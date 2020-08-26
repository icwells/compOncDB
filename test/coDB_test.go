// Performs black box tests on the comparative oncology sql database

package main

import (
	"flag"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"os"
	"strconv"
	"strings"
	"testing"
)

func toFloat(act, exp string) (float64, float64, error) {
	// Converts strings to float
	var e float64
	a, err := strconv.ParseFloat(act, 64)
	if err == nil {
		e, err = strconv.ParseFloat(exp, 64)
	}
	return a, e, err
}

func compareTables(t *testing.T, name string, exp, act *dataframe.Dataframe) {
	// Compares output of equivalent tables
	ac, ar := act.Dimensions()
	ec, er := exp.Dimensions()
	if ac != ec && ar != er {
		t.Errorf("Actual %s dimensions [%d, %d] do not equal expected: [%d, %d]", name, ac, ar, ec, er)
	} else {
		for key := range act.Index {
			for k := range act.Header {
				a, _ := act.GetCell(key, k)
				e, _ := exp.GetCell(key, k)
				if a != e {
					// Make sure error is not due to floating point precision
					af, ef, err := toFloat(a, e)
					if err != nil || af != ef {
						t.Errorf("%s-%s: Actual %s value %s does not equal expected: %s", name, key, k, a, e)
					}
				}
			}
		}
	}
}

//----------------------------------------------------------------------------

func tableToDF(db *dbIO.DBIO, name string) *dataframe.Dataframe {
	// Returns specified table as dataframe
	ret, _ := dataframe.NewDataFrame(0)
	ret.SetHeader(strings.Split(db.Columns[name], ","))
	for _, i := range db.GetTable(name) {
		ret.AddRow(i)
	}
	return ret
}

func TestUpload(t *testing.T) {
	// Compares actual output from table dumps to expected
	flag.Parse()
	exp := getExpectedTables()
	// Get empty database
	c := codbutils.SetConfiguration(*user, true)
	db := dbIO.ReplaceDatabase(c.Host, c.Testdb, *user, *password)
	db.NewTables(c.Tables)
	// Replace column names
	db.GetTableColumns()
	// Upload taxonomy, life history data, denominators
	dbupload.LoadTaxa(db, taxa, true)
	dbupload.LoadLifeHistory(db, lifehistory)
	dbupload.LoadNonCancerTotals(db, denominators)
	// Upload patient data
	dbupload.LoadAccounts(db, uploadfile)
	dbupload.LoadPatients(db, uploadfile, false)
	for k := range db.Columns {
		// Dump all tables for comparison
		if k != "Unmatched" && k != "Update_time" {
			compareTables(t, k, exp[k], tableToDF(db, k))
		}
	}
	// Attempt again to test filtering of existing records
	dbupload.LoadPatients(db, uploadfile, false)
	for k := range db.Columns {
		if k != "Unmatched" && k != "Update_time" {
			compareTables(t, k, exp[k], tableToDF(db, k))
		}
	}
}

func connectToDatabase() *dbIO.DBIO {
	// Manages call to Connect and GetTableColumns
	flag.Parse()
	c := codbutils.SetConfiguration(*user, true)
	db, err := dbIO.Connect(c.Host, c.Testdb, c.User, *password)
	if err != nil {
		os.Exit(1000)
	}
	db.GetTableColumns()
	return db
}

func TestFilterPatients(t *testing.T) {
	// Tests duplicate patient filtering
	db := connectToDatabase()
	exp := getExpectedTables()
	dbupload.LoadPatients(db, uploadfile, false)
	dbupload.LoadPatients(db, uploadfile, true)
	dbextract.AutoCleanDatabase(db)
	for k := range db.Columns {
		if k != "Unmatched" && k != "Update_time" {
			compareTables(t, k, exp[k], tableToDF(db, k))
		}
	}
}

func TestCancerRates(t *testing.T) {
	// Tests taxonomy search output
	var e [][]codbutils.Evaluation
	db := connectToDatabase()
	compareTables(t, "Cancer Rates", getExpectedRates(), dbextract.GetCancerRates(db, 1, false, false, false, e, "species"))
}

func TestSearches(t *testing.T) {
	// Tests taxonomy search output
	db := connectToDatabase()
	cases := newSearchCases(db.Columns)
	for _, i := range cases {
		res, _ := dbextract.SearchColumns(db, i.table, i.eval, false)
		if i.name == "fox" && res.Length() > 0 {
			t.Error("Results returned for gray fox (not present).")
		} else {
			compareTables(t, i.name, i.expected, res)
		}
	}
	// Test searching from file. Given search criteria will only match canis results.
	res, _ := dbextract.SearchDatabase(db, "nil", "nil", searchfile, false)
	compareTables(t, "testSearch", getCanisResults(), res)
}

func TestUpdates(t *testing.T) {
	// Tests dumped tables after update
	db := connectToDatabase()
	exp := getExpectedUpdates()
	dbextract.UpdateEntries(db, updatefile)
	for _, i := range []string{"Patient", "Diagnosis", "Tumor"} {
		compareTables(t, i, exp[i], tableToDF(db, i))
	}
}

func TestDelete(t *testing.T) {
	// Tests delete and outoclean functions
	db := connectToDatabase()
	exp := getCleaned()
	codbutils.DeleteEntries(db, "Patient", "ID", "19")
	dbextract.AutoCleanDatabase(db)
	for k := range exp {
		// Compare all tables to ensure only target data was removed
		compareTables(t, k, exp[k], tableToDF(db, k))
	}
}

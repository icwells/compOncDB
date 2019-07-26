// Performs black box tests on the comparative oncology sql database

package main

import (
	"flag"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"os"
	"strconv"
	//"strings"
	"testing"
)

func compareEntries(actual, expected []string) int {
	// Returns true if both slices are equal
	equal := true
	for idx, i := range actual {
		// Skip randomly assigned IDs
		if idx < len(expected) && i != expected[idx] {
			equal = false
			// Attempt to resolve differences in floating point precision
			a, err := strconv.ParseFloat(i, 64)
			if err == nil {
				var e float64
				e, err = strconv.ParseFloat(expected[idx], 64)
				if err == nil && a == e {
					equal = true
				}
			}
		}
		if equal == false {
			return idx
		}
	}
	return -1
}

func compareTables(t *testing.T, name string, exp, act map[string][]string) {
	// Compares output of equivalent tables
	if len(act) != len(exp) {
		t.Errorf("%s: Actual table length %d does not equal expected: %d", name, len(act), len(exp))
	} else {
		for k := range act {
			if len(act[k]) != len(exp[k]) {
				t.Errorf("%s, %s: Actual length %d does not equal expected: %d", name, k, len(act[k]), len(exp[k]))
				break
			} else {
				idx := compareEntries(act[k], exp[k])
				if idx >= 0 {
					t.Errorf("%s %s-%d: Actual value %s does not equal expected: %s", name, k, idx, act[k][idx], exp[k][idx])
					break
				}
			}
		}
	}
}

//----------------------------------------------------------------------------

func TestUpload(t *testing.T) {
	// Compares actual output from table dumps to expected
	flag.Parse()
	exp := getExpectedTables()
	// Get empty database
	c := codbutils.SetConfiguration(config, *user, true)
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
	dbupload.LoadPatients(db, uploadfile)
	for k := range db.Columns {
		// Dump all tables for comparison
		table := dbupload.ToMap(db.GetTable(k))
		compareTables(t, k, exp[k], table)
	}
}

func connectToDatabase() *dbIO.DBIO {
	// Manages call to Connect and GetTableColumns
	flag.Parse()
	c := codbutils.SetConfiguration(config, *user, true)
	db, err := dbIO.Connect(c.Host, c.Testdb, c.User, *password)
	if err != nil {
		os.Exit(1000)
	}
	db.GetTableColumns()
	return db
}

func TestSearches(t *testing.T) {
	// Tests taxonomy search output
	db := connectToDatabase()
	cases := newSearchCases(db.Columns)
	for _, i := range cases {
		res, _ := dbextract.SearchColumns(db, *user, i.table, i.eval, false, false)
		if i.name == "fox" && len(res) > 0 {
			t.Error("Results returned for gray fox (not present).")
		} else {
			compareTables(t, i.name, i.expected, dbupload.ToMap(res))
		}
	}
}

/*func TestUpdates(t *testing.T) {
	// Tests dumped tables after update
	fmt.Print("\n\tTesting update functions...\n\n")
	db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, true))
	dbextract.UpdateEntries(db, *infile)
	for _, i := range []string{"Patient", "Diagnosis"} {
		table := db.GetTable(i)
		out := fmt.Sprintf("%s%s.csv", *outfile, i)
		iotools.WriteToCSV(out, db.Columns[i], table)
	}

	flag.Parse()
	*indir, _ = iotools.FormatPath(*indir, false)
	files, err := filepath.Glob(*indir + "*.csv")
	if err != nil {
		t.Errorf("Cannot find test files in %s: %v", *indir, err)
	}
	expected := sortInput(files, true)
	actual := sortInput(files, false)
	for k, v := range expected {
		act, ex := actual[k]
		if ex == false {
			t.Errorf("Actual search result %s not found.", k)
		} else {
			compareTables(t, k, v, act)
			// Remove test output
			os.Remove(act)
		}
	}
}*/

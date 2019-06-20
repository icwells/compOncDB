// Assigns input variables to appriate worker functions

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"os"
	"os/exec"
	"path"
	"time"
)

func backup(pw string) {
	// Backup database to local machine
	c := codbutils.SetConfiguration(*config, *user, false)
	fmt.Printf("\n\tBacking up %s database to local machine...\n", c.Database)
	datestamp := time.Now().Format("2006-01-02")
	password := fmt.Sprintf("-p%s", pw)
	host := fmt.Sprintf("-h%s", c.Host)
	res := fmt.Sprintf("--result-file=%s.%s.sql", c.Database, datestamp)
	bu := exec.Command("mysqldump", "-uroot", host, password, res, c.Database)
	err := bu.Run()
	if err == nil {
		fmt.Println("\tBackup complete.")
	} else {
		fmt.Printf("\tBackup failed. %v\n", err)
	}
}

func newDatabase() time.Time {
	// Creates new database and tables
	c := codbutils.SetConfiguration(*config, *user, false)
	db := dbIO.CreateDatabase(c.Host, c.Database, *user)
	db.NewTables(c.Tables)
	return db.Starttime
}

func uploadToDB() time.Time {
	// Uploads infile to given table (all input variables are global)
	if *infile == "nil" {
		fmt.Print("\n\t[Error] Please specify input file. Exiting.\n\n")
		os.Exit(1)
	}
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, false))
	if *taxa == true {
		// Upload taxonomy
		dbupload.LoadTaxa(db, *infile, *common)
	} else if *lh == true {
		// Upload life history table
		dbupload.LoadLifeHistory(db, *infile)
	} else if *den == true {
		// Uplaod denominator table
		dbupload.LoadNonCancerTotals(db, *infile)
	} else if *patient == true {
		// Upload patient data
		dbupload.LoadAccounts(db, *infile)
		dbupload.LoadPatients(db, *infile)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	return db.Starttime
}

func updateDB() time.Time {
	// Updates database with given flags (all input variables are global)
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, false))
	if *total == true {
		dbupload.SpeciesTotals(db)
	} else if *infile != "nil" {
		dbextract.UpdateEntries(db, *infile)
	} else if *column != "nil" && *value != "nil" && *eval != "nil" {
		evaluations := codbutils.SetOperations(db.Columns, *eval)
		e := evaluations[0]
		dbextract.UpdateSingleTable(db, e.Table, *column, *value, e.Column, e.Operator, e.Value)
	} else if *del == true && *eval != "nil" {
		var tables []string
		evaluations := codbutils.SetOperations(db.Columns, *eval)
		e := evaluations[0]
		if *table != "nil" {
			tables = []string{*table}
		} else {
			tables = codbutils.GetTable(db.Columns, e.Column)
		}
		codbutils.DeleteEntries(db, tables, e.Column, e.Value)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	return db.Starttime
}

func extractFromDB() time.Time {
	// Extracts data to outfile/stdout (all input variables are global)
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, false))
	if *dump != "nil" {
		// Extract entire table
		table := db.GetTable(*dump)
		codbutils.WriteResults(*outfile, db.Columns[*dump], table)
	} else if *sum == true {
		summary := dbextract.GetSummary(db)
		codbutils.WriteResults(*outfile, "Field,Total,%\n", summary)
	} else if *cr == true {
		// Extract cancer rates
		header := "Kingdom,Phylum,Class,Orders,Family,Genus,ScientificName,TotalRecords,CancerRecords,CancerRate,"
		header += "AverageAge(months),AvgAgeCancer(months),Male,Female,MaleCancer,FemaleCancer"
		rates := dbextract.GetCancerRates(db, *min, *nec)
		codbutils.WriteResults(*outfile, header, rates)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	return db.Starttime
}

func searchDB() time.Time {
	// Performs search functions on database
	var res [][]string
	var header string
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, false))
	if *eval != "nil" {
		// Search for column/value match
		e := codbutils.SetOperations(db.Columns, *eval)
		res, header = dbextract.SearchColumns(db, *user, *table, e, *count, *infant)
		fmt.Printf("\tFound %d records where %s is %s.\n", len(res), e[0].Column, e[0].Value)
	} else if *taxonomies == true {
		names := codbutils.ReadList(*infile, *col)
		res, header = dbextract.SearchSpeciesNames(db, names)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	if *count == false && len(res) >= 1 {
		codbutils.WriteResults(*outfile, header, res)
	}
	return db.Starttime
}

func testDB() time.Time {
	// Performs test uploads and extractions
	var db *dbIO.DBIO
	if *testsearch == true {
		var terms searchterms
		fmt.Print("\n\tTesting search functions...\n\n")
		db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, true))
		terms.readSearchTerms(db.Columns, *infile, *outfile)
		terms.searchTestCases(db)
	} else if *updates == true {
		fmt.Print("\n\tTesting update functions...\n\n")
		db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, true))
		dbextract.UpdateEntries(db, *infile)
		for _, i := range []string{"Patient", "Diagnosis"} {
			table := db.GetTable(i)
			out := fmt.Sprintf("%s%s.csv", *outfile, i)
			iotools.WriteToCSV(out, db.Columns[i], table)
		}
	} else {
		// Get empty database
		bin, _ := path.Split(*config)
		c := codbutils.SetConfiguration(*config, *user, true)
		db = dbIO.ReplaceDatabase(c.Host, c.Testdb, *user)
		db.NewTables(path.Join(bin, c.Tables))
		// Replace column names
		db.GetTableColumns()
		// Upload taxonomy
		dbupload.LoadTaxa(db, *taxafile, true)
		dbupload.LoadLifeHistory(db, *lifehistory)
		// Uplaod denominator table
		dbupload.LoadNonCancerTotals(db, *noncancer)
		// Upload patient data
		dbupload.LoadAccounts(db, *infile)
		dbupload.LoadPatients(db, *infile)
		fmt.Print("\n\tDumping test tables...\n\n")
		for k := range db.Columns {
			// Dump all tables for comparison
			table := db.GetTable(k)
			out := fmt.Sprintf("%s%s.csv", *outfile, k)
			iotools.WriteToCSV(out, db.Columns[k], table)
		}
	}
	return db.Starttime
}

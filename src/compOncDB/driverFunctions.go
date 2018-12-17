// Assigns input variables to appriate worker functions

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"os"
	"os/exec"
	"strings"
	"time"
)

func backup(pw string) {
	// Backup database to local machine
	fmt.Printf("\n\tBacking up %s database to local machine...\n", DB)
	datestamp := time.Now().Format("2006-01-02")
	password := fmt.Sprintf("-p%s", pw)
	res := fmt.Sprintf("--result-file=%s.%s.sql", DB, datestamp)
	bu := exec.Command("mysqldump", "-uroot", password, res, DB)
	err := bu.Run()
	if err == nil {
		fmt.Println("\tBackup complete.")
	} else {
		fmt.Printf("\tBackup failed. %v\n", err)
	}
}

func connectToDatabase(testdb bool) *dbIO.DBIO {
	// Manages call to Connect and ReadColumns
	d := DB
	if testdb == true {
		d = TDB
	}
	db := dbIO.Connect(d, *user)
	if testdb == false {
		db.ReadColumns(COL, false)
	}
	return db
}

func uploadToDB() time.Time {
	// Uploads infile to given table (all input variables are global)
	if *infile == "nil" {
		fmt.Print("\n\t[Error] Please specify input file. Exiting.\n\n")
		os.Exit(1)
	}
	db := connectToDatabase(false)
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
		dbupload.LoadDiagnoses(db, *infile)
		dbupload.LoadPatients(db, *infile)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	return db.Starttime
}

func updateDB() time.Time {
	// Updates database with given flags (all input variables are global)
	db := connectToDatabase(false)
	if *total == true {
		dbupload.SpeciesTotals(db)
	} else if *infile != "nil" {
		dbextract.UpdateEntries(db, *infile)
	} else if *column != "nil" && *value != "nil" && *eval != "nil" {
		col, op, val := getOperation(*eval)
		tables := getTable(db.Columns, col)
		dbextract.UpdateSingleTable(db, tables[0], *column, *value, col, op, val)
	} else if *del == true && *eval != "nil" {
		if *user == "root" {
			column, _, value := getOperation(*eval)
			tables := getTable(db.Columns, column)
			deleteEntries(db, tables, column, value)
		} else {
			fmt.Print("\n\t[Error] Must be root to delete entries. Exiting.\n\n")
		}
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	return db.Starttime
}

func extractFromDB() time.Time {
	// Extracts data to outfile/stdout (all input variables are global)
	db := connectToDatabase(false)
	if *dump != "nil" {
		if *dump == "Accounts" && *user != "root" {
			fmt.Print("\n\t[Error] Must be root to access Accounts table. Exiting.\n\n")
			os.Exit(1010)
		}
		// Extract entire table
		table := db.GetTable(*dump)
		if *outfile != "nil" {
			iotools.WriteToCSV(*outfile, db.Columns[*dump], table)
		} else {
			printArray(db.Columns[*dump], table)
		}
	} else if *sum == true {
		summary := dbextract.GetSummary(db)
		writeResults(*outfile, "Field,Total,%\n", summary)
	} else if *cr == true {
		// Extract cancer rates
		header := "ScientificName,TotalRecords,CancerRecords,CancerRate,AverageAge(months),AvgAgeCancer(months),Male,Female,MaleCancer,FemaleCancer"
		rates := dbextract.GetCancerRates(db, *min, *nec)
		writeResults(*outfile, header, rates)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	return db.Starttime
}

func searchDB() time.Time {
	// Performs search functions on database
	var res [][]string
	var header string
	db := connectToDatabase(false)
	if *taxon != "nil" {
		// Extract all data for a given species
		var names []string
		if iotools.Exists(*taxon) == true {
			names = readList(*taxon)
		} else {
			// Get single term
			if strings.Contains(*taxon, "_") == true {
				names = []string{strings.Replace(*taxon, "_", " ", -1)}
			} else {
				names = []string{*taxon}
			}
		}
		res, header = dbextract.SearchTaxonomicLevels(db, names, *user, *level, *count, *com)
		fmt.Printf("\tFound %d records where %s is %s.\n", len(res), *level, *taxon)
	} else if *eval != "nil" {
		// Search for column/value match
		column, op, value := getOperation(*eval)
		if *table == "nil" {
			tables := getTable(db.Columns, column)
			res, header = dbextract.SearchColumns(db, tables, *user, column, op, value, *count, *com)
		} else {
			res, header = dbextract.SearchSingleTable(db, *table, *user, column, op, value, *com)
		}
		fmt.Printf("\tFound %d records where %s is %s.\n", len(res), column, value)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	if *count == false && len(res) >= 1 {
		writeResults(*outfile, header, res)
	}
	return db.Starttime
}

func testDB() time.Time {
	// Performs test uploads and extractions
	db := connectToDatabase(true)
	db.NewTables(*tables)
	// Re-read columns without types
	db.ReadColumns(*tables, false)
	if *testsearch == false {
		// Clear existing tables
		for k := range db.Columns {
			db.TruncateTable(k)
		}
		// Upload taxonomy
		dbupload.LoadTaxa(db, *taxafile, true)
		dbupload.LoadLifeHistory(db, *lifehistory)
		// Uplaod denominator table
		dbupload.LoadNonCancerTotals(db, *noncancer)
		// Upload patient data
		dbupload.LoadAccounts(db, *infile)
		dbupload.LoadDiagnoses(db, *infile)
		dbupload.LoadPatients(db, *infile)
		fmt.Print("\n\tDumping test tables...\n\n")
		for k := range db.Columns {
			// Dump all tables for comparison
			table := db.GetTable(k)
			out := fmt.Sprintf("%s%s.csv", *outfile, k)
			iotools.WriteToCSV(out, db.Columns[k], table)
		}
	} else {
		fmt.Print("\n\tTesting search functions...\n\n")
		var terms searchterms
		terms.readSearchTerms(*infile, *outfile)
		terms.searchTestCases(db)
	}
	return db.Starttime
}

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
	"path"
	"strings"
	"time"
)

func backup(pw string) {
	// Backup database to local machine
	c := setConfiguration(false)
	fmt.Printf("\n\tBacking up %s database to local machine...\n", c.database)
	datestamp := time.Now().Format("2006-01-02")
	password := fmt.Sprintf("-p%s", pw)
	host := fmt.Sprintf("-h%s", c.host)
	res := fmt.Sprintf("--result-file=%s.%s.sql", c.database, datestamp)
	bu := exec.Command("mysqldump", "-uroot", host, password, res, c.database)
	err := bu.Run()
	if err == nil {
		fmt.Println("\tBackup complete.")
	} else {
		fmt.Printf("\tBackup failed. %v\n", err)
	}
}

func newDatabase() time.Time {
	// Creates new database and tables
	c := setConfiguration(false)
	db := dbIO.CreateDatabase("", c.database, *user)
	db.NewTables(c.tables)
	return db.Starttime
}

func uploadToDB() time.Time {
	// Uploads infile to given table (all input variables are global)
	if *infile == "nil" {
		fmt.Print("\n\t[Error] Please specify input file. Exiting.\n\n")
		os.Exit(1)
	}
	db := connectToDatabase(setConfiguration(false))
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
	db := connectToDatabase(setConfiguration(false))
	if *total == true {
		dbupload.SpeciesTotals(db)
	} else if *infile != "nil" {
		dbextract.UpdateEntries(db, *infile)
	} else if *column != "nil" && *value != "nil" && *eval != "nil" {
		evaluations := setOperations(*eval)
		e := evaluations[0]
		tables := getTable(db.Columns, e.column)
		dbextract.UpdateSingleTable(db, tables[0], *column, *value, e.column, e.operator, e.value)
	} else if *del == true && *eval != "nil" {
		var tables []string
		evaluations := setOperations(*eval)
		e := evaluations[0]
		if *table != "nil" {
			tables = []string{*table}
		} else {
			tables = getTable(db.Columns, e.column)
		}
		deleteEntries(db, tables, e.column, e.value)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	return db.Starttime
}

func extractFromDB() time.Time {
	// Extracts data to outfile/stdout (all input variables are global)
	db := connectToDatabase(setConfiguration(false))
	if *dump != "nil" {
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
		header := "Kingdom,Phylum,Class,Orders,Family,Genus,ScientificName,TotalRecords,CancerRecords,CancerRate,"
		header += "AverageAge(months),AvgAgeCancer(months),Male,Female,MaleCancer,FemaleCancer"
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
	db := connectToDatabase(setConfiguration(false))
	if *taxon != "nil" {
		// Extract all data for a given species
		var names []string
		if iotools.Exists(*taxon) == true {
			names = readList(*taxon, *col)
		} else {
			// Get single term
			if strings.Contains(*taxon, "_") == true {
				names = []string{strings.Replace(*taxon, "_", " ", -1)}
			} else {
				names = []string{*taxon}
			}
		}
		res, header = dbextract.SearchTaxonomicLevels(db, names, *user, *level, *count, *com, *infant)
		fmt.Printf("\tFound %d records where %s is %s.\n", len(res), *level, *taxon)
	} else if *eval != "nil" {
		// Search for column/value match
		e := setOperations(*eval)
		if *table == "nil" {
			tables := getTable(db.Columns, e[0].column)
			res, header = dbextract.SearchColumns(db, tables, *user, e[0].column, e[0].operator, e[0].value, *count, *com, *infant)
		} else {
			res, header = dbextract.SearchSingleTable(db, *table, *user, e[0].column, e[0].operator, e[0].value, *com, *infant)
		}
		fmt.Printf("\tFound %d records where %s is %s.\n", len(res), e[0].column, e[0].value)
		if len(e) > 1 {
			res = filterSearchResults(header, e[1:], res)
		}
	} else if *taxonomies == true {
		names := readList(*infile, *col)
		res, header = dbextract.SearchSpeciesNames(db, names)
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
	var db *dbIO.DBIO
	if *testsearch == true {
		var terms searchterms
		fmt.Print("\n\tTesting search functions...\n\n")
		db = connectToDatabase(setConfiguration(true))
		terms.readSearchTerms(*infile, *outfile)
		terms.searchTestCases(db)
	} else if *updates == true {
		fmt.Print("\n\tTesting update functions...\n\n")
		db = connectToDatabase(setConfiguration(true))
		dbextract.UpdateEntries(db, *infile)
		for _, i := range []string{"Patient", "Diagnosis"} {
			table := db.GetTable(i)
			out := fmt.Sprintf("%s%s.csv", *outfile, i)
			iotools.WriteToCSV(out, db.Columns[i], table)
		}
	} else {
		// Get empty database
		bin, _ := path.Split(*config)
		c := setConfiguration(true)
		db = dbIO.ReplaceDatabase(c.host, c.testdb, *user)
		db.NewTables(path.Join(bin, c.tables))
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

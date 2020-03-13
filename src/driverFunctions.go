// Assigns input variables to appriate worker functions

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/compOncDB/src/parserecords"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"os"
	"os/exec"
	"time"
)

func backup(pw string) {
	// Backup database to local machine
	c := codbutils.SetConfiguration(*config, *user, false)
	fmt.Printf("\n\tBacking up %s database to local machine...\n", c.Database)
	datestamp := time.Now().Format("2006-01-02")
	user := fmt.Sprintf("-u%s", *user)
	password := fmt.Sprintf("-p%s", pw)
	host := fmt.Sprintf("-h%s", c.Host)
	res := fmt.Sprintf("--result-file=%s.%s.sql", c.Database, datestamp)
	bu := exec.Command("mysqldump", user, host, password, res, c.Database)
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

func parseRecords() time.Time {
	// Parses raw input for unpload to database
	start := time.Now()
	fmt.Print("\n\tProcessing input records...\n")
	ent := parserecords.NewEntries(*service, *infile)
	ent.GetTaxonomy(*taxaFile)
	ent.SortRecords(*debug, *infile, *outfile)
	return start
}

func commandError() {
	// Prints message for invalid input
	fmt.Print("\n\tPlease enter a valid command.\n\n")
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
		codbutils.UpdateTimeStamp(db)
	} else if *lh == true {
		// Upload life history table
		dbupload.LoadLifeHistory(db, *infile)
		codbutils.UpdateTimeStamp(db)
	} else if *den == true {
		// Uplaod denominator table
		dbupload.LoadNonCancerTotals(db, *infile)
		codbutils.UpdateTimeStamp(db)
	} else if *patient == true {
		// Upload patient data
		dbupload.LoadAccounts(db, *infile)
		dbupload.LoadPatients(db, *infile, false)
		codbutils.UpdateTimeStamp(db)
	} else {
		commandError()
	}
	return db.Starttime
}

func updateDB() time.Time {
	// Updates database with given flags (all input variables are global)
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, false))
	if *clean == true {
		dbextract.AutoCleanDatabase(db)
		codbutils.UpdateTimeStamp(db)
	} else if *infile != "nil" {
		dbextract.UpdateEntries(db, *infile)
		codbutils.UpdateTimeStamp(db)
	} else if *column != "nil" && *value != "nil" && *eval != "nil" {
		evaluations := codbutils.SetOperations(db.Columns, *eval)
		e := evaluations[0]
		dbextract.UpdateSingleTable(db, e.Table, *column, *value, e.Column, e.Operator, e.Value)
		codbutils.UpdateTimeStamp(db)
	} else if *del == true && *eval != "nil" {
		evaluations := codbutils.SetOperations(db.Columns, *eval)
		e := evaluations[0]
		if *table == "nil" {
			*table = codbutils.GetTable(db.Columns, e.Column)
		}
		codbutils.DeleteEntries(db, *table, e.Column, e.Value)
		codbutils.UpdateTimeStamp(db)
	} else {
		commandError()
	}
	return db.Starttime
}

func writeDF(table *dataframe.Dataframe) {
	// Writes dataframe to file/screen
	if *count == false && table.Length() >= 1 {
		if *outfile != "nil" {
			table.ToCSV(*outfile)
		} else {
			table.Print()
		}
	}
}

func extractFromDB() time.Time {
	// Extracts data to outfile/stdout (all input variables are global)
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, false))
	if *dump != "nil" {
		// Extract entire table
		table := db.GetTable(*dump)
		codbutils.WriteResults(*outfile, db.Columns[*dump], table)
	} else if *sum {
		summary := dbextract.GetSummary(db)
		codbutils.WriteResults(*outfile, "Field,Total,%\n", summary)
	} else if *cr {
		// Extract cancer rates
		var e []codbutils.Evaluation
		if *eval != "nil" {
			e = codbutils.SetOperations(db.Columns, *eval)
		}
		writeDF(dbextract.GetCancerRates(db, *min, *nec, *lifehist, e))
	} else if *reftaxa {
		writeDF(dbextract.GetReferenceTaxonomy(db))
	} else {
		commandError()
	}
	return db.Starttime
}

func searchDB() time.Time {
	// Performs search functions on database
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, false))
	if *eval != "nil" {
		// Search for column/value match
		e := codbutils.SetOperations(db.Columns, *eval)
		res, msg := dbextract.SearchColumns(db, *table, e, *count, *infant)
		fmt.Print(msg)
		writeDF(res)
	} else if *taxonomies == true {
		names := codbutils.ReadList(*infile, *col)
		writeDF(dbextract.SearchSpeciesNames(db, names))
	} else {
		commandError()
	}
	return db.Starttime
}

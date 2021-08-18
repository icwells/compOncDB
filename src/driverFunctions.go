// Assigns input variables to appriate worker functions

package main

import (
	"bufio"
	"fmt"
	"github.com/icwells/compOncDB/src/cancerrates"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/compOncDB/src/parserecords"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"os"
	"strings"
	"time"
)

func newDatabase() time.Time {
	// Creates new database and tables
	c := codbutils.SetConfiguration(*user, false)
	db := dbIO.CreateDatabase(c.Host, c.Database, *user)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n\tAre you sure you want to initialize a new database? This will erase existing data.")
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(text)
	if text == "y" || text == "yes" {
		fmt.Println("\tInitializing new tables...")
		db.NewTables(c.Tables)
	}
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
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)
	if len(*outfile) > 0 {
		db.BackupDB(*outfile)
	}
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
		dbupload.LoadPatients(db, *infile, false, false)
		codbutils.UpdateTimeStamp(db)
	} else {
		commandError()
	}
	return db.Starttime
}

func writeDF(table *dataframe.Dataframe, output string) {
	// Writes dataframe to file/screen
	if table.Length() >= 1 {
		if output != "nil" && output != "" {
			table.ToCSV(output)
		} else {
			fmt.Println()
			table.Print()
		}
	}
}

func updateDB() time.Time {
	// Updates database with given flags (all input variables are global)
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)
	if len(*outfile) > 0 {
		db.BackupDB(*outfile)
	}
	if *clean {
		dbextract.AutoCleanDatabase(db)
		codbutils.UpdateTimeStamp(db)
	} else if *infile != "nil" {
		dbextract.UpdateEntries(db, *infile)
		codbutils.UpdateTimeStamp(db)
	} else if *column != "nil" && *value != "nil" && *eval != "nil" {
		evaluations := codbutils.SetOperations(db.Columns, *eval)
		e := evaluations[0][0]
		dbextract.UpdateSingleTable(db, e.Table, *column, *value, e.Column, e.Operator, e.Value)
		codbutils.UpdateTimeStamp(db)
	} else if *del && *eval != "nil" {
		evaluations := codbutils.SetOperations(db.Columns, *eval)
		e := evaluations[0][0]
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

func calculateCancerRates() time.Time {
	// Extract cancer rates
	*nec--
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)
	if *pathology {
		prevalence, reports := cancerrates.GetRatesAndRecords(db, *min, *nec, *infant, *lifehist, *keepall, *source, *eval, *location)
		writeDF(prevalence, *outfile)
		writeDF(reports, strings.Replace(*outfile, ".csv", ".Pathology.csv", 1))
	} else {
		writeDF(cancerrates.GetCancerRates(db, *min, *nec, *infant, *lifehist, *wild, *keepall, *source, *eval, *location), *outfile)
	}
	return db.Starttime
}

func searchDB() time.Time {
	// Searches db with given queries
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)
	if *taxonomies == true {
		names := codbutils.ReadList(*infile, *col)
		writeDF(search.SearchSpeciesNames(db, names), *outfile)
	} else if *eval != "nil" || *infile != "nil" {
		// Search for column/value match
		res, msg := search.SearchRecords(db, codbutils.GetLogger(), *eval, *infant, false)
		if msg != "" {
			fmt.Print(msg)
			writeDF(res, *outfile)
		}
	} else if *top {
		writeDF(search.LeaderBoard(db), *outfile)
		codbutils.UpdateTimeStamp(db)
	} else {
		commandError()
	}
	return db.Starttime
}

func extractFromDB() time.Time {
	// Extracts data to outfile/stdout (all input variables are global)
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)
	if *dump != "nil" {
		// Extract entire table
		table := db.GetTable(*dump)
		codbutils.WriteResults(*outfile, db.Columns[*dump], table)
	} else if *dumpdb {
		dbextract.DumpDatabase(db, *outfile)
	} else if *lhsummary {
		writeDF(dbextract.LifeHistorySummary(db, *alltaxa), *outfile)
	} else if *reftaxa {
		writeDF(dbextract.GetReferenceTaxonomy(db), *outfile)
	} else if *sum {
		summary := dbextract.GetSummary(db)
		codbutils.WriteResults(*outfile, "Field,Total,%\n", summary)
	} else {
		commandError()
	}
	return db.Starttime
}

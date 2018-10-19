// This script will manage searching of the comparative oncology database

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	// Global variables
	COL = "tableColumns.txt"
	DB  = "comparativeOncology"
)

var (
	// Kingpin arguments
	app      = kingpin.New("compOncDB", "Comand line-interface for uploading/extrating/manipulating data from the comparative oncology database.")
	user     = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()
	ver      = kingpin.Command("version", "Prints version info and exits.")
	bu       = kingpin.Command("backup", "Backs up database to local machine (Must use root password; output is written to current directory).")
	New      = kingpin.Command("new", "Initializes new tables in new database (database must be initialized manually).")
	column	 = kingpin.Flag("column", "Name of column containing target value (table is automatically determined).").Short('c').Default("nil").String()
	value	 = kingpin.Flag("value", "Name of target value to update.").Short('v').Default("nil").String()
	infile   = kingpin.Flag("infile", "Path to input file (if using).").Short('i').Default("nil").String()
	outfile  = kingpin.Flag("outfile", "Name of output file (writes to stdout if not given).").Short('o').Default("nil").String()

	upload   = kingpin.Command("upload", "Upload data to the database.")
	taxa     = upload.Flag("taxa", "Load taxonomy tables from Kestrel output to update taxonomy table.").Default("false").Bool()
	common   = upload.Flag("common", "Additionally extract common names from Kestrel output to update common name tables.").Default("false").Bool()
	lh       = upload.Flag("lh", "Upload life history info from merged life history table to the database.").Default("false").Bool()
	den      = upload.Flag("den", "Uploads file to denominator table for databases where only cancer records were extracted.").Default("false").Bool()
	patient  = upload.Flag("patient", "Upload patient, account, and diagnosis info from input table to database.").Default("false").Bool()

	update   = kingpin.Command("update", "Update or delete existing records from the database.")
	total	 = update.Flag("count", "Recount species totals and update the Totals table.").Default("false").Bool()
	del		 = update.Flag("delete", "Delete records if column = value.").Default("false").Bool()

	extract  = kingpin.Command("extract", "Extract data from the database and perform optional analyses.")
	dump     = extract.Flag("dump", "Name of table to dump (writes all data from table to output file).").Short('d').Default("nil").String()
	cr       = extract.Flag("cancerRate", "Calculates cancer rates for species with greater than min entries.").Default("false").Bool()
	min      = extract.Flag("min", "Minimum number of entries required for calculations (default = 50).").Short('m').Default("50").Int()
	nec      = extract.Flag("necropsy", "Extract only necropsy records (extracts all matches by default).").Default("false").Bool()

	txn      = "Name of taxonomic unit to extract data for or path to file with single column of units."
	search	 = kingpin.Command("search", "Searches database for matches to given term.")
	taxon	 = search.Flag("taxa", txn).Short('t').Default("nil").String()
	level	 = search.Flag("level", "Taxonomic level of taxon (or entries in taxon file)(default = Species).").Short('l').Default("Species").String()
	com		 = search.Flag("common", "Indicates that common species name was given for taxa.").Default("false").Bool()
	count	 = search.Flag("count", "Returns count of target records instead of printing entire records.").Default("false").Bool()
	eval	 = search.Flag("eval", "Searches life history or totals tables for matches (column operator value; valid operators: <= >= > <).").Default("nil").String()
	table	 = search.Flag("table", "Return matching rows from this table only.").Default("nil").String()
)

func version() {
	fmt.Println("\n\tCompOncDB v0.1 (~) is a package for managing the ASU comparative oncology database.")
	fmt.Println("\n\tCopyright 2018 by Shawn Rupp, Maley Lab, Biodesign Institute, Arizona State University.")
	fmt.Println("\tThis program comes with ABSOLUTELY NO WARRANTY.")
	fmt.Println("\n\tThis is free software, and you are welcome to redistribute it under certain conditions.\n")
	os.Exit(0)
}

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

func uploadToDB() time.Time {
	// Uploads infile to given table (all input variables are global)
	if *infile == "nil" {
		fmt.Println("\n\t[Error] Please specify input file. Exiting.\n")
		os.Exit(1)
	}
	db, _, start := dbIO.Connect(DB, *user)
	col := dbIO.ReadColumns(COL, false)
	defer db.Close()
	if *taxa == true {
		// Upload taxonomy
		loadTaxa(db, col, *infile, *common)
	} else if *lh == true {
		// Upload life history table
		loadLifeHistory(db, col, *infile)
	} else if *den == true {
		// Uplaod denominator table
		loadNonCancerTotals(db, col, *infile)
	} else if *patient == true {
		// Upload patient data
		loadAccounts(db, col, *infile)
		loadDiagnoses(db, col, *infile)
		loadPatients(db, col, *infile)
	} else {
		fmt.Println("\n\tPlease enter a valid command.\n")
	}
	return start
}

func updateDB() time.Time {
	// Updates database with given flags (all input variables are global)
	db, _, start := dbIO.Connect(DB, *user)
	defer db.Close()
	col := dbIO.ReadColumns(COL, false)
	if *total == true {
		speciesTotals(db, col)
	} else if *del == true && *column != "nil" && *value != "nil" {
		if *user == "root" {
			tables := getTable(col, *column)
			deleteEntries(db, col, tables, *column, *value)
		} else {
			fmt.Println("\n\t[Error] Must be root to delete entries. Exiting.\n")
		}
	} else {
		fmt.Println("\n\t[Warning] Update functionality not yet complete.")
		//fmt.Println("\n\tPlease enter a valid command.\n")
	}
	return start
}

func extractFromDB() time.Time {
	// Extracts data to outfile/stdout (all input variables are global)
	col := dbIO.ReadColumns(COL, false)
	db, _, start := dbIO.Connect(DB, *user)
	defer db.Close()	
	if *dump != "nil" {
		if *dump == "Accounts" && *user != "root" {
			fmt.Println("\n\t[Error] Must be root to access Accounts table. Exiting.\n")
			os.Exit(1010)
		}
		// Extract entire table
		table := dbIO.GetTable(db, *dump)
		if *outfile != "nil" {
			iotools.WriteToCSV(*outfile, col[*dump], table)
		} else {
			printArray(col[*dump], table)
		}
	} else if *cr == true {
		// Extract cancer rates
		header := "ScientificName,TotalRecords,CancerRecords,CancerRate,AverageAge(months),AvgAgeCancer(months),Male,Female,MaleCancer,FemaleCancer"
		rates := getCancerRates(db, col, *min, *nec)
		writeResults(*outfile, header, rates)
	} else {
		fmt.Println("\n\tPlease enter a valid command.\n")
	}
	return start
}

func searchDB() time.Time {
	// Performs search functions on database
	var res [][]string
	var header string
	col := dbIO.ReadColumns(COL, false)
	db, _, start := dbIO.Connect(DB, *user)
	defer db.Close()
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
		res, header = searchTaxonomicLevels(db, col, names)	
		fmt.Printf("\tFound %d records where %s is %s.\n", len(res), *level, *taxon)
	} else if *column != "nil" && *value != "nil" {
		// Search for column/value match
		if *table == "nil" {
			tables := getTable(col, *column)
			res, header = searchColumns(db, col, tables)
		} else {
			res, header = searchSingleTable(db, col)
		}
		fmt.Printf("\tFound %d records where %s is %s.\n", len(res), *column, *value)
	} else {
		fmt.Println("\n\tPlease enter a valid command.\n")
	}
	if *count == false && len(res) >= 1 {
		writeResults(*outfile, header, res)
	}
	return start
}

func main() {
	var db *sql.DB
	var start time.Time
	var pw string
	switch kingpin.Parse() {
		case ver.FullCommand():
			version()
		case bu.FullCommand():
			db, pw, start = dbIO.Connect(DB, *user)
			defer db.Close()
			backup(pw)
		case New.FullCommand():
			db, _, start = dbIO.Connect(DB, *user)
			defer db.Close()
			dbIO.NewTables(db, COL)
		case upload.FullCommand():
			start = uploadToDB()
		case update.FullCommand():
			start = updateDB()
		case extract.FullCommand():
			start = extractFromDB()
		case search.FullCommand():
			start = searchDB()
	}
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

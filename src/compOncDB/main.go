// This script will manage searching of the comparative oncology database

package main

import (
	"fmt"
	"github.com/icwells/dbIO"
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
	TDB = "testDataBase"
)

var (
	// Kingpin arguments
	app     = kingpin.New("compOncDB", "Comand line-interface for uploading/extrating/manipulating data from the comparative oncology database.")
	user    = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()
	ver     = kingpin.Command("version", "Prints version info and exits.")
	bu      = kingpin.Command("backup", "Backs up database to local machine (Must use root password; output is written to current directory).")
	New     = kingpin.Command("new", "Initializes new tables in new database (database must be initialized manually).")
	eval    = search.Flag("eval", "Searches tables for matches (table is automatically determined) (column operator value; valid operators: = <= >= > <). ").Default("nil").String()
	infile  = kingpin.Flag("infile", "Path to input file (if using).").Short('i').Default("nil").String()
	outfile = kingpin.Flag("outfile", "Name of output file (writes to stdout if not given).").Short('o').Default("nil").String()

	upload  = kingpin.Command("upload", "Upload data to the database.")
	taxa    = upload.Flag("taxa", "Load taxonomy tables from Kestrel output to update taxonomy table.").Default("false").Bool()
	common  = upload.Flag("common", "Additionally extract common names from Kestrel output to update common name tables.").Default("false").Bool()
	lh      = upload.Flag("lh", "Upload life history info from merged life history table to the database.").Default("false").Bool()
	den     = upload.Flag("den", "Uploads file to denominator table for databases where only cancer records were extracted.").Default("false").Bool()
	patient = upload.Flag("patient", "Upload patient, account, and diagnosis info from input table to database.").Default("false").Bool()

	update = kingpin.Command("update", "Update or delete existing records from the database.")
	total  = update.Flag("count", "Recount species totals and update the Totals table.").Default("false").Bool()
	del    = update.Flag("delete", "Delete records if column = value.").Default("false").Bool()

	extract = kingpin.Command("extract", "Extract data from the database and perform optional analyses.")
	dump    = extract.Flag("dump", "Name of table to dump (writes all data from table to output file).").Short('d').Default("nil").String()
	cr      = extract.Flag("cancerRate", "Calculates cancer rates for species with greater than min entries.").Default("false").Bool()
	min     = extract.Flag("min", "Minimum number of entries required for calculations (default = 50).").Short('m').Default("50").Int()
	nec     = extract.Flag("necropsy", "Extract only necropsy records (extracts all matches by default).").Default("false").Bool()

	search = kingpin.Command("search", "Searches database for matches to given term.")
	taxon  = search.Flag("taxa", "Name of taxonomic unit to extract data for or path to file with single column of units.").Short('t').Default("nil").String()
	level  = search.Flag("level", "Taxonomic level of taxon (or entries in taxon file)(default = Species).").Short('l').Default("Species").String()
	com    = search.Flag("common", "Indicates that common species name was given for taxa.").Default("false").Bool()
	count  = search.Flag("count", "Returns count of target records instead of printing entire records.").Default("false").Bool()
	table  = search.Flag("table", "Return matching rows from this table only.").Default("nil").String()

	test 		= kingpin.Command("test", "Tests database functionality using testDataBase instead of comaprative oncology.")
	tables		= test.Flag("tables", "Path tableColumns.txt file.").String()
	taxafile	= test.Flag("taxonomy", "Path to taxonomy file.").String()
	diagnosis	= test.Flag("diagnosis", "Path to extracted diganoses file.").String()
	lifehistory	= test.Flag("lifehistory", "Path to life history data.").String()
	noncancer	= test.Flag("denominators", "Path to file conataining non-cancer totals.").String()
	testsearch	= test.Flag("search", "Search for matches using above commands.").Default("false").Bool()
)

func version() {
	fmt.Println("\n\tCompOncDB v0.1 (~) is a package for managing the ASU comparative oncology database.")
	fmt.Println("\n\tCopyright 2018 by Shawn Rupp, Maley Lab, Biodesign Institute, Arizona State University.")
	fmt.Println("\tThis program comes with ABSOLUTELY NO WARRANTY.")
	fmt.Print("\n\tThis is free software, and you are welcome to redistribute it under certain conditions.\n\n")
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

func connectToDatabase(testdb bool) *dbIO.DBIO {
	// Manages call to Connect and ReadColumns
	var d string
	if testdb == true {
		d = TDB
	} else {
		d = DB
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
		loadTaxa(db, *infile, *common)
	} else if *lh == true {
		// Upload life history table
		loadLifeHistory(db, *infile)
	} else if *den == true {
		// Uplaod denominator table
		loadNonCancerTotals(db, *infile)
	} else if *patient == true {
		// Upload patient data
		loadAccounts(db, *infile)
		loadDiagnoses(db, *infile)
		loadPatients(db, *infile)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	return db.Starttime
}

func updateDB() time.Time {
	// Updates database with given flags (all input variables are global)
	db := connectToDatabase(false)
	if *total == true {
		speciesTotals(db)
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
	} else if *cr == true {
		// Extract cancer rates
		header := "ScientificName,TotalRecords,CancerRecords,CancerRate,AverageAge(months),AvgAgeCancer(months),Male,Female,MaleCancer,FemaleCancer"
		rates := getCancerRates(db, *min, *nec)
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
		res, header = SearchTaxonomicLevels(db, names)
		fmt.Printf("\tFound %d records where %s is %s.\n", len(res), *level, *taxon)
	} else if *eval != "nil" {
		// Search for column/value match
		column, op, value := getOperation(*eval)
		if *table == "nil" {
			tables := getTable(db.Columns, column)
			res, header = SearchColumns(db, tables, column, op, value)
		} else {
			res, header = SearchSingleTable(db, *table, column, op, value)
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
		loadTaxa(db, *taxafile, true)
		loadLifeHistory(db, *lifehistory)
		// Uplaod denominator table
		loadNonCancerTotals(db, *noncancer)
		// Upload patient data
		loadAccounts(db, *infile)
		loadDiagnoses(db, *infile)
		loadPatients(db, *infile)
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

func main() {
	var start time.Time
	switch kingpin.Parse() {
	case ver.FullCommand():
		version()
	case bu.FullCommand():
		db := connectToDatabase(false)
		start = db.Starttime
		backup(db.Password)
	case New.FullCommand():
		db := connectToDatabase(false)
		start = db.Starttime
		db.NewTables(COL)
	case upload.FullCommand():
		start = uploadToDB()
	case update.FullCommand():
		start = updateDB()
	case extract.FullCommand():
		start = extractFromDB()
	case search.FullCommand():
		start = searchDB()
	case test.FullCommand():
		start = testDB()
	}
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

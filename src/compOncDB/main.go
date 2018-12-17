// Manages input arguments for the comparative oncology database

package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
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
	eval    = kingpin.Flag("eval", "Searches tables for matches (table is automatically determined) (column operator value; valid operators: = <= >= > <). ").Short('e').Default("nil").String()
	infile  = kingpin.Flag("infile", "Path to input file (if using).").Short('i').Default("nil").String()
	outfile = kingpin.Flag("outfile", "Name of output file (writes to stdout if not given).").Short('o').Default("nil").String()

	upload  = kingpin.Command("upload", "Upload data to the database.")
	taxa    = upload.Flag("taxa", "Load taxonomy tables from Kestrel output to update taxonomy table.").Default("false").Bool()
	common  = upload.Flag("common", "Additionally extract common names from Kestrel output to update common name tables.").Default("false").Bool()
	lh      = upload.Flag("lh", "Upload life history info from merged life history table to the database.").Default("false").Bool()
	den     = upload.Flag("den", "Uploads file to denominator table for databases where only cancer records were extracted.").Default("false").Bool()
	patient = upload.Flag("patient", "Upload patient, account, and diagnosis info from input table to database.").Default("false").Bool()

	update = kingpin.Command("update", "Update or delete existing records from the database (see README for upload file template).")
	column = update.Flag("column", "Column to be updated with given value if --eval column == value.").Short('c').Default("nil").String()
	value  = update.Flag("value", "Value to write to column if --eval column == value.").Short('v').Default("nil").String()
	total  = update.Flag("count", "Recount species totals and update the Totals table.").Default("false").Bool()
	del    = update.Flag("delete", "Delete records if column = value.").Default("false").Bool()

	extract = kingpin.Command("extract", "Extract data from the database and perform optional analyses.")
	dump    = extract.Flag("dump", "Name of table to dump (writes all data from table to output file).").Short('d').Default("nil").String()
	sum     = extract.Flag("summarize", "Compiles basic summary statistics of the database.").Default("false").Bool()
	cr      = extract.Flag("cancerRate", "Calculates cancer rates for species with greater than min entries.").Default("false").Bool()
	min     = extract.Flag("min", "Minimum number of entries required for calculations (default = 50).").Short('m').Default("50").Int()
	nec     = extract.Flag("necropsy", "Extract only necropsy records (extracts all matches by default).").Default("false").Bool()

	search = kingpin.Command("search", "Searches database for matches to given term.")
	taxon  = search.Flag("taxa", "Name of taxonomic unit to extract data for or path to file with single column of units.").Short('t').Default("nil").String()
	level  = search.Flag("level", "Taxonomic level of taxon (or entries in taxon file)(default = Species).").Short('l').Default("Species").String()
	com    = search.Flag("common", "Indicates that common species name was given for taxa.").Default("false").Bool()
	count  = search.Flag("count", "Returns count of target records instead of printing entire records.").Default("false").Bool()
	table  = search.Flag("table", "Return matching rows from this table only.").Default("nil").String()

	test        = kingpin.Command("test", "Tests database functionality using testDataBase instead of comaprative oncology.")
	tables      = test.Flag("tables", "Path tableColumns.txt file.").String()
	taxafile    = test.Flag("taxonomy", "Path to taxonomy file.").String()
	diagnosis   = test.Flag("diagnosis", "Path to extracted diganoses file.").String()
	lifehistory = test.Flag("lifehistory", "Path to life history data.").String()
	noncancer   = test.Flag("denominators", "Path to file conataining non-cancer totals.").String()
	testsearch  = test.Flag("search", "Search for matches using above commands.").Default("false").Bool()
	updates  	= test.Flag("update", "Tests update functions.").Default("false").Bool()
)

func version() {
	fmt.Println("\n\tCompOncDB v0.1 (~) is a package for managing the ASU comparative oncology database.")
	fmt.Println("\n\tCopyright 2018 by Shawn Rupp, Maley Lab, Biodesign Institute, Arizona State University.")
	fmt.Println("\tThis program comes with ABSOLUTELY NO WARRANTY.")
	fmt.Print("\n\tThis is free software, and you are welcome to redistribute it under certain conditions.\n\n")
	os.Exit(0)
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

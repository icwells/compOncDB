// Manages input arguments for the comparative oncology database

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	// Kingpin arguments
	app     = kingpin.New("compOncDB", "Command line-interface for uploading/extrating/manipulating data from the comparative oncology database.")
	user    = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()
	config  = kingpin.Flag("config", "Path to config.txt (Default is in utils directory).").Default("config.txt").String()
	eval    = kingpin.Flag("eval", "Searches tables for matches (table is automatically determined) ('column operator value'; valid operators: != = <= >= > <; wrap statement in quotation marks and seperate multiple statements with commas). ").Short('e').Default("nil").String()
	table   = kingpin.Flag("table", "Perform operations on this table only.").Default("nil").String()
	infile  = kingpin.Flag("infile", "Path to input file (if using).").Short('i').Default("nil").String()
	outfile = kingpin.Flag("outfile", "Name of output file (writes to stdout if not given).").Short('o').Default("nil").String()

	ver = kingpin.Command("version", "Prints version info and exits.")
	bu  = kingpin.Command("backup", "Backs up database to local machine (Must use root password; output is written to current directory).")
	New = kingpin.Command("new", "Initializes new tables in new database (database must be initialized manually).")

	parse    = kingpin.Command("parse", "Parse and organize records for upload to the database.")
	service  = parse.Flag("service", "Service database name.").Short('s').Required().String()
	taxaFile = parse.Flag("taxa", "Path to kestrel output.").Short('t').Required().String()
	debug    = parse.Flag("debug", "Adds cancer and code column (if present) for hand checking.").Short('d').Default("false").Bool()

	upload  = kingpin.Command("upload", "Upload data to the database.")
	taxa    = upload.Flag("taxa", "Load taxonomy tables from Kestrel output to update taxonomy table.").Default("false").Bool()
	common  = upload.Flag("common", "Additionally extract common names from Kestrel output to update common name tables.").Default("false").Bool()
	lh      = upload.Flag("lh", "Upload life history info from merged life history table to the database.").Default("false").Bool()
	den     = upload.Flag("den", "Uploads file to denominator table for databases where only cancer records were extracted.").Default("false").Bool()
	patient = upload.Flag("patient", "Upload patient, account, and diagnosis info from input table to database.").Default("false").Bool()

	update = kingpin.Command("update", "Update or delete existing records from the database (see README for upload file template).")
	column = update.Flag("column", "Column to be updated with given value if --eval column == value.").Short('c').Default("nil").String()
	value  = update.Flag("value", "Value to write to column if --eval column == value (only supply one statement).").Short('v').Default("nil").String()
	total  = update.Flag("count", "Recount species totals and update the Totals table.").Default("false").Bool()
	clean  = update.Flag("clean", "Remove extraneous records from the database.").Default("false").Bool()
	del    = update.Flag("delete", "Delete records if column = value.").Default("false").Bool()

	extract  = kingpin.Command("extract", "Extract data from the database and perform optional analyses.")
	dump     = extract.Flag("dump", "Name of table to dump (writes all data from table to output file).").Short('d').Default("nil").String()
	reftaxa  = extract.Flag("reference_taxonomy", "Returns merged common and taxonomy tables.").Short('r').Default("false").Bool()
	sum      = extract.Flag("summarize", "Compiles basic summary statistics of the database.").Default("false").Bool()
	cr       = extract.Flag("cancerRate", "Calculates cancer rates for species with greater than min entries.").Default("false").Bool()
	lifehist = extract.Flag("lifehistory", "Append life history values to cancer rate data.").Default("false").Bool()
	min      = extract.Flag("min", "Minimum number of entries required for calculations.").Short('m').Default("1").Int()
	nec      = extract.Flag("necropsy", "Extract only necropsy records (extracts all matches by default).").Default("false").Bool()

	search     = kingpin.Command("search", "Searches database for matches to given term.")
	count      = search.Flag("count", "Returns count of target records instead of printing entire records.").Default("false").Bool()
	infant     = search.Flag("infant", "Include infant records in results (excluded by default).").Default("false").Bool()
	taxonomies = search.Flag("taxonomies", "Searches for taxonomy matches given column of common/scientific names in a file.").Default("false").Bool()
	col        = search.Flag("names", "Column of input file containing scientific/common species names to search.").Short('n').Default("0").Int()
)

func version() {
	fmt.Println("\n\tCompOncDB v0.3.6 (09/05/19) is a package for managing the ASU comparative oncology database.")
	fmt.Println("\n\tCopyright 2019 by Shawn Rupp, Maley Lab, Biodesign Institute, Arizona State University.")
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
		db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, false))
		start = db.Starttime
		backup(db.Password)
	case New.FullCommand():
		start = newDatabase()
	case parse.FullCommand():
		start = parseRecords()
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

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
	app      = kingpin.New("compOncDB", "Command line-interface for uploading/extrating/manipulating data from the comparative oncology database.")
	eval     = kingpin.Flag("eval", "Searches tables for matches (table is automatically determined) ('column operator value'; valid operators: != = <= >= > < ^; wrap statement in quotation marks and seperate multiple statements with commas; '^' will return match if the column contains the value). ").Short('e').Default("nil").String()
	infant   = kingpin.Flag("infant", "Include infant records in results (excluded by default).").Default("false").Bool()
	infile   = kingpin.Flag("infile", "Path to input file (if using).").Short('i').Default("nil").String()
	min      = kingpin.Flag("min", "Minimum number of entries required for calculations.").Short('m').Default("1").Int()
	outfile  = kingpin.Flag("outfile", "Name of output file (writes to stdout if not given).").Short('o').Default("").String()
	password = kingpin.Flag("password", "Password (for testing or scripting).").Default("").String()
	user     = kingpin.Flag("user", "MySQL username.").Short('u').Default("").String()

	ver = kingpin.Command("version", "Prints version info and exits.")
	bu  = kingpin.Command("backup", "Backs up database to local machine (Specify output directory with '-o' flag).")
	New = kingpin.Command("new", "Initializes new tables in new database (database must be initialized manually).")

	parse    = kingpin.Command("parse", "Parse and organize records for upload to the database.")
	service  = parse.Flag("service", "Service database name.").Short('s').Required().String()
	taxaFile = parse.Flag("taxa", "Path to kestrel output.").Short('t').Required().String()
	debug    = parse.Flag("debug", "Adds cancer and code column (if present) for hand checking.").Short('d').Default("false").Bool()

	verify    = kingpin.Command("verify", "Compares parse output with NLP model predictions. Provide parse records output and new output file with -i and -o.")
	diagnosis = verify.Flag("diagnosis", "Verifies type and location diagnoses only.").Default("false").Bool()
	merge     = verify.Flag("merge", "Merges currated verification results with parse output. Give path to nlp output with -i and path to parse output with -o (it will be overwritten).").Default("false").Bool()
	neoplasia = verify.Flag("neoplasia", "Verifies masspresent diagnosis only.").Default("false").Bool()

	upload  = kingpin.Command("upload", "Upload data to the database. Backs up database if output directory is given with '-o'.")
	taxa    = upload.Flag("taxa", "Load taxonomy tables from Kestrel output to update taxonomy table.").Default("false").Bool()
	common  = upload.Flag("common", "Additionally extract common names from Kestrel output to update common name tables.").Default("false").Bool()
	lh      = upload.Flag("lh", "Upload life history info from merged life history table to the database.").Default("false").Bool()
	den     = upload.Flag("den", "Uploads file to denominator table for databases where only cancer records were extracted.").Default("false").Bool()
	patient = upload.Flag("patient", "Upload patient, account, and diagnosis info from input table to database.").Default("false").Bool()

	update = kingpin.Command("update", "Update or delete existing records from the database (see README for upload file template). Backs up database if output directory is given with '-o'.")
	clean  = update.Flag("clean", "Remove extraneous records from the database.").Default("false").Bool()
	column = update.Flag("column", "Column to be updated with given value if --eval column == value.").Short('c').Default("nil").String()
	del    = update.Flag("delete", "Delete records if column = value.").Default("false").Bool()
	table  = update.Flag("table", "Perform operations on this table only.").Default("nil").String()
	value  = update.Flag("value", "Value to write to column if --eval column == value (only supply one statement).").Short('v').Default("nil").String()

	extract   = kingpin.Command("extract", "Extract data from the database.")
	alltaxa   = extract.Flag("alltaxa", "Summarizes life history table for all species (performs summary for species with records in patient table by default).").Default("false").Bool()
	dump      = extract.Flag("dump", "Name of table to dump (writes all data from table to output file).").Short('d').Default("nil").String()
	dumpdb    = extract.Flag("dump_db", "Extracts entire database into a gzipped tarball of csv files (specify output directory with -o).").Default("false").Bool()
	lhsummary = extract.Flag("lhsummary", "Summarizes life history table.").Default("false").Bool()
	reftaxa   = extract.Flag("reference_taxonomy", "Returns merged common and taxonomy tables.").Short('r').Default("false").Bool()
	sum       = extract.Flag("summarize", "Compiles basic summary statistics of the database.").Default("false").Bool()

	searchdb   = kingpin.Command("search", "Search database for matches to queries.")
	col        = searchdb.Flag("names", "Column of input file containing scientific/common species names to search.").Short('n').Default("0").Int()
	taxonomies = searchdb.Flag("taxonomies", "Searches for taxonomy matches given column of common/scientific names in a file.").Default("false").Bool()

	leader = kingpin.Command("leader", "Calculate neoplasia prevalence leaderboards.")
	typ    = leader.Flag("type", "Returns top 10 species with this tumor type.").Default("").String()

	cancerRates = kingpin.Command("cancerrates", "Calculate neoplasia prevalence for species.")
	keepall     = cancerRates.Flag("keepall", "Keep records without specified tissue when calculating by tissue.").Default("false").Bool()
	lifehist    = cancerRates.Flag("lifehistory", "Append life history values to cancer rate data.").Default("false").Bool()
	location    = cancerRates.Flag("location", "Include tumor location summary for each species.").Default("").String()
	noavgage    = cancerRates.Flag("noavgage", "Will not return average age columns in output file.").Default("false").Bool()
	nosexcol    = cancerRates.Flag("nosexcol", "Will not return male/female specific columns in output file.").Default("false").Bool()
	notaxacol   = cancerRates.Flag("notaxacol", "Will not return Kingdom-Genus columns in output file.").Default("false").Bool()
	nec         = cancerRates.Flag("necropsy", "2: Extract only necropsy records, 1: extract all records by default, 0: extract non-necropsy records.").Default("2").Int()
	pathology   = cancerRates.Flag("pathology", "Additionally extract pathology records for target species.").Default("false").Bool()
	source      = cancerRates.Flag("source", "Zoo/institute records to calculate prevalence with; all: use all records, approved: used zoos approved for publication, aza: use only AZA member zoos, zoo: use only zoos.").Short('z').Default("approved").String()
	tissue      = cancerRates.Flag("tissue", "Include tumor tissue type summary for each species (supercedes location analysis).").Default("").String()
	wild        = cancerRates.Flag("wild", "Return results for wild records only (returns non-wild only by default).").Default("false").Bool()

	newuser  = kingpin.Command("newuser", "Adds new user to database. Must performed on the server using root password.")
	admin    = newuser.Flag("admin", "Grant all privileges to user. Also allows remote MySQL access.").Default("false").Bool()
	username = newuser.Flag("username", "MySQL username for new user. Password will be be set to this name until it is updated.").Default("").String()
)

func version() {
	fmt.Println("\n\tCompOncDB v0.9.1 is a package for managing the ASU comparative oncology database.")
	fmt.Println("\n\tCopyright 2021 by Shawn Rupp, Maley Lab, Biodesign Institute, Arizona State University.")
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
		db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)
		db.BackupDB(*outfile)
		start = db.Starttime
	case New.FullCommand():
		start = newDatabase()
	case parse.FullCommand():
		start = parseRecords()
	case verify.FullCommand():
		start = verifyDiagnoses()
	case upload.FullCommand():
		start = uploadToDB()
	case update.FullCommand():
		start = updateDB()
	case extract.FullCommand():
		start = extractFromDB()
	case searchdb.FullCommand():
		start = searchDB()
	case leader.FullCommand():
		start = leaderBoards()
	case cancerRates.FullCommand():
		start = calculateCancerRates()
	case newuser.FullCommand():
		start = addNewUser()
	}
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

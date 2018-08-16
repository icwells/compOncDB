// This script will manage searching of the comparative oncology database

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/Songmu/prompter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/exec"
	"time"
)

var (
	// Global arguemnts
	COL = "tableColumns.txt"
	DB  = "comparativeOncology"
)

var (
	// Kingpin arguments
	app      = kingpin.New("compOncDB", "Comand line-interface for uploading/extrating/manipulating data from the comparative oncology database.")
	user     = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()
	ver      = kingpin.Command("version", "Print version info and exit.")
	bu       = kingpin.Command("backup", "Backs up database to local machine (Must use root password).")
	New      = kingpin.Command("new", "Initializes new tables in new database (database must be initialized manually).")

	upload   = kingpin.Command("upload", "Upload data to the database.")
	common   = upload.Flag("common", "Additionally extract common names from Kestrel output to update common name tables.").Default("false").Bool()
	taxa     = upload.Flag("taxa", "Load taxonomy tables from Kestrel output to update taxonomy table.").Default("false").Bool()
	lh       = upload.Flag("lh", "Upload life history info from merged life history table to the database.").Default("false").Bool()
	accounts = upload.Flag("accounts", "Extract account info from input file and update database.").Default("false").Bool()
	diag     = upload.Flag("diagnosis", "Extract diagnosis info from input file and update database.").Default("false").Bool()
	patient  = upload.Flag("patient", "Uploads patient info from input table to database.").Default("false").Bool()
	infile   = upload.Arg("infile", "Path to input file.").Required().String()

	extract  = kingpin.Command("extract", "Extracts data from the database and performs optional analyses.")
	dump     = extract.Flag("dump", "Name of table to dump (writes all data from table to output file).").Short('d').Default("nil").String()
	cr       = extract.Flag("cancerRate", "Calculates cancer rates for species with gtreater than min entries.").Default("false").Bool()
	min      = extract.Flag("min", "Minimum number of entries required for calculations (default = 50).").Short('m').Default("50").Int()
	nec      = extract.Flag("necropsy", "Extract only necropsy records (extracts all matches by default).").Default("false").Bool()
	outfile  = extract.Arg("outfile", "Name of output file.").Default("nil").String()
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
	dump := exec.Command("mysqldump", "-uroot", password, res, DB)
	err := dump.Run()
	if err == nil {
		fmt.Println("\tBackup complete.")
	} else {
		fmt.Printf("\tBackup failed. %v\n", err)
	}
}

func connect(user string) (*sql.DB, string, time.Time) {
	// Attempts to connect to sql database. Returns db instance.
	// Prompt for password
	pw := prompter.Password("\n\tEnter MySQL password")
	// Begin recording time after password input
	start := time.Now()
	db, err := sql.Open("mysql", user+":"+pw+"@/"+DB)
	if err != nil {
		fmt.Printf("\n\t[Error] Connecting to database: %v", err)
		os.Exit(2)
	}
	if err = db.Ping(); err != nil {
		fmt.Printf("\n\t[Error] Cannot connect to database: %v", err)
	}
	return db, pw, start
}

func main() {
	var db *sql.DB
	var start time.Time
	var pw string
	switch kingpin.Parse() {
		case ver.FullCommand():
			version()
		case bu.FullCommand():
			db, _, start = connect(*user)
			defer db.Close()
			backup(pw)
		case New.FullCommand():
			db, _, start = connect(*user)
			defer db.Close()
			dbIO.NewTables(db, COL)
		case upload.FullCommand():
			// Upload infile to given table
			col := dbIO.ReadColumns(COL, false)
			db, _, start = connect(*user)
			defer db.Close()
			if *taxa == true {
				// Upload taxonomy
				loadTaxa(db, col, *infile, *common)
			} else if *lh == true {
				// Upload life history table
				loadLifeHistory(db, col, *infile)
			} else if *accounts == true {
				// Upload account info
				loadAccounts(db, col, *infile)
			} else if *diag == true {
				loadDiagnoses(db, col, *infile)
			} else if *patient == true {
				// Upload patient data
				loadPatients(db, col, *infile)
			}
		case extract.FullCommand():
			// Extract data to outfile/stdout
			col := dbIO.ReadColumns(COL, false)
			db, _, start = connect(*user)
			defer db.Close()	
			if *dump != "nil" {
				// Extract entire table
				table := dbIO.GetTable(db, *dump)
				iotools.WriteToCSV(*outfile, col[*dump], table)
			} else if *cr == true {
				header := "ScientificName,TotalRecords,CancerRecords,CancerRate,AverageAge(months),AvgAgeCancer(months),Male:Female\n"
				rates := getCancerRates(db, col, *min, *nec)
				iotools.WriteToCSV(*outfile, header, rates)
			}
	}
	fmt.Printf("\n\tFinished. Runtime: %s\n\n", time.Since(start))
}

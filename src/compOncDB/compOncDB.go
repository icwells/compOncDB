// This script will manage searching of the comparative oncology database

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/icwells/go-tools/iotools"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/exec"
	"time"
)

func version() {
	fmt.Println("\n\tCompOncDV v0.1 (~) is a package for managing the ASU comparative oncology database.")
	fmt.Println("\n\tCopyright 2018 by Shawn Rupp, Maley Lab, Biodesign Institute, Arizona State University.")
	fmt.Println("\tThis program comes with ABSOLUTELY NO WARRANTY.")
	fmt.Println("\n\tThis is free software, and you are welcome to redistribute it under certain conditions.\n")
	os.Exit(0)
}

func backup(DB, pw string) {
	// Backup database to local machine
	fmt.Printf("\n\tBacking up %s database to local machine...\n", DB)
	datestamp := time.Now().Format("2006-01-02")
	password := fmt.Sprintf("-p%s", pw)
	res := fmt.Sprintf("--result-file=%s.%s.sql", DB, datestamp)
	dump := exec.Command("mysqldump","-uroot", password, res, DB)
	err := dump.Run()
	if err == nil {
		fmt.Println("\tBackup complete.")
	} else {
		fmt.Printf("\tBackup failed. %v\n", err)
	}
}

func connect(DB, user, pw string) *sql.DB {
	// Attempts to connect to sql database. Returns db instance.
	db, err := sql.Open("mysql", user+":"+pw+"@/"+DB)
	if err != nil {
		fmt.Printf("\n\t[Error] Connecting to database: %v", err)
		os.Exit(2)
	}
	if err = db.Ping(); err != nil {
		fmt.Printf("\n\t[Error] Cannot connect to database: %v", err)
	}
	return db
}

func main() {
	COL := "tableColumns.txt"
	DB := "comparativeOncology"
	var (
		user      = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()
		ver       = kingpin.Flag("version", "Print version info and exit.").Short('v').Default("false").Bool()
		bu        = kingpin.Flag("backup", "Backs up database to local machine (Must use root password).").Default("false").Bool()
		New       = kingpin.Flag("new", "Initializes new tables in new database (database must be initialized manually).").Default("false").Bool()
		taxa      = kingpin.Flag("taxa", "Load taxonomy tables from Kestrel output to update taxonomy and common name tables.").Default("false").Bool()
		accounts  = kingpin.Flag("accounts", "Extract account info from input file and update database.").Default("false").Bool()
		diag      = kingpin.Flag("diagnosis", "Extract diagnosis info from input file and update database.").Default("false").Bool()
		upload    = kingpin.Flag("upload", "Uploads patient info from input table to database.").Default("false").Bool()
		dump      = kingpin.Flag("dump", "Name of table to dump (writes all data from table to output file).").Short('d').Default("nil").String()
		infile    = kingpin.Flag("infile", "Path to input file.").Short('i').Default("nil").String()
		outfile   = kingpin.Flag("outfile", "Name of output file.").Short('o').Default("nil").String()
	)
	kingpin.Parse()
	if *ver == true {
		version()
	}
	// Prompt for password
	password := prompter.Password("\n\tEnter MySQL password")
	// Begin recording time after password input
	start := time.Now()
	db := connect(DB, *user, password)
	defer db.Close()
	if *bu == true {
		backup(DB, password)
	} else {
		col := dbIO.ReadColumns(COL, false)
		if *New == true {
			dbIO.NewTables(db, COL)
		} else if *dump != "nil" {
			// Extract entire table
			if *outfile == "nil" {
				fmt.Println("\n\t[Error] Please specify output file. Exiting.\n")
				os.Exit(1)
			}
			table := dbIO.GetTable(db, *dump)
			iotools.WriteToCSV(*outfile, col[*dump], table)
		} else if *taxa == true {
			// Upload taxonomy
			LoadTaxa(db, col, *infile)
		} else if *accounts == true {
			// Upload account info
			LoadAccounts(db, col, *infile)
		} else if *diag == true {
			LoadDiagnoses(db, col, *infile)
		} else if *upload == true {
			// Upload patient data
			LoadPatients(db, col, *infile)
		}
	}
	fmt.Printf("\n\tFinished. Runtime: %s\n\n", time.Since(start))
}

// This script will manage searching of the comparative oncology database

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/go-sql-driver/mysql"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strconv"
	"strings"
	"time"
)

func version() {
	fmt.Println("\n\tCompOncDV v0.1 (~) is a package for managing the ASU comparative oncology database.")
	fmt.Println("\n\tCopyright 2018 by Shawn Rupp, Maley Lab, Biodesign Institute, Arizona State University.")
	fmt.Println("\tThis program comes with ABSOLUTELY NO WARRANTY.")
	fmt.Prinln("\n\tThis is free software, and you are welcome to redistribute it under certain conditions.\n")
	os.Exit(0)
}

/*func backup(DB string) {
	// Backup database to local machine
	fmt.Println("\n\tBacking up %s database to local machine...", DB)
	dump := exec.Command("mysqldump", fmt.Sprintf("-u root -p --result-file=%s.%s.sql '%s'", DB, time.Date(), DB)
	err := dump.Run()
	if err == nil {
		fmt.Println("\tBackup complete.\n")
	} else {
		fmt.Println("\tBackup failed.\n")
	}
}*/

func connect(DB, user, pw string) *DB {
	// Attempts to connect to sql database. Returns db instance.
	if pw == nil {
		// Prompt for password
		pw = prompter.Password("Enter MySQL password: ")
	}
	db, err := sql.Open("mysql", user+":"+pw+"@/"+DB)
	if err != nil {
		fmt.Fprintf("\n\t[Error] Connecting to database: %v", err)
		os.Exit(2)
	}
	if err = db.Ping(); err != nil {
		fmt.Fprintf("\n\t[Error] Bad database connection: %v", err)
	}
	return db
}

func main() {
	start := time.Now()
	COL := "tableColumns.txt"
	DB := "comparativeOncology"
	var (
		user    = kingpin.Flag("u", "MySQL username (default is root)").Default("root").String()
		pw      = kingpin.Flag("p", "MySQL password").String()
		ver     = kingpin.Flag("v", "Print version info").Default("false").Bool()
		bu      = kingpin.Flag("backup", "Backs up database to local machine").String()
		New     = kingpin.Flag("new", "Initializes new tables in new database (database must be made manually).").Default("false").Bool()
		dump    = kingpin.Flag("dump", "Name of table to dump (writes all data from table to output file).").Default("").String()
		infile  = kingpin.Flag("i", "Path to input file.").String()
		outfile = kingpin.Flag("o", "Name of output file.").String()
		//cpu     = kingpin.Flag("t", "Number of threads (default = 1).").Default("1").Int()
	)
	kingpin.Parse()
	db := connect(DB, *user, *pw)
	defer db.Close()
	if *bu == true {
		backup(DB)
	} else if *New == true {
		dbIO.NewTables(db, COL)
	} else if *dump != nil {
		// Extract entire table
		if *outfile == nil {
			fmt.Println("\n\t[Error] Please specify output file. Exiting.\n")
			os.Exit(1)
		}
		col := dbIO.ReadColumns(COL, false)
		table := dbIO.GetTable(db, *dump)
		printCSV(*outfile, col[*dump], table)
	} else if *infile != nil {
		// Upload csv
		col := dbIO.ReadColumns(COL, false)

	}
	fmt.Println("\n\tFinished. Runtime: %s", time.Since(start))
}

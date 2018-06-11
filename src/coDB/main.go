// This script will manage searching of the comparative oncology database

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/icwells/go-tools/iotools"
	"github.com/Songmu/prompter"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strconv"
	"strings"
)

func connect(DB, user, pw string) *DB {
	// Attempts to connect to sql database. Returns db instance.
	if len(pw) <= 0 {
		// Prompt for password
		pw = prompter.Password("Enter MySQL password: ")
	}
	db, err := sql.Open("mysql", user + ":" + pw + "@/" + DB)
	if err != nil {
		fmt.Fprintf("\n\t[Error] Connecting to database: %v", err)
		os.Exit(1)
	}
	if err = db.Ping(); err != nil {
		fmt.Fprintf("\n\t[Error] Bad database connection: %v", err)
	}
	return db
}


func main() {
	DB := "comparativeOncology"
	var (
		user	= kingpin.Flag("u", "MySQL username").Required().String()
		pw		= kingpin.Flag("p", "MySQL password").Default("").String()
		New		= kingpin.Flag("new", "Initializes new tables in new database (database must be made manually).").Default("false").Boolean()
		dump	= kingpin.Flag("dump", "Name of table to dump (writes all data from table to output file).").Default("").String()
		infile	= kingpin.Flag("i", "Path to input file.").Default("").String()
		outfile = kingpin.Flag("o", "Name of output file.").Default("").String()
		//cpu     = kingpin.Flag("t", "Number of threads (default = 1).").Default("1").Int()
	)
	kingpin.Parse()
	db := connect(DB, *user, *pw)
	defer db.Close()

}

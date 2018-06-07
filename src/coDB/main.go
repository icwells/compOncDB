// This script will manage searching of the comparative oncology database

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strconv"
	"strings"
)

func connect(DB string) *DB {
	// Attempts to connect to sql database. Returns db instance.
	db, err := sql.Open("mysql", DB)
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
		outfile = kingpin.Flag("o", "Name of output file.").Required().String()
		cpu     = kingpin.Flag("t", "Number of threads (default = 1).").Default("1").Int()
		n       = kingpin.Flag("n", "Number of generations to simulate (default = 20).").Default("20").Int()
	)
	kingpin.Parse()
	db := connect(DB)
	defer db.Close()

}

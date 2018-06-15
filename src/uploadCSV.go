// This script will upload csvs to the comparative oncology database

package main

import (
	"github.com/icwells/go-tools/iotools"
	//"os"
	//"strconv"
	//"strings"
)

type Columns struct (
	ID		int
	name	int
	kingdom	int
	phylum	int
	class	int
	order	int
	family	int
	genus	int
	species	int
	date	int
	diag	int
	
)

/*func (c *Columns) getColumns(line) {
	// Returns struct of column numbers with given fields

}

func parseLine(db *DB, c Columns, line string) {
	// Sorts data from line into appopriate tables

}*/

func readCSV(db *DB, infile string) {
	// Reads csv files and sends data to apprriate table
	first = true
	var col Columns
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			parseLine(db, col, line)
		} else {
			col.getColumns(line)
			first = false
		}
	}
}

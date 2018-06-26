// This script will upload patient data to the comparative oncology database

package main

import (
	"bufio"
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"strconv"
	"strings"
)


func uploadTable(db *sql.DB, col map[string]string, p Patient, d Diagnosis, t TumorRelation, s Source) {
	// Uploads unique account entries with random ID number
	var acc [][]string
	for _, i := range accounts {
		// Add unique taxa ID
		count++
		c := strconv.Itoa(count)
		acc = append(acc, []string{c, i})
	}
	if len(acc) > 0 {
		vals, l := dbIO.FormatSlice(acc)
		dbIO.UpdateDB(db, "Accounts", col["Accounts"], vals, l)
	}
}

func getCodes() int {
	// Returns SMALLINT code type for true/false/NA/unspecified

}

func extractPatients(infile string, count int) {
	// Assigns patient data to appropriate structs for sorting later
	first := true
	var p Patient
	var d Diagnosis
	var t TumorRelation
	var s Source
	var c Columns
	var col int
	fmt.Printf("\n\tExtracting accounts from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			count++

		} else {
			c.setIndeces(line)
			first = false
		}
	}
}

func LoadPatients(db *sql.DB, col map[string]string, infile string) {
	// Loads unique patient info to appropriate tables
	m := dbIO.GetMax(db, "Patient", "ID")
	tumor := dbIO.GetTable(db, "Tumor")
	acc := dbIO.GetTable(db, "Accounts")
	meta := dbIO.GetTable(db, "Metastasis")
	species := dbIO.GetColumns(db, "Taxonomy", []string{"taxa_id", "Species"})
	p, d, t, s := extractPatients(infile, m)
	p, d, t, s := sortPatients(p, d, t, s, tumor, acc, meta, species)
	uploadAccounts(db, col, p, d, t, s)
}

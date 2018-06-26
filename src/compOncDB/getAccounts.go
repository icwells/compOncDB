// This srcipt will summarize diagnosis and account data from database files
// and upload them the comparative oncology database

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

func uploadAccounts(db *sql.DB, col map[string]string, accounts []string, count int) {
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

func getIndex(line string) int {
	// Assigns column indeces to struct
	var col int
	s := strings.Split(line, ",")
	for idx, i := range s {
		if i == "Owner" || i == "Client" {
			col = idx
			break
		}
	}
	return col
}

func extractAccounts(infile string, acc []string) []string {
	// Extracts accounts from input file
	first := true
	var accounts []string
	var col int
	fmt.Printf("\n\tExtracting accounts from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			s := strings.Split(line, ",")
			a := strings.Trim(s[col], " \n\t")
			if strarray.InSliceStr(acc, a) == false && strarray.InSliceStr(accounts, a) == false {
				// Add unique occurances
				accounts = append(accounts, a)
			}
		} else {
			col = getIndex(line)
			first = false
		}
	}
	return accounts
}

func LoadAccounts(db *sql.DB, col map[string]string, infile string) {
	// Loads unique entries into comparative oncology metastatis, tumor, and account tables
	m := dbIO.GetMax(db, "Accounts", "account_id")
	acc := dbIO.GetColumnText(db, "Accounts", "Name")
	accounts := extractAccounts(infile, acc)
	uploadAccounts(db, col, accounts, m)
}

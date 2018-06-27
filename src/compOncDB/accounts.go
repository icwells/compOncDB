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

func uploadAccounts(db *sql.DB, col map[string]string, accounts map[string][]string, count int) {
	// Uploads unique account entries with random ID number
	var acc [][]string
	for k, v := range accounts {
		for _, i := range v {
			// Add unique taxa ID
			count++
			c := strconv.Itoa(count)
			acc = append(acc, []string{c, i})
		}
	}
	if len(acc) > 0 {
		vals, l := dbIO.FormatSlice(acc)
		dbIO.UpdateDB(db, "Accounts", col["Accounts"], vals, l)
	}
}

func getIndex(line string) (int, int) {
	// Assigns column indeces to struct
	var a, c int
	s := strings.Split(line, ",")
	for idx, i := range s {
		if i == "Owner" || i == "Client" {
			c = idx
		} else if i == "Account" {
			a = idx
		}
	}
	return a, c
}

func tableToMap(t [][]string) map[string][]string {
	// Converts extracted table to map for easier sorting
	m := make(map[string][]string)
	for _, i := range t {
		if strarray.InMapSli(m, i[0]) == true {
			if strarray.InSliceStr(m[i[0]], i[1]) == false {
				// Add new submitter name
				m[i[0]] = append(m[i[0]], i[1])
			}
		} else {
			m[i[0]] = []string{i[1]}
		}
	}
	return m
}

func extractAccounts(infile string, table [][]string) map[string][]string {
	// Extracts accounts from input file
	first := true
	accounts := make(map[string][]string)
	var col int
	acc := tableToMap(table)
	fmt.Printf("\n\tExtracting accounts from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			pass := false
			s := strings.Split(line, ",")
			account := strings.Trim(s[a], " \n\t")
			client := strings.Trim(s[c], " \n\t")
			// Determine if entry is unique
			rep := strarray.InMapSli(accounts, account)
			if rep == false {
				pass = true
			} else if rep == true && strarray.InSliceStr(accounts[account], client) == false {
				pass = true
			} else if strarray.InMapSli(acc, account) == true && strarray.InSliceStr(acc[account], client) == false {
				pass = true
			}
			if pass == true {
				// Add unique occurances
				accounts = append(accounts, []string{account, client})
			}
		} else {
			a, c = getIndex(line)
			first = false
		}
	}
	return accounts
}

func LoadAccounts(db *sql.DB, col map[string]string, infile string) {
	// Loads unique entries into comparative oncology metastatis, tumor, and account tables
	m := dbIO.GetMax(db, "Accounts", "account_id")
	acc := dbIO.GetColumns(db, "Accounts", []string{"Name", "submitter_name"})
	accounts := extractAccounts(infile, acc)
	uploadAccounts(db, col, accounts, m)
}

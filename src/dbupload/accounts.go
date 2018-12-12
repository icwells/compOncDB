// This srcipt will summarize diagnosis and account data from database files
// and upload them the comparative oncology database

package dbupload

import (
	"bufio"
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"strconv"
	"strings"
)

func uploadAccounts(db *dbIO.DBIO, accounts map[string][]string, count int) {
	// Uploads unique account entries with random ID number
	var acc [][]string
	for k, v := range accounts {
		for _, i := range v {
			// Add unique taxa ID
			count++
			c := strconv.Itoa(count)
			acc = append(acc, []string{c, k, i})
		}
	}
	if len(acc) > 0 {
		vals, l := dbIO.FormatSlice(acc)
		db.UpdateDB("Accounts", vals, l)
	}
}

func extractAccounts(infile string, table [][]string) map[string][]string {
	// Extracts accounts from input file
	first := true
	accounts := make(map[string][]string)
	acc := ToMap(table)
	fmt.Printf("\n\tExtracting accounts from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		s := strings.Split(line, ",")
		if first == false && len(s) == 17 {
			pass := false
			account := strings.Trim(s[15], " \n\t")
			client := strings.Trim(s[16], " \n\t")
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
				accounts[account] = append(accounts[account], client)
			}
		} else {
			first = false
		}
	}
	return accounts
}

func LoadAccounts(db *dbIO.DBIO, infile string) {
	// Loads unique entries into comparative oncology metastatis, tumor, and account tables
	m := db.GetMax("Accounts", "account_id")
	acc := db.GetColumns("Accounts", []string{"Account", "submitter_name"})
	accounts := extractAccounts(infile, acc)
	uploadAccounts(db, accounts, m)
}

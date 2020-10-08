// This srcipt will summarize diagnosis and account data from database files
// and upload them the comparative oncology database

package dbupload

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"log"
	"strconv"
	"strings"
)

type accounts struct {
	db       *dbIO.DBIO
	count    int
	logger   *log.Logger
	acc, neu map[string][]string
}

func newAccounts(db *dbIO.DBIO) *accounts {
	// Returns new account struct
	a := new(accounts)
	a.db = db
	a.acc = codbutils.ToMap(a.db.GetColumns("Accounts", []string{"Account", "submitter_name"}))
	a.count = a.db.GetMax("Accounts", "account_id") + 1
	a.logger = codbutils.GetLogger()
	a.neu = make(map[string][]string)
	return a
}

func (a *accounts) uploadAccounts() {
	// Uploads unique account entries with random ID number
	if len(a.neu) > 0 {
		vals, l := dbIO.FormatMap(a.neu)
		a.db.UpdateDB("Accounts", vals, l)
	}
}

func (a *accounts) extractAccounts(infile string) {
	// Extracts accounts from input file
	var col map[string]int
	var l int
	first := true
	a.logger.Printf("Extracting accounts from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := iotools.GetScanner(f)
	for input.Scan() {
		s := strings.Split(string(input.Text()), ",")
		if first == false && len(s) == l {
			pass := false
			account := strings.TrimSpace(s[col["Account"]])
			client := strings.TrimSpace(s[col["Submitter"]])
			// Determine if entry is unique
			row, ex := a.neu[account]
			if ex == false {
				pass = true
			} else if strarray.InSliceStr(row, client) == false {
				pass = true
			} else if _, e := a.acc[account]; e == true && strarray.InSliceStr(a.acc[account], client) == false {
				pass = true
			}
			if pass == true {
				// Add unique occurances
				a.neu[account] = []string{strconv.Itoa(a.count), account, client}
				a.count++
			}
		} else {
			col = iotools.GetHeader(s)
			l = len(s)
			first = false
		}
	}
}

func LoadAccounts(db *dbIO.DBIO, infile string) {
	// Loads unique entries into comparative oncology metastatis, tumor, and account tables
	a := newAccounts(db)
	a.extractAccounts(infile)
	a.uploadAccounts()
}

// This srcipt will summarize diagnosis and account data from database files
// and upload them the comparative oncology database

package dbupload

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/simpleset"
	"log"
	"strconv"
)

type accounts struct {
	db        *dbIO.DBIO
	count     int
	logger    *log.Logger
	neu       map[string][]string
	submitter *simpleset.Set
}

func newAccounts(db *dbIO.DBIO) *accounts {
	// Returns new account struct
	a := new(accounts)
	a.db = db
	a.count = a.db.GetMax("Accounts", "account_id") + 1
	a.logger = codbutils.GetLogger()
	a.neu = make(map[string][]string)
	a.submitter = simpleset.FromStringSlice(a.db.GetColumnText("Accounts", "submitter_name"))
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
	a.logger.Printf("Extracting accounts from %s\n", infile)
	reader, col := iotools.YieldFile(infile, true)
	l := len(col)
	for s := range reader {
		if len(s) == l {
			client := s[col["Submitter"]]
			if ex, _ := a.submitter.InSet(client); !ex {
				a.submitter.Add(client)
				// Add unique occurances
				a.neu[client] = []string{strconv.Itoa(a.count), client}
				a.count++
			}
		}
	}
}

func LoadAccounts(db *dbIO.DBIO, infile string) {
	// Loads unique entries into comparative oncology metastatis, tumor, and account tables
	a := newAccounts(db)
	a.extractAccounts(infile)
	a.uploadAccounts()
}

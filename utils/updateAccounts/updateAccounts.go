// Re-formates submitter names and adds zoo and institute columns

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/clusteraccounts"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"strconv"
	"time"
)

type updater struct {
	db                            *dbIO.DBIO
	source, accounts, newaccounts map[string][]string
	keys                          map[string]map[string]string
}

func (u *updater) setSources(key, aid, zoo, inst string) {
	// Stores new account id, and source types by source id
	// Get source ids corresponding to old account id
	sids, ex := u.source[key]
	if ex == true {
		for _, i := range sids {
			// Store in source maps by field type
			u.keys["Zoo"][i] = zoo
			u.keys["Institute"][i] = inst
			u.keys["account_id"][i] = aid
		}
	}
}

func (u *updater) setAccounts() {
	// Calls clusteraccounts and assigns old ids to keys map
	count := 1
	a := clusteraccounts.NewAccounts("")
	for _, v := range u.accounts {
		a.Queries.Add(v[1])
	}
	neu := a.ResolveAccounts()
	for k, v := range u.accounts {
		n, ex := neu[v[1]]
		if ex == true {
			id := strconv.Itoa(count)
			row, e := u.newaccounts[n[0]]
			if e == false {
				// Assign novel entries to map
				u.newaccounts[n[0]] = []string{id, v[0], n[0]}
				count++
			} else {
				// Get existing account id
				id = row[0]
			}
			u.setSources(k, id, n[1], n[2])
		}
	}
}

func newUpdater(db *dbIO.DBIO) *updater {
	// Initializes updater struct
	u := new(updater)
	u.db = db
	u.keys = make(map[string]map[string]string)
	u.keys["Zoo"] = make(map[string]string)
	u.keys["Institute"] = make(map[string]string)
	u.keys["account_id"] = make(map[string]string)
	u.newaccounts = make(map[string][]string)
	u.source = dbupload.ToMap(u.db.GetColumns("Source", []string{"account_id", "ID"}))
	u.accounts = dbupload.ToMap(u.db.GetTable("Accounts"))
	u.setAccounts()
	return u
}

//----------------------------------------------------------------------------------

func (u *updater) updateSources() {
	// Adds zoo and institute columns to source table
	fmt.Println("\tUpdataing sources table...")
	_ = u.db.UpdateColumns("Source", "ID", u.keys)
}

func (u *updater) updateAccounts() {
	// Uploads newaccounts to database
	var table [][]string
	fmt.Println("\tUploading new accounts table...")
	for _, v := range u.newaccounts {
		table = append(table, v)
	}
	u.db.TruncateTable("Accounts")
	vals, l := dbIO.FormatSlice(table)
	u.db.UpdateDB("Accounts", vals, l)
}

func main() {
	start := time.Now()
	fmt.Println("\n\tUpdating account values in database...")
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration("config.txt", "smrupp", false))
	u := newUpdater(db)
	u.updateAccounts()
	u.updateSources()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

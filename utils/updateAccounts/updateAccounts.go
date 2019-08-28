// Re-formates submitter names and adds zoo and institute columns

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/compOncDB/src/clusteraccounts"
	"github.com/icwells/dbIO"
	"strconv"
	"time"
)

type updater struct {
	db			*dbIO.DBIO
	source		map[string]string
	accounts, newaccounts, keys	map[string][]string
}

func (u *updater) setAccounts() {
	// Calls clusteraccounts and assigns old ids to keys map
	count := 0
	a := clusteraccounts.NewAccounts("")
	for k, v := range u.accounts {
		a.queries.Add(v[1])
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
				id = row[0]
			}
			// Store old id: [zoo, inst, new id]
			u.keys[k] = append(n[1:], id)
		}
	}
}

func newUpdater(db *dbIO.DBIO) *updater {
	// Initializes updater struct
	u := new(updater)
	u.db = db
	u.keys = make(map[string][]string)
	u.newaccounts = make(map[string][]string)
	u.source = dbupload.EntryMap(u.db.GetColumns("Source", []string{"account_id", "ID"}))
	u.accounts = dbupload.ToMap(u.db.GetTable("Accounts"))
	u.setAccounts()
	return u
}

func (u *updater) updateSources() {
	// Adds zoo and institute columns to source table
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

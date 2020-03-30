// Re-formates submitter names and adds zoo and institute columns

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/clusteraccounts"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"math"
	"os"
	"strconv"
	"time"
)

type updater struct {
	db                                  *dbIO.DBIO
	keys, source, accounts, newaccounts map[string][]string
	newsource                           [][]string
}

func (u *updater) setSources(key, aid, zoo, inst string) {
	// Stores new account id, and source types with source id
	// Get source ids corresponding to old account id
	sids, ex := u.keys[key]
	if ex == true {
		for _, i := range sids {
			row, e := u.source[i]
			if e == true {
				// Append patient id, service name, zoo/inst status, and account id
				u.newsource = append(u.newsource, []string{i, row[0], zoo, inst, aid})
			}
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
	u.newaccounts = make(map[string][]string)
	u.keys = dbupload.ToMap(u.db.GetColumns("Source", []string{"account_id", "ID"}))
	u.source = dbupload.ToMap(u.db.GetTable("Source"))
	u.accounts = dbupload.ToMap(u.db.GetTable("Accounts"))
	u.setAccounts()
	return u
}

func (u *updater) toSlice(table string) [][]string {
	// Converts map to slice
	var ret [][]string
	var m map[string][]string
	if table == "Source" {
		m = u.source
	} else {
		m = u.accounts
	}
	for k, v := range m {
		row := append([]string{k}, v...)
		ret = append(ret, row)
	}
	return ret
}

func (u *updater) restoreTable(table string) {
	// Returns table to original state
	var end int
	fmt.Printf("\tRestoring original %s table...", table)
	s := u.toSlice(table)
	d := u.getDenominator(s)
	length := len(s)
	l := int(math.Floor(float64(length) / float64(d)))
	idx := 0
	u.db.TruncateTable(table)
	for i := 0; i < d; i++ {
		// Get end index
		if idx+l > length {
			end = length
		} else {
			end = l + idx
		}
		vals, ln := dbIO.FormatSlice(s[idx:end])
		u.db.UpdateDB(table, vals, ln)
	}
}

//----------------------------------------------------------------------------------

func (u *updater) getDenominator(s [][]string) int {
	// Returns denominator for subsetting upload slice (size in bytes / 16Mb)
	max := 10000000.0
	size := 0
	for _, row := range s {
		for _, i := range row {
			size += len([]byte(i))
		}
	}
	return int(math.Ceil(float64(size*8) / max))
}

func (u *updater) updateSources() {
	// Adds zoo and institute columns to source table
	var end int
	fmt.Println("\tUpdating sources table...")
	d := u.getDenominator(u.newsource)
	length := len(u.newsource)
	l := int(math.Floor(float64(length) / float64(d)))
	idx := 0
	u.db.TruncateTable("Source")
	for i := 0; i < d; i++ {
		fmt.Printf("\r\tPerforming update %d of %d...", i+1, d)
		// Get end index
		if idx+l > length {
			end = length
		} else {
			end = l + idx
		}
		vals, ln := dbIO.FormatSlice(u.newsource[idx:end])
		res := u.db.UpdateDB("Source", vals, ln)
		if res == 1 {
			idx = end
		} else {
			u.restoreTable("Source")
			os.Exit(1)
		}
	}
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
	res := u.db.UpdateDB("Accounts", vals, l)
	if res == 0 {
		u.restoreTable("Accounts")
	}
}

func main() {
	start := time.Now()
	fmt.Println("\n\tUpdating account values in database...")
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration("config.txt", "smrupp", false))
	u := newUpdater(db)
	u.updateSources()
	u.updateAccounts()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

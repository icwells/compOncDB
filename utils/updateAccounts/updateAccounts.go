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
	db                            *dbIO.DBIO
	source, accounts, newaccounts map[string][]string
	newsource                     [][]string
}

func (u *updater) setSources(key, aid, zoo, inst string) {
	// Stores new account id, and source types with source id
	// Get source ids corresponding to old account id
	sids, ex := u.source[key]
	if ex == true {
		for _, i := range sids {
			u.newsource = append(u.newsource, []string{i, zoo, inst, aid})
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
	u.source = dbupload.ToMap(u.db.GetColumns("Source", []string{"account_id", "ID"}))
	u.accounts = dbupload.ToMap(u.db.GetTable("Accounts"))
	u.setAccounts()
	return u
}

//----------------------------------------------------------------------------------

func (u *updater) getDenominator() int {
	// Returns denominator for subsetting upload slice (size in bytes / 16Mb)
	max := 2000000.0
	size := 0
	for _, row := range u.newsource {
		for _, i := range row {
			size += len([]byte(i))
		}
	}
	return int(math.Ceil(float64(size*8) / max))
}

func (u *updater) subsetMap(start, end int) map[string]map[string]string {
	// Subsets maps from newsource to upload to table
	ret := make(map[string]map[string]string)
	ret["Zoo"] = make(map[string]string)
	ret["Institute"] = make(map[string]string)
	ret["account_id"] = make(map[string]string)
	for _, i := range u.newsource[start:end] {
		// Store in source maps by field type
		ret["Zoo"][i[0]] = i[1]
		ret["Institute"][i[0]] = i[2]
		ret["account_id"][i[0]] = i[3]
	}
	return ret
}

func (u *updater) updateSources() {
	// Adds zoo and institute columns to source table
	var end int
	fmt.Println("\tUpdating sources table...")
	d := u.getDenominator()
	length := len(u.newsource)
	l := int(math.Floor(float64(length) / float64(d)))
	idx := 0
	for i := 0; i < d; i++ {
		fmt.Printf("\r\tPerforming update %d of %d...", i+1, d)
		// Get end index
		if idx+l > length {
			end = length
		} else {
			end = l + idx
		}
		res := u.db.UpdateColumns("Source", "ID", u.subsetMap(idx, end))
		if res == true {
			idx = end
		} else {
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
	u.db.UpdateDB("Accounts", vals, l)
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

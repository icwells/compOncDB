// Tests update accounts package

package main

import (
	"flag"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"os"
	"strconv"
	"testing"
)

var (
	password = flag.String("password", "", "MySQL password.")
	user     = flag.String("user", "", "MySQL username.")
)

func getAccounts() map[string][]string {
	// Returns map of account data
	ret := make(map[string][]string)
	ret["1"] = []string{"X520", "XYZ"}
	ret["2"] = []string{"A16", "Kv Zoo"}
	return ret
}

func getSource() map[string][]string {
	// Returns map of account data
	ret := make(map[string][]string)
	ret["1"] = []string{"NWZP", "0", "0", "1"}
	ret["2"] = []string{"NWZP", "0", "0", "1"}
	ret["3"] = []string{"NWZP", "0", "0", "1"}
	ret["4"] = []string{"NWZP", "0", "0", "1"}
	ret["5"] = []string{"NWZP", "0", "0", "1"}
	ret["6"] = []string{"NWZP", "0", "0", "1"}
	ret["7"] = []string{"NWZP", "0", "0", "1"}
	ret["8"] = []string{"NWZP", "1", "0", "2"}
	ret["9"] = []string{"NWZP", "1", "0", "2"}
	ret["10"] = []string{"NWZP", "1", "0", "2"}
	ret["11"] = []string{"NWZP", "1", "0", "2"}
	ret["12"] = []string{"NWZP", "1", "0", "2"}
	ret["13"] = []string{"NWZP", "1", "0", "2"}
	ret["14"] = []string{"NWZP", "1", "0", "2"}
	ret["15"] = []string{"NWZP", "1", "0", "2"}
	ret["16"] = []string{"NWZP", "1", "0", "2"}
	ret["17"] = []string{"NWZP", "1", "0", "2"}
	ret["18"] = []string{"NWZP", "1", "0", "2"}
	//ret["19"] = []string{"NWZP", "1", "0", "2"}
	return ret
}

func connectToDatabase() *dbIO.DBIO {
	// Manages call to Connect and GetTableColumns
	flag.Parse()
	c := codbutils.SetConfiguration(*user, true)
	db, err := dbIO.Connect(c.Host, c.Testdb, c.User, *password)
	if err != nil {
		os.Exit(1000)
	}
	db.GetTableColumns()
	return db
}

func compareEntries(actual, expected []string) int {
	// Returns true if both slices are equal
	equal := true
	for idx, i := range actual {
		// Skip randomly assigned IDs
		if idx < len(expected) && i != expected[idx] {
			equal = false
			// Attempt to resolve differences in floating point precision
			a, err := strconv.ParseFloat(i, 64)
			if err == nil {
				var e float64
				e, err = strconv.ParseFloat(expected[idx], 64)
				if err == nil && a == e {
					equal = true
				}
			}
		}
		if equal == false {
			return idx
		}
	}
	return -1
}

func compareTables(t *testing.T, name string, exp, act map[string][]string) {
	// Compares output of equivalent tables
	if len(act) != len(exp) {
		t.Errorf("%s: Actual table length %d does not equal expected: %d", name, len(act), len(exp))
	} else {
		for k := range act {
			if len(act[k]) != len(exp[k]) {
				t.Errorf("%s, %s: Actual length %d does not equal expected: %d", name, k, len(act[k]), len(exp[k]))
				break
			} else {
				idx := compareEntries(act[k], exp[k])
				if idx >= 0 {
					t.Errorf("%s %s-%d: Actual value %s does not equal expected: %s", name, k, idx, act[k][idx], exp[k][idx])
					break
				}
			}
		}
	}
}

func TestUpdateAccounts(t *testing.T) {
	// Runs white and black box tests
	acc := getAccounts()
	u := newUpdater(connectToDatabase())
	compareTables(t, "origninal Accounts", acc, u.accounts)
	u.updateAccounts()
	u.updateSources()
	accounts := dbupload.ToMap(u.db.GetTable("Accounts"))
	sources := dbupload.ToMap(u.db.GetTable("Source"))
	compareTables(t, "updated Accounts", u.accounts, accounts)
	compareTables(t, "Source", getSource(), sources)
}

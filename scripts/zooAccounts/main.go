// Merges accounts records based pn zoo name.

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/simpleset"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"time"
)

var user = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Required().String()

type account struct {
	accounts *simpleset.Set
	id       string
}

func newAccount(id string) *account {
	// Returns initialized struct
	a := new(account)
	a.accounts = simpleset.NewStringSet()
	a.id = id
	return a
}

func (a *account) addAccount(name string) {
	// Adds name to accounts set
	a.accounts.Add(name)
}

func (a *account) getAccounts() string {
	// Returns sources as comma seperated string
	a.accounts.Pop(a.id)
	return strings.Join(a.accounts.ToStringSlice(), ",")
}

type zoos struct {
	accounts map[string]*account
	db       *dbIO.DBIO
	names    *simpleset.Set
	source   []string
}

func newZoos() *zoos {
	// Returns new struct
	z := new(zoos)
	z.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	z.accounts = make(map[string]*account)
	z.names = simpleset.NewStringSet()
	for _, i := range z.db.GetRows("Source", "Zoo", "1", "Zoo,account_id") {
		z.source = append(z.source, i[0])
	}
	return z
}

func (z *zoos) setAccounts() {
	// Stores patient ids by taxa
	fmt.Println("\n\tStoring account IDs by species...")
	for _, i := range z.db.GetRows("Accounts", "account_id", strings.Join(z.source, ","), z.db.Columns["Accounts"]) {
		id := i[0]
		name := i[2]
		if _, ex := z.accounts[name]; !ex {
			z.accounts[name] = newAccount(id)
		} else {
			z.accounts[name].addAccount(id)
		}
	}
}

func (z *zoos) updateAccounts() {
	// Replaces redundant account ids in account table
	fmt.Println("Updating Accounts table...")
	count := 0
	for _, i := range z.accounts {
		count++
		z.db.DeleteRows("Accounts", "account_id", i.accounts.ToStringSlice())
		z.db.UpdateRow("Accounts", "Account", "NA", "account_id", "=", i.id)
		fmt.Printf("\r\tUpdated %d of %d account ids.", count, len(z.accounts))
	}
	fmt.Println()
}

func (z *zoos) updateSource() {
	// Updates redundant account ids in source table
	fmt.Println("Updating Source table...")
	count := 0
	for _, i := range z.accounts {
		count++
		if err := z.db.UpdateRow("Source", "account_id", i.id, "account_id", "IN", i.getAccounts()); !err {
			fmt.Println("\t[Warning] Failed to update Source.")
		} else {
			fmt.Printf("\r\tUpdated %d of %d account ids.", count, len(z.accounts))
		}
	}
	fmt.Println()
}

func main() {
	start := time.Now()
	kingpin.Parse()
	z := newZoos()
	z.setAccounts()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

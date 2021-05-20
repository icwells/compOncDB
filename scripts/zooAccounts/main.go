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

func (a *account) length() int {
	// Returns length of accounts set
	return a.accounts.Length()
}

func (a *account) getAccounts() string {
	// Returns sources as comma seperated string
	a.accounts.Pop(a.id)
	return strings.Join(a.accounts.ToStringSlice(), ",")
}

type zoos struct {
	accounts map[string]*account
	db       *dbIO.DBIO
	source   []string
}

func newZoos() *zoos {
	// Returns new struct
	z := new(zoos)
	z.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	z.accounts = make(map[string]*account)
	for _, i := range z.db.GetRows("Source", "Zoo", "1", "account_id,Zoo") {
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
	for k, v := range z.accounts {
		// Remove records without multiple ids
		if v.length() == 0 {
			delete(z.accounts, k)
		}
	}
}

func (z *zoos) update(table, command string) {
	// Executes command
	cmd, err := z.db.DB.Prepare(command)
	if err != nil {
		fmt.Printf("\n\t[Error] Preparing command for %s: %v\n", table, err)
	} else {
		_, err = cmd.Exec()
		cmd.Close()
		if err != nil {
			fmt.Printf("\n\t[Error] Executing command on %s: %v\n", table, err)
		}
	}
}


func (z *zoos) updateAccounts() {
	// Replaces redundant account ids in account table
	fmt.Println("Updating Accounts IDs...")
	count := 0
	for _, i := range z.accounts {
		count++
		z.update("Source", fmt.Sprintf("UPDATE Source SET account_id = %s WHERE account_id IN (%s);", i.id, i.getAccounts()))
		z.update("Accounts", fmt.Sprintf("DELETE FROM Accounts WHERE account_id IN (%s);",i.getAccounts()))
		z.update("Accounts", fmt.Sprintf("UPDATE Accounts SET Account = 'NA' WHERE account_id = %s;", i.id))
		fmt.Printf("\r\tUpdated %d of %d account ids.", count, len(z.accounts))
	}
	fmt.Println()
}

func main() {
	start := time.Now()
	kingpin.Parse()
	z := newZoos()
	z.setAccounts()
	z.updateAccounts()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

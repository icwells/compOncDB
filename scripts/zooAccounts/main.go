// Merges accounts records based on zoo name.

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"gopkg.in/alecthomas/kingpin.v2"
	"strconv"
	"time"
)

var user = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Required().String()

type account struct {
	id   int
	name string
	old  []int
}

func newAccount(id int, name string) *account {
	// Returns initialized struct
	a := new(account)
	a.id = id
	a.name = name
	return a
}

func (a *account) addID(id int) {
	// Adds old id to replace
	if id < a.id {
		// Keep lowest number
		a.old = append(a.old, a.id)
		a.id = id
	} else {
		a.old = append(a.old, id)
	}
}

type zoos struct {
	accounts map[string]*account
	db       *dbIO.DBIO
	orig     int
}

func newZoos() *zoos {
	// Returns new struct
	z := new(zoos)
	z.accounts = make(map[string]*account)
	z.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	return z
}

func (z *zoos) setAccounts() {
	// Stores patient ids by taxa
	fmt.Println("\n\tStoring account IDs by species...")
	for _, i := range z.db.GetTable("Accounts") {
		z.orig++
		id, _ := strconv.Atoi(i[0])
		name := i[1]
		if _, ex := z.accounts[name]; !ex {
			z.accounts[name] = newAccount(id, name)
		} else {
			// Store ids
			z.accounts[name].addID(id)
		}
	}
	fmt.Printf("\tFound %d unique submitter names from a total of %d.\n", len(z.accounts), z.orig)
	for k, v := range z.accounts {
		if len(v.old) == 0 {
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
	var count int
	fmt.Println("\tUpdating Accounts IDs...")
	for _, v := range z.accounts {
		count++
		for _, i := range v.old {
			z.update("Accounts", fmt.Sprintf("UPDATE Source SET account_id = %d WHERE account_id = %d;", v.id, i))
			z.update("Accounts", fmt.Sprintf("DELETE FROM Accounts WHERE account_id = %d;", i))
		}
		//z.update("Accounts", fmt.Sprintf("UPDATE Accounts SET account_id = %d WHERE submitter_name =\"%s\";", v.id, v.name))
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

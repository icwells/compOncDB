// Prints name and id of zoos which need to have approval updated.

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"time"
)

var (
	infile = kingpin.Flag("infile", "Path to input file.").Short('i').Default("nil").String()
	user   = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type approval struct {
	approved []string
	db       *dbIO.DBIO
	na       []string
	rejected []string
	sources  map[string][]string
}

func newApproval() *approval {
	// Returns new struct
	a := new(approval)
	a.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	fmt.Println("\n\tInitializing struct...")
	cmd := "SELECT DISTINCT(Accounts.submitter_name),Source.account_id,Source.Approved FROM Source JOIN Accounts ON Accounts.account_id = Source.account_id;"
	a.sources = codbutils.ToMap(a.db.Execute(cmd))
	return a
}

func (a *approval) setApproved() {
	// Stores zoo approvals from input file
	fmt.Println("\tReading input file...")
	input, header := iotools.YieldFile(*infile, true)
	for i := range input {
		approved := strings.ToLower(i[header["Approved"]])
		if strings.Contains(approved, "yes") || approved == "1" {
			a.approved = append(a.approved, i[header["Zoo"]])
		} else if approved == "no" || approved == "0" {
			a.rejected = append(a.rejected, i[header["Zoo"]])
		}
	}
}

func (a *approval) checkApprovals() {
	// Determines if any approvals need to be updated
	for idx, list := range [][]string{a.approved, a.rejected} {
		if idx == 0 {
			fmt.Println("\n\tAccounts to approve:")
		} else {
			fmt.Println("\n\tAccounts to remove:")
		}
		var found bool
		for _, i := range list {
			if val, ex := a.sources[i]; ex {
				if val[1] != "1" {
					fmt.Printf("\t\t%s %s\n", val[0], i)
					found = true
				}
			}
		}
		if !found {
			fmt.Println("\t\tNone")
		}
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	a := newApproval()
	a.setApproved()
	a.checkApprovals()
	fmt.Printf("\n\tFinished. Runtime: %s\n\n", time.Since(start))
}

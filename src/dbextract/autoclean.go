// Removes entries form sub-tables if master entry has been deleted

package dbextract

import (
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/strarray"
)

type cleaner struct {
	db     *dbIO.DBIO
	pids   strarray.Set
	pchild strarray.Set
	tids   strarray.Set
	tchild strarray.Set
	aids   strarray.Set
	achild strarray.Set
}

func (c *cleaner) setIDs(col string, tables []string) strarray.Set {
	// Returns set of ids from tables
	ret := strarray.NewSet()
	for _, i := range tables {
		ret.Extend(c.db.GetColumnText(i, col))
	}
	return ret
}

func newCleaner(db *dbIO.DBIO) *cleaner {
	// Initializes struct
	var c cleaner
	c.db = db
	c.pids = c.setIDs("ID", []string{"Patient"})
	c.pchild = c.setIDs("ID", []string{"Diagnosis", "Tumor", "Source"})
	c.tids = c.setIDs("taxa_id", []string{"Taxonomy", "Denominators", "Patient", "Life_history"})
	// Skip totals table since it will be overwritten
	c.tchild = c.setIDs("taxa_id", []string{"Common"})
	c.aids = c.setIDs("account_id", []string{"Accounts"})
	c.achild = c.setIDs("account_id", []string{"Source"})
	return &c
}

func (c *cleaner) cleanTables(col string, tables []string, parent, child strarray.Set) {
	// Removes records from target tables if id is present in child but not parent
	var rm []string
	for _, i := range child.ToSlice() {
		if parent.InSet(i) == false {
			// Record extraneous ids
			rm = append(rm, i)
		}
	}
	if len(rm) > 0 {
		for _, i := range tables {
			c.db.DeleteRows(i, col, rm)
		}
	}
}

func AutoCleanDatabase(db *dbIO.DBIO) {
	// Cleans database and recalcutates species totals
	c := newCleaner(db)
	c.cleanTables("ID", []string{"Diagnosis", "Tumor", "Source"}, c.pids, c.pchild)
	c.cleanTables("taxa_id", []string{"Common"}, c.tids, c.tchild)
	c.cleanTables("account_id", []string{"Source"}, c.tids, c.tchild)
	dbupload.SpeciesTotals(db)
}

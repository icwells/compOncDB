// Removes entries form sub-tables if master entry has been deleted

package dbextract

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/simpleset"
)

type cleaner struct {
	db     *dbIO.DBIO
	pids   *simpleset.Set
	pchild *simpleset.Set
	tids   *simpleset.Set
	tchild *simpleset.Set
	aids   *simpleset.Set
	achild *simpleset.Set
}

func (c *cleaner) setIDs(col string, tables []string) *simpleset.Set {
	// Returns set of ids from tables
	ret := simpleset.NewStringSet()
	for _, t := range tables {
		for _, i := range c.db.GetColumnText(t, col) {
			ret.Add(i)
		}
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

func (c *cleaner) cleanTables(col string, tables []string, parent, child *simpleset.Set) {
	// Removes records from target tables if id is present in child but not parent
	var rm []string
	for _, i := range child.ToStringSlice() {
		if ex, _ := parent.InSet(i); !ex {
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
	count := dbupload.FilterPatients(db)
	if count > 0 {
		codbutils.GetLogger().Printf("Removed %d duplicate patient records.\n", count)
	}
	c := newCleaner(db)
	c.cleanTables("ID", []string{"Diagnosis", "Tumor", "Source"}, c.pids, c.pchild)
	c.db.OptimizeTables()
	c.cleanTables("taxa_id", []string{"Common"}, c.tids, c.tchild)
	c.cleanTables("account_id", []string{"Source"}, c.tids, c.tchild)
}

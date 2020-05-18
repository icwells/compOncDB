// Updates taxonmy entries in place

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"strconv"
	"strings"
)

type tableupdate struct {
	table  string
	target string
	values map[string]map[string]string
}

func newTableUpdate(table, target string) *tableupdate {
	// Initializes new update struct
	t := new(tableupdate)
	t.table = table
	t.target = target
	t.values = make(map[string]map[string]string)
	return t
}

func (t *tableupdate) add(id string, col string, val string) {
	// Adds value to values
	if _, ex := t.values[col]; ex == false {
		t.values[col] = make(map[string]string)
	}
	t.values[col][id] = val
}

func (t *tableupdate) updateTable(db *dbIO.DBIO) {
	// Uploads contents to database
	pass := db.UpdateColumns(t.table, t.target, t.values)
	if pass == false {
		fmt.Printf("\t[Warning] Failed to upload to %s.\n", t.table)
	} else {
		fmt.Printf("\tUpdated %d records in %s.\n", len(t.values), t.table)
	}
}

//----------------------------------------------------------------------------

type updater struct {
	db      *dbIO.DBIO
	delim   string
	target  string
	header  map[int]string
	columns map[string]string
	tables  map[string]*tableupdate
	taxa    map[string]string
}

func newUpdater(db *dbIO.DBIO) updater {
	// Initializes new update struct
	var u updater
	u.db = db
	u.header = make(map[int]string)
	u.columns = db.Columns
	u.tables = make(map[string]*tableupdate)
	u.taxa = dbupload.EntryMap(u.db.GetColumns("Taxonomy", []string{"taxa_id", "Species"}))
	return u
}

func (u *updater) setHeader(line string) {
	// Stores input file columns to database tables and columns
	u.delim, _ = iotools.GetDelim(line)
	for idx, i := range strings.Split(line, u.delim) {
		i = strings.TrimSpace(i)
		if len(i) > 0 {
			if idx == 0 {
				// Store target column for identification
				u.target = i
			} else {
				if strings.ToUpper(i) == "ID" {
					u.header[idx] = "ID"
				} else if strings.Contains(i, "_") == true {
					u.header[idx] = strings.ToLower(i)
				} else {
					u.header[idx] = strarray.TitleCase(i)
				}
			}
		}
	}
}

func (u *updater) checkTaxaID(id string) string {
	// Replaces scentific name with taxa id
	if _, err := strconv.Atoi(id); err != nil {
		tid, ex := u.taxa[id]
		if ex == true {
			id = tid
		}
	}
	return id
}

func (u *updater) evaluateRow(row []string) {
	// Assigns row values to substruct
	id := strings.TrimSpace(row[0])
	if len(id) >= 1 {
		for k, v := range u.header {
			if k < len(row) {
				val := strings.TrimSpace(row[k])
				if len(val) >= 1 {
					table := codbutils.GetTable(u.columns, v)
					if _, ex := u.tables[table]; ex == false {
						// Initialize new struct
						u.tables[table] = newTableUpdate(table, u.target)
					}
					// Add new value
					u.tables[table].add(u.checkTaxaID(id), v, val)
				}
			}
		}
	}
}

func (u *updater) getUpdateFile(infile string) {
	// Returns map of data to be updated
	first := true
	fmt.Println("\n\tReading input file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Text()))
		if first == false {
			u.evaluateRow(strings.Split(line, u.delim))
		} else {
			u.setHeader(line)
			first = false
		}
	}
}

func (u *updater) updateTables() {
	// Updates database with all identified values
	fmt.Println("\tUpdating tables...")
	for _, v := range u.tables {
		if v.table == "Accounts" {
			fmt.Print("\n\t[Warning] Skipping Accounts table.\n\n")
		} else {
			v.updateTable(u.db)
		}
	}
}

//----------------------------------------------------------------------------

func UpdateEntries(db *dbIO.DBIO, infile string) {
	// Updates taxonomy entries
	u := newUpdater(db)
	u.getUpdateFile(infile)
	u.updateTables()
}

func UpdateSingleTable(db *dbIO.DBIO, table, column, value, target, op, key string) {
	// Updates single table
	fmt.Printf("\n\tUpdating %s...\n", table)
	c := db.UpdateRow(table, column, value, target, op, key)
	if c == false {
		fmt.Printf("\t[Warning] Failed to upload to %s.\n", table)
	} else {
		fmt.Printf("\tSuccessfully updated %s.\n", table)
	}
}

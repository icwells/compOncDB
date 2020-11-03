// Updates taxonmy entries in place

package dbextract

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"log"
	"strconv"
	"strings"
)

type tableupdate struct {
	logger *log.Logger
	table  string
	target string
	total  int
	values map[string]map[string]string
}

func newTableUpdate(logger *log.Logger, table, target string) *tableupdate {
	// Initializes new update struct
	t := new(tableupdate)
	t.logger = logger
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
	t.total++
}

func (t *tableupdate) updateTable(db *dbIO.DBIO) {
	// Uploads contents to database
	pass := db.UpdateColumns(t.table, t.target, t.values)
	if pass == false {
		t.logger.Printf("[Warning] Failed to upload to %s.\n", t.table)
	} else {
		t.logger.Printf("Updated %d records in %s.\n", t.total, t.table)
	}
}

//----------------------------------------------------------------------------

type Updater struct {
	columns map[string]string
	db      *dbIO.DBIO
	delim   string
	header  map[int]string
	logger  *log.Logger
	tables  map[string]*tableupdate
	target  string
	taxa    map[string]string
}

func NewUpdater(db *dbIO.DBIO) Updater {
	// Initializes new update struct
	var u Updater
	u.db = db
	u.columns = db.Columns
	u.header = make(map[int]string)
	u.logger = codbutils.GetLogger()
	u.tables = make(map[string]*tableupdate)
	u.taxa = codbutils.EntryMap(u.db.GetColumns("Taxonomy", []string{"taxa_id", "Species"}))
	return u
}

func (u *Updater) ClearTables() {
	// Empties tables map
	u.tables = make(map[string]*tableupdate)
}

func (u *Updater) SetHeader(line string) {
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

func (u *Updater) checkTaxaID(id string) string {
	// Replaces scentific name with taxa id
	if _, err := strconv.Atoi(id); err != nil {
		tid, ex := u.taxa[id]
		if ex == true {
			id = tid
		}
	}
	return id
}

func (u *Updater) EvaluateRow(row []string) {
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
						u.tables[table] = newTableUpdate(u.logger, table, u.target)
					}
					// Add new value
					u.tables[table].add(u.checkTaxaID(id), v, val)
				}
			}
		}
	}
}

func (u *Updater) getUpdateFile(infile string) {
	// Returns map of data to be updated
	first := true
	u.logger.Println("Reading input file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Text()))
		if first == false {
			u.EvaluateRow(strings.Split(line, u.delim))
		} else {
			u.SetHeader(line)
			first = false
		}
	}
}

func (u *Updater) UpdateTables() {
	// Updates database with all identified values
	u.logger.Println("Updating tables...")
	for _, v := range u.tables {
		if v.table == "Accounts" {
			u.logger.Print("[Warning] Skipping Accounts table.\n\n")
		} else {
			v.updateTable(u.db)
		}
	}
}

//----------------------------------------------------------------------------

func UpdateEntries(db *dbIO.DBIO, infile string) {
	// Updates taxonomy entries
	u := NewUpdater(db)
	u.getUpdateFile(infile)
	u.UpdateTables()
}

func UpdateSingleTable(db *dbIO.DBIO, table, column, value, target, op, key string) {
	// Updates single table
	logger := codbutils.GetLogger()
	logger.Printf("Updating %s...\n", table)
	c := db.UpdateRow(table, column, value, target, op, key)
	if c == false {
		logger.Printf("[Warning] Failed to upload to %s.\n", table)
	} else {
		logger.Printf("Successfully updated %s.\n", table)
	}
}

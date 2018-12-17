// Updates taxonmy entries in place

package dbextract

import (
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"os"
	"strings"
)

type tableupdate struct {
	table  string
	target string
	values map[string][]string
}

func newTableUpdate(table, target string) *tableupdate {
	// Initializes new update struct
	t := new(tableupdate)
	t.table = table
	t.target = target
	t.values = make(map[string][]string)
	return t
}

func (t *tableupdate) add(id string, row []string) {
	// Adds row to values
	t.values[id] = row
}

func (t *tableupdate) updateTable(db *dbIO.DBIO) {
	// Uploads contents to database
	c := db.UpdateRows(t.table, t.target, t.values)
	if c == 0 {
		fmt.Printf("\t[Warning] Failed to upload to %s.\n", t.table)
	} else {
		fmt.Printf("\tUpdated %d records in %s.\n", c, t.table)
	}
}

//----------------------------------------------------------------------------

type updater struct {
	columns map[string][]int
	col     map[string]string
	tables  map[string]*tableupdate
}

func newUpdater(col map[string]string) updater {
	// Initializes new update struct
	var u updater
	u.columns = make(map[string][]int)
	u.col = col
	u.tables = make(map[string]*tableupdate)
	return u
}

func (u *updater) formatHeader(row []string) []string {
	// Ensures proper formatting of input header values
	var ret []string
	for _, i := range row {
		if strings.ToUpper(i) == "ID" {
			i = "ID"
		} else if strings.Contains(i, "_") == true {
			i = strings.ToLower(i)
		} else {
			i = strings.Title(i)
		}
		ret = append(ret, strings.TrimSpace(i))
	}
	return ret
}

func (u *updater) setColumns(row []string) {
	// Correlates input file columns to database tables and columns
	keep := false
	for k, v := range u.col {
		head := strings.Split(v, ",")
		// Initialize new column header and fill (missing values have an index of -1)
		for _, i := range head {
			// Store file header index in index of database table column
			ind := strarray.SliceIndex(row, i)
			u.columns[k] = append(u.columns[k], ind)
			if ind > 0 {
				keep = true
			}
		}
		if keep == false {
			// Remove empty tables
			delete(u.columns, k)
		} else {
			// Initialize struct
			u.tables[k] = newTableUpdate(k, row[0])
		}
	}
}

func (u *updater) evaluateRow(row []string) {
	// Stores row in appriate substructs
	id := row[0]
	if len(id) >= 1 {
		for k, v := range u.columns {
			var line []string
			for _, i := range v {
				if i <= 0 {
					// Skip empty/ID fields
					line = append(line, "")
				} else if len(row[i]) < 1 {
					line = append(line, "")
				} else {
					line = append(line, row[i])
				}
			}
			u.tables[k].add(id, line)
		}
	}
}

func (u *updater) getUpdateFile(infile string) {
	// Returns map of data to be updated
	var d string
	first := true
	fmt.Println("\n\tReading input file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Text()))
		if first == false {
			u.evaluateRow(strings.Split(line, d))
		} else {
			d = iotools.GetDelim(line)
			u.setColumns(strings.Split(line, d))
			first = false
		}
	}
	fmt.Println(u.columns)
	os.Exit(0)
}

func (u *updater) updateTables(db *dbIO.DBIO) {
	// Updates database with all identified values
	fmt.Println("\tUpdating tables...")
	for _, v := range u.tables {
		v.updateTable(db)
	}
}

//----------------------------------------------------------------------------

func UpdateEntries(db *dbIO.DBIO, infile string) {
	// Updates taxonomy entries
	u := newUpdater(db.Columns)
	u.getUpdateFile(infile)
	u.updateTables(db)
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

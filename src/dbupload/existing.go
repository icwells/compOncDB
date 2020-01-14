// Defines struct for identifying redundant

package dbupload

import (
	"github.com/icwells/dbIO"
)

type Entry struct {
	Age  string
	Taxa string
	Date string
}

func NewEntry(row []string) *Entry {
	// Returns empty struct
	e := new(Entry)
	e.Age = row[0]
	e.Taxa = row[1]
	e.Date = row[2]
	return e
}

type Existing struct {
	Entries   map[string]map[string]*Entry
	Unmatched map[string]*Entry
}

func (e *Existing) setUnmatched(db *dbIO.DBIO) {
	// Populates unmathced map
	for _, i := range db.GetColumns("Unmatched", []string{"sourceID", "age", "name", "date"}) {
		e.Unmatched[i[0]] = NewEntry(i[1:])
	}
}

func (e *Existing) setEntries(db *dbIO.DBIO) {
	// Populates entries map
	accounts := EntryMap(db.GetColumns("Source", []string{"account_id", "ID"}))
	e.Entries["-1"] = make(map[string]*Entry)
	for _, v := range accounts {
		// Initialize inner maps
		e.Entries[v] = make(map[string]*Entry)
	}
	for _, i := range db.GetColumns("Patient", []string{"source_id", "Age", "taxa_id", "Date"}) {
		id := i[0]
		acc, ex := accounts[id]
		if !ex {
			acc = "-1"
		}
		e.Entries[acc][id] = NewEntry(i[1:])
	}
}

func NewExisting(db *dbIO.DBIO) *Existing {
	// Populates and returns struct
	e := new(Existing)
	e.Entries = make(map[string]map[string]*Entry)
	e.Unmatched = make(map[string]*Entry)
	if db != nil {
		e.setEntries(db)
		e.setUnmatched(db)
	}
	return e
}

func (e *Existing) Exists(acc, id, age, taxa, date string) bool {
	// Returns true if given record is in map
	if acc == "" && taxa == "" {
		// Check unmatched table
		if row, ex := e.Unmatched[id]; ex {
			if row.Age == age && row.Date == date {
				return true
			}
		}
	} else {
		// Check patient table
		if _, exists := e.Entries[acc]; exists {
			if row, ex := e.Entries[acc][id]; ex {
				if row.Age == age && row.Taxa == taxa && row.Date == date {
					return true
				}
			}
		}
	}
	return false
}

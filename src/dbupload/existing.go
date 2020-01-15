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
	db        *dbIO.DBIO
	Entries   map[string]map[string]*Entry
	Unmatched map[string]*Entry
	accounts  map[string]string
	ids       map[string][]string
}

func (e *Existing) setUnmatched() {
	// Populates unmathced map
	for _, i := range e.db.GetColumns("Unmatched", []string{"sourceID", "age", "name", "date"}) {
		e.Unmatched[i[0]] = NewEntry(i[1:])
	}
}

func (e *Existing) setEntries() {
	// Populates entries map
	e.Entries["-1"] = make(map[string]*Entry)
	for _, v := range e.accounts {
		// Initialize inner maps
		e.Entries[v] = make(map[string]*Entry)
	}
	for _, i := range e.db.GetColumns("Patient", []string{"source_id", "Age", "taxa_id", "Date"}) {
		id := i[0]
		acc, ex := e.accounts[id]
		if !ex {
			acc = "-1"
		}
		e.Entries[acc][id] = NewEntry(i[1:])
	}
}

func NewExisting(db *dbIO.DBIO) *Existing {
	// Populates and returns struct
	e := new(Existing)
	e.db = db
	e.Entries = make(map[string]map[string]*Entry)
	e.Unmatched = make(map[string]*Entry)
	if db != nil {
		e.accounts = EntryMap(db.GetColumns("Source", []string{"account_id", "ID"}))
		e.setEntries()
		e.setUnmatched()
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

func (e *Existing) setIDs() []string {
	// Stores map of source ids to uids
	var ret []string
	e.ids = make(map[string][]string)
	for _, i := range e.db.GetColumns("Patient", []string{"ID", "source_id", "Age", "taxa_id", "Date"}) {
		if acc, ex := e.accounts[i[1]]; ex {
			if e.Exists(acc, i[1], i[2], i[3], i[4]) {
				e.ids[i[1]] = append(e.ids[i[1]], i[0])
			}
		}
	}
	for k, v := range e.ids {
		if len(v) == 1 {
			delete(e.ids, k)
		} else {
			ret = append(ret, v[1:]...)
		}
	}
	return ret
}

func FilterPatients(db *dbIO.DBIO) int {
	// Removes duplicate patient entries
	e := NewExisting(db)
	rm := e.setIDs()
	e.db.DeleteRows("Patient", "ID", rm)
	return len(e.ids)
}

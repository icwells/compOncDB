// This script defines a struct for identifying duplicate entries

package parserecords

import (
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"strings"
)

type duplicates struct {
	count   int
	ids     map[string][]string
	records map[string]map[string]record
	reps    map[string][]string
}

func newDuplicates() duplicates {
	// Makes duplicates maps
	var d duplicates
	d.ids = make(map[string][]string)
	d.records = make(map[string]map[string]record)
	d.reps = make(map[string][]string)
	return d
}

func (d *duplicates) add(id, source string) {
	// Stores all id combinations in ids and repeats in reps
	if strings.ToUpper(id) != "NA" && strings.ToUpper(source) != "NA" {
		row, ex := d.ids[source]
		if ex == true {
			if strarray.InSliceStr(row, id) == true {
				d.count++
				r, e := d.reps[source]
				if e == true {
					if strarray.InSliceStr(r, id) == false {
						// Add to rep id slice
						d.reps[source] = append(d.reps[source], id)
					}
				} else {
					// Make new rep entry
					d.reps[source] = []string{id}
				}
			} else {
				// Add to id slice
				d.ids[source] = append(d.ids[source], id)
			}
		} else {
			// Make new id map entry
			d.ids[source] = []string{id}
		}
	}
}

func (e *entries) getDuplicates(infile string) {
	// Identifies duplicate entries
	first := true
	e.dups = newDuplicates()
	e.dupsPresent = true
	e.logger.Println("Identifying duplicate entries...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		if first == false {
			s := strings.Split(line, e.d)
			if len(s) > e.col.max {
				if e.col.patient >= 0 {
					pid := s[e.col.patient]
					acc := s[e.col.submitter]
					e.dups.add(pid, acc)
				}
			}
		} else {
			e.parseHeader(line)
			first = false
		}
	}
	e.logger.Printf("Found %d duplicate patients.\n", e.dups.count)
}

func (e *entries) resolveDuplicates(rec record) {
	// Determines whether to replace existing record with input record
	mp, exists := e.dups.records[rec.submitter]
	if exists == true {
		row, ex := mp[rec.patient]
		if ex == true {
			if row.massPresent != "1" {
				if rec.massPresent == "1" {
					// Only replace is stored record is not a cancer record and new one is
					e.dups.records[rec.submitter][rec.patient] = rec
				}
			}
		} else {
			// Store new patient
			e.dups.records[rec.submitter][rec.patient] = rec
		}
	} else {
		// Make new map
		e.dups.records[rec.submitter] = make(map[string]record)
		e.dups.records[rec.submitter][rec.patient] = rec
	}
}

func (e *entries) inDuplicates(rec record) bool {
	// Returns true if rec is a duplicate record
	ret := false
	row, ex := e.dups.reps[rec.submitter]
	if ex == true {
		if strarray.InSliceStr(row, rec.patient) == true {
			ret = true
		}
	}
	return ret
}

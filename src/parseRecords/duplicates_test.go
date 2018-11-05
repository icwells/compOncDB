// Performs white box tests on the duplicates stuct's methods

package main

import (
	"fmt"
	"github.com/icwells/go-tools/strarray"
	"testing"
)

type source struct {
	source	string
	id		string
	mass	string
	repeat	bool
	indups	bool
	resolve	bool
}

func newSources() []source {
	// Returns a slice of sources to test
	return []source{
		{"KV Zoo", "102", "1", false, true, true},
		{"KV Zoo", "95", "0", false, false, false},
		{"KV Zoo", "102", "0", true, true, false},
		{"XY Zoo", "95", "0", false, true, false},
		{"XY Zoo", "95", "1", true, true, true},
		{"AB", "300", "1", false, true, true},
		{"AB", "300", "1", true, true, false},
	}
}

func TestAdd(t *testing.T) {
	// Tests add method
	e := newEntries("service")
	sources := newSources()
	for _, i := range sources {
		l := len(e.dups.reps)
		e.dups.add(i.id, i.source)
		if i.repeat == false {
			// Make sure reps map has not changed and new info is added to ids map
			if len(e.dups.reps) != l {
				t.Error("Length of reps map changed with novel entry.")
			}
			row, ex := e.dups.ids[i.source]
			if ex == false {
				t.Error("New source not added to ids map.")
			} else if strarray.InSliceStr(row, i.id) == false {
				t.Error("New id not added to ids map.")
			}
		} else {
			// Make sure new info is added to reps map
			if len(e.dups.reps) != l + 1 {
				t.Error("Length of reps map did not change with novel entry.")
			}
			row, ex := e.dups.reps[i.source]
			if ex == false {
				t.Error("New source not added to reps map.")
			} else if strarray.InSliceStr(row, i.id) == false {
				t.Error("New id not added to reps map.")
			}			
		}
	}
}

func TestResolveDuplicates(t *testing.T) {
	// Tests resolve duplicates method
	e := newEntries("service")
	sources := newSources()
	for _, i := range sources {
		// Load maps
		e.dups.add(i.id, i.source)
	}
	for _, i := range sources {
		r := newRecord()
		r.patient = i.id
		r.submitter = i.source
		r.massPresent = i.mass
		e.resolveDuplicates(r)
		if i.resolve == true {
			//Compare each record as updated
			msg := compareRecords(e.dups.records[i.source][i.id], r)
			if len(msg) > 1 {
				t.Error(msg)
			}
		}
	}
}

func TestInDuplicates(t *testing.T) {
	// Tests inDuplicates method
	e := newEntries("service")
	sources := newSources()
	for _, i := range sources {
		// Load maps
		e.dups.add(i.id, i.source)
	}
	for _, i := range sources {
		r := newRecord()
		r.patient = i.id
		r.submitter = i.source
		actual := e.inDuplicates(r)
		if actual != i.indups {
			msg := fmt.Sprintf("inDuplucates returned %v were %v was expected.", actual, i.indups)
			t.Error(msg)
		}
	}
}

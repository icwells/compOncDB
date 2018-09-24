// This script defines a struct for identifying duplicate entries

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"os"
	"strings"
)

type duplicates struct {
	ids		map[string][]string
	reps	map[string][]string
}

func newDuplicates() duplicates {
	// Makes duplicates maps
	var d duplicates
	d.ids = make(map[string][]string)
	d.reps = make(map[string][]string)
}

func (d *duplicates) add(id, source string) {
	// Stores all id combinations in ids and repeats in reps
	row, ex := ids[source]
	if ex == true {
		if strarray.InSliceSli(row, id) == true {
			_, e := reps[source]
			if e == true {
				if strarray.InSliceSli(d.reps, id) == false {
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

func (e *entries) getDuplicates(infile string) {
	// Identifies duplicate entries
	first := true
	e.dups = newDuplicates()
	fmt.Println("\tIdentifying duplicate entries...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			s := strings.Split(line, d)
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
}

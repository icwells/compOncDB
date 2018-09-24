// This script defines a struct for identifying duplicate entries

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"os"
	"strings"
)

type duplicates struct {
	ids		map[string][]string
	reps	map[string][]string
}

func (e *entries) getDuplicates(infile string) {
	// Identifies duplicate entries
	first := true
	e.diag = make(map[string][]string)
	e.diagPresent = true
	fmt.Println("\tIdentifying duplicate entries...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			s := strings.Split(line, d)
			if len(s) > e.col.max {
				id := s[e.col.id]
			}
		} else {
			e.parseHeader(line)
			first = false
		}
	}
}

// This script defines a struct for managing comparative oncology records

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"os"
	"strings"
)

type entries struct {
	d           string
	col         columns
	service     string
	taxa        map[string]string
	diag        map[string][]string
	match		matcher
	dups		duplicates
	taxaPresent bool
	diagPreset  bool
	dupsPresent	bool
}

func newEntries(service string) entries {
	// Initializes new struct
	var e entries
	e.service = service
	e.col = newColumns()
	e.taxaPresent = false
	e.diagPResent = false
	e.dupsPresent = false
	return e
}

func (e *entries) parseHeader(header string) {
	// Stores column numbers and delimiter from header
	e.d = getDelim(header)
	head := strings.Split(header, e.d)
	e.col.setColumns(head)
}

func (e *entries) getTaxonomy(infile string) {
	// Reads in map of species
	var d string
	first := true
	e.taxa = make(map[string]string)
	e.taxaPresent = true
	fmt.Println("\tReading taxonomy file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			s := strings.Split(line, d)
			// Store binomial for each search term
			e.taxa[s[0]] = s[8]
		} else {
			d = getDelim(line)
			first = false
		}
	}
}

func (e *entries) getDiagnosis(infile) {
	// Reads in diagnosis data
	var d string
	first := true
	e.diag = make(map[string][]string)
	e.diagPresent = true
	fmt.Println("\tReading diagnosis file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			s := strings.Split(line, d)
			// Store daignosis data by ids
			e.diag[s[0]] = s[1:]
		} else {
			d = getDelim(line)
			first = false
		}
	}
}

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
	taxa        map[string][]string
	diag        map[string][]string
	match       matcher
	dups        duplicates
	dupsPresent bool
}

func newEntries(service string) entries {
	// Initializes new struct
	var e entries
	e.service = service
	e.col = newColumns()
	e.match = newMatcher()
	e.dups = newDuplicates()
	e.dupsPresent = false
	return e
}

func (e *entries) parseHeader(header string) {
	// Stores column numbers and delimiter from header
	e.d = iotools.GetDelim(header)
	e.col.setColumns(strings.Split(header, e.d))
}

func (e *entries) getOutputFile(outfile, header string) *os.File {
	// Opends file for appending
	var f *os.File
	if iotools.Exists(outfile) == true {
		f = iotools.AppendFile(outfile)
	} else {
		f = iotools.CreateFile(outfile)
		f.WriteString(header)
	}
	return f
}

func (e *entries) getTaxonomy(infile string) {
	// Reads in map of species
	var d string
	col := make(map[string]int)
	first := true
	e.taxa = make(map[string][]string)
	fmt.Println("\tReading taxonomy file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		if first == false {
			s := strings.Split(line, d)
			e.taxa[s[0]] = []string{s[col["Family"]], s[col["Genus"]], s[col["Species"]]}
		} else {
			d = iotools.GetDelim(line)
			for idx, i := range strings.Split(line, d) {
				col[strings.TrimSpace(i)] = idx
			}
			first = false
		}
	}
}

func (e *entries) getDiagnosis(infile string) {
	// Reads in diagnosis data
	var d string
	first := true
	e.diag = make(map[string][]string)
	fmt.Println("\tReading diagnosis file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		if first == false {
			s := strings.Split(line, d)
			// Store daignosis data by ids
			e.diag[s[0]] = s[1:]
		} else {
			d = iotools.GetDelim(line)
			first = false
		}
	}
}

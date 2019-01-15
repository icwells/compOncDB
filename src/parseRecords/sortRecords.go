// This script defines functions for sorting entries data

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"strings"
)

func subsetLine(idx int, line []string) string {
	// Returns line[idx]/NA
	ret := "NA"
	if idx > -1 && idx < len(line) {
		ret = strings.TrimSpace(line[idx])
		if len(ret) <= 0 {
			ret = "NA"
		}
	}
	return ret
}

func (e *entries) sortLine(line []string) (record, bool) {
	// Returns formatted string and true if it should be written
	write := false
	rec := newRecord()
	var idx int
	if e.col.common >= 0 {
		// Get common name if present
		idx = e.col.common
	} else if e.col.species >= 0 {
		idx = e.col.species
	} else {
		printFatal("Cannot determine species column", 20)
	}
	if len(line) >= e.col.max && len(line[idx]) >= 3 && strings.ToUpper(line[idx]) != "N/A" {
		// Proceed if line is properly formatted and species is present and no NA
		id := subsetLine(e.col.id, line)
		rec.setID(id)
		if id != "NA" {
			row, ex := e.diag[rec.id]
			if ex == true {
				// Assign diagnosis data if id is present in map
				rec.setDiagnosis(row)
			}
		}
		// Replace entry with scientific name
		sp, ex := e.taxa[line[idx]]
		if ex == true {
			rec.setSpecies(sp)
		}
		rec.setDate(subsetLine(e.col.date, line))
		rec.setComments(subsetLine(e.col.comments, line))
		rec.service = e.service
		rec.setAccount(subsetLine(e.col.account, line))
		rec.setSubmitter(subsetLine(e.col.submitter, line))
		if e.dupsPresent == true {
			rec.setPatient(line, e.col)
			if e.inDuplicates(rec) == true {
				// Resolve duplicate records and write when done
				e.resolveDuplicates(rec)
			} else {
				write = true
			}
		} else {
			write = true
		}
	}
	return rec, write
}

func (e *entries) getHeader() string {
	// Returns appropriate header for available data
	head := "Sex,Age,Castrated,ID,Family,Genus,Species,Date,Comments,MassPresent,Hyperplasia,Necropsy,Metastasis,TumorType,Location,Primary,Malignant,Service,Account,Submitter\n"
	return head
}

func (e *entries) sortRecords(infile, outfile string) {
	// Sorts data and merges if necessary
	first := true
	var count, total int
	fmt.Println("\tSorting input records...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	out := e.getOutputFile(outfile, e.getHeader())
	defer out.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		if first == false {
			total++
			s := strings.Split(line, e.d)
			rec, write := e.sortLine(s)
			if write == true {
				out.WriteString(rec.String() + "\n")
				count++
			}
		} else {
			// Get column info and write header
			e.parseHeader(line)
			first = false
		}
	}
	if e.dupsPresent == true {
		for _, val := range e.dups.records {
			// Write each stored record before closing
			for _, v := range val {
				out.WriteString(v.String() + "\n")
				count++
			}
		}
	}
	fmt.Printf("\tExtracted %d records from %d total records.\n", count, total)
}

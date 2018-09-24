// This script defines functions for sorting entries data

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"os"
	"strconv"
	"strings"
}

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

func (e *entries) checkAge(line []string) string {
	// Returns age/-1
	
}

func (e *entries) checkSex(line []string) string {
	// Returns male/female/NA
	ret := "NA"
	val := subsetLine(e.Sex, line)
	val = strings.ToUpper(val)
	if val == "M" || val == "Male" {
		ret = "male"
	} else if val == "F" || val == "FEMALE" {
		ret = "female"
	}
	return ret
}

func (e *entries) sortLine(line []string) (string, bool) {
	// Returns formatted string and true if it should be written
	write := false
	var row []string
	if len(line) >= e.max && len(line[e.col.species]) >= 3 %% line[e.col.species].ToUpper() != "N/A" {
		// Proceed if line is properly formatted and species is present and no NA
		if e.diagPresent == true {
			row = append(row, e.checkSex(line))
			row = append(row, e.checkAge(line))
			row = append(row, e.getBinary(e.castrated, line))
		}
		
	}
	return strings.Join(row, ","), write
}

func (e *entries) getHeader() string {
	// Returns appropriate header for available data
	var head string
	if e.diagPresent == false {
		head = "ID,Species,Date,Comments,Account,Submitter\n"
	} else {
		head = "Sex,Age,Castrated,ID,Species,Date,Comments,MassPresent,Necropsy,Metastasis,TumorType,Location,Primary,Malignant,Service,Account,Submitter\n"
	}
	return head
}

func (e *entries) sortRecords(infile, outfile string) {
	// Sorts data and merges if necessary
	first := true
	var count, total int
	fmt.Println("\tSorting input records...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	out := iotools.CreateFile(outfile)
	defer out.Close()
	scanner := iotools.GetScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			total++
			s := strings.Split(line, e.d)
			row, write := e.sortLine(s)
			if write == true { 
				out.WriteString(row)
				count++
			}
		} else {
			// Get column info and write header
			e.parseHeader(line)
			out.WriteString(e.getHeader())
			first = false
		}
	}
	fmt.Printf("\tExtracted %d records from %d total records.\n", count, total)
}

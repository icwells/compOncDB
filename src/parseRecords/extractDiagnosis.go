// This script will extract diagnosis information from a given input file

package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func countNA(row []string) (bool, bool) {
	// Determines if any or all fields have been identified
	var found, complete bool
	count := 0
	l := len(row)
	for _, i := range row {
		if i == "NA" {
			count++
		}
	}
	if count < l {
		found = true
		if count == 0 {
			complete = true
		}
	}
	return found, complete
}

func (e *entries) getDiagnosis(line, age string, cancer bool) []string {
	// Examines line for each diagnosis case
	var row []string
	prim := "N"
	if e.match.infantRecords(line) == true {
		age = "0"
	} else if age == "NA" {
		// Try to extract age if it's not given
		age = e.match.getAge(line)
	}
	row = append(row, age)
	row = append(row, e.match.getMatch(e.match.sex, line))
	row = append(row, e.match.getCastrated(line))
	row = append(row, e.match.getLocation(line, cancer))
	t := e.match.getType(line, cancer)
	row = append(row, t)
	row = append(row, e.match.binaryMatch(e.malignant, line, "benign"))
	met := e.match.binaryMatch(e.metastasis, line, "")
	if met == "N" && t != "NA" {
		// Store yes for primary if a tumor was found but no metastasis
		prim = "Y"
	} else {
		if e.match.getMatch(e.match.primary, line) != "NA" {
			prim = "Y"
		}
	}
	row = append(row, prim)
	row = append(row, met)
	row = append(row, e.match.binaryMatch(e.match.necropsy, line, "biopsy")
	return row
}

func (e *entries) checkAge(line []string, idx int) string {
	// Returns age from column if given
	ret := "NA"
	if e.col.days >= 0 {
		age, err := strconv.ParseFloat(line[e.col.days])
		if err == nil {
			// convert days to months
			age = age / 30.0
			ret = strconv.FormatFloat(age, 'f', -1, 64)
		}
	} else if e.col.age >= 0 {
		age, err := strconv.ParseFloat(line[e.col.age])
		if err == nil {
			// convert years to months
			age = age * 12.0
			ret = strconv.FormatFloat(age, 'f', -1, 64)
		}
	}
	return ret	
}

func (e *entries) parseLine(line []string) ([]string, bool, bool) {
	// Extracts diagnosis info from line
	var row []string
	cancer := true
	if len(line) > e.col.id {
		id := line[e.col.id]
		age := e.checkAge(line, e.col.age)
		// Remove ID and join line
		line = append(line[:e.col.id], line[e.col.id+1:]...)
		str := strings.Join(line, " ")
		if e.service == "NWZP" && strings.Contains(line[e.col.code], "8") == false {
			cancer = false
		}
		row = getDiagnosis(str, age, cancer)
	}
	found, complete := countNA(row)
	return row, found, complete
}

func (e *entries) extractDiagnosis(dict, infile, outfile string) {
	// Get diagnosis information using regexp struct
	var res [][]string
	var count, total, complete int
	first := true
	e.match = newMatcher(dict)
	head := "ID,Age(months),Sex,Castrated,Location,Type,Malignant,PrimaryTumor,Metastasis,Necropsy\n"
	fmt.Println("\tExtracting diagnosis data...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Input() {
		line := string(input.Text())
		if first == false {
			total++
			s := strings.Split(line, e.d)
			row, found, com := e.parseLine(s)
			if found == true {
				res = append(res, row)
				count++
				if com == true {
					complete++
				}
			}
		} else {
			e.parseHeader(line)
			first = false
		}
	}
	fmt.Printf("\tFound data for %d of %d records.\n", count, total)
	fmt.Printf("\tFound complete information for %d records.\n", complete)
	iotools.WriteToCSV(outfile, header, res)
}

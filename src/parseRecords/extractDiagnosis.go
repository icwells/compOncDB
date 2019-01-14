// This script will extract diagnosis information from a given input file

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"strconv"
	"strings"
)

func countNA(row []string) (bool, bool) {
	// Determines if any or all fields have been identified
	var found, complete bool
	count := 0
	l := len(row) - 2
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

func (e *entries) parseDiagnosis(line, age string, cancer, necropsy bool) []string {
	// Examines line for each diagnosis case
	var row []string
	prim := "N"
	line = strings.ToLower(line)
	if e.match.infantRecords(line) == true {
		age = "0"
	} else if age == "NA" {
		// Try to extract age if it's not given
		age = e.match.getAge(line)
	}
	if ch, _ := strconv.ParseFloat(age, 64); ch < 0.0 {
		// Make sure values aren't below 0
		age = "0"
	}
	row = append(row, age)
	row = append(row, e.match.getMatch(e.match.sex, line))
	row = append(row, e.match.getCastrated(line))
	row = append(row, e.match.getLocation(line, cancer))
	t, mal := e.match.getType(line, cancer)
	row = append(row, t)
	met := e.match.binaryMatch(e.match.metastasis, line)
	if met == "Y" {
		// Assume malignancy if metastasis is detected
		mal = "Y"
	}
	row = append(row, mal)
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
	if necropsy == true {
		row = append(row, "Y")
	} else {
		row = append(row, e.match.getNecropsy(line))
	}
	return row
}

func (e *entries) checkAge(line []string) string {
	// Returns age from column if given
	ret := "NA"
	if e.col.days >= 0 {
		age, err := strconv.ParseFloat(line[e.col.days], 64)
		if err == nil {
			// Convert days to months
			age = age / 30.0
			ret = strconv.FormatFloat(age, 'f', -1, 64)
		}
	} else if e.col.age >= 0 {
		age, err := strconv.ParseFloat(line[e.col.age], 64)
		if err == nil {
			// Convert years to months
			age = age * 12.0
			ret = strconv.FormatFloat(age, 'f', -1, 64)
		}
	}
	return ret
}

func (e *entries) parseLine(line []string) ([]string, bool, bool) {
	// Extracts diagnosis info from line
	var row []string
	var necropsy, found, complete bool
	cancer := true
	idx := e.col.id
	if e.service == "NWZP" && e.col.code > idx {
		// Get larger index
		idx = e.col.code
	}
	if len(line) > idx {
		id := line[e.col.id]
		age := e.checkAge(line)
		if e.service == "NWZP" {
			// Get neoplasia and euthnasia codes from NWZP
			cancer = strings.Contains(line[e.col.code], "8")
			necropsy = strings.Contains(line[e.col.code], "14")
		}
		// Remove ID and join line
		line = append(line[:e.col.id], line[e.col.id+1:]...)
		str := strings.Join(line, " ")
		row = e.parseDiagnosis(str, age, cancer, necropsy)
		// Prepend id
		row = append([]string{id}, row...)
		found, complete = countNA(row)
	}
	return row, found, complete
}

func (e *entries) extractDiagnosis(infile, outfile string) {
	// Get diagnosis information using regexp struct
	var res [][]string
	var count, total, complete int
	first := true
	head := "ID,Age(months),Sex,Castrated,Location,Type,Malignant,PrimaryTumor,Metastasis,Necropsy"
	fmt.Println("\n\tExtracting diagnosis data...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
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
	iotools.WriteToCSV(outfile, head, res)
}

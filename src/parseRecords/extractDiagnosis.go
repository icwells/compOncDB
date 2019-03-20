// This script will extract diagnosis information from a given input file

package main

import (
	"strconv"
	"strings"
)

func countNA(r *record) (bool, bool) {
	// Determines if any or all fields have been identified
	var found, complete bool
	count := 0
	targets := []string{r.age, r.sex, r.castrated, r.location, r.tumorType, r.malignant, r.metastasis, r.necropsy}
	for _, i := range targets {
		if i == "NA" || i == "-1" {
			count++
		}
	}
	if count < len(targets) {
		found = true
		if count == 0 {
			complete = true
		}
	}
	return found, complete
}

func (e *entries) parseDiagnosis(rec *record, line string, cancer, necropsy bool) {
	// Examines line for each diagnosis case
	line = strings.ToLower(line)
	if e.match.infantRecords(line) == true {
		rec.age = "0.0"
	} else if rec.age == "-1" {
		// Try to extract age if it's not given
		rec.age = e.match.getAge(line)
	}
	if ch, _ := strconv.ParseFloat(rec.age, 64); ch < -1.0 {
		// Make sure values aren't below 0
		rec.age = "-1"
	}
	rec.sex = e.match.getMatch(e.match.sex, line)
	rec.castrated = e.match.getCastrated(line)
	rec.tumorType, rec.location, rec.malignant = e.match.getTumor(line, cancer)
	rec.metastasis = e.match.binaryMatch(e.match.metastasis, line)
	if rec.metastasis == "1" {
		// Assume malignancy if metastasis is detected
		rec.malignant = "1"
	}
	if rec.tumorType != "NA" {
		// Only check for primary tumor if a tumor was found
		if rec.metastasis == "0" && strings.Count(rec.tumorType, ";") == 0 && strings.Count(rec.location, ";") == 0 {
			// Store yes for primary if a tumor was found but no metastasis
			rec.primary = "1"
		} else if e.match.getMatch(e.match.primary, line) != "NA" {
			rec.primary = "1"
		}
	}
	if necropsy == true {
		rec.necropsy = "1"
	} else {
		rec.necropsy = e.match.getNecropsy(line)
	}
}

func (e *entries) checkAge(line []string) string {
	// Returns age from column if given
	ret := "-1"
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

func (e *entries) parseLine(rec *record, line []string) {
	// Extracts diagnosis info from line
	var necropsy bool
	cancer := true
	rec.age = e.checkAge(line)
	if e.service == "NWZP" {
		// Get neoplasia and euthnasia codes from NWZP
		cancer = strings.Contains(line[e.col.code], "8")
		necropsy = strings.Contains(line[e.col.code], "14")
	}
	// Remove ID and join line (make copy to preserve column indeces)
	row := make([]string, len(line))
	copy(row, line)
	row = append(row[:e.col.id], row[e.col.id+1:]...)
	str := strings.Join(line, " ")
	e.parseDiagnosis(rec, str, cancer, necropsy)
	found, com := countNA(rec)
	if found == true {
		e.found++
		if com == true {
			e.complete++
		}
	}
}

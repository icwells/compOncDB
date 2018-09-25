// This script defines a struct for handling single pathology records

package main

import (
	"strconv"
	"strings"
)

func checkString(val string) {
	// Returns NA if string is malformed
	if len(val) <= 0 {
		val = "NA"
	} else if val == "na" || strings.ToUpper(val) == "N/A" {
		val = "NA"
	}
	return val
}

func checkBinary(val string) string {
	// Returns binary options as -1/0/1
	ret := "-1"
	val = strings.ToUpper(val)
	if val == "Y" || val == "YES" {
		ret = "1"
	} else if val == "N" || val == "NO" {
		ret = "0"
	}
	return ret
}

type record struct {
	sex         string
	age         string
	castrated   string
	id          string
	species     string
	date        string
	comments    string
	massPresent string
	necropsy    string
	metastasis  string
	tumorType   string
	location    string
	primary     string
	malignant   string
	service     string
	account     string
	submitter   string
	patient		string
}

func (r *record) String() {
	// Returns formatted string
	var row []string
	row = append(row, r.sex)
	row = append(row, r.age)
	row = append(row, r.castrated)
	row = append(row, r.id)
	row = append(row, r.species)
	row = append(row, r.date)
	row = append(row, r.comments)
	row = append(row, r.massPresent)
	row = append(row, r.necropsy)
	row = append(row, r.metastasis)
	row = append(row, r.tumorType)
	row = append(row, r.location)
	row = append(row, r.primary)
	row = append(row, r.malignant)
	row = append(row, r.service)
	row = append(row, r.account)
	row = append(row, r.submitter)
	for idx, i := range row {
		// Make sure values are present
		if len(i) <= 0 {
			row = append(row[:idx], row[idx+1:]...)
		}
	}
	return strings.Join(row, ",")
}

func (r *record) setPatient(line []string, c columns) {
	// Attempts to identify patient id
	if c.patient >= 0 {
		r.patient = checkString(line[c.patient])
	} else if c.id >= 0 {
		r.patient = checkString(line[c.id])
	}
}

func (r *record) setAccount(val string) {
	//Store account/NA
	r.account = checkString(val)
}

func (r *record) setSubmitter(val string) {
	//Store submitter/NA
	r.submitter = checkString(val)
}

func (r *record) setDate(val string) {
	//Store date/NA
	r.date = checkString(val)
}

func (r *record) setComments(val string) {
	//Store comments/NA
	r.comments = checkString(val)
}

func (r *record) setLocation(val string) {
	// Store location/NA
	r.location = checkString(val)
}

func (r *record) setType(val string) {
	// Store type/NA and masspresent
	r.tumorType = checkString(val)
	if r.tumorType != "NA" {
		r.massPresent = "1"
	} else {
		r.massPresent = "0"
	}
}

func (r *record) setID(val string) {
	// Makes sure ID is an int
	if _, err := strconv.Atoi(val); err == nil {
		r.id = val
	} else {
		r.id = "-1"
	}
}

func (r *record) setAge(val string) {
	// Returns age/-1
	if strings.ToUpper(val) == "NA" || len(val) > 7 {
		// Set -1 if age is too long (age would be impossible)
		r.age = "-1"
	} else if _, err := strconv.parseFloat(val, 64); err == nil {
		r.age = val
	} else {
		r.age = "-1"
	}
}

func (r *record) setSex(val string) {
	// Returns male/female/NA
	val = strings.ToUpper(val)
	if val == "M" || val == "Male" {
		r.sex = "male"
	} else if val == "F" || val == "FEMALE" {
		r.sex = "female"
	} else {
		r.sex = "NA"
	}
}

func (r *record) setDiagnosis(row []string) {
	// Stores and formats input from diagnosis
	r.setAge(row[0])
	r.setSex(row[1])
	r.setCastrated(checkBinary(row[2]))
	r.setLocation(row[3])
	r.setType(row[4])
	r.setMalignant(checkBinary(row[5]))
	r.setPrimary(checkBinary(row[6]))
	r.setMetastasis(checkBinary(row[7]))
	r.setNecropsy(checkBinary(row[8]))
}

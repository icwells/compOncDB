// This script defines a struct for handling single pathology records

package main

import (
	"bytes"
	"strings"
)

func checkString(val string) string {
	// Returns NA if string is malformed
	v := strings.ToLower(val)
	if len(val) <= 0 {
		val = "NA"
	} else if v == "na" || v == "n/a" {
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
	genus       string
	species     string
	name        string
	date        string
	comments    string
	massPresent string
	hyperplasia string
	necropsy    string
	metastasis  string
	tumorType   string
	location    string
	primary     string
	malignant   string
	service     string
	account     string
	submitter   string
	patient     string
	cancer      string
	code        string
}

func newRecord() record {
	// Returns record with default values
	var r record
	r.sex = "NA"
	r.age = "-1"
	r.castrated = "-1"
	r.id = "NA"
	r.genus = "NA"
	r.species = "NA"
	r.date = "NA"
	r.comments = "NA"
	r.massPresent = "0"
	r.hyperplasia = "0"
	r.necropsy = "-1"
	r.metastasis = "-1"
	r.tumorType = "NA"
	r.location = "NA"
	r.primary = "0"
	r.malignant = "-1"
	r.service = "NA"
	r.account = "NA"
	r.submitter = "NA"
	r.patient = "NA"
	r.cancer = "N"
	r.code = "NA"
	return r
}

func (r *record) String(debug bool) string {
	// Returns formatted string
	buffer := bytes.NewBufferString(r.sex)
	buffer.WriteByte(',')
	buffer.WriteString(r.age)
	buffer.WriteByte(',')
	buffer.WriteString(r.castrated)
	buffer.WriteByte(',')
	buffer.WriteString(r.id)
	buffer.WriteByte(',')
	buffer.WriteString(r.genus)
	buffer.WriteByte(',')
	buffer.WriteString(r.species)
	buffer.WriteByte(',')
	buffer.WriteString(r.name)
	buffer.WriteByte(',')
	buffer.WriteString(r.date)
	buffer.WriteByte(',')
	buffer.WriteString(r.comments)
	buffer.WriteByte(',')
	buffer.WriteString(r.massPresent)
	buffer.WriteByte(',')
	buffer.WriteString(r.hyperplasia)
	buffer.WriteByte(',')
	buffer.WriteString(r.necropsy)
	buffer.WriteByte(',')
	buffer.WriteString(r.metastasis)
	buffer.WriteByte(',')
	buffer.WriteString(r.tumorType)
	buffer.WriteByte(',')
	buffer.WriteString(r.location)
	buffer.WriteByte(',')
	buffer.WriteString(r.primary)
	buffer.WriteByte(',')
	buffer.WriteString(r.malignant)
	buffer.WriteByte(',')
	buffer.WriteString(r.service)
	buffer.WriteByte(',')
	buffer.WriteString(r.account)
	buffer.WriteByte(',')
	buffer.WriteString(r.submitter)
	if debug == true {
		buffer.WriteByte(',')
		buffer.WriteString(r.cancer)
		buffer.WriteByte(',')
		buffer.WriteString(r.code)
	}
	return buffer.String()
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
	r.submitter = val
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
	// Store type/NA and hyperplasia
	r.tumorType = checkString(val)
	if r.tumorType == "hyperplasia" || r.tumorType == "neoplasia" {
		r.hyperplasia = "1"
	}
}

func (r *record) setSpecies(t []string) {
	// Stores family, genus, and species
	r.genus = t[0]
	r.species = t[1]
}

func (r *record) setID(val string) {
	// Stores ID as string
	r.id = checkString(val)
}

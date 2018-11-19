// This script defines a struct for handling single pathology records

package main

import (
	"bytes"
	"strconv"
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
	patient     string
}

func newRecord() record {
	// Returns record with default values
	var r record
	r.sex = "NA"
	r.age = "-1"
	r.castrated = "-1"
	r.id = "NA"
	r.species = "NA"
	r.date = "NA"
	r.comments = "NA"
	r.massPresent = "0"
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
	return r
}

func (r *record) String() string {
	// Returns formatted string
	buffer := bytes.NewBufferString(r.sex)
	buffer.WriteByte(',')
	buffer.WriteString(r.age)
	buffer.WriteByte(',')
	buffer.WriteString(r.castrated)
	buffer.WriteByte(',')
	buffer.WriteString(r.id)
	buffer.WriteByte(',')
	buffer.WriteString(r.species)
	buffer.WriteByte(',')
	buffer.WriteString(r.date)
	buffer.WriteByte(',')
	buffer.WriteString(r.comments)
	buffer.WriteByte(',')
	buffer.WriteString(r.massPresent)
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
	terms := map[string]string{"Animal Clinic": "A. C.", "Animal Hospital": "A. H.", "Veterinary Clinic": "V. C.",
"Veterinary Hospital": "V. H.", "Veterinary Services": "V. S.", "Pet Vet": "P. V.", "International": "Intl ", "Animal": "Anim "}
	val = checkString(val)
	if val != "NA" {
		// Resolve abbreviations
		for k, v := range terms {
			var alt string
			if strings.Contains(v, ".") == false {
				// Add trailing period
				alt = strings.Replace(v, " ", ".", 1)
			} else {
				// Remove space
				alt = strings.Replace(v, " ", "", 1)
			}
			if strings.Contains(val, v) == true {
				val = strings.Replace(val, v, k, 1)
				break
			} else if strings.Contains(val, alt) == true {
				val = strings.Replace(val, alt, k, 1)
				break
			}
		}
	}
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
	// Store type/NA and masspresent
	r.tumorType = checkString(val)
}

func (r *record) setID(val string) {
	// Stores ID as string
	r.id = checkString(val)
}

func (r *record) setAge(val string) {
	// Returns age/-1
	if strings.ToUpper(val) == "NA" || len(val) > 7 {
		// Set -1 if age is too long (age would be impossible)
		r.age = "-1"
	} else if _, err := strconv.ParseFloat(val, 64); err == nil {
		r.age = val
	} else {
		r.age = "-1"
	}
}

func (r *record) setSex(val string) {
	// Returns male/female/NA
	val = strings.ToUpper(val)
	if val == "M" || val == "MALE" {
		r.sex = "male"
	} else if val == "F" || val == "FEMALE" {
		r.sex = "female"
	} else {
		r.sex = "NA"
	}
}

func (r *record) setMassPresent() {
	// Determines if mass is present by evaluating diagnosis settings
	if r.tumorType != "NA" || r.malignant == "1" || r.metastasis == "1" || r.primary == "1" {
		r.massPresent = "1"
	}
}

func (r *record) setDiagnosis(row []string) {
	// Stores and formats input from diagnosis
	r.setAge(row[0])
	r.setSex(row[1])
	r.castrated = checkBinary(row[2])
	r.setLocation(row[3])
	r.setType(row[4])
	r.malignant = checkBinary(row[5])
	r.primary = checkBinary(row[6])
	r.metastasis = checkBinary(row[7])
	r.necropsy = checkBinary(row[8])
	r.setMassPresent()
}

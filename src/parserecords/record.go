// This script defines a struct for handling single pathology records

package parserecords

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var DIGIT = regexp.MustCompile(`([0-9]*[.])?[0-9]+`)

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
	year        string
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
	zoo         string
	aza         string
	institute   string
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
	r.year = "-1"
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
	r.zoo = "-1"
	r.aza = "-1"
	r.institute = "-1"
	r.patient = "NA"
	r.cancer = "N"
	r.code = "NA"
	return r
}

func (r *record) String(debug bool) string {
	// Returns formatted string
	columns := []string{r.sex, r.age, r.castrated, r.id, r.genus, r.species, r.name, r.date, r.year, r.comments}
	columns = append(columns, []string{r.massPresent, r.hyperplasia, r.necropsy, r.metastasis, r.tumorType, r.location, r.primary, r.malignant}...)
	columns = append(columns, []string{r.service, r.account, r.submitter, r.zoo, r.aza, r.institute}...)
	if debug {
		columns = append(columns, []string{r.cancer, r.code}...)
	}
	return strings.Join(columns, ",")
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

func (r *record) setSubmitter(v []string) {
	//Store submitter/NA
	if len(v) == 4 {
		r.submitter = v[0]
		r.zoo = v[1]
		r.aza = v[2]
		r.institute = v[3]
	}
}

func (r *record) setDate(val string) {
	//Store date/NA
	r.date = checkString(val)
}

func (r *record) setYear(val string) {
	//Stores year in 4 digit format
	if val = checkString(val); val != "NA" {
		year := DIGIT.FindString(val)
		fmt.Println(year)
		if len(year) == 2 {
			if y, _ := strconv.Atoi(year); y > 50 {
				year = "19" + year
			} else {
				year = "20" + year
			}

		}
		if len(year) == 4 {
			r.year = year
		}
	}
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

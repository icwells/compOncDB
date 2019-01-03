// This script will perform white box tests on parseRecords diagnosis functions

package main

import (
	"strconv"
	"testing"
)

func TestCountNA(t *testing.T) {
	// Tests coutn NA method
	nas := []struct {
		row      []string
		found    bool
		complete bool
	}{
		{[]string{"1", "12", "male", "Y", "Liver", "neoplasm", "N", "N", "N", "Y"}, true, true},
		{[]string{"2", "12", "female", "NA", "NA", "NA", "NA", "NA", "NA", "NA"}, true, false},
		{[]string{"3", "12", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA"}, false, false},
		{[]string{"NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA"}, false, false},
	}
	for _, i := range nas {
		found, complete := countNA(i.row)
		if found != i.found || complete != i.complete {
			t.Errorf("countNA returned %v, %v instead of %v, %v.", found, complete, i.found, i.complete)
		}
	}
}

func getAge(n int) string {
	// Returns formatted age
	return strconv.FormatFloat(float64(n)*12.0, 'f', -1, 64)
}

func getDays(n int) string {
	// Return formatted age from days
	return strconv.FormatFloat(float64(n)/30.0, 'f', -1, 64)
}

func TestCheckAge(t *testing.T) {
	// Tests checkAge for age in days and years
	e := newEntries("service")
	ages := []struct {
		row  []string
		idx  int
		age  string
		days string
	}{
		{[]string{"9", "ABC", "A16", "A99", "7", "1-Dec", "PACIFIC FORD TURTLE", "99-121", "M/F", "", "A99-7"}, 0, getAge(9), getDays(9)},
		{[]string{"f1212351", "Bongo", "skin biopsy:  squamous cell carcinoma, in situ", "", "Female", "2"}, 5, getAge(2), getDays(2)},
		{[]string{"6254519", "KV Zoo", "39179", "Arctic Fox", "6211", "", "", "Female"}, 5, "NA", "NA"},
	}
	for _, i := range ages {
		e.col.age = i.idx
		actual := e.checkAge(i.row)
		if actual != i.age {
			t.Errorf("Returned incorect value %s from age column. Expected %s.", actual, i.row[i.idx])
		}
	}
	for _, i := range ages {
		// Repeat with days column
		e.col.days = i.idx
		actual := e.checkAge(i.row)
		if actual != i.days {
			t.Errorf("Returned incorect value %s from days column. Expected %s.", actual, i.row[i.idx])
		}
	}
}

/*func TestParseDiagnosis(t *testing.T) {
	// Tests parseDiagnosis with matches from matcher_test.go
	e := newEntries("service")
	matches := newMatches()
	for _, i := range matches {
		row := e.parseDiagnosis(strings.ToLower(i.line), "NA", true, false)
		if row[0] != i.age && i.infant == false {
			t.Errorf("Actual age %s does not equal expected %s", row[0], i.age)
		} else if row[1] != i.sex {
			t.Errorf("Actual sex %s does not equal expected %s", row[1], i.sex)
		} else if row[2] != i.castrated {
			t.Errorf("Actual neuter value %s does not equal expected %s", row[2], i.castrated)
		} else if row[3] != i.location {
			t.Errorf("Actual location %s does not equal expected %s", row[3], i.location)
		} else if row[4] != i.typ {
			t.Errorf("Actual type %s does not equal expected %s", row[4], i.typ)
		} else if row[5] != i.malignant {
			t.Errorf("Actual malignant value %s does not equal expected %s", row[5], i.malignant)
		} else if row[6] != i.primary {
			if row[4] != "NA" && row[7] == "N" {
				// Skip if function determined a single tumor to be primary
				if row[6] != "Y" {
					t.Errorf("Actual primary %s does not equal expected %s", row[6], i.primary)
				}
			} else {
				t.Errorf("Actual primary %s does not equal expected %s", row[6], i.primary)
			}
		} else if row[7] != i.metastasis {
			t.Errorf("Actual metastasis value %s does not equal expected %s", row[7], i.metastasis)
		} else if row[8] != i.necropsy {
			t.Errorf("%s: Actual necropsy value %s does not equal expected %s", i.line, row[8], i.necropsy)
		}
	}
}*/

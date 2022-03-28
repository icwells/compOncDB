// This script will perform white box tests on parseRecords diagnosis functions

package parserecords

import (
	"github.com/icwells/compOncDB/src/diagnoses"
	"strconv"
	"strings"
	"testing"
)

func toRecord(row []string) *record {
	// Sorts slice into record struct
	r := newRecord()
	r.age = row[1]
	r.sex = row[2]
	r.castrated = row[3]
	r.location = row[4]
	r.tumorType = row[5]
	r.malignant = row[6]
	r.primary = row[7]
	r.metastasis = row[8]
	r.necropsy = row[9]
	return &r
}

func TestCountNA(t *testing.T) {
	// Tests coutn NA method
	nas := []struct {
		row      []string
		found    bool
		complete bool
	}{
		{[]string{"1", "12", "male", "Y", "Liver", "neoplasm", "N", "N", "N", "Y"}, true, true},
		{[]string{"2", "12", "female", "NA", "NA", "NA", "NA", "NA", "NA", "NA"}, true, false},
		{[]string{"3", "12", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA"}, true, false},
		{[]string{"NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA"}, false, false},
	}
	for _, i := range nas {
		found, complete := countNA(toRecord(i.row))
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
	e := NewEntries("service", "")
	ages := []struct {
		row  []string
		idx  int
		age  string
		days string
	}{
		{[]string{"9", "ABC", "A16", "A99", "7", "1-Dec", "PACIFIC FORD TURTLE", "99-121", "M/F", "", "A99-7"}, 0, getAge(9), getDays(9)},
		{[]string{"f1212351", "Bongo", "skin biopsy:  squamous cell carcinoma, in situ", "", "Female", "2"}, 5, getAge(2), getDays(2)},
		{[]string{"6254519", "KV Zoo", "39179", "Arctic Fox", "6211", "", "", "Female"}, 5, "-1", "-1"},
	}
	for _, i := range ages {
		e.col.age = i.idx
		actual := e.checkAge(i.row)
		if actual != i.age {
			t.Errorf("Returned incorect value %s from age column. expected: %s.", actual, i.age)
		}
	}
	for _, i := range ages {
		// Repeat with days column
		e.col.days = i.idx
		actual := e.checkAge(i.row)
		if actual != i.days {
			t.Errorf("Returned incorect value %s from days column. expected: %s.", actual, i.days)
		}
	}
}

func TestParseDiagnosis(t *testing.T) {
	// Tests parseDiagnosis with matches from matcher_test.go
	e := NewEntries("service", "")
	matches := diagnoses.NewMatches()
	for _, i := range matches {
		rec := newRecord()
		e.parseDiagnosis(&rec, strings.ToLower(i.Line), true, false)
		if rec.age != i.Age && i.Infant == false {
			t.Errorf("Actual age %s does not equal expected: %s", rec.age, i.Age)
		} else if rec.sex != i.Sex {
			t.Errorf("Actual sex %s does not equal expected: %s", rec.sex, i.Sex)
		} else if rec.castrated != i.Castrated {
			t.Errorf("Actual neuter value %s does not equal expected: %s", rec.castrated, i.Castrated)
		} else if rec.location != i.Location {
			t.Errorf("Actual location %s does not equal expected: %s", rec.location, i.Location)
		} else if rec.tumorType != i.Typ {
			t.Errorf("Actual type %s does not equal expected: %s", rec.tumorType, i.Typ)
		} else if rec.malignant != i.Malignant {
			t.Errorf("Actual malignant value %s does not equal expected: %s", rec.malignant, i.Malignant)
		} else if rec.primary != i.Primary {
			/*if rec.tumorType != "NA" && rec.metastasis == "0" {
				// Skip if function determined a single tumor to be primary
				if rec.primary != "1" {
					t.Errorf("Actual primary %s does not equal expected: %s", rec.primary, i.Primary)
				}
			} else {*/
				t.Errorf("Actual primary %s does not equal expected: %s", rec.primary, i.Primary)
			//}
		} else if rec.metastasis != i.Metastasis {
			t.Errorf("Actual metastasis value %s does not equal expected: %s", rec.metastasis, i.Metastasis)
		} else if rec.necropsy != i.Necropsy {
			t.Errorf("%s: Actual necropsy value %s does not equal expected: %s", i.Line, rec.necropsy, i.Necropsy)
		}
	}
}

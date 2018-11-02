// This script will perform white box tests on parseRecords diagnosis functions

package main

import (
	"fmt"
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
			msg := fmt.Sprintf("countNA returned %v, %v instead of %v, %v.", found, complete, i.found, i.complete)
			t.Error(msg)
		}
	}
}

func getAge(n int) string {
	// Returns formatted age
	return strconv.FormatFloat(float64(n * 12), 'f', -1, 64)
}

func getDays(n int) string {
	// Return formatted age from days
	return strconv.FormatFloat(float64(n / 30.0), 'f', -1, 64)
}

func TestCheckAge(t *testing.T) {
	// Tests checkAge for age in days and years
	e := newEntries("service")
	ages := []struct {
		row		[]string
		idx		int
		age		string
		days	string
	}{
		{[]string{"9", "ABC", "A16", "A99", "7", "1-Dec", "PACIFIC FORD TURTLE", "99-121", "M/F", "", "A99-7"}, 0, getAge(9), getDays(9)},
		{[]string{"f1212351", "Bongo", "", "", "skin biopsy:  squamous cell carcinoma, in situ", "", "Female", "2"}, 7, getAge(2), getDays(2)},
		{[]string{"6254519", "KV Zoo", "39179", "Arctic Fox", "6211", "", "", "Female"}, 5, "NA", "NA"},
	}
	for _, i := range ages {
		e.col.age = i.idx
		actual := e.checkAge(i.row)
		if actual != i.age {
			msg := fmt.Sprintf("Returned incorect value %s from age column. Expected %s.", actual, i.row[i.idx])
			t.Error(msg)
		}
		// Repeat with days column
		e.col.age++
		e.col.days = i.idx
		actual = e.checkAge(i.row)
		if actual != i.days {
			msg := fmt.Sprintf("Returned incorect value %s from days column. Expected %s.", actual, i.row[i.idx])
			t.Error(msg)
		}
	}
}

func TestParseDiagnosis(t *testing.T) {
	// Tests parseDiagnosis

}

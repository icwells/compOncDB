// PErforms white box tests on the records struct's methods

package main

import (
	"fmt"
	"testing"
)

func TestSubsetLine(t *testing.T) {
	// Tests subsetLine (in sortRecords)
	line := []string{"Cheetah ", "", " lion", " Heloderma suspectum"}
	matches := []struct{
		idx			int
		expected	string
	} {
		{-1, "NA"},
		{0, "Cheetah"},
		{1, "NA"},
		{2, "lion"},
		{3, "Heloderma suspectum"},
	}
	for _, i := range matches {
		actual := subsetLine(i.idx, line)
		if actual != i.expected {
			msg := fmt.Sprintf("Actual %s does not equal expected: %s.", actual, i.expected)
			t.Error(msg)
		}
	}
}

func TestCheckString(t *testing.T) {
	// Tests checkString for NA determining
	matches := []struct {
		value		string
		expected	string
	} {
		{"", "NA"},
		{"yes", "yes"},
		{"N", "N"},
		{"n/a", "NA"},
		{"Na", "NA"},
	}
	for _, i := range matches {
		actual := checkString(i.value)
		if actual != i.expected {
			msg := fmt.Sprintf("Actual string value %s does not equal expected %s", actual, i.expected)
			t.Error(msg)
		}
	}
}

func TestCheckBinary(t *testing.T) {
	// Tests binary option conversion 
	matches := []struct {
		value		string
		expected	string
	} {
		{"Y", "1"},
		{"yes", "1"},
		{"N", "0"},
		{"nO", "0"},
		{"Na", "-1"},
		{"test", "-1"},
	}
	for _, i := range matches {
		actual := checkBinary(i.value)
		if actual != i.expected {
			msg := fmt.Sprintf("Actual binary option %s does not equal expected %s", actual, i.expected)
			t.Error(msg)
		}
	}
}

func compareRecords(a, e record) string {
	// Returns error message/empty string
	var msg string
	if a.sex != e.sex {
		msg = fmt.Sprintf("Actual sex %s does not equal expected: %s.", a.sex, e.sex)
	} else if a.age != e.age {
		msg = fmt.Sprintf("Actual age %s does not equal expected: %s.", a.age, e.age)
	} else if a.castrated != e.castrated {
		msg = fmt.Sprintf("Actual neuter value %s does not equal expected: %s.", a.castrated, e.castrated)
	} else if a.id != e.id {
		msg = fmt.Sprintf("Actual ID %s does not equal expected: %s.", a.id, e.id)
	} else if a.species != e.species {
		msg = fmt.Sprintf("Actual species %s does not equal expected: %s.", a.species, e.species)
	} else if a.date != e.date {
		msg = fmt.Sprintf("Actual date %s does not equal expected: %s.", a.date, e.date)
	} else if a.comments != e.comments {
		msg = fmt.Sprintf("Actual comments %s do not equal expected: %s.", a.comments, e.comments)
	} else if a.massPresent != e.massPresent {
		msg = fmt.Sprintf("Actual mass present value %s does not equal expected: %s.", a.massPresent, e.massPresent)
	} else if a.necropsy != e.necropsy {
		msg = fmt.Sprintf("Actual necropsy value %s does not equal expected: %s.", a.necropsy, e.necropsy)
	} else if a.metastasis != e.metastasis {
		msg = fmt.Sprintf("Actual metastasis value %s does not equal expected: %s.", a.metastasis, e.metastasis)
	} else if a.tumorType != e.tumorType {
		msg = fmt.Sprintf("Actual tumor type %s does not equal expected: %s.", a.tumorType, e.tumorType)
	} else if a.location != e.location {
		msg = fmt.Sprintf("Actual location %s does not equal expected: %s.", a.location, e.location)
	} else if a.primary != e.primary {
		msg = fmt.Sprintf("Actual primary tumor value %s does not equal expected: %s.", a.primary, e.primary)
	} else if a.malignant != e.malignant {
		msg = fmt.Sprintf("Actual malignant value %s does not equal expected: %s.", a.malignant, e.malignant)
	} else if a.service != e.service {
		msg = fmt.Sprintf("Actual service %s does not equal expected: %s.", a.service, e.service)
	} else if a.account != e.account {
		msg = fmt.Sprintf("Actual account %s does not equal expected: %s.", a.account, e.account)
	} else if a.submitter != e.submitter {
		msg = fmt.Sprintf("Actual submitter %s does not equal expected: %s.", a.submitter, e.submitter)
	} else if a.patient != e.patient {
		msg = fmt.Sprintf("Actual patient %s does not equal expected: %s.", a.patient, e.patient)
	}
	return msg
}

func testRecords(rows [][]string) []record {
	// Returns struct of record test cases
	var ret []record
	for idx := range rows {
		r := newRecord()
		switch (idx) {
			case 0:
				r.age = "12"
				r.sex = "male"
				r.castrated = "1"
				r.location = "Spleen"
				r.tumorType = "Carcinoma"
				r.massPresent = "1"
				r.malignant = "1"
				r.primary = "-1"
				r.metastasis = "1"
				r.necropsy = "1"
			case 1:
				r.age = "-1"
				r.sex = "female"
				r.castrated = "-1"
				r.location = "NA"
				r.tumorType = "NA"
				r.massPresent = "0"
				r.malignant = "-1"
				r.primary = "-1"
				r.metastasis = "-1"
				r.necropsy = "-1"
			case 2:
				r.age = "16.5"
				r.sex = "NA"
				r.castrated = "0"
				r.location = "liver"
				r.tumorType = "sarcoma"
				r.massPresent = "1"
				r.malignant = "0"
				r.primary = "1"
				r.metastasis = "0"
				r.necropsy = "0"
		}
		ret = append(ret, r)
	}
	return ret
}

func TestSetDiagnosis(t *testing.T) {
	// Tests setDiagnosis (and setAge, setSex, setType, and setLocation)
	matches := [][]string{
		{"12", "male", "Y", "Spleen", "Carcinoma", "Y", "NA", "Y", "Y"},
		{"12345678", "f", "NA", "NA", "NA", "NA", "NA", "NA", "NA"},
		{"16.5", "NA", "N", "liver", "sarcoma", "N", "Y", "N", "N"},
	}
	expected := testRecords(matches)
	for idx, i := range matches {
		actual := newRecord()
		actual.setDiagnosis(i)
		msg := compareRecords(actual, expected[idx])
		if len(msg) > 1 {
			t.Error(msg)
		}
	}
}

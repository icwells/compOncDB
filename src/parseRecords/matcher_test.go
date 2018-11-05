// This script will perform white box tests on the mathcer struct's methods

package main

import (
	"fmt"
	"testing"
)

type matches struct {
	line 		string
	location	string
	typ			string
	infant		bool
	age			string
	sex			string
	castrated	string
	malignant	string
	metastasis	string
	primary		string
	necropsy	string
}

func newMatches() []matches {
	// Initializes test matches
	line1 := "spinal neoplasia, biopsy; castration helps to resolve the situation since it is somewhat hormonal dependent, Female, 2 years old"
	line2 := "cause of death: single Malignant liver carcinoma; retarded growth has also been reported. 37 month old male"
	line3 := "metastatis lymphoma, infant, 30 days, not castrated, "
	line4 := "spayed female gray fox, "
	return []matches {
		{line1, "spinal cord", "neoplasia", false, "24.0000", "female", "Y", "NA", "NA", "N", "N"},
		{line2, "liver", "carcinoma", false, "37", "male", "NA", "Y", "NA", "Y", "Y"},
		{line3, "lymph", "lymphoma", true, "1.0000", "NA", "N", "NA", "Y", "N", "NA"},
		{line4, "NA", "NA", false, "NA", "female", "Y", "NA", "NA", "N", "NA"},
	}
}

func TestGetType(t *testing.T) {
	// Tests getMatch method
	m := newMatcher()
	matches := newMatches()
	for _, i := range matches {
		actual := m.getType(i.line, true)
		if actual != i.typ {
			msg := fmt.Sprintf("Actual type %s does not equal expected: %s.", actual, i.typ)
			t.Error(msg)
		}
	}
}

func TestGetLocation(t *testing.T) {
	// Tests getMatch method
	m := newMatcher()
	matches := newMatches()
	for _, i := range matches {
		actual := m.getLocation(i.line, true)
		if actual != i.location {
			msg := fmt.Sprintf("Actual location %s does not equal expected: %s.", actual, i.location)
			t.Error(msg)
		}
	}
}

func TestGetCastrated(t *testing.T) {
	// Tests getMatch method
	m := newMatcher()
	matches := newMatches()
	for _, i := range matches {
		actual := m.getCastrated(i.line)
		if actual != i.castrated {
			msg := fmt.Sprintf("Actual neuter value %s does not equal expected: %s.", actual, i.castrated)
			t.Error(msg)
		}
	}
}

func TestInfantRecords(t *testing.T) {
	// Tests getMatch method
	m := newMatcher()
	matches := newMatches()
	for _, i := range matches {
		actual := m.infantRecords(i.line)
		if actual != i.infant {
			msg := fmt.Sprintf("Actual infant record value %v does not equal expected: %v.", actual, i.infant)
			t.Error(msg)
		}
	}
}

func TestGetAge(t *testing.T) {
	// Tests getMatch method
	m := newMatcher()
	matches := newMatches()
	for _, i := range matches {
		actual := m.getAge(i.line)
		if actual != i.age {
			msg := fmt.Sprintf("Actual age %s does not equal expected: %s.", actual, i.age)
			t.Error(msg)
		}
	}
}

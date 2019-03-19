// This script will perform white box tests on the mathcer struct's methods

package main

import (
	"testing"
)

type matches struct {
	line       string
	location   string
	typ        string
	infant     bool
	age        string
	sex        string
	castrated  string
	malignant  string
	metastasis string
	primary    string
	necropsy   string
}

func newMatches() []matches {
	// Initializes test matches
	line1 := "spinal neoplasia, biopsy; castration helps to resolve the situation since it is somewhat hormonal dependent, Female, 2 years old"
	line2 := "cause of death: single Malignant liver carcinoma; retarded growth has also been reported. 37 month old male"
	line3 := "metastatis lymphoma, infant, 30 days, not castrated, "
	line4 := "spayed female gray fox, "
	return []matches{
		{line1, "spinal cord", "neoplasia", false, "24", "female", "1", "-1", "-1", "0", "0"},
		{line2, "liver", "carcinoma", false, "37", "male", "-1", "1", "-1", "1", "1"},
		{line3, "lymph nodes", "lymphoma", true, "1", "NA", "0", "1", "1", "0", "-1"},
		{line4, "NA", "NA", false, "-1", "female", "1", "-1", "-1", "0", "-1"},
	}
}

func TestGetTypes(t *testing.T) {
	// Tests getMatch method
	m := newMatcher()
	matches := newMatches()
	for _, i := range matches {
		f := newTumorFinder()
		m.getTypes(&f, i.line)
		typ, _, mal := f.toStrings()
		if typ != i.typ {
			t.Errorf("Actual type %s does not equal expected: %s.", typ, i.typ)
		} else if mal != i.malignant {
			t.Errorf("Actual malignant value %s does not equal expected: %s.", mal, i.malignant)
		}
	}
}

func TestGetLocations(t *testing.T) {
	// Tests getMatch method
	m := newMatcher()
	matches := newMatches()
	for _, i := range matches {
		f := newTumorFinder()
		m.getTypes(&f, i.line)
		m.getLocations(&f, i.line)
		_, loc, _ := f.toStrings()
		if loc != i.location {
			t.Errorf("Actual location %s does not equal expected: %s.", loc, i.location)
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
			t.Errorf("Actual neuter value %s does not equal expected: %s.", actual, i.castrated)
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
			t.Errorf("Actual infant record value %v does not equal expected: %v.", actual, i.infant)
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
			t.Errorf("Actual age %s does not equal expected: %s.", actual, i.age)
		}
	}
}

// This script will perform white box tests on the mathcer struct's methods

package diagnoses

import (
	//"strings"
	"testing"
)

type distcases struct {
	line, typ, key, match string
	dist                  int
}

func newDistCases() []distcases {
	// Returns slice of test cases for setDistance
	line1 := "spinal neoplasia, biopsy; castration helps to resolve the situation since it is somewhat hormonal dependent, Female, 2 years old"
	line2 := "cause of death: single Malignant liver carcinoma; retarded growth has also been reported. 37 month old male"
	line3 := "metastatis lymphoma, infant, 30 days, not castrated, "
	line4 := "spayed female gray fox cutaneous malignant melanoma, "
	return []distcases{
		{line1, "neoplasia", "spinal cord", "spinal", 1},
		{line2, "carcinoma", "liver", "liver", 1},
		{line3, "lymphoma", "lymph nodes", "lymphoma", 0},
		{line4, "melanoma", "skin", "cutaneous", 11},
	}
}

/*func TestSetDistance(t *testing.T) {
	cases := newDistCases()
	for _, c := range cases {
		a := newTumorType(c.typ, strings.Index(c.line, c.typ), len(c.line))
		a.setDistance(c.key, c.match, c.line)
		if a.locations[c.key] != c.dist {
			t.Errorf("Actual distance %d does not equal expected: %d", a.locations[c.key], c.dist)
		}
		a.setLocation()
		if a.location != c.key {
			t.Errorf("Actual set location %s does not equal expected: %s", a.location, c.key)
		}
	}
}*/

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

func TestTumor(t *testing.T) {
	// Tests getMatch method
	m := NewMatcher()
	matches := newMatches()
	for _, i := range matches {
		typ, loc, mal := m.GetTumor(i.line, "NA", true)
		if typ != i.typ {
			t.Errorf("Actual type %s does not equal expected: %s.", typ, i.typ)
		} else if loc != i.location {
			t.Errorf("Actual location %s does not equal expected: %s.", loc, i.location)
		} else if mal != i.malignant {
			t.Errorf("Actual malignant value %s does not equal expected: %s.", mal, i.malignant)
		}
	}
}

func TestGetCastrated(t *testing.T) {
	// Tests getMatch method
	m := NewMatcher()
	matches := newMatches()
	for _, i := range matches {
		actual := m.GetCastrated(i.line)
		if actual != i.castrated {
			t.Errorf("Actual neuter value %s does not equal expected: %s.", actual, i.castrated)
		}
	}
}

func TestInfantRecords(t *testing.T) {
	// Tests getMatch method
	m := NewMatcher()
	matches := newMatches()
	for _, i := range matches {
		actual := m.InfantRecords(i.line)
		if actual != i.infant {
			t.Errorf("Actual infant record value %v does not equal expected: %v.", actual, i.infant)
		}
	}
}

func TestGetAge(t *testing.T) {
	// Tests getMatch method
	m := NewMatcher()
	matches := newMatches()
	for _, i := range matches {
		actual := m.GetAge(i.line)
		if actual != i.age {
			t.Errorf("Actual age %s does not equal expected: %s.", actual, i.age)
		}
	}
}

// This script will perform white box tests on the mathcer struct's methods

package diagnoses

import (
	"strings"
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

func TestSetDistance(t *testing.T) {
	cases := newDistCases()
	for _, c := range cases {
		a := newTumorHit(c.typ, strings.Index(c.line, c.typ), len(c.line))
		a.setDistance(c.key, c.match, c.line)
		if a.locations[c.key] != c.dist {
			t.Errorf("Actual distance %d does not equal expected: %d", a.locations[c.key], c.dist)
		}
		a.setLocation()
		if a.location != c.key {
			t.Errorf("Actual set location %s does not equal expected: %s", a.location, c.key)
		}
	}
}

func getTissueLocations() []string {
	// Returns slice of known tissues
	return []string{"fibrous", "myxomatous tissue", "fat", "notochord", "smooth muscle", "striated muscle", "peripheral nerve sheath", "meninges", "blood", "cartilage", "synovium", "bone", "bone marrow", "lymph nodes", "spleen", "mast cell", "dendritic cell", "pigment cell", "skin", "hair follicle", "gland", "mammary", "glial cell", "nerve cell", "pnet", "neuroepithelial", "spinal cord", "brain", "pituitary gland", "parathyroid gland", "thyroid", "adrenal medulla", "adrenal cortex", "pancreas", "stomach", "carotid body", "neuroendocrine", "testis", "prostate", "ovary", "vulva", "uterus", "kidney", "bladder", "liver", "bile duct", "gall bladder", "stomach", "small intestine", "colon", "esophagus", "oral", "duodenum", "abdomen", "iris", "pupil", "larynx", "trachea", "lung", "nose", "transitional epithelium", "mesothelium", "heart", "widespread"}
}

func TestNewMatcher(t *testing.T) {
	m := NewMatcher()
	if len(m.location) == 0 {
		t.Error("Matcher locations were not read from file.")
	} else {
		for _, i := range getTissueLocations() {
			if _, ex := m.location[i]; ex == false {
				t.Errorf("%s not found in location map.", i)
			}
		}
	}
}

func TestTumor(t *testing.T) {
	// Tests getMatch method
	m := NewMatcher()
	matches := NewMatches()
	for _, i := range matches {
		typ, loc, mal := m.GetTumor(i.Line, i.Sex, true)
		if typ != i.Typ {
			t.Errorf("Actual type %s does not equal expected: %s.", typ, i.Typ)
		} else if loc != i.Location {
			t.Errorf("Actual location %s does not equal expected: %s.", loc, i.Location)
		} else if mal != i.Malignant {
			t.Errorf("Actual malignant value %s does not equal expected: %s.", mal, i.Malignant)
		}
	}
}

func TestGetCastrated(t *testing.T) {
	m := NewMatcher()
	matches := NewMatches()
	for _, i := range matches {
		actual := m.GetCastrated(i.Line)
		if actual != i.Castrated {
			t.Errorf("Actual neuter value %s does not equal expected: %s.", actual, i.Castrated)
		}
	}
}

func TestInfantRecords(t *testing.T) {
	m := NewMatcher()
	matches := NewMatches()
	for _, i := range matches {
		actual := m.InfantRecords(i.Line)
		if actual != i.Infant {
			t.Errorf("Actual infant record value %v does not equal expected: %v.", actual, i.Infant)
		}
	}
}

func TestGetAge(t *testing.T) {
	m := NewMatcher()
	matches := NewMatches()
	for _, i := range matches {
		actual := m.GetAge(i.Line)
		if actual != i.Age {
			t.Errorf("Actual age %s does not equal expected: %s.", actual, i.Age)
		}
	}
}

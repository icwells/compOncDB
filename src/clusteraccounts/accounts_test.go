// Test functions for accounts struct

package clusteraccounts

import (
	"testing"
)

func getTestQueries() map[string][]string {
	// Returns map of expected queries
	return map[string][]string{
		"Payson Animal Hospital":        {"payson animal hospital", "Payson A. H.", "Payson anim. Hospital"},
		"Phoenix Veterinary Services":   {"phoenix v.s.", "phoenix v. s."},
		"Phoenix Veterinarian Services": {"phoenix Veterinarian services"},
	}
}

func getTestPool() map[string][]string {
	// Returns map of expected pool
	return map[string][]string{
		"Payson Animal Hospital":        {"Payson Animal Hospital", "Payson Animal Hospital", "Payson Animal Hospital"},
		"Phoenix Veterinary Services":   {"Phoenix Veterinary Services", "Phoenix Veterinary Services"},
		"Phoenix Veterinarian Services": {"Phoenix Veterinarian Services"},
	}
}

func getTestScores() map[string]int {
	// Returns map of expected scores
	return map[string]int{
		"Payson Animal Hospital":        5,
		"Phoenix Veterinary Services":   5,
		"Phoenix Veterinarian Services": 4,
	}
}

func getTestTerms() map[string]string {
	// Returns map of expected terms
	return map[string]string{
		"payson animal hospital":        "Payson Animal Hospital",
		"Payson A. H.":                  "Payson Animal Hospital",
		"Payson anim. Hospital":         "Payson Animal Hospital",
		"phoenix v.s.":                  "Phoenix Veterinary Services",
		"phoenix v. s.":                 "Phoenix Veterinary Services",
		"phoenix Veterinarian services": "Phoenix Veterinary Services",
	}
}

func TestSetScores(t *testing.T) {
	// Sets setScores (fuzzymatch and checkSpelling)
	var s []string
	expected := getTestScores()
	a := NewAccounts("")
	a.queries = getTestQueries()
	a.pool = getTestPool()
	for k := range a.pool {
		s = append(s, k)
	}
	a.setScores(s)
	for k, v := range expected {
		a, ex := a.scores[k]
		if ex == false {
			t.Errorf("Expected key %s not present in scores map.", k)
		} else if v != a {
			t.Errorf("Actual score %d does not equal expected: %d.", a, v)
		}
	}
}

func TestSetQueries(t *testing.T) {
	// Tests setQueries (by extension: checkAmpersand, checkPeriods, and checkAbbreviations)
	expected := getTestQueries()
	a := NewAccounts("")
	var s []string
	for _, v := range expected {
		s = append(s, v...)
	}
	a.setQueries(s, false)
	for k, v := range a.queries {
		e, ex := expected[k]
		if ex == false {
			t.Errorf("Actual queries key %s not present in expected map.", k)
		} else {
			for idx, i := range v {
				if i != e[idx] {
					t.Errorf("Actual queries value %s does not equal expected: %s.", i, e[idx])
				}
			}
		}
	}
}

/*func TestClusterNames(t *testing.T) {
	a := newAccounts()
	a.queries = getTestQueries()
	a.clusterNames()
	t.Error(a.pool)
}

func TestSetTerms(t *testing.T) {
	expected := getTestTerms()
	a := newAccounts()
	a.queries = getTestQueries()
	a.pool = getTestPool()
}*/

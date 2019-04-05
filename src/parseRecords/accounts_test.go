// Test functions for accounts struct

package main

import (
	"testing"
)
func getTestQueries() map[string][]string {
	// Returns map of expected queries
	return map[string][]string{
		"Payson Animal Hospital": {"payson animal hospital", "Payson A. H.", "Payson anim. Hospital"},
		"Phoenix Veterinary Services": {"phoenix v.s.", "phoenix v. s."},
		"Phoenix Veterinarian Services": {"phoenix Veterinarian services"},
	}
}

func TestSetQueries(t *testing.T) {
	// Tests set queries (by extension: checkAmpersand and checkAbbreviations)
	expected := getTestQueries()
	a := newAccounts()
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

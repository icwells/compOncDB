// Test functions for accounts struct

package clusteraccounts

import (
	"github.com/icwells/simpleset"
	"strings"
	"testing"
)

func getTestTerms() map[string][]string {
	// Returns map of expected terms
	return map[string][]string{
		"payson animal hospital":        []string{"Payson Animal Hospital", "0", "0", "0"},
		"Payson A. H.":                  []string{"Payson Animal Hospital", "0", "0", "0"},
		"Payson anim. Hospital":         []string{"Payson Animal Hospital", "0", "0", "0"},
		"phoenix v.s.":                  []string{"Phoenix Veterinary Services", "0", "0", "1"},
		"phoenix v. s.":                 []string{"Phoenix Veterinary Services", "0", "0", "1"},
		"phoenix Veterinarian services": []string{"Phoenix Veterinarian Services", "0", "0", "1"},
		"matt":                          []string{"Matt", "0", "0", "0"},
		" zoo; Phoenix ":                []string{"Phoenix Zoo", "1", "1", "0"},
		" Phoenix zoo ":                 []string{"Phoenix Zoo", "1", "1", "0"},
		"tuscon aquarium":               []string{"Tuscon Aquarium", "1", "0", "0"},
		"wildlife rescue center":        []string{"Wildlife Rescue Center", "0", "0", "1"},
		"wildlfe rescue center":         []string{"Wildlife Rescue Center", "0", "0", "1"},
		"lemur Institute":               []string{"Lemur Institute", "0", "0", "1"},
	}
}

func getCorpus(terms map[string][]string) []string {
	// Creates corpus from test terms
	s := simpleset.NewStringSet()
	for _, v := range terms {
		for _, i := range strings.Split(v[0], " ") {
			if i != "Payson" && i != "wildlfe rescue center" {
				s.Add(i)
			}
		}
	}
	return s.ToStringSlice()
}

func TestCheckAbbreviations(t *testing.T) {
	// Tests corpus and checkAbbreviations (by extension: checkAmpersand, checkPeriods, and checkCaps)
	expected := getTestTerms()
	corpus := getCorpus(expected)
	delete(expected, "wildlfe rescue center")
	a := NewAccounts("")
	for k, v := range expected {
		ch := make(chan string)
		go a.checkAbbreviations(ch, k)
		act := <-ch
		if act != v[0] {
			t.Errorf("Actual formatted value %s does not equal expected: %s.", act, v[0])
		}
		a.terms = append(a.terms, newTerm(k, act))
	}
	for _, i := range corpus {
		if ex, _ := a.corpus.InSet(i); !ex {
			t.Errorf("Expected value %s not in accounts corpus.", i)
		}
	}
}

func TestResolveAccounts(t *testing.T) {
	// Tests spelling correction and clustering
	c := []string{"corrected name", "zoo", "aza", "institute"}
	expected := getTestTerms()
	a := NewAccounts("")
	a.zoos = []string{"Phoenix Zoo"}
	for k := range expected {
		a.Queries.Add(k)
	}
	act := a.ResolveAccounts()
	for k, v := range expected {
		if _, ex := act[k]; ex == false {
			t.Errorf("Expected term %s not in actual accounts map.", k)
		} else {
			for idx, i := range v {
				if i != act[k][idx] {
					t.Errorf("Actual %s column value for %s %s does not equal expected: %s", c[idx], k, act[k][idx], i)
				}
			}
		}
	}
}

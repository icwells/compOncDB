// Tests species search functions

package search

import (
	"strings"
	"testing"
)

func testSpeciesSearcher() speciesSearcher {
	// Initilaizes struct for testing
	var s speciesSearcher
	s.species = map[string]string{
		"Coyote":                   "1",
		"Canis latrans":            "1",
		"Wolf":                     "2",
		"Canis lupus":              "2",
		"Gray fox":                 "3",
		"Urocyon cinereoargenteus": "3",
	}
	s.taxa = map[string][]string{
		"1": []string{"Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis latrans", "http://eol.org/api/hierarchy_entries/1.0.xml?id=52440711"},
		"2": []string{"Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus", "http://eol.org/api/hierarchy_entries/1.0.xml?id=52624675"},
		"3": []string{"Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Urocyon", "Urocyon cinereoargenteus", "http://eol.org/api/hierarchy_entries/1.0.xml?id=52578011"},
	}
	for k := range s.species {
		s.list = append(s.list, k)
	}
	return s
}

func TestGetTaxonomy(t *testing.T) {
	ch := make(chan []string)
	s := testSpeciesSearcher()
	input := map[string]string{
		"COYOTE":      "1",
		"canis lupus": "2",
		"Gray Fox":    "3",
		"Wolf":        "2",
		"fox":         "",
	}
	for k, v := range input {
		go s.getTaxonomy(ch, k)
		row := <-ch
		if len(v) == 1 {
			if len(row) == 0 {
				t.Errorf("No result returned for %s.", k)
			} else if row[0] != k {
				t.Errorf("Incorrect query %s returned for %s.", row[0], k)
			} else if s.species[row[1]] != v {
				t.Errorf("Incorrect species %s returned for %s.", row[1], k)
			} else {
				actual := strings.Join(row[2:], ",")
				if actual != strings.Join(s.taxa[v], ",") {
					t.Errorf("Expected taxonomy for %s does not equal actual: %s", k, strings.Join(row, ","))
				}
			}
		} else if len(row) > 0 {
			t.Errorf("Unexpected results returned for %s.", k)
		}
	}
}

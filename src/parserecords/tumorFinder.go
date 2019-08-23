// Contians tumorFinder struct and methods for parseRecords

package parserecords

import (
	"strconv"
	"strings"
)

type tumorType struct {
	match     string
	index     int
	length    int
	locations map[string]int
	location  string
}

func newTumorType(m string, i, l int) *tumorType {
	// Initializes new struct
	var t tumorType
	t.match = m
	t.index = i
	t.length = l
	t.locations = make(map[string]int)
	t.location = "NA"
	return &t
}

func (t *tumorType) setDistance(k, m, line string) {
	// Stores location hit and distance from type
	var dist int
	idx := strings.Index(line, m)
	if idx == t.index {
		dist = 0
	} else if idx > t.index {
		// Location index - type index + length of type
		dist = idx - (t.index + len(t.match))
	} else {
		// Type index - (location index + location length)
		dist = t.index - (idx + len(m))
	}
	t.locations[k] = dist
}

func (t *tumorType) setLocation() {
	// Determines location with shortest distance from type
	var loc string
	min := t.length
	for k, v := range t.locations {
		if v < min {
			if k == "abdomen" || k == "fat" {
				// Only keep vague terms if there are no other potential matches
				loc = k
			} else {
				min = v
				t.location = k
				if min <= 1 {
					// Accept neighboring word
					break
				}
			}
		}
	}
	if t.location == "NA" && len(loc) > 0 {
		t.location = loc
	}
}

//----------------------------------------------------------------------------

type tumorFinder struct {
	types     map[string]*tumorType
	malignant string
}

func newTumorFinder() tumorFinder {
	// Initializes new struct
	var t tumorFinder
	t.types = make(map[string]*tumorType)
	t.malignant = "-1"
	return t
}

func (t *tumorFinder) toStrings() (string, string, string) {
	// Returns values as strings
	if len(t.types) == 0 {
		return "NA", "NA", t.malignant
	} else {
		var types, locations []string
		for k, v := range t.types {
			types = append(types, k)
			locations = append(locations, v.location)
		}
		return strings.Join(types, ";"), strings.Join(locations, ";"), t.malignant
	}
}

//----------------------------------------------------------------------------

func (m *matcher) getLocations(t *tumorFinder, line string) {
	// Searches line preceding type index for locations
	for key := range t.types {
		for k, v := range m.location {
			// Search for matches in words between previous and current match
			match := m.getMatch(v, line)
			if match != "NA" {
				t.types[key].setDistance(k, match, line)
			}
		}
		t.types[key].setLocation()
	}
}

func (m *matcher) setMalignant(t *tumorFinder, line string) {
	// Sets malignant value for tumorFinder
	for i := range t.types {
		for k := range m.types {
			// Get sub-map
			if _, ex := m.types[k][i]; ex == true {
				vm, _ := strconv.Atoi(m.types[k][i].malignant)
				tm, _ := strconv.Atoi(t.malignant)
				if vm > tm {
					// Malignant > non-malignant > NA
					t.malignant = m.types[k][i].malignant
				}
				break
			}
		}
	}
	if t.malignant == "-1" {
		t.malignant = m.getMalignancy(line)
	}
}

func (m *matcher) getTypes(t *tumorFinder, line string) {
	// Returns types from map
	for key := range m.types {
		found := false
		var term, typ string
		for k, v := range m.types[key] {
			match := m.getMatch(v.expression, line)
			if match != "NA" {
				if key == "other" || k != key {
					// Keep specific diagnosis terms in struct
					t.types[k] = newTumorType(match, strings.Index(line, match), len(line))
					found = true
				} else {
					// Store potentially overlapping terms
					term = match
					typ = k
				}
			}
		}
		if found == false && len(typ) > 1 {
			t.types[typ] = newTumorType(term, strings.Index(line, term), len(line))
		}
	}
}

func (m *matcher) getTumor(line string, cancer bool) (string, string, string) {
	// Returns type, location, and malignancy
	t := newTumorFinder()
	if cancer == true {
		m.getTypes(&t, line)
		if len(t.types) > 0 {
			m.setMalignant(&t, line)
			m.getLocations(&t, line)
		}
	}
	return t.toStrings()
}

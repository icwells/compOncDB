// Contians tumorFinder struct and methods for parseRecords

package diagnoses

import (
	"strconv"
	"strings"
)

type tumorHit struct {
	match     string
	index     int
	length    int
	locations map[string]int
	location  string
}

func newTumorHit(m string, i, l int) *tumorHit {
	// Initializes new struct
	var t tumorHit
	t.match = m
	t.index = i
	t.length = l
	t.locations = make(map[string]int)
	t.location = "NA"
	return &t
}

func (t *tumorHit) setDistance(k, m, line string) {
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

func (t *tumorHit) setLocation() {
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
	types     map[string]*tumorHit
	malignant string
}

func newTumorFinder() *tumorFinder {
	// Initializes new struct
	var t tumorFinder
	t.types = make(map[string]*tumorHit)
	t.malignant = "-1"
	return &t
}

func (t *tumorFinder) checkKeys(name string, idx int) {
	// Removes incomplete tumor name matches
	for key := range t.types {
		if strings.Contains(name, key) {
			delete(t.types, key)
		}
	}
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

func (m *Matcher) setMalignant(t *tumorFinder, line string) {
	// Sets malignant value for tumorFinder; searches for match if malignant and benign for type are both or both false
	malignant := -1
	for key := range t.types {
		mal := -1
		if key == "widespread" {
			mal = 1
		} else if m.types[key].benign && !m.types[key].malignant {
			mal = 0
		} else if !m.types[key].benign && m.types[key].malignant {
			mal = 1
		} else if m.types[key].benign && m.types[key].malignant {
			mal, _ = strconv.Atoi(m.GetMalignancy(line))
		} else if !m.types[key].benign && !m.types[key].malignant {
			mal, _ = strconv.Atoi(m.GetMalignancy(line))
		}
		if mal > malignant {
			malignant = mal
			if malignant == 1 {
				break
			}
		}
	}
	t.malignant = strconv.Itoa(malignant)
}

func (m *Matcher) searchLocation(t *tumorFinder, line, key, i, sex string) {
	// Searches for a match to given location
	match := m.GetMatch(m.location[i], line)
	if match != "NA" {
		if match == "interstitial" {
			if sex == "male" {
				i = "testis"
			} else if sex == "female" {
				i = "ovary"
			}
		}
		t.types[key].setDistance(i, match, line)
	}
}

func (m *Matcher) getLocations(t *tumorFinder, line, sex string) {
	// Searches line for locations of matches
	for key := range t.types {
		locations := m.types[key].locations.ToSlice()
		if len(locations) > 0 {
			// Search for matches in known locations
			for _, i := range locations {
				m.searchLocation(t, line, key, i, sex)
			}
		} else {
			// Search for match to any location
			for k := range m.location {
				m.searchLocation(t, line, key, k, sex)
			}
		}
		t.types[key].setLocation()
	}
}

func (m *Matcher) getTypes(t *tumorFinder, line string) {
	// Returns types from map
	for k, v := range m.types {
		match := m.GetMatch(v.expression, line)
		if match != "NA" {
			idx := strings.Index(line, match)
			t.checkKeys(k, idx)
			t.types[k] = newTumorHit(match, idx, len(line))
		}
	}
}

func (m *Matcher) GetTumor(line, sex string, cancer bool) (string, string, string) {
	// Returns type, location, and malignancy
	t := newTumorFinder()
	if cancer == true {
		m.getTypes(t, line)
		if len(t.types) > 0 {
			m.getLocations(t, line, sex)
			m.setMalignant(t, line)
		}
	}
	return t.toStrings()
}

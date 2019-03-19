// Contians tumorFinder strcut and methods for parseRecords

package main

import (
	"strconv"
	"strings"
)

type tumorFinder struct {
	matches   []string
	types     []string
	locations []string
	malignant string
}

func newTumorFinder() tumorFinder {
	// Initializes new struct
	var t tumorFinder
	t.malignant = "-1"
	return t
}

func (t *tumorFinder) toStrings() (string, string, string) {
	// Returns values as strings
	if len(t.types) == 0 {
		t.types = []string{"NA"}
	}
	if len(t.locations) == 0 {
		t.locations = []string{"NA"}
	}
	return strings.Join(t.types, ";"), strings.Join(t.locations, ";"), t.malignant
}

func (t *tumorFinder) getSearchIndeces(idx int, line string) (int, int) {
	// Returns indeces of last match and next match
	start, end := 0, len(line)
	if idx > 0 {
		last := idx - 1
		start = strings.Index(line, t.matches[last]) + len(t.matches[last])
	}
	if start < len(line) {
		// Include type as it might be informative (i.e. lymphoma)
		end = strings.Index(line, t.matches[idx]) + len(t.matches[idx])
	}
	if start > end {
		// Reset illogical results
		start, end = -1, -1
	}
	return start, end
}

//----------------------------------------------------------------------------

func (m *matcher) getTypes(t *tumorFinder, line string) {
	// Returns location from map
	for k, v := range m.types {
		match := m.getMatch(v.expression, line)
		if match != "NA" {
			t.matches = append(t.matches, match)
			t.types = append(t.types, k)
			vm, _ := strconv.Atoi(v.malignant)
			tm, _ := strconv.Atoi(t.malignant)
			if vm > tm {
				// Malignant > non-malignant > NA
				t.malignant = v.malignant
			}
		}
	}
	if t.malignant == "-1" && len(t.types) > 0 {
		t.malignant = m.getMalignancy(line)
	}
}

func (m *matcher) getLocations(t *tumorFinder, line string) {
	// Searches line preceding type index for locations
	for idx, _ := range t.matches {
		loc := "NA"
		start, end := t.getSearchIndeces(idx, line)
		if start >= 0 && end < len(line) {
			for k, v := range m.location {
				// Search for matches in words between previous and current match
				match := m.getMatch(v, line[start:end])
				if match != "NA" {
					loc = k
					if loc != "widespread" && loc != "other" {
						// Break if descriptive match is found
						break
					}
				}
			}
		}
		// Append one location for each type
		t.locations = append(t.locations, loc)
	}
}

func (m *matcher) getTumor(line string, cancer bool) (string, string, string) {
	// Returns type, location, and malignancy
	t := newTumorFinder()
	if cancer == true {
		m.getTypes(&t, line)
		if len(t.types) > 0 {
			m.getLocations(&t, line)
		}
	}
	return t.toStrings()
}

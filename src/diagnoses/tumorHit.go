// Defines tumor hit struct

package diagnoses

import (
	"github.com/icwells/go-tools/strarray"
	"strings"
)

func FilterMultipleHits(typ, loc string) (string, string) {
	// Removes duplicate hits and hits missing either location/type
	d := ";"
	if strings.Contains(typ, d) {
		var rm []int
		l := strings.Split(loc, d)
		t := strings.Split(typ, d)
		for idx, i := range t {
			if i == "NA" || l[idx] == "NA" {
				if len(rm) < len(t)-1 {
					// Only remove if there is at least one record left
					rm = append(rm, idx)
				}
			} else if idx > 0 && strarray.InSliceStr(t[:idx], i) && strarray.InSliceStr(l[:idx], l[idx]) {
				// Remove duplicate diagnoses
				rm = append(rm, idx)
			}
		}
		for _, i := range rm {
			t = strarray.DeleteSliceIndex(t, i)
			l = strarray.DeleteSliceIndex(l, i)
		}
		typ = strings.Join(t, d)
		loc = strings.Join(l, d)
	}
	return typ, loc
}

type tumorHit struct {
	end       int
	index     int
	indeces   map[string]int
	length    int
	locations map[string]int
	location  string
	match     string
}

func newTumorHit(m string, i, l int) *tumorHit {
	// Initializes new struct
	var t tumorHit
	t.match = strings.ToLower(m)
	t.index = i
	t.indeces = make(map[string]int)
	t.length = l
	t.locations = make(map[string]int)
	t.location = "NA"
	return &t
}

func (t *tumorHit) setDistance(k, m, line string) {
	// Stores location hit and distance from type
	var dist int
	k = strings.ToLower(k)
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
	t.indeces[k] = idx
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
		if t.indeces[loc] > t.index {
			t.end = t.indeces[loc] + len(loc)
		}
	} else if strings.Contains(t.match, "sarcoma") && t.location == "NA" {
		t.location = "sarcoma"
		t.end = t.index + len(t.match)
	}
}

// Contians tumorFinder struct and methods for parseRecords

package diagnoses

import (
	"github.com/icwells/go-tools/strarray"
	"sort"
	"strconv"
	"strings"
)

var (
	FEMALE = []string{"ovary", "vulva", "uterus", "oviduct"}
	MALE   = []string{"testis", "prostate"}
)

type tumorFinder struct {
	hits      []*tumorHit
	malignant string
	types     map[string]*tumorHit
}

func newTumorFinder() *tumorFinder {
	// Initializes new struct
	var t tumorFinder
	t.types = make(map[string]*tumorHit)
	t.malignant = "-1"
	return &t
}

func (t *tumorFinder) checkKeys(name string, idx int) bool {
	// Removes incomplete tumor name matches
	ret := true
	for key := range t.types {
		if strings.Contains(name, key) {
			delete(t.types, key)
		} else if strings.Contains(key, name) {
			ret = false
			break
		}
	}
	return ret
}

func (t *tumorFinder) Len() int {
	return len(t.hits)
}

func (t *tumorFinder) Less(i, j int) bool {
	return t.hits[i].end < t.hits[j].end
}

func (t *tumorFinder) Swap(i, j int) {
	t.hits[i], t.hits[j] = t.hits[j], t.hits[i]
}

func (t *tumorFinder) sortHits() {
	// Returns slice if hits ordered by end index
	for _, v := range t.types {
		t.hits = append(t.hits, v)
	}
	sort.Sort(t)
}

func (t *tumorFinder) subsetLine(line string, start, end int) string {
	// Slices line between start and end
	var ret string
	if start < end && end < len(line) {
		ret = line[start : end+1]
	}
	return ret
}

func (t *tumorFinder) getGrowthType(i *tumorHit) (string, string) {
	// Returns neoplasia and hyperplasia values
	neoplasia, hyperplasia := "1", "0"
	if i.match == "hyperplasia" {
		neoplasia = "0"
		hyperplasia = "1"
	}
	return neoplasia, hyperplasia
}

func (t *tumorFinder) SplitStrings(line string) [][]string {
	// Splits line so that each piece contains one tumor diagnosis
	var ret [][]string
	if len(t.types) == 0 {
		ret = append(ret, []string{line, "0", "0", "NA", "NA"})
	} else {
		var start int
		t.sortHits()
		for _, i := range t.hits[:len(t.hits)-1] {
			if s := t.subsetLine(line, start, i.end); s != "" {
				neoplasia, hyperplasia := t.getGrowthType(i)
				row := []string{s, neoplasia, hyperplasia, i.match, i.location}
				ret = append(ret, row)
				start = i.end
			}
		}
		i := t.hits[len(t.hits)-1]
		if s := t.subsetLine(line, start, len(line)-1); s != "" {
			neoplasia, hyperplasia := t.getGrowthType(i)
			ret = append(ret, []string{s, neoplasia, hyperplasia, i.match, i.location})
		}
	}
	return ret
}

func (t *tumorFinder) toStrings(tissues map[string]string) (string, string, string, string) {
	// Returns values as strings
	if len(t.types) == 0 {
		return "NA", "NA", "NA", t.malignant
	} else {
		var types, tissue, locations []string
		for k, v := range t.types {
			types = append(types, k)
			locations = append(locations, v.location)
			if val, ex := tissues[v.location]; ex {
				tissue = append(tissue, val)
			} else {
				tissue = append(tissue, "NA")
			}
		}
		return strings.Join(types, ";"), strings.Join(tissue, ";"), strings.Join(locations, ";"), t.malignant
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
	pass := true
	match := m.GetMatch(m.location[i], line)
	if match != "NA" {
		if match == "interstitial" {
			if sex == "male" {
				i = "testis"
			} else if sex == "female" {
				i = "ovary"
			}
		} else if strarray.InSliceStr(FEMALE, i) && sex == "male" {
			pass = false
		} else if strarray.InSliceStr(MALE, i) && sex == "female" {
			pass = false
		}
		if pass {
			t.types[key].setDistance(i, match, line)
		}
	}
}

func (m *Matcher) getLocations(t *tumorFinder, line, sex string) {
	// Searches line for locations of matches
	for key := range t.types {
		// Search for matches in known locations
		for _, i := range m.types[key].locations.ToStringSlice() {
			m.searchLocation(t, line, key, i, sex)
		}
		// Search for match to any location
		for k := range m.location {
			m.searchLocation(t, line, key, k, sex)
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
			if t.checkKeys(k, idx) {
				t.types[k] = newTumorHit(match, idx, len(line))
			}
		}
	}
}

func (m *Matcher) SplitOnTumors(line, sex string) [][]string {
	// Splits input so that each piece contains only one tumor type and location
	t := newTumorFinder()
	m.getTypes(t, line)
	if len(t.types) > 0 {
		m.getLocations(t, line, sex)
	}
	return t.SplitStrings(line)
}

func (m *Matcher) GetTumor(line, sex string, cancer bool) (string, string, string, string) {
	// Returns type, location, and malignancy
	t := newTumorFinder()
	if cancer == true {
		m.getTypes(t, line)
		if len(t.types) > 0 {
			m.getLocations(t, line, sex)
			m.setMalignant(t, line)
		}
	}
	return t.toStrings(m.tissues)
}

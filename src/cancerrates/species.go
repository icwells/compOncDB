// Species level struct for storing records

package cancerrates

import (
	"github.com/icwells/simpleset"
	"sort"
	"strings"
)

func emptySlice(n int) []string {
	// Returns empty slice of length n
	var ret []string
	for i := 0; i < n; i++ {
		ret = append(ret, "")
	}
	return ret
}

type species struct {
	grandtotal  int
	id          string
	infancy     float64
	lifehistory []string
	location    string
	locations   *simpleset.Set
	taxonomy    []string
	tissue      *record
	tissues     map[string]*record
	total       *record
}

func newSpecies(id, location string, taxonomy []string) *species {
	// Return new species struct
	s := new(species)
	s.id = id
	s.location = location
	s.locations = simpleset.NewStringSet()
	s.taxonomy = taxonomy
	s.tissue = newRecord()
	s.tissues = make(map[string]*record)
	s.total = newRecord()
	s.setLocations()
	return s
}

func (s *species) setLocations() {
	// Initializes location set
	if strings.Contains(s.location, ",") {
		s.location = strings.Replace(s.location, ",", ";", -1)
	}
	if strings.Contains(s.location, ";") {
		for _, i := range strings.Split(s.location, ";") {
			s.locations.Add(i)
			s.tissues[i] = newRecord()
		}
	} else {
		s.locations.Add(s.location)
	}
}

func (s *species) tissueSlice(name string, r *record) []string {
	// Formats rows for specific tissues
	ret := []string{s.id}
	ret = append(ret, emptySlice(len(s.taxonomy))...)
	ret = append(ret, name)
	ret = append(ret, r.calculateRates(s.total.total)...)
	if len(s.lifehistory) > 0 {
		ret = append(ret, emptySlice(len(s.lifehistory))...)
	}
	return ret
}

func (s *species) toSlice() [][]string {
	// Formats cancer rates and returns row for tissue and total
	var ret [][]string
	total := append([]string{s.id}, s.taxonomy...)
	total = append(total, "all")
	total = append(total, s.total.calculateRates(-1)...)
	if len(s.lifehistory) > 0 {
		total = append(total, s.lifehistory...)
	}
	ret = append(ret, total)
	if s.location != "" {
		ret = append(ret, s.tissueSlice(s.location, s.tissue))
		for k, v := range s.tissues {
			if v.grandtotal > 0 {
				ret = append(ret, s.tissueSlice(k, v))
			}
		}
	}
	return ret
}

func (s *species) highestMalignancy(mal string) string {
	// Returns highest malignacy code
	if strings.Contains(mal, ";") {
		m := strings.Split(mal, ";")
		sort.Strings(m)
		return m[len(m)-1]
	}
	return mal
}

func (s *species) checkLocation(mal, loc string) (bool, string) {
	// Returns true if s.location is in loc
	if loc != "" {
		if strings.Contains(loc, ";") {
			m := strings.Split(mal, ";")
			for idx, i := range strings.Split(loc, ";") {
				if ex, _ := s.locations.InSet(i); ex {
					return true, m[idx]
				}
			}
		} else if ex, _ := s.locations.InSet(loc); ex {
			return true, mal
		}
	}
	return false, ""
}

func (s *species) addCancer(age float64, sex, nec, mal, loc, service, aid string) {
	// Adds cancer measures
	s.total.cancerMeasures(age, sex, s.highestMalignancy(mal), service)
	if eq, m := s.checkLocation(mal, loc); eq {
		// Add all measures for target tissue
		s.tissue.cancerMeasures(age, sex, m, service)
		s.tissue.nonCancerMeasures(age, sex, nec, service, aid)
		if _, ex := s.tissues[loc]; ex {
			// Add to specific location
			s.tissues[loc].cancerMeasures(age, sex, m, service)
			s.tissues[loc].nonCancerMeasures(age, sex, nec, service, aid)
		}
	}
}

func (s *species) addNonCancer(age float64, sex, nec, service, aid string) {
	// Adds non-cancer measures
	s.total.nonCancerMeasures(age, sex, nec, service, aid)
}

func (s *species) addDenominator(d int) {
	// Adds denominator to records
	//s.tissue.addTotal(d)
	s.total.addTotal(d)
}

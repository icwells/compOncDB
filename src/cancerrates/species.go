// Species level struct for storing records

package cancerrates

import (
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
	taxonomy    []string
	tissue      *record
	total       *record
}

func newSpecies(id, location string, taxonomy []string) *species {
	// Return new species struct
	s := new(species)
	s.id = id
	s.location = location
	s.taxonomy = taxonomy
	s.tissue = newRecord()
	s.total = newRecord()
	return s
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
		tissue := []string{s.id}
		tissue = append(tissue, emptySlice(len(s.taxonomy))...)
		tissue = append(tissue, s.location)
		tissue = append(tissue, s.tissue.calculateRates(s.total.total)...)
		if len(s.lifehistory) > 0 {
			tissue = append(tissue, emptySlice(len(s.lifehistory))...)
		}
		ret = append(ret, tissue)
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
				if i == s.location {
					return true, m[idx]
				}
			}
		} else if loc == s.location {
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
	}
}

func (s *species) addNonCancer(age float64, sex, nec, service, aid string) {
	// Adds non-cancer measures
	s.total.nonCancerMeasures(age, sex, nec, service, aid)
}

func (s *species) addDenominator(d int) {
	// Adds denominator to records
	s.tissue.addTotal(d)
	s.total.addTotal(d)
}

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

type Species struct {
	denominator  int
	Grandtotal   int
	id           string
	ids          *simpleset.Set
	infancy      float64
	lifehistory  []string
	Location     string
	locationprop float64
	notissue     int
	other        *Record
	taxonomy     []string
	tissue       *Record
	total        *Record
}

func newSpecies(id, location string, taxonomy []string) *Species {
	// Return new species struct
	s := new(Species)
	s.id = id
	s.ids = simpleset.NewStringSet()
	s.Location = location
	s.locationprop = 0.05
	s.other = newRecord()
	s.taxonomy = taxonomy
	s.tissue = newRecord()
	s.total = newRecord()
	return s
}

func (s *Species) tissueSlice(name string, r *Record) []string {
	// Formats rows for specific tissues
	ret := []string{s.id}
	ret = append(ret, emptySlice(len(s.taxonomy))...)
	ret = append(ret, name)
	ret = append(ret, r.calculateRates(s.denominator, -1)...)
	if len(s.lifehistory) > 0 {
		ret = append(ret, emptySlice(len(s.lifehistory))...)
	}
	return ret
}

func (s *Species) ToSlice(keepall bool) [][]string {
	// Formats cancer rates and returns row for tissue and total
	var ret [][]string
	if keepall || s.Location == "" || s.tissue.total > 0 {
		// Keep records with target tissue or at least 5% of records have locations
		s.denominator = s.total.total - s.notissue
		total := append([]string{s.id}, s.taxonomy...)
		if s.Location != "" {
			total = append(total, "all")
			total = append(total, s.total.calculateRates(s.total.total, s.notissue)...)
		} else {
			// Omit location column
			total = append(total, s.total.calculateRates(-1, s.notissue)...)
		}
		if len(s.lifehistory) > 0 {
			total = append(total, s.lifehistory...)
		}
		ret = append(ret, total)
		if s.Location != "" {
			ret = append(ret, s.tissueSlice(s.Location, s.tissue))
			ret = append(ret, s.tissueSlice("Other", s.other))
		}
	}
	return ret
}

func (s *Species) highestMalignancy(mal string) string {
	// Returns highest malignacy code
	if strings.Contains(mal, ";") {
		m := strings.Split(mal, ";")
		sort.Strings(m)
		return m[len(m)-1]
	}
	return mal
}

func (s *Species) checkLocation(mal, loc string) (bool, string) {
	// Returns true if s.location is in loc
	if loc != "" && loc != "NA" {
		if strings.Contains(loc, ";") {
			m := strings.Split(mal, ";")
			for idx, i := range strings.Split(loc, ";") {
				if s.Location == i {
					if idx < len(m) {
						return true, m[idx]
					}
				}
			}
		} else if s.Location == loc {
			return true, mal
		}
	}
	if strings.Contains(mal, ";") {
		ret := "-1"
		for _, i := range strings.Split(mal, ";") {
			if i == "1" {
				return false, "1"
			} else if i == "0" {
				ret = "0"
			}
		}
		return false, ret
	}
	return false, mal
}

func (s *Species) addCancer(allrecords bool, age, sex, nec, mal, loc, service, aid string) {
	// Adds cancer measures
	s.total.cancerMeasures(allrecords, age, sex, s.highestMalignancy(mal), service)
	eq, m := s.checkLocation(mal, loc)
	if eq {
		// Add all measures for target tissue
		s.tissue.cancerMeasures(allrecords, age, sex, m, service)
		s.tissue.nonCancerMeasures(allrecords, age, sex, nec, service, aid)
	} else {
		s.other.cancerMeasures(allrecords, age, sex, m, service)
		s.other.nonCancerMeasures(allrecords, age, sex, nec, service, aid)
	}
}

func (s *Species) addNonCancer(allrecords bool, age, sex, nec, service, aid, id string) {
	// Adds non-cancer measures
	s.total.nonCancerMeasures(allrecords, age, sex, nec, service, aid)
	s.Grandtotal = s.total.grandtotal
	s.ids.Add(id)
}

func (s *Species) addDenominator(masspresent, loc string) {
	// Adds to notissue if no reported location
	if masspresent == "1" {
		if loc == "NA" || loc == "" {
			s.notissue++
		} else if strings.Contains(loc, ";") {
			var found bool
			for _, i := range strings.Split(loc, ";") {
				if i != "NA" {
					found = true
					break
				}
			}
			if !found {
				s.notissue++
			}
		}
	}
}

func (s *Species) AddTissue(v *Species) {
	// Adds v.tissue to s.tissue
	s.tissue.Add(v.tissue)
}

func (s *Species) Copy() *Species {
	// Returns deep copy of struct
	ret := newSpecies(s.id, s.Location, s.taxonomy)
	ret.Grandtotal = s.Grandtotal
	ret.infancy = s.infancy
	ret.lifehistory = s.lifehistory
	ret.tissue = s.tissue.Copy()
	ret.total = s.total.Copy()
	return ret
}

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
	infancy      float64
	lifehistory  []string
	Location     string
	locationprop float64
	locations    *simpleset.Set
	notissue     int
	taxonomy     []string
	tissue       *Record
	tissues      map[string]*Record
	total        *Record
}

func newSpecies(id, location string, taxonomy []string) *Species {
	// Return new species struct
	s := new(Species)
	s.id = id
	s.Location = location
	s.locationprop = 0.05
	s.locations = simpleset.NewStringSet()
	s.taxonomy = taxonomy
	s.tissue = newRecord()
	s.tissues = make(map[string]*Record)
	s.total = newRecord()
	s.setLocations()
	return s
}

func (s *Species) setLocations() {
	// Initializes location set
	if strings.Contains(s.Location, ",") {
		s.Location = strings.Replace(s.Location, ",", ";", -1)
	}
	if strings.Contains(s.Location, ";") {
		for _, i := range strings.Split(s.Location, ";") {
			s.locations.Add(i)
			s.tissues[i] = newRecord()
		}
	} else {
		s.locations.Add(s.Location)
	}
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

func (s *Species) ToSlice() [][]string {
	// Formats cancer rates and returns row for tissue and total
	var ret [][]string
	if s.Location == "" || s.tissue.total > 0{
		// Keep records with target tissue or at least 5% of records have locations
		s.denominator = s.total.total - s.notissue
		total := append([]string{s.id}, s.taxonomy...)
		total = append(total, "all")
		total = append(total, s.total.calculateRates(-1, s.notissue)...)
		if len(s.lifehistory) > 0 {
			total = append(total, s.lifehistory...)
		}
		ret = append(ret, total)
		if s.Location != "" {
			ret = append(ret, s.tissueSlice(s.Location, s.tissue))
			for k, v := range s.tissues {
				if v.grandtotal > 0 {
					ret = append(ret, s.tissueSlice(k, v))
				}
			}
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

func (s *Species) addCancer(allrecords bool, age, sex, nec, mal, loc, service, aid string) {
	// Adds cancer measures
	s.total.cancerMeasures(allrecords, age, sex, s.highestMalignancy(mal), service)
	if eq, m := s.checkLocation(mal, loc); eq {
		// Add all measures for target tissue
		s.tissue.cancerMeasures(allrecords, age, sex, m, service)
		s.tissue.nonCancerMeasures(allrecords, age, sex, nec, service, aid)
		if _, ex := s.tissues[loc]; ex {
			// Add to specific location
			s.tissues[loc].cancerMeasures(allrecords, age, sex, m, service)
			s.tissues[loc].nonCancerMeasures(allrecords, age, sex, nec, service, aid)
		}
	}
}

func (s *Species) addNonCancer(allrecords bool, age, sex, nec, service, aid string) {
	// Adds non-cancer measures
	s.total.nonCancerMeasures(allrecords, age, sex, nec, service, aid)
	s.Grandtotal = s.total.grandtotal
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
	for k, v := range s.tissues {
		ret.tissues[k] = v.Copy()
	}
	ret.total = s.total.Copy()
	return ret
}

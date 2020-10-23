// Species level struct for storing records

package cancerrates

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
	total = append(total, s.total.calculateRates()...)
	if len(s.lifehistory) > 0 {
		total = append(total, s.lifehistory...)
	}
	ret = append(ret, total)
	if s.location != "" {
		tissue := []string{s.id}
		tissue = append(tissue, emptySlice(len(s.taxonomy))...)
		tissue = append(tissue, s.location)
		tissue = append(tissue, s.tissue.calculateRates()...)
		if len(s.lifehistory) > 0 {
			tissue = append(tissue, emptySlice(len(s.lifehistory))...)
		}
		ret = append(ret, tissue)
	}
	return ret
}

func (s *species) addCancer(age float64, sex, mal, loc, service, aid string) {
	// Adds cancer measures
	s.total.cancerMeasures(age, sex, mal, service)
	if loc == s.location {
		// Add all measures for target tissue
		s.tissue.cancerMeasures(age, sex, mal, service)
		if service != "MSU" {
			// Add to total and grandtotal
			s.tissue.addTotal(1)
			s.tissue.age++
			if sex == "male" {
				s.tissue.male++
			} else if sex == "female" {
				s.tissue.female++
			}
		} else {
			// Increment grand total
			s.tissue.grandtotal++
		}
		s.tissue.sources.Add(aid)
	}
}

func (s *species) addNonCancer(age float64, sex, service, aid string) {
	// Adds non-cancer measures
	if service != "MSU" {
		// Add to total and grandtotal
		s.total.addTotal(1)
		s.total.age++
		if sex == "male" {
			s.total.male++
		} else if sex == "female" {
			s.total.female++
		}
	} else {
		// Increment grand total
		s.total.grandtotal++
	}
	s.total.sources.Add(aid)
}

func (s *species) addDenominator(d int) {
	// Adds denominator to records
	s.tissue.addTotal(d)
	s.total.addTotal(d)
}

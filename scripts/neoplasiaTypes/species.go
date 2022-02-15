// Defines species struct for storing neoplasia by type

package main

import (
	"strconv"
)

type species struct {
	benign    map[string]int
	btotal    int
	cancer    int
	malignant map[string]int
	mtotal    int
	taxonomy  []string
	total     int
}

func newSpecies(taxonomy []string) *species {
	// Returns initialized struct
	s := new(species)
	s.benign = make(map[string]int)
	s.malignant = make(map[string]int)
	s.taxonomy = taxonomy
	return s
}

func (s *species) addNeoplasia(malignant int) {
	// Added to neoplasia and benign/malignant totals
	s.cancer++
	if malignant == 1 {
		s.mtotal++
	} else {
		s.btotal++
	}
}

func (s *species) addType(malignant int, typ string) {
	// Adds record to appropriate type map
	if malignant == 1 {
		if _, ex := s.malignant[typ]; !ex {
			s.malignant[typ] = 0
		}
		s.malignant[typ]++
	} else {
		if _, ex := s.benign[typ]; !ex {
			s.benign[typ] = 0
		}
		s.benign[typ]++
	}
}

func (s *species) divide(n, d int) string {
	// Divides n by d and returns formatted string
	if n != 0 && d != 0 {
		return strconv.FormatFloat(float64(n)/float64(d), 'f', 3, 64)
	}
	return "0"
}

func (s *species) getRates(l []string, m map[string]int) []string {
	// Returns benign/malignant rates as string slice
	var ret []string
	for _, i := range l {
		if v, ex := m[i]; ex {
			ret = append(ret, s.divide(v, s.cancer))
		} else {
			ret = append(ret, "")
		}
	}
	return ret
}

func (s *species) toSlice(malignant, benign []string) []string {
	// Returns values as string slice
	var ret []string
	ret = append(ret, s.taxonomy...)
	ret = append(ret, strconv.Itoa(s.total))
	ret = append(ret, strconv.Itoa(s.cancer))
	ret = append(ret, s.divide(s.cancer, s.total))
	ret = append(ret, strconv.Itoa(s.mtotal))
	ret = append(ret, s.divide(s.mtotal, s.cancer))
	ret = append(ret, strconv.Itoa(s.btotal))
	ret = append(ret, s.divide(s.btotal, s.cancer))
	ret = append(ret, s.getRates(malignant, s.malignant)...)
	ret = append(ret, s.getRates(benign, s.benign)...)
	return ret
}

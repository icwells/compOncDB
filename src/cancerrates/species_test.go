// Tests species struct

package cancerrates

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"testing"
)

func canidTaxa() [][]string {
	// Returns taxonomies for records
	var ret [][]string
	ret = append(ret, []string{"1", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus"})
	ret = append(ret, []string{"1", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis latrans"})
	ret = append(ret, []string{"3", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Vulpes", "Vulpes vulpes"})
	return ret
}

func locationSlice() [][]string {
	// Return slice of expected location values
	var ret [][]string
	ret = append(ret, []string{"100", "5", "5", "0.05", "2", "0.02", "0.40", "3", "0.03", "0.60", "10.00", "10.00", "2", "3", "2", "3", "1", "0"})
	ret = append(ret, []string{"110", "5", "5", "0.05", "2", "0.02", "0.40", "3", "0.03", "0.60", "10.00", "10.00", "2", "3", "2", "3", "1", "0"})
	return ret
}

func locationRecords() []*record {
	// Returns slice of records for testing
	var ret []*record
	ret = append(ret, setRecord([]float64{100, 5, 50.0, 2, 3, 5, 50.0, 2, 3, 2, 3, 1, 10, 4, 6}))
	ret = append(ret, setRecord([]float64{110, 5, 50.0, 2, 3, 5, 50.0, 2, 3, 2, 3, 1, 10, 4, 6}))
	return ret
}

func getExpectedSpecies() [][][]string {
	// Returns test slices
	var ret [][][]string
	taxa := canidTaxa()
	loc := locationSlice()
	for idx, i := range getExpectedRecords() {
		var sp [][]string
		r := append(taxa[idx], "all")
		sp = append(sp, append(r, i...))
		if idx < 2 {
			r = append([]string{taxa[idx][0], "", "", "", "", "", "", ""})
			r = append(r, "liver")
			r = append(r, loc[idx]...)
			sp = append(sp, r)
		}
		ret = append(ret, sp)
	}
	return ret
}

func setSpecies(taxa []string, r *record) *species {
	// Returns test struct
	s := newSpecies(taxa[0], "", taxa[1:])
	s.total = r
	return s
}

func getTestSpecies() []*species {
	// Returns slice of test structs
	var ret []*species
	taxa := canidTaxa()
	loc := locationRecords()
	for idx, i := range testRecords() {
		s := setSpecies(taxa[idx], i)
		if idx < 2 {
			s.location = "liver"
			s.tissue = loc[idx]
		}
		ret = append(ret, s)
	}
	return ret
}

func TestToSlice(t *testing.T) {
	head := codbutils.CancerRateHeader()
	expected := getExpectedSpecies()
	for ind, s := range getTestSpecies() {
		act := s.toSlice()
		exp := expected[ind]
		if len(act) != len(exp) {
			t.Errorf("%d: Actual number of rows %d does not equal expected: %d", ind, len(act), len(exp))
			break
		}
		for row, a := range act {
			if len(a) != len(exp[row]) {
				t.Errorf("%d: Actual length %d does not equal expected: %d", ind, len(a), len(exp[row]))
				break
			}
			for idx, i := range a {
				if i != exp[row][idx] {
					t.Errorf("%d %d: Actual slice value %s %s does not equal expected: %s", ind, row, head[idx], i, exp[row][idx])
					break
				}
			}
		}
	}
}

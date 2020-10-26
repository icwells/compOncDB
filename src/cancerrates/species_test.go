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

func setSpecies(taxa []string, r *record) *species {
	// Returns test struct
	s := newSpecies(taxa[0], "", taxa[1:])
	s.location = "liver"
	s.total = r
	return s
}

func getSpecies() []*species {
	// Returns species structs for testing
	var ret []*species
	taxa := canidTaxa()
	for idx, i := range testRecords() {
		ret = append(ret, setSpecies(taxa[idx], i))
	}
	return ret
}

func TestAddDenominators(t *testing.T) {
	sp := getSpecies()
	input := [][]int{
		{100, 200, 150},
		{50, 160, 150},
		{0, 50, 50},
	}
	for idx, i := range input {
		s := sp[idx]
		s.addDenominator(i[0])
		if s.total.grandtotal != i[1] {
			t.Errorf("%d: Actual grand total %d does not equal %d.", idx, s.total.grandtotal, i[1])
		} else if s.total.total != i[2] {
			t.Errorf("%d: Actual total %d does not equal %d.", idx, s.total.total, i[2])
		} else if s.tissue.grandtotal != i[0] {
			t.Errorf("%d: Actual grand total %d does not equal %d.", idx, s.tissue.grandtotal, i[0])
		} else if s.tissue.total != i[0] {
			t.Errorf("%d: Actual total %d does not equal %d.", idx, s.tissue.total, i[0])
		}
	}
}

func addRow(s *record, age float64, sex, mal, loc, service, aid string) {
	// Adds values to struct
	s.grandtotal++
	s.allcancer++
	s.sources.Add(aid)
	if mal == "1" {
		s.maltotal++
	} else {
		s.bentotal++
	}
	if service != "MSU" {
		s.total++
		s.cancer++
		s.age += age
		s.cancerage += age
		if sex == "male" {
			s.male++
			s.malecancer++
		} else {
			s.female++
			s.femalecancer++
		}
		if mal == "1" {
			s.malignant++
		} else {
			s.benign++
		}
	}
}

func compareRecords(t *testing.T, a, e *record) {
	// Campares values in structs
	if a.age != e.age {
		t.Errorf("Actual age %f does not equal %f.", a.age, e.age)
	} else if a.allcancer != e.allcancer {
		t.Errorf("Actual allcancer %d does not equal %d.", a.allcancer, e.allcancer)
	} else if a.benign != e.benign {
		t.Errorf("Actual benign %d does not equal %d.", a.benign, e.benign)
	} else if a.bentotal != e.bentotal {
		t.Errorf("Actual bentotal %d does not equal %d.", a.bentotal, e.bentotal)
	} else if a.cancer != e.cancer {
		t.Errorf("Actual cancer %d does not equal %d.", a.cancer, e.cancer)
	} else if a.cancerage != e.cancerage {
		t.Errorf("Actual cancerage %f does not equal %f.", a.cancerage, e.cancerage)
	} else if a.female != e.female {
		t.Errorf("Actual female %d does not equal %d.", a.female, e.female)
	} else if a.femalecancer != e.femalecancer {
		t.Errorf("Actual femalecancer %d does not equal %d.", a.femalecancer, e.femalecancer)
	} else if a.grandtotal != e.grandtotal {
		t.Errorf("Actual bentotal %d does not equal %d.", a.grandtotal, e.grandtotal)
	} else if a.male != e.male {
		t.Errorf("Actual male %d does not equal %d.", a.male, e.male)
	} else if a.malecancer != e.malecancer {
		t.Errorf("Actual malecancer %d does not equal %d.", a.malecancer, e.malecancer)
	} else if a.malignant != e.malignant {
		t.Errorf("Actual malignant %d does not equal %d.", a.malignant, e.malignant)
	} else if a.maltotal != e.maltotal {
		t.Errorf("Actual maltotal %d does not equal %d.", a.maltotal, e.maltotal)
	} else if a.necropsy != e.necropsy {
		t.Errorf("Actual necropsy %d does not equal %d.", a.necropsy, e.necropsy)
	} else if a.sources.Length() != e.sources.Length() {
		t.Errorf("Actual sources %d does not equal %d.", a.sources.Length(), e.sources.Length())
	} else if a.total != e.total {
		t.Errorf("Actual total %d does not equal %d.", a.total, e.total)
	}
}

func TestAddMeasures(t *testing.T) {
	// Tests addNonCancer and addCancer
	sp := getSpecies()[0]
	input := []struct {
		age                         float64
		sex, mal, loc, service, aid string
	}{
		{10.0, "male", "1", "liver", "NWZP", "13"},
		{50.0, "female", "0", "kidney", "ZEPS", "2"},
		{5.0, "male", "0", "liver", "MSU", "2"},
	}
	for _, i := range input {
		s := sp.total.Copy()
		l := sp.tissue.Copy()
		addRow(s, i.age, i.sex, i.mal, i.loc, i.service, i.aid)
		if i.loc == "liver" {
			addRow(l, i.age, i.sex, i.mal, i.loc, i.service, i.aid)
		}
		sp.addNonCancer(i.age, i.sex, i.service, i.aid)
		sp.addCancer(i.age, i.sex, i.mal, i.loc, i.service, i.aid)
		compareRecords(t, sp.total, s)
		compareRecords(t, sp.tissue, l)
	}
}

//----------------------------------------------------------------------------

func locationSlice() [][]string {
	// Return slice of expected location values
	var ret [][]string
	ret = append(ret, []string{"100", "5", "5", "0.10", "2", "0.04", "0.40", "3", "0.06", "0.60", "10.00", "10.00", "2", "3", "2", "3", "1", "0"})
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

func getTestSpecies() []*species {
	// Returns slice of test structs
	var ret []*species
	taxa := canidTaxa()
	loc := locationRecords()
	for idx, i := range testRecords() {
		s := setSpecies(taxa[idx], i)
		if idx < 2 {
			s.tissue = loc[idx]
		} else {
			s.location = ""
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

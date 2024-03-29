// Tests species struct

package cancerrates

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"strconv"
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

func setSpecies(taxa []string, r *Record) *Species {
	// Returns test struct
	s := newSpecies(taxa[0], "liver", taxa[1:])
	//s.location = "liver"
	s.total = r
	return s
}

func getSpecies() []*Species {
	// Returns species structs for testing
	var ret []*Species
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
		s.total.addTotal(i[0])
		if s.total.grandtotal != i[1] {
			t.Errorf("%d: Actual grand total %d does not equal %d.", idx, s.total.grandtotal, i[1])
		} else if s.total.total != i[2] {
			t.Errorf("%d: Actual total %d does not equal %d.", idx, s.total.total, i[2])
			/*} else if s.tissue.grandtotal != i[0] {
				t.Errorf("%d: Actual grand total %d does not equal %d.", idx, s.tissue.grandtotal, i[0])
			} else if s.tissue.total != i[0] {
				t.Errorf("%d: Actual total %d does not equal %d.", idx, s.tissue.total, i[0])*/
		}
	}
}

func TestHighestMalignancy(t *testing.T) {
	// Tests highest malignancy determination
	sp := getSpecies()[0]
	input := [][]string{
		{"-1", "-1"},
		{"0;-1", "0"},
		{"0;1;-1", "1"},
		{"-1;0;1;0", "1"},
	}
	for _, i := range input {
		act := sp.highestMalignancy(i[0])
		if act != i[1] {
			t.Errorf("Actual highest malignancy %s does not equal %s.", act, i[1])
		}
	}
}

func TestCheckLocation(t *testing.T) {
	// Tests location comparison
	sp := getSpecies()[0]
	input := []struct {
		mal string
		loc string
		exp bool
		m   string
	}{
		{"0", "liver", true, "0"},
		{"0;-1", "ovary;mammary", false, "0"},
		{"0;1;-1", "testis;liver;kidney", true, "1"},
		{"1", "livr", false, "1"},
		{"-1", "oral", false, "-1"},
	}
	for _, i := range input {
		act, m := sp.checkLocation(i.mal, i.loc)
		if act != i.exp {
			t.Errorf("Actual result %v for %s does not equal %v.", act, i.loc, i.exp)
		} else if m != i.m {
			t.Errorf("Actual malignant value %s does not equal %s.", m, i.m)
		}
	}
}

func addRow(s *Record, age, sex, nec, mal, loc, service, aid string) {
	// Adds values to struct
	a, _ := strconv.ParseFloat(age, 64)
	s.grandtotal++
	s.allcancer++
	s.sources.Add(aid)
	if mal == "1" {
		s.maltotal++
	} else if mal == "0" {
		s.bentotal++
	}
	if service != "MSU" {
		s.total++
		s.cancer++
		s.age += a
		s.agetotal++
		s.cancerage += a
		s.catotal++
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
		if nec == "1" {
			s.necropsy++
		}
	}
}

func compareRecords(t *testing.T, a, e *Record) {
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
		age, sex, nec, mal, loc, service, aid string
	}{
		{"10.0", "male", "0", "1", "liver", "NWZP", "13"},
		{"50.0", "female", "1", "0", "kidney", "ZEPS", "2"},
		{"5.0", "male", "0", "0", "liver", "MSU", "2"},
	}
	for _, i := range input {
		var allrecords bool
		if i.service != "MSU" {
			allrecords = true
		}
		s := sp.total.Copy()
		l := sp.tissue.Copy()
		addRow(s, i.age, i.sex, i.nec, i.mal, i.loc, i.service, i.aid)
		if i.loc == "liver" {
			addRow(l, i.age, i.sex, i.nec, i.mal, i.loc, i.service, i.aid)
		}
		sp.addNonCancer(allrecords, i.age, i.sex, i.nec, i.service, i.aid, "1")
		sp.addCancer(allrecords, i.age, i.sex, i.nec, i.mal, i.loc, i.service, i.aid)
		compareRecords(t, sp.total, s)
		compareRecords(t, sp.tissue, l)
	}
}

//----------------------------------------------------------------------------

func locationSlice() [][]string {
	// Return slice of expected location values
	var ret [][]string
	//"RecordsWithDenominators", "NeoplasiaDenominator", "NeoplasiaWithDenominators", "NeoplasiaPrevalence"
	ret = append(ret, []string{"5", "50", "5", "0.1", "-", "10", "2", "0.04", "0.4", "3", "0.06", "0.6", "-", "10.00", "10.00", "-", "2", "2", "0", "3", "3", "0", "-", "100", "10", "1", "0", ""})
	ret = append(ret, []string{"5", "100", "5", "0.05", "-", "10", "2", "0.02", "0.4", "3", "0.03", "0.6", "-", "10.00", "10.00", "-", "2", "2", "0", "3", "3", "0", "-", "110", "10", "1", "0", ""})
	return ret
}

func locationRecords() []*Record {
	// Returns slice of records for testing
	var ret []*Record
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
		r := taxa[idx]
		if idx < 2 {
			r = append(r, "all")
			sp = append(sp, append(r, i...))
			// Append location slice
			r = append([]string{taxa[idx][0], "", "", "", "", "", "", ""})
			r = append(r, "liver")
			r = append(r, loc[idx]...)
			sp = append(sp, r)
		} else {
			// Skip location, denominator, and notissue column
			sp = append(sp, append(r, i[1:len(i)-1]...))
		}
		ret = append(ret, sp)
	}
	return ret
}

func getTestSpecies() []*Species {
	// Returns slice of test structs
	var ret []*Species
	taxa := canidTaxa()
	loc := locationRecords()
	for idx, i := range testRecords() {
		s := setSpecies(taxa[idx], i)
		if idx < 2 {
			s.tissue = loc[idx]
		} else {
			s.Location = ""
		}
		ret = append(ret, s)
	}
	return ret
}

func TestToSlice(t *testing.T) {
	head := codbutils.CancerRateHeader(true, false, true, true, true)
	expected := getExpectedSpecies()
	for ind, s := range getTestSpecies() {
		s.denominator = s.total.total
		act := s.ToSlice(false, true, true)
		if len(act) == 3 {
			act = act[:2]
		}
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

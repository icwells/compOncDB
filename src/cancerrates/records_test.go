// Performs white box tests on various methods in the compOncDB package

package cancerrates

import (
	//"github.com/icwells/compOncDB/src/codbutils"
	"testing"
)

func TestAvgAge(t *testing.T) {
	// Tests avgAge method (in speciesTotals script)
	ages := []struct {
		num      float64
		den      int
		expected float64
	}{
		{-1.1, 15, -1.0},
		{12.8, 0, -1.0},
		{12.0, 4, 3.0},
		{6.0, 8, 0.75},
	}
	for _, i := range ages {
		actual := avgAge(i.num, i.den)
		if actual != i.expected {
			t.Errorf("Actual age %f does not equal expected: %f", actual, i.expected)
		}
	}
}

/*func canidTaxa() ([]string, []string) {
	// Returns taxonomies for records
	canis := []string{"Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis"}
	vulpes := []string{"Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Vulpes"}
	return canis, vulpes
}

func setRecord(taxonomy []string, v []float64) *Record {
	// Initilaizes new testing record
	r := NewRecord()
	r.setTaxonomy(taxonomy)
	r.Total = int(v[0])
	r.Age = v[1]
	r.Male = int(v[2])
	r.Female = int(v[3])
	r.Cancer = int(v[4])
	r.Cancerage = v[5]
	r.Malecancer = int(v[6])
	r.Femalecancer = int(v[7])
	r.Malignant = int(v[8])
	r.Benign = int(v[9])
	r.Necropsy = int(v[10])
	r.grandtotal = r.Total
	r.allcancer = r.Cancer
	r.maltotal = r.Malignant
	r.bentotal = r.Benign
	return r
}

func testRecords() []*Record {
	// Returns slice of records for testing
	canis, vulpes := canidTaxa()
	var ret []*Record
	ret = append(ret, setRecord(append(canis, "Canis lupus"), []float64{100, 1000.0, 50, 50, 25, 250.0, 15, 10, 5, 10, 20}))
	ret = append(ret, setRecord(append(canis, "Canis latrans"), []float64{110, 900.0, 50, 70, 30, 300.0, 12, 18, 3, 5, 5}))
	ret = append(ret, setRecord(append(vulpes, "Vulpes vulpes"), []float64{50, 600.0, 25, 35, 0, 0.0, 50, 0, 0, 0, 0}))
	return ret
}

func TestAdd(t *testing.T) {
	r := NewRecord()
	var exp []*Record
	exp = append(exp, setRecord([]string{""}, []float64{100, 1000.0, 50, 50, 25, 250.0, 15, 10, 5, 10, 20}))
	exp = append(exp, setRecord([]string{""}, []float64{210, 1900.0, 100, 120, 55, 550.0, 27, 28, 8, 15, 25}))
	exp = append(exp, setRecord([]string{""}, []float64{260, 2500.0, 125, 155, 55, 550.0, 77, 28, 8, 15, 25}))
	for idx, i := range testRecords() {
		r.Add(i)
		r.grandtotal = r.Total
		if r.Total != exp[idx].Total {
			t.Errorf("%d: Total %d does not equal expected: %d", idx, r.Total, exp[idx].Total)
		} else if r.Age != exp[idx].Age {
			t.Errorf("%d: Age %f does not equal expected: %f", idx, r.Age, exp[idx].Age)
		} else if r.Male != exp[idx].Male {
			t.Errorf("%d: Male %d does not equal expected: %d", idx, r.Male, exp[idx].Male)
		} else if r.Female != exp[idx].Female {
			t.Errorf("%d: Female %d does not equal expected: %d", idx, r.Female, exp[idx].Female)
		} else if r.Cancer != exp[idx].Cancer {
			t.Errorf("%d: Cancer %d does not equal expected: %d", idx, r.Cancer, exp[idx].Cancer)
		} else if r.Cancerage != exp[idx].Cancerage {
			t.Errorf("%d: Cancerage %f does not equal expected: %f", idx, r.Cancerage, exp[idx].Cancerage)
		} else if r.Malecancer != exp[idx].Malecancer {
			t.Errorf("%d: Malecancer %d does not equal expected: %d", idx, r.Malecancer, exp[idx].Malecancer)
		} else if r.Femalecancer != exp[idx].Femalecancer {
			t.Errorf("%d: Femalecancer %d does not equal expected: %d", idx, r.Femalecancer, exp[idx].Femalecancer)
		} else if r.Malignant != exp[idx].Malignant {
			t.Errorf("%d: Malignant %d does not equal expected: %d", idx, r.Malignant, exp[idx].Malignant)
		} else if r.Benign != exp[idx].Benign {
			t.Errorf("%d: Benign %d does not equal expected: %d", idx, r.Benign, exp[idx].Benign)
		} else if r.Necropsy != exp[idx].Necropsy {
			t.Errorf("%d: Necropsy %d does not equal expected: %d", idx, r.Necropsy, exp[idx].Necropsy)
		} else if r.grandtotal != exp[idx].grandtotal {
			t.Errorf("%d: grandtotal %d does not equal expected: %d", idx, r.grandtotal, exp[idx].grandtotal)
		} else if r.allcancer != exp[idx].allcancer {
			t.Errorf("%d: allcancer %d does not equal expected: %d", idx, r.allcancer, exp[idx].allcancer)
		} else if r.maltotal != exp[idx].maltotal {
			t.Errorf("%d: maltotal %d does not equal expected: %d", idx, r.maltotal, exp[idx].maltotal)
		} else if r.bentotal != exp[idx].bentotal {
			t.Errorf("%d: bentotal %d does not equal expected: %d", idx, r.bentotal, exp[idx].bentotal)
		}
	}
}

func getExpectedRecords() [][]string {
	// Return slice of expected values
	var expected [][]string
	canis, vulpes := canidTaxa()
	wolf := append(canis, []string{"Canis lupus", "100", "25", "0.25", "5", "0.05", "0.20", "10", "0.10", "0.40", "10.00", "10.00", "50", "50", "15", "10", "20", "0"}...)
	coyote := append(canis, []string{"Canis latrans", "110", "30", "0.27", "3", "0.03", "0.10", "5", "0.05", "0.17", "8.18", "10.00", "50", "70", "12", "18", "5", "0"}...)
	fox := append(vulpes, []string{"Vulpes vulpes", "50", "0", "0.00", "0", "0.00", "0.00", "0", "0.00", "0.00", "12.00", "NA", "25", "35", "50", "0", "0", "0"}...)
	expected = append(expected, wolf)
	expected = append(expected, coyote)
	return append(expected, fox)
}

func TestCalculateRates(t *testing.T) {
	// Tests calculateRates method
	head := codbutils.CancerRateHeader("")
	expected := getExpectedRecords()
	for ind, r := range testRecords() {
		actual := r.CalculateRates("", "", false)
		for idx, i := range actual {
			if i != expected[ind][idx] {
				t.Errorf("%d: Actual calculated rate %s %s does not equal expected: %s", ind, head[idx+1], i, expected[ind][idx])
			}
		}
	}
}*/

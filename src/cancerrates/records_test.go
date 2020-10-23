// Performs white box tests on various methods in the compOncDB package

package cancerrates

import (
	"github.com/icwells/compOncDB/src/codbutils"
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
}*/

func setRecord(v []float64) *record {
	// Initilaizes new testing record
	r := newRecord()
	r.grandtotal = int(v[0])
	r.total = int(v[1])
	r.age = v[2]
	r.male = int(v[3])
	r.female = int(v[4])
	r.cancer = int(v[5])
	r.cancerage = v[6]
	r.malecancer = int(v[7])
	r.femalecancer = int(v[8])
	r.malignant = int(v[9])
	r.benign = int(v[10])
	r.necropsy = int(v[11])
	r.allcancer = int(v[12])
	r.maltotal = int(v[13])
	r.bentotal = int(v[14])
	return r
}

func testRecords() []*record {
	// Returns slice of records for testing
	var ret []*record
	ret = append(ret, setRecord([]float64{100, 50, 1000.0, 50, 50, 25, 250.0, 15, 10, 5, 10, 20, 30, 6, 12}))
	ret = append(ret, setRecord([]float64{110, 100, 900.0, 50, 70, 30, 300.0, 12, 18, 3, 5, 5, 35, 5, 8}))
	ret = append(ret, setRecord([]float64{50, 50, 600.0, 25, 35, 0, 0.0, 50, 0, 0, 0, 0, 0, 0, 0}))
	return ret
}

func getExpectedRecords() [][]string {
	// Return slice of expected values
	var expected [][]string
	expected = append(expected, []string{"100", "50", "25", "0.25", "5", "0.05", "0.20", "10", "0.10", "0.40", "20.00", "10.00", "50", "50", "15", "10", "20", "0"})
	expected = append(expected, []string{"110", "100", "30", "0.27", "3", "0.03", "0.14", "5", "0.05", "0.23", "9.00", "10.00", "50", "70", "12", "18", "5", "0"})
	expected = append(expected, []string{"50", "50", "0", "0.00", "0", "0.00", "0.00", "0", "0.00", "0.00", "12.00", "NA", "25", "35", "50", "0", "0", "0"})
	return expected
}

func TestCalculateRates(t *testing.T) {
	// Tests calculateRates method
	h := codbutils.NewHeaders()
	head := h.Rates[1:]
	expected := getExpectedRecords()
	for ind, r := range testRecords() {
		actual := r.calculateRates()
		if len(actual) != len(expected[ind]) {
			t.Errorf("%d: Actual length %d does not equal expected: %d", ind, len(actual), len(expected[ind]))
			break
		}
		for idx, i := range actual {
			if i != expected[ind][idx] {
				t.Errorf("%d: Actual calculated rate %s %s does not equal expected: %s", ind, head[idx], i, expected[ind][idx])
			}
		}
	}
}

/*func TestAdd(t *testing.T) {
	r := newRecord()
	var exp []*record
	exp = append(exp, setRecord([]float64{100, 1000.0, 50, 50, 25, 250.0, 15, 10, 5, 10, 20}))
	exp = append(exp, setRecord([]float64{210, 1900.0, 100, 120, 55, 550.0, 27, 28, 8, 15, 25}))
	exp = append(exp, setRecord([]float64{260, 2500.0, 125, 155, 55, 550.0, 77, 28, 8, 15, 25}))
	for idx, i := range testRecords() {
		r.Add(i)
		r.grandtotal = r.total
		if r.total != exp[idx].total {
			t.Errorf("%d: Total %d does not equal expected: %d", idx, r.total, exp[idx].total)
		} else if r.age != exp[idx].age {
			t.Errorf("%d: Age %f does not equal expected: %f", idx, r.age, exp[idx].age)
		} else if r.male != exp[idx].male {
			t.Errorf("%d: Male %d does not equal expected: %d", idx, r.male, exp[idx].male)
		} else if r.female != exp[idx].female {
			t.Errorf("%d: Female %d does not equal expected: %d", idx, r.female, exp[idx].female)
		} else if r.cancer != exp[idx].cancer {
			t.Errorf("%d: Cancer %d does not equal expected: %d", idx, r.cancer, exp[idx].cancer)
		} else if r.cancerage != exp[idx].cancerage {
			t.Errorf("%d: Cancerage %f does not equal expected: %f", idx, r.cancerage, exp[idx].cancerage)
		} else if r.malecancer != exp[idx].malecancer {
			t.Errorf("%d: Malecancer %d does not equal expected: %d", idx, r.malecancer, exp[idx].malecancer)
		} else if r.femalecancer != exp[idx].femalecancer {
			t.Errorf("%d: Femalecancer %d does not equal expected: %d", idx, r.femalecancer, exp[idx].femalecancer)
		} else if r.malignant != exp[idx].malignant {
			t.Errorf("%d: Malignant %d does not equal expected: %d", idx, r.malignant, exp[idx].malignant)
		} else if r.benign != exp[idx].benign {
			t.Errorf("%d: Benign %d does not equal expected: %d", idx, r.benign, exp[idx].benign)
		} else if r.necropsy != exp[idx].necropsy {
			t.Errorf("%d: Necropsy %d does not equal expected: %d", idx, r.necropsy, exp[idx].necropsy)
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
}*/

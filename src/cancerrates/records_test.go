// Performs white box tests on various methods in the compOncDB package

package cancerrates

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"testing"
)

/*func TestAvgAge(t *testing.T) {
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
}*/

func setRecord(v []float64) *Record {
	// Initilaizes new testing record
	r := newRecord()
	r.grandtotal = int(v[0])
	r.total = int(v[1])
	r.age = v[2]
	r.agetotal = r.total
	r.male = int(v[3])
	r.female = int(v[4])
	r.cancer = int(v[5])
	r.cancerage = v[6]
	r.catotal = r.cancer
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

func testRecords() []*Record {
	// Returns slice of records for testing
	var ret []*Record
	ret = append(ret, setRecord([]float64{100, 50, 1000.0, 50, 50, 25, 250.0, 15, 10, 5, 10, 20, 30, 6, 12}))
	ret = append(ret, setRecord([]float64{110, 100, 900.0, 50, 70, 30, 300.0, 12, 18, 3, 5, 5, 35, 5, 8}))
	ret = append(ret, setRecord([]float64{50, 50, 600.0, 25, 35, 0, 0.0, 50, 0, 0, 0, 0, 0, 0, 0}))
	return ret
}

func getExpectedRecords() [][]string {
	// Return slice of expected values
	var expected [][]string
	//"RecordsWithDenominators", "NeoplasiaDenominator", "NeoplasiaWithDenominators", "NeoplasiaPrevalence"
	expected = append(expected, []string{"50", "50", "25", "0.5", "-", "18", "5", "0.1", "0.3333333333333333", "10", "0.2", "0.6666666666666666", "-", "20.00", "10.00", "50", "15", "0", "50", "10", "0", "-", "100", "30", "20", "0", "0"})
	expected = append(expected, []string{"100", "100", "30", "0.3", "-", "13", "3", "0.03", "0.38461538461538464", "5", "0.05", "0.6153846153846154", "-", "9.00", "10.00", "50", "12", "0", "70", "18", "0", "-", "110", "35", "5", "0", "0"})
	expected = append(expected, []string{"50", "50", "0", "0", "-", "0", "0", "0", "NA", "0", "0", "NA", "-", "12.00", "NA", "25", "50", "0", "35", "0", "0", "-", "50", "0", "0", "0", "0"})
	return expected
}

func TestCalculateRates(t *testing.T) {
	// Tests calculateRates method
	h := codbutils.NewHeaders()
	head := append(h.Neoplasia, h.Malignancy...)
	head = append(head, h.AgeSex...)
	head = append(head, h.RatesTail...)
	expected := getExpectedRecords()
	for ind, r := range testRecords() {
		actual := r.calculateRates(r.total, 0)
		if len(actual) != len(expected[ind]) {
			t.Errorf("%d: Actual length %d does not equal expected: %d", ind, len(actual), len(expected[ind]))
			break
		}
		for idx, i := range actual {
			if i != expected[ind][idx] {
				t.Errorf("%d: Actual calculated rate %s %s does not equal expected: %s", ind, head[idx], i, expected[ind][idx])
				break
			}
		}
	}
}

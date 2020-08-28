// Performs white box tests on various methods in the compOncDB package

package dbextract

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

func canidTaxa() ([]string, []string) {
	// Returns taxonomies for records
	canis := []string{"Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis"}
	vulpes := []string{"Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Vulpes"}
	return canis, vulpes
}

func testRecords() []Record {
	// Returns slice of records for testing
	canis, vulpes := canidTaxa()
	return []Record{
		{append(canis, "Canis lupus"), 100, 1000.0, 50, 50, 25, 250.0, 15, 10, 5, 20, nil},
		{append(canis, "Canis latrans"), 110, 900.0, 50, 70, 30, 300.0, 12, 18, 3, 5, nil},
		{append(vulpes, "Vulpes vulpes"), 50, 600.0, 25, 35, 0, 0.0, 50, 0, 0, 0, nil},
	}
}

func TestAdd(t *testing.T) {
	r := NewRecord()
	exp := []*Record{
		{[]string{""}, 100, 1000.0, 50, 50, 25, 250.0, 15, 10, 5, 20, nil},
		{[]string{""}, 210, 1900.0, 100, 120, 55, 550.0, 27, 28, 8, 25, nil},
		{[]string{""}, 260, 2500.0, 125, 155, 55, 550.0, 77, 28, 8, 25, nil},
	}
	for idx, v := range testRecords() {
		i := &v
		r.Add(i)
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
		} else if r.Necropsy != exp[idx].Necropsy {
			t.Errorf("%d: Necropsy %d does not equal expected: %d", idx, r.Necropsy, exp[idx].Necropsy)
		}
	}
}

func getExpectedRecords() [][]string {
	// Return slice of expected values
	var expected [][]string
	canis, vulpes := canidTaxa()
	wolf := append(canis, []string{"Canis lupus", "100", "25", "0.25", "5", "0.05", "0.20", "10.00", "10.00", "50", "50", "15", "10", "20"}...)
	coyote := append(canis, []string{"Canis latrans", "110", "30", "0.27", "3", "0.03", "0.10", "8.18", "10.00", "50", "70", "12", "18", "5"}...)
	fox := append(vulpes, []string{"Vulpes vulpes", "50", "0", "0.00", "0", "0.00", "0.00", "12.00", "NA", "25", "35", "50", "0", "0"}...)
	expected = append(expected, wolf)
	expected = append(expected, coyote)
	return append(expected, fox)
}

func TestCalculateRates(t *testing.T) {
	// Tests calculateRates method
	head := codbutils.CancerRateHeader("taxa_id", "")
	expected := getExpectedRecords()
	//rec := testRecords()
	for ind, r := range testRecords() {
		actual := r.CalculateRates([]string{}, false)
		for idx, i := range actual {
			if i != expected[ind][idx] {
				t.Errorf("%d: Actual calculated rate %s %s does not equal expected: %s", ind, head[idx+1], i, expected[ind][idx])
			}
		}
	}
}

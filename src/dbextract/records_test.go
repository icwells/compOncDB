// Performs white box tests on various methods in the compOncDB package

package dbextract

import (
	//"strconv"
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

func testRecords() []Record {
	// Returns slice of records for testing
	canis, vulpes := canidTaxa()
	return []Record{
		{append(canis, "Canis lupus"), 100, 1000.0, 50, 50, 25, 250.0, 15, 10, nil},
		{append(canis, "Canis latrans"), 110, 900.0, 50, 70, 30, 300.0, 12, 18, nil},
		{append(vulpes, "Vulpes vulpes"), 50, 600.0, 25, 35, 0, 0.0, 50, 0, nil},
	}
}

func canidTaxa() ([]string, []string) {
	// Returns taxonomies for records
	canis := []string{"Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis"}
	vulpes := []string{"Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Vulpes"}
	return canis, vulpes
}

func getExpectedRecords() [][]string {
	// Return slice of expected values
	var expected [][]string
	canis, vulpes := canidTaxa()
	wolf := append(canis, []string{"Canis lupus", "100", "25", "0.25", "10.00", "10.00", "50", "50", "15", "10"}...)
	coyote := append(canis, []string{"Canis latrans", "110", "30", "0.27", "8.18", "10.00", "50", "70", "12", "18"}...)
	fox := append(vulpes, []string{"Vulpes vulpes", "50", "0", "0.00", "12.00", "NA", "25", "35", "50", "0"}...)
	expected = append(expected, wolf)
	expected = append(expected, coyote)
	return append(expected, fox)
}

func TestCalculateRates(t *testing.T) {
	// Tests calculateRates method
	expected := getExpectedRecords()
	rec := testRecords()
	for ind, r := range rec {
		actual := r.CalculateRates("", false)
		for idx, i := range actual {
			if i != expected[ind][idx] {
				t.Error(idx)
				t.Errorf("Actual calculated rate %s does not equal expected: %s", i, expected[ind][idx])
			}
		}
	}
}

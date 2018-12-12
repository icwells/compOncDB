// Tests cancerRates functions

package dbextract

import(
	"testing"
)

func TestCalculateRates(t *testing.T) {
	// Tests calculateRates method
	rec := testRecords()
	expected := [][]string{
		{"Canis lupus", "100", "25", "0.25", "1000.00", "250.00", "50", "50", "15", "10"},
		{"Canis latrans", "110", "30", "0.27", "900.00", "300.00", "50", "70", "12", "18"},
		{"Vulpes vulpes", "50", "0", "0.00", "600.00", "0.00", "25", "35", "0", "0"},
	}
	for ind, r := range rec {
		actual := r.calculateRates()
		for idx, i := range actual {
			if i != expected[ind][idx] {
				t.Errorf("Actual calculated rate %s does not equal expected: %s", i, expected[ind][idx])
			}
		}
	}
}

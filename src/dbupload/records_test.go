// Performs white box tests on various methods in the compOncDB package

package dbupload

import (
	"strconv"
	"testing"
)

func TestAvgAge(t *testing.T) {
	// Tests avgAge method (in speciesTotals script)
	ages := []struct {
		num      float64
		den      int
		expected string
	}{
		{-1.1, 15, "-1"},
		{12.8, 0, "-1"},
		{12.0, 4, "3"},
		{6.0, 8, "0.75"},
	}
	for _, i := range ages {
		actual := avgAge(i.num, i.den)
		if actual != i.expected {
			t.Errorf("Actual age %s does not equal expected: %s", actual, i.expected)
		}
	}
}

func testRecords() []Record {
	// Returns slice of records for testing
	return []Record{
		{"Canis lupus", 6.2, 105, 1000.0, 50, 50, 25, 250.0, 100, 15, 10},
		{"Canis latrans", 5.8, 120, 900.0, 50, 70, 30, 300.0, 110, 12, 18},
		{"Vulpes vulpes", 5.0, 60, 600.0, 25, 35, 0, 0.0, 50, 0, 0},
	}
}

func TestToSlice(t *testing.T) {
	// Tests toSlice method (in speciesTotals script)
	rec := testRecords()
	expected := [][]string{
		{"1", "105", "10", "100", "50", "50", "25", "10", "15", "10"},
		{"2", "120", "8.181818181818182", "110", "50", "70", "30", "10", "12", "18"},
		{"3", "60", "12", "50", "25", "35", "0", "-1", "0", "0"},
	}
	for ind, r := range rec {
		id := strconv.Itoa(ind + 1)
		actual := r.ToSlice(id)
		for idx, i := range actual {
			if i != expected[ind][idx] {
				t.Errorf("Actual slice value %s does not equal expected: %s", i, expected[ind][idx])
			}
		}
	}
}

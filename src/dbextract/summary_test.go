// Performs white box tests on summary functions

package dbextract

import (
	"strconv"
	"testing"
)

func TestGetRow(t *testing.T) {
	// Tests getRows output
	expected := []struct {
		name    string
		num     int
		den     int
		percent string
		length  int
	}{
		{"total", 100, 0, "", 2},
		{"cancer", 50, 100, "50.00%", 3},
		{"male", 25, 100, "25.00%", 3},
		{"benign", 10, 200, "5.00%", 3},
	}
	for _, i := range expected {
		actual := getRow(i.name, i.num, i.den)
		if len(actual) != i.length {
			t.Errorf("Actual length %d does not equal expected: %d", len(actual), i.length)
		} else if actual[0] != i.name {
			t.Errorf("Actual name %s does not equal expected: %s", actual[0], i.name)
		} else if actual[1] != strconv.Itoa(i.num) {
			t.Errorf("Actual total %s does not equal expected: %d", actual[1], i.num)
		} else if i.length == 3 && actual[2] != i.percent {
			t.Errorf("Actual percent %s does not equal expected: %s", actual[2], i.percent)
		}
	}
}

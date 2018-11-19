// Performs black box tests on parseRecords output

package parseRecords_test

import (
	"flag"
	"github.com/icwells/go-tools/iotools"
	"strings"
	"testing"
)

var (
	exp = flag.String("expected", "", "Path to expected output.")
	act = flag.String("actual", "", "Path to actual output.")
)

func getInput(file string) map[string][]string {
	// Returns input test file as a map of string slices
	ret := make(map[string][]string)
	f := iotools.OpenFile(file)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		s := strings.Split(string(scanner.Text()), ",")
		ret[s[0]] = s[1:]
	}
	return ret
}

func TestExtractDiagnosis(t *testing.T) {
	// Compares output of parseRecords extract with expected output
	flag.Parse()
	header := []string{"ID", "Age(months)", "Sex", "Castrated", "Location", "Type", "Malignant", "PrimaryTumor", "Metastasis", "Necropsy"}
	expected := getInput(*exp)
	actual := getInput(*act)
	if len(actual) != len(expected) {
		t.Errorf("Actual length %d does not equal expected: %d", len(actual), len(expected))
	}
	for key, line := range actual {
		for idx, i := range line {
			if i != expected[key][idx] {
				t.Errorf("%s: Actual %s value %s does not equal expected: %s", key, header[idx+1], i, expected[key][idx])
			}
		}
	}
}

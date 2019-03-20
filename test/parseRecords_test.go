// Performs black box tests on parseRecords output

package coDB_test

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

func getInput(file string, col int) map[string][]string {
	// Returns input test file as a map of string slices
	ret := make(map[string][]string)
	f := iotools.OpenFile(file)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		s := strings.Split(string(scanner.Text()), ",")
		ret[s[col]] = s
	}
	return ret
}

func TestMergeRecords(t *testing.T) {
	// Compares output of parseRecords merge with expected output
	flag.Parse()
	header := []string{"Sex", "Age", "Castrated", "ID", "Genus", "Species", "Name", "Date", "Comments", "MassPresent", "Hyperplasia",
		"Necropsy", "Metastasis", "TumorType", "Location", "Primary", "Malignant", "Service", "Account", "Submitter"}
	expected := getInput(*exp, 3)
	actual := getInput(*act, 3)
	if len(actual) != len(expected) {
		t.Errorf("Actual length %d does not equal expected: %d", len(actual), len(expected))
	}
	for key, line := range actual {
		for idx, i := range line {
			if i != expected[key][idx] {
				t.Errorf("%s: Actual %s value %s does not equal expected: %s", key, header[idx], i, expected[key][idx])
			}
		}
	}
}

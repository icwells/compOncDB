// Performs black box tests on parseRecords output

package coDB_test

import (
	"flag"
	"github.com/icwells/compOncDB/src/parserecords"
	"github.com/icwells/go-tools/iotools"
	"os"
	"strings"
	"testing"
)

var (
	user = flag.String("user", "", "MySQL username.")
	service = "NWZP"
	taxa = "input/taxonomies.csv"
	lifehistory = 
	infile = "input/testInput.csv"
	parseout = "merged.csv"
	uploadfile = "input/testUpload.csv"
	config = "config.txt"
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

func TestParseRecords(t *testing.T) {
	// Compares output of parseRecords with expected output
	header := []string{"Sex", "Age", "Castrated", "ID", "Genus", "Species", "Name", "Date", "Comments", "MassPresent", "Hyperplasia",
		"Necropsy", "Metastasis", "TumorType", "Location", "Primary", "Malignant", "Service", "Account", "Submitter", "Zoo", "Institute"}
	// Parse test file
	ent := parserecords.NewEntries(service)
	ent.GetTaxonomy(taxa)
	ent.SortRecords(false, infile, parseout)
	// Compare actual output with expected
	expected := getInput(uploadfile, 3)
	actual := getInput(parseout, 3)
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
	os.Remove(parseout)
}

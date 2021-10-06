// Performs black box tests on parseRecords output

package main

import (
	"flag"
	"github.com/icwells/compOncDB/src/parserecords"
	"github.com/icwells/go-tools/dataframe"
	"os"
	"testing"
)

var (
	denominators = "input/testDenominators.csv"
	infile       = "input/testInput.csv"
	lifehistory  = "input/testLifeHistories.csv"
	parseout     = "merged.csv"
	password     = flag.String("password", "", "MySQL password.")
	searchfile   = "input/testSearch.csv"
	service      = "NWZP"
	taxa         = "input/taxonomies.csv"
	updatefile   = "input/testUpdate.csv"
	uploadfile   = "input/testUpload.csv"
	user         = flag.String("user", "", "MySQL username.")
)

func TestParseRecords(t *testing.T) {
	// Compares output of parseRecords with expected output
	ent := parserecords.NewEntries(service, infile)
	ent.GetTaxonomy(taxa)
	ent.SortRecords(false, infile, parseout)
	// Compare actual output with expected
	expected, _ := dataframe.FromFile(uploadfile, 3)
	actual, err := dataframe.FromFile(parseout, 3)
	if err != nil {
		t.Error(err)
	} else {
		ac, ar := actual.Dimensions()
		ec, er := expected.Dimensions()
		if ac != ec && ar != er {
			t.Errorf("Actual dimensions [%d, %d] does not equal expected: [%d, %d]", ac, ar, ec, er)
		}
		for key := range actual.Index {
			for k := range actual.Header {
				a, _ := actual.GetCell(key, k)
				e, _ := expected.GetCell(key, k)
				if a != e {
					t.Errorf("%s: Actual %s value %s does not equal expected: %s", key, k, a, e)
				}
			}
		}
		os.Remove(parseout)
	}
}

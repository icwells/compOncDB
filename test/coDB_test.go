// Performs black box tests on the comparative oncology sql database

package coDB_test

import (
	"flag"
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func sortInput(files []string, expected bool) map[string]string {
	// Returns sorted actual or expected files
	ret := make(map[string]string)
	for _, i := range files {
		base := iotools.GetFileName(i)
		if expected == true {
			// Remove test prefix from map key
			base = strings.Replace(base, "test", "", 1)
		}
		ret[base] = i
	}
	return ret
}

func loadTable(file string) map[string][]string {
	// Returns table as a map of string slices
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

func compareTables(t *testing.T, name, exp, act string) {
	// Compares output of equivalent tables
	expected := loadTable(exp)
	actual := loadTable(act)
	if len(actual) != len(expected) {
		t.Errorf("%s: Actual length %d does not equal expected: %d", name, len(actual), len(expected))
	} else {
		for k, v := range actual {
			for idx, i := range v {
				err := false
				// Attempt to resolve differences in floating point precision
				a, er := strconv.ParseFloat(i, 64)
				if er == nil {
					var e float64
					e, er = strconv.ParseFloat(expected[k][idx], 64)
					if er == nil && a != e {
						err = true
					}
				} else if i != expected[k][idx] {
					err = true
				}
				if err == true && strings.Contains(name, "Tumor") == false {
					t.Errorf("%s %d: Actual value %s does not equal expected: %s", name, idx + 1, i, expected[k][idx])
				}
			}
		}
	}
}

func TestDumpTables(t *testing.T) {
	// Compares actual output from table dumps to expected
	flag.Parse()
	*indir, _ = iotools.FormatPath(*indir, false)
	files, err := filepath.Glob(*indir + "*.csv")
	if err != nil {
		fmt.Printf("\n\t[Error] Cannot find test files in %s: %v", *indir, err)
		os.Exit(10)
	}
	expected := sortInput(files, true)
	actual := sortInput(files, false)
	for k, v := range expected {
		act, ex := actual[k]
		if ex == false {
			t.Errorf("Actual table %s not found.", k)
		} else {
			compareTables(t, k, v, act)
			// Remove test output
			//os.Remove(act)
		}
	}
}

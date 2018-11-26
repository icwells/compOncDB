// Performs black box tests on the comparative oncology sql database

package coDB_test

import (
	"flag"
	"github.com/icwells/go-tools/iotools"
	//"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var (
	indir = flag.String("indir", "", "Path to output directory with test data to compare.")
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

func loadTable(file string) [][]string {
	// Returns table as a map of string slices
	first := true
	var col []int
	var ret [][]string
	f := iotools.OpenFile(file)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		s := strings.Split(string(scanner.Text()), ",")
		if first == false {
			for _, c := range col {
				// Remove randomly assigned id entries
				var head []string
				if c == 1 {
					head = []string{s[0]}
				} else {
					head = s[:c]
				s = append(head, s[c+1:]...)
				}
			}
			ret = append(ret, s)
		} else {
			for idx, i := range s {
				if i == "ID" || strings.Contains(i, "_id") == true {
					col = append(col, idx)
				}
			}
		}
	}
	return ret
}

func compareEntries(actual, expected []string) (bool, int) {
	// Returns true if both slices are equal
	ret := true
	var index int
	for idx, i := range actual {
		if i != expected[idx] {
			ret = false
			// Attempt to resolve differences in floating point precision
			a, err := strconv.ParseFloat(i, 64)
			if err == nil {
				var e float64
				e, err = strconv.ParseFloat(expected[idx], 64)
				if err == nil && a == e {
					ret = true
				}
			}
		}
		if ret == false {
			index = idx
			break
		}
	}
	return ret, index
}

func compareTables(t *testing.T, name, exp, act string) {
	// Compares output of equivalent tables
	expected := loadTable(exp)
	actual := loadTable(act)
	if len(actual) != len(expected) {
		t.Errorf("%s: Actual length %d does not equal expected: %d", name, len(actual), len(expected))
	} else {
		for k, v := range actual {
			equal := false
			var idx int
			for _, val := range expected {
				// Ignore randomly assigned IDs and compare to all entries
				if len(v) == len(val) {
					equal, idx = compareEntries(v, val)
				}
				if equal == true { 
					break
				}
			}
			if equal == false {
				t.Errorf("%s %d: Actual value %s does not equal expected: %s", name, idx, actual[k][idx], expected[k][idx])
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
		t.Errorf("Cannot find test files in %s: %v", *indir, err)
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

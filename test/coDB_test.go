// Performs black box tests on the comparative oncology sql database

package coDB_test

import (
	"flag"
	"github.com/icwells/go-tools/iotools"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var (
	indir  = flag.String("indir", "", "Path to output directory with test data to compare.")
	level  = flag.String("level", "", "Empty field for taxonomic level search.")
	tables = flag.String("tables", "", "Path tableColumns.txt file.")
)

type idtrimmer struct {
	columns	[]int
}

func (t *idtrimmer) setColumns(row []string) {
	// Stores id column indeces
	for idx, i := range row {
		if i == "ID" || strings.Contains(i, "_id") == true {
			t.columns = append(t.columns, idx)
		}
	}
}

func (t *idtrimmer) trimColumns(row []string) []string {
	// Removes randomly generated id numbers from column
	for _, c := range t.columns {
		// Remove randomly assigned id entries
		var head []string
		if c == 1 {
			head = []string{row[0]}
		} else {
			head = row[:c]
		row = append(head, row[c+1:]...)
		}
	}
	return row
}

//----------------------------------------------------------------------------

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
	var ret [][]string
	var trim idtrimmer
	f := iotools.OpenFile(file)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		s := strings.Split(string(scanner.Text()), ",")
		if first == false {
			s = trim.trimColumns(s)
			ret = append(ret, s)
		} else {
			trim.setColumns(s)
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
			os.Remove(act)
		}
	}
}

//----------------------------------------------------------------------------

func TestSearches(t *testing.T) {
	// Tests taxonomy search output
	flag.Parse()
	*indir, _ = iotools.FormatPath(*indir, false)
	files, err := filepath.Glob(*indir + "*.csv")
	if err != nil {
		t.Errorf("Cannot find test files in %s: %v", *indir, err)
	}
	expected := sortInput(files, true)
	actual := sortInput(files, false)
	if iotools.Exists(*indir + "gray_fox.csv") == true {
		t.Error("Empty result saved to file.")
	}
	for k, v := range expected {
		act, ex := actual[k]
		if ex == false {
			t.Errorf("Actual search result %s not found.", k)
		} else {
			compareTables(t, k, v, act)
			// Remove test output
			os.Remove(act)
		}
	}
}

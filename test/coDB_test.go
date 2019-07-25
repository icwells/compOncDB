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

type filepair struct {
	t			*testing.T
	name		string
	columns		map[string]int
	ids			[]int
	expected	map[string][]string
	actual		map[string][]string
	index		int
	value		string
	pass		bool
}

func (f *filepair) setColumns(header []string) {
	// Stores column value:index pairs
	for idx, i := range header {
		f.columns[i] = idx
		i := strings.TrimSpace(i)
		if i == "ID" || strings.Contains(i, "_id") == true {
			f.ids = append(f.ids, idx)
		}
	}
}

func (f *filepair) loadTable(file string) map[string][]string {
	// Returns table as a map of string slices
	first := true
	ret := make(map[string][]string)
	fl := iotools.OpenFile(file)
	defer fl.Close()
	scanner := iotools.GetScanner(fl)
	for scanner.Scan() {
		s := strings.Split(string(scanner.Text()), ",")
		if first == false {
			// Store map entry without id columns
			ret[s[0]] = s
		} else {
			if len(f.columns) == 0 {
				f.setColumns(s)
			}
			first = false
		}
	}
	return ret
}

func newFilePair(t *testing.T, name, e, a string) filepair {
	// Initializes struct and reads in data from input files
	var f filepair
	f.t = t
	f.name = name
	f.index = -1
	f.pass = true
	f.columns = make(map[string]int)
	f.expected = f.loadTable(e)
	f.actual = f.loadTable(a)
	return f
}

//----------------------------------------------------------------------------

func (f *filepair) isID(idx int) bool {
	// Returns true if column idx is in f.ids
	for _, i := range f.ids {
		if i == idx {
			return true
		}
	}
	return false
}

func (f *filepair) compareEntries(actual, expected []string)  {
	// Returns true if both slices are equal
	equal := true
	for idx, i := range actual {
		// Skip randomly assigned IDs
		if idx < len(expected) {
			if f.isID(idx) == false && i != expected[idx] {
				equal = false
				// Attempt to resolve differences in floating point precision
				a, err := strconv.ParseFloat(i, 64)
				if err == nil {
					var e float64
					e, err = strconv.ParseFloat(expected[idx], 64)
					if err == nil && a == e {
						equal = true
					}
				}
			}
			if equal == false {
				f.pass = false
				f.index = idx
				f.value = expected[idx]
				break
			}
		}
	}
}

func (f *filepair) compareRows(k string) {
	// Compares actual[idx] to expected entries
	row := f.actual[k]
	e := f.expected[k]
	if len(row) != len(f.columns) {
		f.t.Errorf("%s %s: Actual line length %d does not equal expected: %d", f.name, k, len(row), len(f.columns))
	} else {
		f.compareEntries(row, e)
	}
	if f.pass == false {
		f.t.Errorf("%s %s-%d: Actual value %s does not equal expected: %s", f.name, k, f.index, row[f.index], f.value)
	}
}

func compareTables(t *testing.T, name, exp, act string) {
	// Compares output of equivalent tables
	f := newFilePair(t, name, exp, act)
	if len(f.actual) != len(f.expected) {
		f.t.Errorf("%s: Actual file length %d does not equal expected: %d", name, len(f.actual), len(f.expected))
	} else {
		for k := range f.actual {
			f.compareRows(k)
			if f.pass == false {
				// Report 1 error per file
				break
			}
		}
	}
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
			//os.Remove(act)
		}
	}
}

func TestUpdates(t *testing.T) {
	// Tests dumped tables after update 
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
			t.Errorf("Actual search result %s not found.", k)
		} else {
			compareTables(t, k, v, act)
			// Remove test output
			os.Remove(act)
		}
	}
}

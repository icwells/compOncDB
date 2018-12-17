// Defines a struct for managing test searches

package main

import (
	"bytes"
	"fmt"
	"dbextract"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"strings"
)

type testcase struct {
	level  string
	common bool
	column string
	op     string
	value  string
	table  string
}

func (c *testcase) String() string {
	// Returns formatted string
	ret := bytes.NewBufferString(c.level)
	if c.common == true {
		ret.WriteString(", true, ")
	} else {
		ret.WriteString(", false, ")
	}
	ret.WriteString(c.op)
	ret.WriteString(", ")
	ret.WriteString(c.value)
	ret.WriteString(", ")
	ret.WriteString(c.table)
	return ret.String()
}

type searchterms struct {
	outdir string
	cases  []*testcase
}

func (s *searchterms) setCase(row []string) {
	// Adds new test case to slice
	set := false
	c := new(testcase)
	for idx, i := range row {
		i = strings.TrimSpace(i)
		if idx == 0 && len(i) >= 1 {
			c.level = i
			set = true
		} else if idx == 1 {
			if strings.ToLower(i) == "true" {
				c.common = true
			}
			set = true
		} else if idx == 2 && len(i) >= 1 {
			if strings.Contains(i, "=") == true || strings.Contains(i, ">") == true || strings.Contains(i, "<") == true {
				c.column, c.op, c.value = getOperation(i)
			} else {
				c.value = i
			}
			set = true
		} else if idx == 3 && len(i) >= 1 {
			c.table = i
			set = true
		}
	}
	if set == true {
		// Skip empty lines
		s.cases = append(s.cases, c)
	}
}

func (s *searchterms) readSearchTerms(infile, outdir string) {
	// Loads test cases from file
	s.outdir = outdir
	first := true
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		if first == false {
			row := strings.Split(string(scanner.Text()), ",")
			s.setCase(row)
		} else {
			first = false
		}
	}
}

func (s *searchterms) searchTestCases(db *dbIO.DBIO) {
	// Assigns each case to an appriate search type
	for _, c := range s.cases {
		var res [][]string
		var header string
		outfile := "nil"
		if s.outdir != "nil" {
			outfile = fmt.Sprintf("%s%s.csv", s.outdir, strings.Replace(c.value, " ", "_", 1))
		}
		if len(c.column) >= 1 {
			if len(c.table) >= 1 {
				// Perform single table search
				res, header = dbextract.SearchSingleTable(db, c.table, *user, c.column, c.op, c.value, false)
			} else {
				// Perform column search
				tables := getTable(db.Columns, c.column)
				res, header = dbextract.SearchColumns(db, tables, *user, c.column, c.op, c.value, false, false)
			}
		} else {
			// Perform taxonomy search
			res, header = dbextract.SearchTaxonomicLevels(db, []string{c.value}, *user, c.level, false, c.common)
		}
		if len(res) >= 1 {
			writeResults(outfile, header, res)
		}
	}
}

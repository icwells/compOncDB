// Defines a struct for managing test searches

package main

import (
	"bytes"
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"strconv"
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
	ret.WriteString(", ")
	if c.common == true {
		ret.WriteString("true")
	} else {
		ret.WriteString("false")
	}
	ret.WriteString(", ")
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

func (s *searchterms) addCase(c *testcase) {
	// Appends new test case
	s.cases = append(s.cases, c)
}

func (s *searchterms) setCase(row []string) {
	// Adds new test case to slice
	c := new(testcase)
	for idx, i := range row {
		i = strings.TrimSpace(i)
		if idx == 0 && len(i) >= 1 {
			c.level = i
		} else if idx == 1 && len(i) >= 1 {
			c.common = true
		} else if idx == 2 && len(i) >= 1 {
			c.column = i
		} else if idx == 3 && len(i) >= 1 {
			c.op = strings.Replace(i, "'", "", -1)
		} else if idx == 4 && len(i) >= 1 {
			c.value = i
		} else if idx == 5 && len(i) >= 1 {
			c.table = i
		}
	}
	s.addCase(c)
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
		fmt.Println(c)
		outfile := fmt.Sprintf("%s%s.csv", s.outdir, strings.Replace(c.value, " ", "_", 1))
		if _, er := strconv.Atoi(c.column); er == nil {
			if len(c.table) >= 2 {
				// Perform single table search
				res, header = SearchSingleTable(db, c.table, c.column, c.op, c.value)
			} else {
				// Perform column search
				tables := getTable(db.Columns, c.column)
				res, header = SearchColumns(db, tables, c.column, c.op, c.value)
			}
		} else {
			// Perform taxonomy search
			*com = c.common
			*level = c.level
			fmt.Println(*com, *level)
			res, header = SearchTaxonomicLevels(db, []string{c.value})
		}
		if len(res) >= 1 {
			writeResults(outfile, header, res)
		}
	}
}

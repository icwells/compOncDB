// Defines a struct for managing test searches

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"strings"
)

type searchterms struct {
	outdir string
	cases  []codbutils.Evaluation
}

func (s *searchterms) setCase(columns map[string]string, row []string) {
	// Adds new test case to slice
	set := false
	var e []codbutils.Evaluation
	for idx, i := range row {
		i = strings.TrimSpace(i)
		if idx == 0 && len(i) >= 1 {
			e = codbutils.SetOperations(columns, i)
			set = true
		}
	}
	if set == true {
		// Skip empty lines
		s.cases = append(s.cases, e[0])
	}
}

func (s *searchterms) readSearchTerms(columns map[string]string, infile, outdir string) {
	// Loads test cases from file
	s.outdir = outdir
	first := true
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		if first == false {
			row := strings.Split(string(scanner.Text()), ",")
			s.setCase(columns, row)
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
			outfile = fmt.Sprintf("%s%s.csv", s.outdir, strings.Replace(c.Value, " ", "_", 1))
		}
		if c.Table == "Life_history" {
			// Perform single table search
			res, header = dbextract.SearchSingleTable(db, c.Table, *user, c.Column, c.Operator, c.Value, false)
		} else {
			// Perform column search
			res, header = dbextract.SearchColumns(db, *user, []codbutils.Evaluation{c}, false, false)
		}
		if len(res) >= 1 {
			codbutils.WriteResults(outfile, header, res)
		}
	}
}

// This script contains methods for searching tumor tables

package main

import (
	"github.com/icwells/dbIO"
	"strings"
)

type searcher struct {
	db			*dbIO.DBIO
	user		string
	tables		[]string
	column		string
	value		string
	operator	string
	short		bool
	common		bool
	res			[][]string
	ids			[]string
	taxaids		[]string
	header		string
	na			[]string
}

func newSearcher(db *dbIO.DBIO, tables []string, column, op, value string) *searcher {
	// Assigns starting values to searcher
	s := new(searcher)
	s.db = db
	s.user = *user
	s.tables = tables
	s.column = column
	s.value = value
	s.operator = op
	s.common = *common
	s.na = []string{"NA", "NA", "NA", "NA", "NA", "NA", "NA"}
	return s
}

func (s *searcher) getIDs() {
	// Gets ids from target table and get patient records
	ids := s.db.GetRows(s.tables[0], s.column, s.value, "ID")
	for _, i := range ids {
		// Convert to single slice
		s.ids = append(s.ids, i[0])
	}
	s.res = s.db.GetRows("Patient", "ID", strings.Join(s.ids, ","), "*")
}

func (s *searcher) setIDs() {
	// Sets IDs from s.res (ID must be in first column)
	for _, i := range s.res {
		s.ids = append(s.ids, i[0])
	}
}

func (s *searcher) setTaxaIDs() {
	// Stores taxa ids from patient results
	for _, i := range s.res {
		s.taxaids = append(s.taxaids, i[4])
	}
}

func (s *searcher) appendSource() {
	// Appends data from source table to res
	m := toMap(s.db.GetRows("Source", "ID", strings.Join(s.ids, ","), "*"))
	for idx, i := range s.res {
		row , ex := m[i[0]]
		if ex == true {
			s.res[idx] = append(i, row...)
		} else {
			s.res[idx] = append(i, s.na[:2]...)
		}
	}
}

func (s *searcher) appendTaxonomy() {
	// Appends raxonomy to s.res
	taxa := s.getTaxonomy(s.taxaids, true)
	for idx, i := range s.res {
		// Apppend taxonomy to records
		taxonomy, ex := taxa[i[4]]
		if ex == true {
			s.res[idx] = append(i, taxonomy...)
		} else {
			s.res[idx] = append(i, s.na...)
		}
	}
}

func (s *searcher) appendDiagnosis() {
	// Appends data from tumor and tumor relation tables
	d := toMap(s.db.GetRows("Diagnosis", "ID", strings.Join(s.ids, ","), "*"))
	t := s.getTumor()
	for idx, i := range s.res {
		// Concatenate tables
		id := i[0]
		diag, ex := d[id]
		if ex == true {
			i = append(i, diag...)
		} else {
			i = append(i, s.na[:4]...)
		}
		tumor, e := t[id]
		if e == true {
			i = append(i, tumor...)
		} else {
			i = append(i, s.na[:5]...)
		}
		s.res[idx] = i
	}
}

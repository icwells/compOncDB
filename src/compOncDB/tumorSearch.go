// This script contains methods for searching tumor tables

package main

import (
	"database/sql"
	"dbIO"
	"strings"
)

func getTumorRecords(ch chan []string, db *sql.DB, id string, tumor map[string][]string) {
	// Returns tumor information for given id
	var loc, typ, mal, prim []string
	rows := dbIO.GetRows(db, "Tumor_relation", "ID", id, "*")
	for _, i := range rows {
		j, ex := tumor[i[1]]
		if ex == true {
			if i[2] == "1" {
				// Prepend primary tumor
				prim = append([]string{i[2]}, prim...)
				mal = append([]string{i[3]}, mal...)
				typ = append([]string{j[0]}, typ...)
				loc = append([]string{j[1]}, loc...)
			} else {
				// Append tumor type and location
				prim = append(prim, i[2])
				mal = append(mal, i[3])
				typ = append(typ, j[0])
				loc = append(loc, j[1])
			}
		}
	}
	diag := []string{strings.Join(prim, ";"), strings.Join(mal, ";"), strings.Join(typ, ";"), strings.Join(loc, ";")}
	ch <- diag
}

func (s *searcher) getTumor() map[string][]string {
	// Returns map of tumor data from patient ids
	ch := make(chan []string)
	// {id: [types], [locations]}
	rec := make(map[string][]string)
	tumor := toMap(dbIO.GetTable(s.db, "Tumor"))
	for _, id := range s.ids {
		// Get records for each patient concurrently
		go getTumorRecords(ch, s.db, id, tumor)
		ret := <-ch
		if len(ret) >= 1 {
			rec[id] = ret
		}
	}
	return rec
}

func (s *searcher) getTumorRelation(ch chan [][]string, row []string) {
	// Returns matching tumor relation entries
	var ret [][]string
	table := dbIO.GetRows(s.db, "tumor_relation", "tumor_id", row[0], "*")
	for _, i := range table {
		res := append(i, row[1:]...)
		ret = append(ret, res)
	}
	ch <- ret
}

func (s * searcher) searchTumor() {
	// Finds matches in tumor tables
	s.header = "ID,tumor_id,primary_tumor,Malignant,Type,Location"
	if s.tables[0] == "tumor_relation" {
		s.searchPairedTables(1)
	} else if s.tables[0] == "Tumor" {
		ch := make(chan [][]string)
		t := dbIO.GetRows(s.db, s.tables[0], s.column, s.value, "*")
		for _, i := range t {
			go s.getTumorRelation(ch, i)
			ret := <-ch
			s.res = append(s.res, ret...)
		}
	}
	if s.short == false {
		s.header = "ID,Sex,Age,Castrated,taxa_id,source_id,Species,Date,Comments,tumor_id,primary_tumor,Malignant,Type,Location"
		s.getPatients()
	}
}

// This script contains functions for searching tables for a given column/value combination

package main

import (
	"fmt"
	"github.com/icwells/dbIO"
	"os"
	"strings"
)

func getTumorRecords(ch chan []string, db *dbIO.DBIO, rows [][]string, tumor map[string][]string) {
	// Returns tumor information for given id
	var loc, typ, mal, prim []string
	for _, i := range rows {
		j, ex := tumor[i[1]]
		if ex == true {
			if i[2] == "1" && len(prim) > 0 {
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

func (s *searcher) tumorMap() map[string][][]string {
	// Converts tumor_relation table to map
	ret := make(map[string][][]string)
	rows := s.db.GetRows("Tumor_relation", "ID", strings.Join(s.ids, ","), "*")
	for _, i := range rows{
		_, ex := ret[i[0]]
		if ex == true {
			ret[i[0]] = append(ret[i[0]], i)
		} else {
			ret[i[0]] = [][]string{i}
		}
	}
	return ret
}

func (s *searcher) getTumor() map[string][]string {
	// Returns map of tumor data from patient ids
	ch := make(chan []string)
	// {id: [types], [locations]}
	rec := make(map[string][]string)
	tumor := toMap(s.db.GetTable("Tumor"))
	tr := s.tumorMap()
	for _, id := range s.ids {
		// Get records for each patient concurrently (may be multiple tumor relation records for an id)
		go getTumorRecords(ch, s.db, tr[id], tumor)
		ret := <-ch
		if len(ret) >= 1 {
			rec[id] = ret
		}
	}
	return rec
}

func (s *searcher) searchTumor() {
	// Gets IDs from tumor ids
	var tumorids []string
	tids := s.db.GetRows(s.tables[0], s.column, s.value, "ID")
	for _, i := range tids {
		// Convert to single slice
		tumorids = append(tumorids, i[0])
	}
	ids := s.db.GetRows("Tumor_relation", "tumor_id", strings.Join(tumorids, ","), "ID")
	for _, i := range ids {
		s.ids = append(s.ids, i[0])
	}
	s.res = s.db.GetRows("Patient", "ID", strings.Join(s.ids, ","), "*")
}
	
func (s *searcher) searchAccounts() {
	// Searches source tables
	if s.user != "root" {
		fmt.Print("\n\t[Error] Must be root to access Accounts table. Exiting.\n\n")
		os.Exit(1010)
	}
	var accounts []string
	target := s.value
	if s.column != "account_id" {
		// Get target account IDs
		aids := s.db.GetRows(s.tables[0], s.column, s.value, "ID")
		for _, i := range aids {
			accounts = append(accounts, i[0])
		}
		target = strings.Join(accounts, ",")
	}
	// Get target patient IDs
	ids := s.db.GetRows("Source", "account_id", target, "ID")
	for _, i := range ids {
		s.ids = append(s.ids, i[0])
	}
	s.res = s.db.GetRows("Patient", "ID", strings.Join(s.ids, ","), "*")
}

func (s *searcher) searchTaxaIDs() {
	// Searches for matches in any table with taxa_ids as primary key
	tids := s.db.EvaluateRows(s.tables[0], s.column, s.operator, s.value, "taxa_id")
	for _, i := range tids {
		s.taxaids = append(s.taxaids, i[0])
	}
	s.res = s.db.GetRows("Patient", "taxa_id", strings.Join(s.taxaids, ","), "*")
	s.setIDs()
}

func (s *searcher) searchPatient() {
	// Searches any match that include the patient table
	s.res = s.db.GetRows(s.tables[0], s.column, s.value, "*")
	s.setIDs()
}

func (s *searcher) assignSearch() {
	// Runs appropriate search based on input
	// Store standardized header
	s.header = "ID,Sex,Age,Castrated,taxa_id,source_id,Species,Date,Comments,"
	s.header = s.header + "Kingdom,Phylum,Class,Order,Family,Genus,Masspresent,Necropsy,Metastasis,"
	s.header = s.header + "primary_tumor,Malignant,Type,Location,service_name,account_id"
	switch s.tables[0] {
		// Start with potential mutliple entries
		case "Patient":
			s.searchPatient()
		case "Source":
			s.getIDs()
		case "Tumor_relation":
			s.getIDs()
		case "Taxonomy":
			s.searchTaxaIDs()
		case "Common":
			s.searchTaxaIDs()
		case "Life_history":
			s.searchTaxaIDs()
		case "Totals":
			s.searchTaxaIDs()
		case "Diagnosis":
			s.getIDs()
		case "Tumor":
			s.searchTumor()
		case "Accounts":
			s.searchAccounts()
	}
	if *count == false {
		// res and ids must be set first
		s.setTaxaIDs()
		s.appendTaxonomy()
		s.appendDiagnosis()
		s.appendSource()
	}
}

func SearchColumns(db *dbIO.DBIO, tables []string, column, op, value string) ([][]string, string) {
	// Determines search procedure
	fmt.Printf("\tSearching for records with %s in column %s...\n", value, column)
	s := newSearcher(db, tables, column, op, value)
	s.assignSearch()
	return s.res, s.header
}

func SearchSingleTable(db *dbIO.DBIO, table, column, op, value string) ([][]string, string) {
	// Returns results from single table
	fmt.Printf("\tSearching table %s for records with %s in column %s...\n", table, value, column)
	s := newSearcher(db, []string{table}, column, op, value)
	s.header = s.db.Columns[table]
	s.res = s.db.GetRows(table, s.column, s.value, "*")
	return s.res, s.header
}

// This script contains functions for searching tables for a given column/value combination

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	"os"
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
	// Gets IDs from tumor ids
	var tumorids []string
	tids = dbIO.GetRows(s.db, s.tables[0], s.column, s.value, "ID")
	for _, i := range tids {
		// Convert to single slice
		tumorids = append(tumorids, i[0])
	}
	ids = dbIO.GetRows(s.db, "Tumor_relation", "tumor_id", strings.Join(tumorids, ","), "ID")
	for _, i := range ids {
		s.ids = append(s.ids, i[0])
	}
	s.res := dbIO.GetRows(s.db, "Patient", "ID", strings.Join(s.ids, ","), "*")
}

func (s *searcher) appendDiagnosis() {
	// Appends data from tumor and tumor relation tables
	d := toMap(dbIO.GetRows(s.db, "Diagnosis", "ID", strings.Join(s.ids, ","), "*"))
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
	
func (s *searcher) searchAccounts() {
	// Searches source tables
	switch s.column {
		case "account_id":
			if s.user == "root" {
				// Return both tables
				s.searchPairedTables(2)
			} else {
				if s.tables[0] == "Source" {
					// Return single table
					s.res = dbIO.GetRows(s.db, s.tables[0], s.column, s.value, "*")
				} else {
					fmt.Println("\n\t[Error] Must be root to access Accounts table. Exiting.\n")
					os.Exit(99)
				}
			}
		default:
			if s.tables[0] == "Accounts" && s.user != "root" {
				fmt.Println("\n\t[Error] Must be root to access Accounts table. Exiting.\n")
				os.Exit(99)
			} else {
				s.res = dbIO.GetRows(s.db, s.tables[0], s.column, s.value, "*")
			}
	}
}

func (s *searcher) searchPatient() {
	// Searches any match that include the patient table
	s.res = dbIO.GetRows(s.db, s.tables[0], s.column, s.value, "*")
	s.setIDs()
}

func (s *searcher) assignSearch() {
	// Runs appropriate search based on input
	// Store standardized header
	s.header = "ID,Sex,Age,Castrated,taxa_id,source_id,Species,Date,Comments,"
	s.header = s.header + "Kingdom,Phylum,Class,Order,Family,Genus,Species,"
	s.header = s.header + "Masspresent,Necropsy,Metastasis,primary_tumor,Malignant,Type,Location,service_name,account_id"
	switch s.tables[0] {
		// Start with potential mutliple entries
		case "Patient":
			s.searchPatient()
		case "Source":
			s.getIDs()
		case "Tumor_relation":
			s.getIDs()
		//case "Taxonomy":
			//s.searchTaxonomy()
		//case "Common":
			//s.searchTaxonomy()
		//case "Life_history":
			//s.searchLifeHistory()
		//case "Totals":
			//s.searchTotals()
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

func searchColumns(db *sql.DB, col map[string]string, tables []string) ([][]string, string) {
	// Determines search procedure
	fmt.Printf("\tSearching for records with %s in column %s...\n", *value, *column)
	s := newSearcher(db, col, tables)
	s.assignSearch()
	return s.res, s.header
}

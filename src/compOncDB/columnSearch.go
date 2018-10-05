// This script contains functions for searching tables for a given column/value combination

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/icwells/go-tools/strarray"
	"os"
	"strings"
)

type searcher struct {
	db		*sql.DB
	user	string
	columns	map[string]string
	tables	[]string
	column	string
	value	string
	short	bool
	scour	bool
	res		[][]string
	header	string
}

func newSearcher(db *sql.DB, col map[string]string, tables []string) searcher {
	// Assigns starting values to searcher
	s := new(searcher)
	s.db = db
	s.user = *user
	s.columns = col
	s.tables = tables
	s.column = *column
	s.value = *value
	s.short = *short
	s.scour = *scour
	return s
}

func (s *searcher) searchPairedTables(c int) {
	// Cancatentes results from paired tables, c indicates id column
	s.res := dbIO.GetRows(s.db, s.tables[0], s.column, s.value, "*")
	// Get entire table in case value is not in both tables
	m := toMap(dbIO.GetTable(s.db, s.tables[1]))
	for idx, i := range s.res {
		row, ex := m[i[c]]
		if ex == true {
			s.res[idx] = append(i, row...)
		}
	}
}

func (s *searcher) getIDs() []string {
	// Returns IDs from s.res
	var ret []string
	for _, i := range s.res {
		ret = append(ret, i[0])
	}
	return ret
}

func (s *searcher) patientMap() map[string][]string {
	// Adds to map of patients in blocks
	var res [][]string
	ids := s.getIDs()
	l := float64(len(ids))
	d := math.Ceil(l / 50000.0)
	idx := int(l/d)
	ind := 0
	for i := 0; i >= d; i++ {
		if ind+idx > int(l) {
			// Get last less than idx rows
			idx = int(l) - ind + 1
		}
		vals := strings.Join(ids[ind:idx], ",")
		r := dbIO.GetRows(s.db, "Patient", "ID", vals, "*")
		res = append(res, r...)
		ind = ind + idx
	}
	return toMap(res)
}

func (s *searcher) getPatients() {
	// Adds patient records to s.res
	m := s.patientMap()
	for idx, i := range s.res {
		row, ex := m[i[0]]
		if ex == true {
			// Append existing record without ID
			s.res[idx] = append(row, i[1:])
		}
	}
}

func (s *searcher) getTumorRelation(ch chan [][]string, row []string) {
	// Returns matching tumor relation entries
	var ret [][]string
	table := dbIO.GetRows(s.db, "tumor_relation", "tumor_id", row[0], "*")
	for _, i := range table {
		res := append(i, row[1:])
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
	
func (s *searcher) searchSource() {
	// Searches source tables
	s.header = "ID,service_name,account_id"
	switch s.column {
		case "account_id":
			if s.user == "root" {
				// Return both tables
				s.header = s.header + ",Account,submitter_name"
				s.searchPairedTables(2)
			} else {
				if s.tables[0] == "Source" {
					// Return single table
					s.res = dbIO.GetRows(s.db, s.table[0], s.column, s.value, "*")
				} else {
					fmt.Println("\n\t[Error] Must be root to access Accounts table. Exiting.\n")
					os.Exit(99)
				}
			}
		default:
			if s.tables[0] == "Source" && s.user != "root" {
				fmt.Println("\n\t[Error] Must be root to access Accounts table. Exiting.\n")
				os.Exit(99)
			} else {
				s.res = dbIO.GetRows(s.db, s.table[0], s.column, s.value, "*")
			}
	}
	if s.short == false {
		s.header = "ID,Sex,Age,Castrated,taxa_id,source_id,Species,Date,Comments" + strings.Replace(s.header, "ID,", ",")
		s.getPatients()
	}
}

func (s *searcher) searchPatient() {
	// Searches any match that include the patient table
	switch s.column {
		case "ID":
			
		case "taxa_id":

		default:
	}
}

func (s *searcher) assignSearch() {
	// Runs appropriate search based on input
	switch s.table[0] {
		// Start with potential mutliple entries
		case "Patient":
			s.searchPatient()
		case "Source":
			s.searchSource()
		case "Tumor_relation":
			s.searchTumor()
		case "Taxonomy":
			//s.searchTaxonomy()
		case "Common":
			//s.searchTaxonomy()
		case "Life_history":
			//s.searchLifeHistory()
		case "Totals":
			//s.searchTotals()
		case "Diagnosis":
			//s.searchDiagnosis()
		case "Tumor":
			s.searchTumor()
		case "Accounts":
			s.searchSource()
	}
}

func searchColumns(db *sql.DB, col map[string]string, tables []string) ([][]string, string) {
	// Determines search procedure
	fmt.Printf("\tSearching for records with %s in column %s...", *value, *column)
	s := newSearcher(db, col, tables)
	s.assignSearch()
	return s.res, r.header
}

// This script contains functions for searching tables for a given column/value combination

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	"os"
	"strings"
)

var (
	na := []string{"NA", "NA", "NA", "NA", "NA", "NA", "NA"}
)

type searcher struct {
	db		*sql.DB
	user	string
	columns	map[string]string
	tables	[]string
	column	string
	value	string
	short	bool
	common	bool
	res		[][]string
	ids		[]string
	taxaids	[]string
	header	string
}

func newSearcher(db *sql.DB, col map[string]string, tables []string) *searcher {
	// Assigns starting values to searcher
	s := new(searcher)
	s.db = db
	s.user = *user
	s.columns = col
	s.tables = tables
	s.column = *column
	s.value = *value
	s.short = *short
	s.common = *common
	return s
}

func (s *searcher) searchPairedTables(c int) {
	// Cancatentes results from paired tables, c indicates id column
	s.res = dbIO.GetRows(s.db, s.tables[0], s.column, s.value, "*")
	// Get entire table in case value is not in both tables
	m := toMap(dbIO.GetTable(s.db, s.tables[1]))
	for idx, i := range s.res {
		row, ex := m[i[c]]
		if ex == true {
			s.res[idx] = append(i, row...)
		}
	}
}

func (s *searcher) setIDs() {
	// Sets IDs from s.res
	for _, i := range s.res {
		s.ids = append(s.ids, i[0])
	}
}

func (s *searcher) setTaxaIDs() {
	// Stores taxa ids from res
	for _, i := range s.res {
		s.taxaids = append(s.taxaids, i[4])
	}
}

func (s *searcher) appendSource() {
	// Appends data from source table to res
	m := toMap(dbIO.GetRows(s.db, "Source", "ID", strings.Join(s.ids, ","), "*")
	for idx, i := range s.res {
		row , ex := m[i[0]]
		if ex == true {
			s.res[idx] = append(i, row...)
		} else {
			s.res[idx] = append(i, na[:2]...)
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
			s.res[idx] = append(i, na...)
		}
	}
}

func (s *searcher) getPatients() {
	// Prepends patient records to s.res
	s.setIDs()
	m := toMap(dbIO.GetRows(s.db, "Patient", "ID", strings.Join(s.ids, ","), "*"))
	for idx, i := range s.res {
		row, ex := m[i[0]]
		if ex == true {
			// Append existing record without ID
			s.res[idx] = append(row, i[1:]...)
		}
	}
}

// ---------------------------------------------------------------------------
	
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
					s.res = dbIO.GetRows(s.db, s.tables[0], s.column, s.value, "*")
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
				s.res = dbIO.GetRows(s.db, s.tables[0], s.column, s.value, "*")
			}
	}
	if s.short == false {
		s.header = "ID,Sex,Age,Castrated,taxa_id,source_id,Species,Date,Comments" + strings.Replace(s.header, "ID,", ",", 1)
		s.getPatients()
	}
}

func (s *searcher) searchDiagnosis() {
	// Returns diagnosis entires
	s.header = "ID,MassPresent,Necopsy,Metastasis"
	s.res = dbIO.GetRows(s.db, s.tables[0], s.column, s.value, "*")
	if s.short == false {
		s.header = "ID,Sex,Age,Castrated,taxa_id,source_id,Species,Date,Comments" + strings.Replace(s.header, "ID,", ",", 1)
		s.getPatients()
	}
}

func (s *searcher) searchPatient() {
	// Searches any match that include the patient table
	s.res = dbIO.GetRows(s.db, s.tables[0], s.column, s.value, "*")
	s.setIDs()
	s.setTaxaIDs()
	s.appendTaxonomy()
	if s.short == false {
		var t map[string][]string
		d := toMap(dbIO.GetRows(s.db, "Diagnosis", "ID", strings.Join(s.ids, ","), "*"))
		if *count == false {
		// Skip if not needed since this is the most time consuming step
			t = s.getTumor()
		}
		for idx, i := range s.res {
			// Concatenate tables
			row := i
			id := i[0]
			diag, ex := d[id]
			if ex == true {
				row = append(row, diag...)
			} else {
				row = append(row, []string{"NA", "NA", "NA"}...)
			}
			tumor, e := t[id]
			if e == true {
				row = append(row, tumor...)
			} else {
				row = append(row, []string{"NA", "NA", "NA", "NA"}...)
			}
			s.res[idx] = row
		}
		s.appendSource()
	}
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
			s.searchSource()
		case "Tumor_relation":
			s.searchTumor()
		//case "Taxonomy":
			//s.searchTaxonomy()
		//case "Common":
			//s.searchTaxonomy()
		//case "Life_history":
			//s.searchLifeHistory()
		//case "Totals":
			//s.searchTotals()
		case "Diagnosis":
			s.searchDiagnosis()
		case "Tumor":
			s.searchTumor()
		case "Accounts":
			s.searchSource()
	}
}

func searchColumns(db *sql.DB, col map[string]string, tables []string) ([][]string, string) {
	// Determines search procedure
	fmt.Printf("\tSearching for records with %s in column %s...\n", *value, *column)
	s := newSearcher(db, col, tables)
	s.assignSearch()
	return s.res, s.header
}

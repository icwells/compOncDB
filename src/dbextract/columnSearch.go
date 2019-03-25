// This script contains functions for searching tables for a given column/value combination

package dbextract

import (
	"fmt"
	"github.com/icwells/dbIO"
	"os"
	"strings"
)

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

func (s *searcher) assignSearch(count bool) {
	// Runs appropriate search based on input
	switch s.tables[0] {
	// Start with potential mutliple entries
	case "Patient":
		s.searchPatient()
	case "Source":
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
	case "Accounts":
		s.searchAccounts()
	}
	if len(s.res) >= 1 {
		// res and ids must be set first
		s.setTaxaIDs()
		if s.infant == false {
			s.filterInfantRecords()
		}
		s.appendDiagnosis()
		s.appendTaxonomy()
		s.appendSource()
	}
}

func SearchColumns(db *dbIO.DBIO, tables []string, user, column, op, value string, count, com, inf bool) ([][]string, string) {
	// Determines search procedure
	fmt.Printf("\tSearching for records with '%s' in column %s...\n", value, column)
	s := newSearcher(db, tables, user, column, op, value, com, inf)
	s.assignSearch(count)
	return s.res, s.header
}

func SearchSingleTable(db *dbIO.DBIO, table, user, column, op, value string, com, inf bool) ([][]string, string) {
	// Returns results from single table
	fmt.Printf("\tSearching table %s for records where %s %s %s...\n", table, column, op, value)
	s := newSearcher(db, []string{table}, user, column, op, value, com, inf)
	// Overwrite standard header
	s.header = s.db.Columns[table]
	s.res = s.db.EvaluateRows(table, s.column, s.operator, s.value, "*")
	if s.infant == false && table == "Patient" {
		s.setIDs()
		s.setTaxaIDs()
		s.filterInfantRecords()
	}
	return s.res, s.header
}

// This script contains functions for searching tables for a given column/value combination

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/strarray"
	//"os"
	"strings"
)

/*func (s *searcher) searchAccounts(e codbutils.Evaluation) {
	// Searches source tables
	if s.user != "root" {
		fmt.Print("\n\t[Error] Must be root to access Accounts table. Exiting.\n\n")
		os.Exit(1010)
	}
	var accounts []string
	target := s.value
	if s.column != "account_id" {
		// Get target account IDs
		aids := s.db.GetRows(e.Table, e.Column, e.Value, "ID")
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
	s.res = dbupload.ToMap(s.db.GetRows("Patient", "ID", strings.Join(s.ids, ","), "*"))
}*/

func (s *searcher) setPatient() {
	// Reads all patient records with ids in s.ids
	s.res = dbupload.ToMap(s.db.GetRows("Patient", "ID", strings.Join(s.ids, ","), "*"))
}

func (s *searcher) submitEvaluation(e codbutils.Evaluation) []string {
	// Gets ids matching evaluation criteria
	var ret []string
	ids := s.db.EvaluateRows(e.Table, e.Column, e.Operator, e.Value, e.ID)
	for _, i := range ids {
		// Convert to string slice
		ret = append(ret, i[0])
	}
	return ret
}

func (s *searcher) filterIDs(target, match []string) []string {
	// Removes target which are not present in ids slice
	var ret []string
	for _, i := range target {
		if strarray.InSliceStr(match, i) {
			ret = append(ret, i)
		}
	}
	return ret
}

func (s *searcher) assignSearch(eval []codbutils.Evaluation) {
	// Runs appropriate search based on input
	var taxa, patients []codbutils.Evaluation
	for _, i := range eval {
		// Sort by id type
		if i.ID == "ID" {
			patients = append(patients, i)
		} else {
			taxa = append(taxa, i)
		}
	}
	for idx, i := range taxa {
		// Populate taxaids and filter with additional criteria
		if idx == 0 {
			s.taxaids = s.submitEvaluation(i)
		} else {
			s.taxaids = s.filterIDs(s.taxaids, s.submitEvaluation(i))
		}
	}
	// Populate patient ids
	if len(s.taxaids) > 0 {
		s.getIDs("Patient", "taxa_id", strings.Join(s.taxaids, ","))
	} else if len(patients) > 0 {
		s.ids = s.submitEvaluation(patients[0])
		if len(patients) > 1 {
			patients = patients[1:]
		} else {
			patients = nil
		}
	}
	if len(patients) > 0 {
		// Filter patient ids by additional criteria
		for _, i := range patients {
			s.ids = s.filterIDs(s.ids, s.submitEvaluation(i))
		}
	}
	// Store patient results and update taxaids
	s.setPatient()
	s.setTaxaIDs()
}

func (s *searcher) searchSingleTable(table string) {
	// Stores value from single table
	var ids string
	typ := "taxa_id"
	s.header = s.db.Columns[table]
	if table == "Patient" || !strings.Contains(s.header, typ) {
		typ = "ID"
		ids = strings.Join(s.ids, ",")
	} else {
		ids = strings.Join(s.taxaids, ",")
	}
	s.res = dbupload.ToMap(s.db.GetRows(table, typ, ids, "*"))
}

func SearchColumns(db *dbIO.DBIO, user, table string, eval []codbutils.Evaluation, count, inf bool) ([][]string, string) {
	// Determines search procedure
	fmt.Println("\tSearching for matching records...")
	s := newSearcher(db, user, inf)
	s.assignSearch(eval)
	if len(s.res) >= 1 {
		if s.infant == false {
			s.filterInfantRecords()
		}
		if table != "" {
			// Return results from single table
			s.searchSingleTable(table)
		} else if count == false {
			// res and ids must be set first
			s.appendDiagnosis()
			s.appendTaxonomy()
			s.appendSource()
		}
	}
	return s.toSlice(), s.header
}

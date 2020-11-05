// This script contains functions for searching tables for a given column/value combination

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/strarray"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"log"
	"sort"
	"strings"
)

func (s *searcher) setErr(e codbutils.Evaluation) {
	// Stores error message if no match is found for given evalutation
	s.msg = fmt.Sprintf("Found 0 records where %s is %s.", e.Column, e.Value)
	matches := fuzzy.RankFindFold(e.Value, s.db.GetColumnText(e.Table, e.Column))
	if matches.Len() > 0 {
		sort.Sort(matches)
		if matches[0].Target != e.Value {
			s.msg += fmt.Sprintf(" Did you mean %s?", matches[0].Target)
		}
	}
	s.msg += "\n"
}

func (s *searcher) setPatient() {
	// Reads all patient records with ids in s.ids
	if len(s.ids) > 0 {
		s.res = codbutils.ToMap(s.db.GetRows("Patient", "ID", strings.Join(s.ids, ","), "*"))
	} else if len(s.taxaids) > 0 {
		s.res = codbutils.ToMap(s.db.GetRows("Patient", "taxa_id", strings.Join(s.taxaids, ","), "*"))
	}
}

func (s *searcher) submitEvaluation(e codbutils.Evaluation) []string {
	// Gets ids matching evaluation criteria
	var ret []string
	var ids [][]string
	if e.Operator == "^" {
		ids = s.db.ColumnContains(e.Table, e.Column, e.Value, e.ID)
	} else {
		ids = s.db.EvaluateRows(e.Table, e.Column, e.Operator, e.Value, e.ID)
	}
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
	s.res = codbutils.ToMap(s.db.GetRows(table, typ, ids, "*"))
}

func (s *searcher) searchPatientIDs(patients []codbutils.Evaluation) {
	// Populate patient ids
	if len(s.taxaids) > 0 {
		s.getIDs("Patient", "taxa_id", strings.Join(s.taxaids, ","))
	} else if len(patients) > 0 {
		s.ids = s.submitEvaluation(patients[0])
		if len(s.ids) == 0 {
			s.setErr(patients[0])
		} else if len(patients) > 1 {
			patients = patients[1:]
		} else {
			patients = nil
		}
	}
	if s.msg == "" && len(patients) > 0 {
		// Filter patient ids by additional criteria
		for _, i := range patients {
			s.ids = s.filterIDs(s.ids, s.submitEvaluation(i))
			if len(s.ids) == 0 {
				s.setErr(i)
				break
			}
		}
	}
}

func (s *searcher) searchTaxaIDs(taxa []codbutils.Evaluation) {
	// Populates taxaids and filter with additional criteria
	for idx, i := range taxa {
		if idx == 0 {
			s.taxaids = s.submitEvaluation(i)
		} else {
			s.taxaids = s.filterIDs(s.taxaids, s.submitEvaluation(i))
		}
		if len(s.taxaids) == 0 {
			s.setErr(i)
			break
		}
	}
}

func (s *searcher) searchJoin(id, table string, eval []codbutils.Evaluation) []string {
	// Searches for given id type with given evaluations
	var join, where, ret []string
	target := fmt.Sprintf("%s.%s", table, id)
	// Subset by taxa ids
	if len(s.taxaids) > 0 {
		where = append(where, fmt.Sprintf("%s.taxa_id IN (%s)", table, strings.Join(s.taxaids, ",")))
	}
	for _, i := range eval {
		join = append(join, fmt.Sprintf("JOIN %s %s ON %s = %s.%s", table, i.Table, target, i.Table, id))
		if i.Operator == "^" {
			where = append(where, fmt.Sprintf("INSTR(%s.%s, '%s') > 0", i.Table, i.Column, i.Value))
		} else {
			where = append(where, fmt.Sprintf("%s.%s %s %s", i.Table, i.Column, i.Operator, i.Value))
		}

	}
	cmd := fmt.Sprintf("SELECT %s FROM %s", target, table)
	if len(eval) > 1 {
		cmd = fmt.Sprintf("%s %s ", cmd, strings.Join(join, " "))
	}
	cmd = fmt.Sprintf("%s WHERE %s;", cmd, strings.Join(where, " AND "))
	s.logger.Println(cmd)
	for _, i := range s.db.Execute(cmd) {
		// Convert to string slice
		ret = append(ret, i[0])
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
		} else if i.ID == "taxa_id" {
			taxa = append(taxa, i)
		}
	}
	if len(taxa) > 0 {
		//s.taxaids = s.searchJoin("taxa_id", "Taxonomy", taxa)
		s.searchTaxaIDs(taxa)
	}
	if s.msg == "" {
		//s.ids = s.searchJoin("ID", "Patient", patients)
		s.searchPatientIDs(patients)
		// Store patient results and update taxaids
		s.setPatient()
		s.setTaxaIDs()
	}
}

func columnSearch(db *dbIO.DBIO, logger *log.Logger, table string, eval []codbutils.Evaluation, inf bool) *searcher {
	// Determines search procedure
	s := newSearcher(db, logger, inf)
	s.assignSearch(eval)
	if len(s.res) >= 1 {
		if s.infant == false {
			s.filterInfantRecords()
		}
		if table != "" && table != "nil" {
			// Return results from single table
			s.searchSingleTable(table)
		} else {
			// res and ids must be set first
			s.appendDiagnosis()
			s.appendTaxonomy()
			s.appendSource()
		}
	}
	return s
}

func SearchColumns(db *dbIO.DBIO, logger *log.Logger, table string, eval [][]codbutils.Evaluation, inf bool) (*dataframe.Dataframe, string) {
	// Wraps calls to columnSearch
	var ret *dataframe.Dataframe
	logger.Println("Searching for matching records...")
	for idx, i := range eval {
		s := columnSearch(db, logger, table, i, inf)
		res := s.toDF()
		if s.msg != "" {
			logger.Print(s.msg)
		} else {
			logger.Printf("Found %d records where %s.\n", res.Length(), i[0].String())
		}
		if idx == 0 {
			ret = res
		} else if res.Length() > 0 {
			ret.Extend(res)
		}
	}
	return ret, fmt.Sprintf("\tFound %d records matching search criteria.\n", ret.Length())
}

func SearchDatabase(db *dbIO.DBIO, table, eval, infile string, infant bool) (*dataframe.Dataframe, string) {
	// Directs queries to appropriate functions
	var e [][]codbutils.Evaluation
	logger := codbutils.GetLogger()
	if eval != "nil" && eval != "" {
		// Search for column/value match
		e = codbutils.SetOperations(db.Columns, eval)
	} else if infile != "nil" && infile != "" {
		e = codbutils.OperationsFromFile(db.Columns, infile)
	}
	return SearchColumns(db, logger, table, e, infant)
}

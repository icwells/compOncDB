// This script contains functions for searching tables for a given column/value combination

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/simpleset"
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
	if s.ids.Length() > 0 {
		s.res = codbutils.ToMap(s.db.GetRows("Patient", "ID", strings.Join(s.ids.ToStringSlice(), ","), "*"))
	} else if s.taxaids.Length() > 0 {
		s.res = codbutils.ToMap(s.db.GetRows("Patient", "taxa_id", strings.Join(s.taxaids.ToStringSlice(), ","), "*"))
	}
	s.setTaxaIDs()
}

func (s *searcher) submitEvaluation(e codbutils.Evaluation) <-chan string {
	// Gets ids matching evaluation criteria
	ch := make(chan string)
	var ids [][]string
	if e.Operator == "^" {
		ids = s.db.ColumnContains(e.Table, e.Column, e.Value, e.ID)
	} else {
		ids = s.db.EvaluateRows(e.Table, e.Column, e.Operator, e.Value, e.ID)
	}
	go func() {
		for _, i := range ids {
			// Yield string
			ch <- i[0]
		}
		close(ch)
	}()
	return ch
}

func (s *searcher) filterIDs(ids *simpleset.Set, e codbutils.Evaluation) *simpleset.Set {
	// Removes target ids which are not present in ids slice
	ret := simpleset.NewStringSet()
	for i := range s.submitEvaluation(e) {
		if ex, _ := ids.InSet(i); ex {
			ret.Add(i)
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
		ids = strings.Join(s.ids.ToStringSlice(), ",")
	} else {
		ids = strings.Join(s.taxaids.ToStringSlice(), ",")
	}
	s.res = codbutils.ToMap(s.db.GetRows(table, typ, ids, "*"))
}

func (s *searcher) searchPatientIDs(patients []codbutils.Evaluation) {
	// Populate patient ids
	s.setIDs()
	if s.msg == "" && len(patients) > 0 {
		// Filter patient ids by additional criteria
		for _, i := range patients {
			s.ids = s.filterIDs(s.ids, i)
			if s.ids.Length() == 0 {
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
			for i := range s.submitEvaluation(i) {
				s.taxaids.Add(i)
			}
		} else {
			s.taxaids = s.filterIDs(s.taxaids, i)
		}
		if s.taxaids.Length() == 0 {
			s.setErr(i)
			break
		}
	}
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
		// Store patient results
		s.setPatient()
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

// This script contains methods for searching tumor tables

package search

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

func TumorMap(db *dbIO.DBIO) map[string][]string {
	// Returns map of all tumor entries per ID in 2d slice
	ret := make(map[string][]string)
	for _, row := range db.GetTable("Tumor") {
		id := row[0]
		if _, ex := ret[id]; !ex {
			// Add new entry
			ret[id] = row[1:]
		} else {
			// Add new entry to existing cells
			for idx, i := range row[1:] {
				ret[id][idx] += ";" + i
			}
		}
	}
	return ret
}

type searcher struct {
	db      *dbIO.DBIO
	header  []string
	logger  *log.Logger
	metadata string
	msg     string
	res     map[string][]string
}

func newSearcher(db *dbIO.DBIO, logger *log.Logger) *searcher {
	// Assigns starting values to searcher
	s := new(searcher)
	// Add default header
	s.db = db
	//s.header = strings.Join(codbutils.RecordsHeader(), ",")
	s.header = strings.Split(s.db.Columns["Records"], ",")
	s.logger = logger
	s.res = make(map[string][]string)
	return s
}

func (s *searcher) clearSearcher() {
	// Removes records from previous search
	for k := range s.res {
		delete(s.res, k)
	}
}

func (s *searcher) toDF() *dataframe.Dataframe {
	// Converts res map to dataframe
	ret, _ := dataframe.NewDataFrame(0)
	ret.SetHeader(s.header)
	for _, v := range s.res {
		if err := ret.AddRow(v); err != nil {
			panic(err)
		}
	}
	if s.metadata != "" {
		ret.SetMetaData(s.metadata)
	}
	return ret
}

func (s *searcher) toSlice() [][]string {
	// Converts res map to slice
	var ret [][]string
	for _, v := range s.res {
		ret = append(ret, v)
	}
	return ret
}

func (s *searcher) setErr(e codbutils.Evaluation) {
	// Stores error message if no match is found for given evalutation
	s.msg = fmt.Sprintf("Found 0 records where %s is %s.", e.Column, e.Value)
	if e.Operator != "^" {
		// Skip the 'in' command since results would be illogical
		matches := fuzzy.RankFindFold(e.Value, s.db.GetColumnText(e.Table, e.Column))
		if matches.Len() > 0 {
			sort.Sort(matches)
			if matches[0].Target != e.Value {
				s.msg += fmt.Sprintf(" Did you mean %s?", matches[0].Target)
			}
		}
	}
	s.msg += "\n"
}

func (s *searcher) setMetaData(eval []codbutils.Evaluation) {
	// Stores search options as string
	var m []string
	if len(s.metadata) == 0 {
		m = append(m, codbutils.GetTimeStamp())
	} else {
		// Store metadata for multiple searches
		m = append(m, s.metadata + ", ,")
	}
	if len(eval) > 0 {
		for _, i := range eval {
			m = append(m, i.String())
		}
	}
	s.metadata = strings.Join(m, ",")
}

func (s *searcher) replaceNull(row []string) []string {
	// Replaces Null values with NA
	for idx, i := range row {
		if i == "" || i == "NULL" {
			row[idx] = "NA"
		}
	}
	return row
}

func (s *searcher) formatCommand(e []codbutils.Evaluation) (string, []codbutils.Evaluation) {
	// Formats sql command
	cmd := strings.Builder{}
	cmd.WriteString("SELECT * FROM Records")
	for idx, i := range e {
		if idx == 0 {
			cmd.WriteString(" WHERE ")
		} else {
			cmd.WriteString(" AND ")
		}
		cmd.WriteByte(' ')
		cmd.WriteString(i.String())
	}
	cmd.WriteByte(';')
	return cmd.String(), e
}

func (s *searcher) getRecords(eval []codbutils.Evaluation, lh bool) {
	// Gets matching records from view
	idx := strarray.SliceIndex(s.header, "female_maturity")
	pid := strarray.SliceIndex(s.header, "primary_tumor")
	tid := strarray.SliceIndex(s.header, "Type")
	lid := strarray.SliceIndex(s.header, "Location")
	cmd, e := s.formatCommand(eval)
	rows := s.db.Execute(cmd)
	if len(rows) == 0 {
		s.setErr(e[0])
	} else {
		for _, i := range rows {
			id := i[0]
			if _, ex := s.res[id]; !ex {
				if !lh {
					// Drop life history data
					s.res[id] = s.replaceNull(i[:idx])
				} else {
					s.res[id] = s.replaceNull(i)
				}
			} else {
				// Merge multiple tumor records
				if i[pid] == "1" {
					// Store highest primary tumor value
					s.res[id][pid] = i[pid]
				}
				s.res[id][tid] += ";" + i[tid]
				s.res[id][lid] += ";" + i[lid]
			}
		}
		if !lh {
			// Remove life history from header
			s.header = s.header[:idx]
		}
		s.setMetaData(eval)
	}
}

func (s *searcher) setEvaluations(eval string, inf bool) []codbutils.Evaluation {
	// Formats evaluations
	if len(eval) > 0 {
		eval += ","
	}
	// Add evaluation to remove infant records
	if inf {
		eval += "Infant = 1"
	} else {
		eval += "Infant != 1"
	}
	return codbutils.RecordsEvaluations(s.db.Columns, eval)
}

func SearchRecords(db *dbIO.DBIO, logger *log.Logger, eval string, inf, lh bool) (*dataframe.Dataframe, string) {
	// Searches for matching records
	s := newSearcher(db, logger)
	s.logger.Println("Searching for matching records...")
	s.getRecords(s.setEvaluations(eval, inf), lh)
	ret := s.toDF()
	if s.msg == "" {
		s.msg = fmt.Sprintf("\tFound %d records matching search criteria.\n", ret.Length())
	}
	return ret, s.msg
}

func SearchColumns(db *dbIO.DBIO, logger *log.Logger, eval [][]codbutils.Evaluation, inf, lh bool) (*dataframe.Dataframe, string) {
	// Wraps calls from server to getRecords
	var ret *dataframe.Dataframe
	s := newSearcher(db, logger)
	s.logger.Println("Searching for matching records...")
	for idx, i := range eval {
		s.getRecords(i, lh)
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
		s.clearSearcher()
	}
	return ret, fmt.Sprintf("\tFound %d records matching search criteria.\n", ret.Length())
}

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

func (s *searcher) toDF() *dataframe.Dataframe {
	// Converts res map to dataframe
	ret, _ := dataframe.NewDataFrame(0)
	ret.SetHeader(s.header)
	for k, v := range s.res {
		row := append([]string{k}, v...)
		if err := ret.AddRow(row); err != nil {
			panic(strings.Join(s.header, " "))
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
	for k, v := range s.res {
		row := append([]string{k}, v...)
		ret = append(ret, row)
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

func (s *searcher) setMetaData(eval string, inf bool) {
	// Stores search options as string
	var m []string
	m = append(m, codbutils.GetTimeStamp())
	if eval != "" && eval != "nil" {
		m = append(m, eval)
	}
	m = append(m, fmt.Sprintf("KeepInfantRecords=%v", inf))
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

func (s *searcher) formatCommand(eval string, inf bool) (string, []codbutils.Evaluation) {
	// Formats sql command
	cmd := strings.Builder{}
	if inf {
		// Add evaluation to remove infant records
		if len(eval) > 0 {
			eval += ","
		}
		eval += "Infant != 1"
	}
	e := codbutils.RecordsEvaluations(s.db.Columns, eval)
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

func (s *searcher) getRecords(eval string, inf, lh bool) {
	// Gets matching records from view
	idx := strarray.SliceIndex(s.header, "female_maturity")
	pid := strarray.SliceIndex(s.header, "primary_tumor")
	tid := strarray.SliceIndex(s.header, "Type")
	lid := strarray.SliceIndex(s.header, "Location")
	cmd, e := s.formatCommand(eval, inf)
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
	}
}

func SearchRecords(db *dbIO.DBIO, logger *log.Logger, eval string, inf, lh bool) (*dataframe.Dataframe, string) {
	// Wraps calls to columnSearch
	logger.Println("Searching for matching records...")
	s := newSearcher(db, logger)
	s.getRecords(eval, inf, lh)
	ret := s.toDF()
	if s.msg == "" {
		s.msg = fmt.Sprintf("\tFound %d records matching search criteria.\n", ret.Length())
	}
	return ret, s.msg
}

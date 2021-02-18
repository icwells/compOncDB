// This script contains methods for searching tumor tables

package search

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/simpleset"
	"log"
	"strings"
)

type searcher struct {
	db      *dbIO.DBIO
	header  string
	ids     *simpleset.Set
	logger  *log.Logger
	msg     string
	na      []string
	res     map[string][]string
	taxa    map[string][]string
	taxaids *simpleset.Set
}

func newSearcher(db *dbIO.DBIO, logger *log.Logger) *searcher {
	// Assigns starting values to searcher
	s := new(searcher)
	// Add default header
	s.db = db
	s.header = strings.Join(codbutils.RecordsHeader(), ",")
	s.ids = simpleset.NewStringSet()
	s.logger = logger
	s.na = []string{"NA", "NA", "NA", "NA", "NA", "NA", "NA"}
	s.res = make(map[string][]string)
	s.taxa = make(map[string][]string)
	s.taxaids = simpleset.NewStringSet()
	return s
}

func (s *searcher) toDF() *dataframe.Dataframe {
	// Converts res map to dataframe
	ret, _ := dataframe.NewDataFrame(0)
	ret.SetHeader(strings.Split(s.header, ","))
	for k, v := range s.res {
		row := append([]string{k}, v...)
		ret.AddRow(row)
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

func (s *searcher) setIDs() {
	// Stores initial ids set
	s.ids.Clear()
	if s.taxaids.Length() > 0 {
		for _, i := range s.db.GetRows("Patient", "taxa_id", strings.Join(s.taxaids.ToStringSlice(), ","), "ID") {
			s.ids.Add(i[0])
		}
	} else {
		// Get all ids
		for _, i := range s.db.GetColumnText("Patient", "ID") {
			s.ids.Add(i)
		}
	}
}

func (s *searcher) setTaxaIDs() {
	// Stores taxa ids from patient results
	s.taxaids.Clear()
	for _, v := range s.res {
		s.taxaids.Add(v[5])
	}
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

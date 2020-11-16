// This script contains methods for searching tumor tables

package dbextract

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/simpleset"
	"log"
	"strconv"
	"strings"
)

type searcher struct {
	db      *dbIO.DBIO
	header  string
	ids     *simpleset.Set
	infant  bool
	logger  *log.Logger
	msg     string
	na      []string
	res     map[string][]string
	taxa    map[string][]string
	taxaids *simpleset.Set
}

func newSearcher(db *dbIO.DBIO, logger *log.Logger, inf bool) *searcher {
	// Assigns starting values to searcher
	s := new(searcher)
	// Add default header
	s.db = db
	s.header = strings.Join(codbutils.RecordsHeader(), ",")
	s.ids = simpleset.NewStringSet()
	s.infant = inf
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

func getColumn(idx int, table [][]string) []string {
	// Stores values in set and returns slice of unique entries
	tmp := simpleset.NewStringSet()
	for _, i := range table {
		if i[idx] != "-1" {
			tmp.Add(i[idx])
		}
	}
	return tmp.ToStringSlice()
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
		s.taxaids.Add(v[3])
	}
}

func (s *searcher) filterInfantRecords() {
	// Removes infant records from search results
	// In summary.go
	ages := GetMinAges(s.db, s.taxaids.ToStringSlice())
	// Filter results
	for k, v := range s.res {
		if len(v) >= 4 {
			min, ex := ages[v[3]]
			if ex == true {
				age, err := strconv.ParseFloat(v[1], 64)
				if err == nil && age <= min {
					// Remove infant record
					delete(s.res, k)
					s.ids.Pop(k)
				}
			}
		}
	}
	// Update ids
	s.setTaxaIDs()
}

func (s *searcher) appendSource() {
	// Appends data from source table to res
	m := codbutils.ToMap(s.db.GetRows("Source", "ID", strings.Join(s.ids.ToStringSlice(), ","), "*"))
	for k, v := range s.res {
		row, ex := m[k]
		if ex == true {
			s.res[k] = append(v, row...)
		} else {
			s.res[k] = append(v, s.na[:3]...)
		}
	}
}

func (s *searcher) getTaxonomy() {
	// Stores taxonomy (ids must be set first)
	s.taxa = codbutils.ToMap(s.db.GetRows("Taxonomy", "taxa_id", strings.Join(s.taxaids.ToStringSlice(), ","), "taxa_id,Kingdom,Phylum,Class,Orders,Family,Genus,Species"))
}

func (s *searcher) appendTaxonomy() {
	// Appends taxonomy to s.res
	if len(s.taxa) == 0 && s.taxaids.Length() > 0 {
		s.getTaxonomy()
	}
	for k, v := range s.res {
		// Apppend taxonomy to records
		taxonomy, ex := s.taxa[v[3]]
		if ex == false {
			taxonomy = s.na
		}
		s.res[k] = append(s.res[k], taxonomy...)
	}
}

func (s *searcher) appendDiagnosis() {
	// Appends data from tumor and tumor relation tables
	d := codbutils.ToMap(s.db.GetRows("Diagnosis", "ID", strings.Join(s.ids.ToStringSlice(), ","), "*"))
	t := codbutils.ToMap(s.db.GetRows("Tumor", "ID", strings.Join(s.ids.ToStringSlice(), ","), "*"))
	for k := range s.res {
		// Concatenate tables
		diag, ex := d[k]
		if ex == false {
			diag = s.na[:4]
		}
		tumor, exists := t[k]
		if exists == false {
			tumor = s.na[:4]
		}
		diag = append(diag, tumor...)
		s.res[k] = append(s.res[k], diag...)
	}
}

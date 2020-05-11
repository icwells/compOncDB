// This script contains methods for searching tumor tables

package dbextract

import (
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/simpleset"
	"strconv"
	"strings"
)

type searcher struct {
	db      *dbIO.DBIO
	infant  bool
	res     map[string][]string
	taxa    map[string][]string
	ids     []string
	taxaids []string
	header  string
	na      []string
	msg     string
}

func newSearcher(db *dbIO.DBIO, inf bool) *searcher {
	// Assigns starting values to searcher
	s := new(searcher)
	s.res = make(map[string][]string)
	s.taxa = make(map[string][]string)
	// Add default header
	s.header = "ID,Sex,Age,Castrated,taxa_id,source_id,source_name,Date,Year,Comments,"
	s.header = s.header + "Masspresent,Hyperplasia,Necropsy,Metastasis,primary_tumor,Malignant,Type,Location,"
	s.header = s.header + "Kingdom,Phylum,Class,Orders,Family,Genus,Species,service_name,Zoo,Aza,Institute,account_id"
	s.db = db
	s.infant = inf
	s.na = []string{"NA", "NA", "NA", "NA", "NA", "NA", "NA"}
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

func (s *searcher) getIDs(table, column, value string) {
	// Gets ids from target table and get patient records
	s.ids = getColumn(0, s.db.GetRows(table, column, value, "ID"))
}

func (s *searcher) setIDs() {
	// Sets IDs from s.res (ID must be in first column)
	for k := range s.res {
		s.ids = append(s.ids, k)
	}
}

func (s *searcher) setTaxaIDs() {
	// Stores taxa ids from patient results
	s.taxaids = getColumn(4, s.toSlice())
}

func (s *searcher) filterInfantRecords() {
	// Removes infant records from search results
	// In summary.go
	ages := getMinAges(s.db, s.taxaids)
	// Filter results
	for k, v := range s.res {
		if len(v) >= 4 {
			min, ex := ages[v[3]]
			if ex == true {
				age, err := strconv.ParseFloat(v[1], 64)
				if err == nil && age <= min {
					// Remove infant record
					delete(s.res, k)
				}
			}
		}
	}
	// Update ids
	s.setIDs()
	s.setTaxaIDs()
}

func (s *searcher) appendSource() {
	// Appends data from source table to res
	m := dbupload.ToMap(s.db.GetRows("Source", "ID", strings.Join(s.ids, ","), "*"))
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
	s.taxa = dbupload.ToMap(s.db.GetRows("Taxonomy", "taxa_id", strings.Join(s.taxaids, ","), "taxa_id,Kingdom,Phylum,Class,Orders,Family,Genus,Species"))
}

func (s *searcher) appendTaxonomy() {
	// Appends raxonomy to s.res
	if len(s.taxa) == 0 && len(s.taxaids) > 0 {
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
	d := dbupload.ToMap(s.db.GetRows("Diagnosis", "ID", strings.Join(s.ids, ","), "*"))
	t := dbupload.ToMap(s.db.GetRows("Tumor", "ID", strings.Join(s.ids, ","), "*"))
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

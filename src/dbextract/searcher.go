// This script contains methods for searching tumor tables

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/strarray"
	"strconv"
	"strings"
)

type searcher struct {
	db       *dbIO.DBIO
	user     string
	tables   []string
	column   string
	value    string
	operator string
	common   bool
	infant   bool
	res      map[string][]string
	ids      []string
	taxaids  []string
	header   string
	na       []string
}

func newSearcher(db *dbIO.DBIO, tables []string, user, column, op, value string, com, inf bool) *searcher {
	// Assigns starting values to searcher
	s := new(searcher)
	s.res = make(map[string][]string)
	// Add default header
	s.header = "ID,Sex,Age,Castrated,taxa_id,source_id,Date,Comments,"
	s.header = s.header + "Masspresent,Hyperplasia,Necropsy,Metastasis,primary_tumor,Malignant,Type,Location,"
	s.header = s.header + "Kingdom,Phylum,Class,Orders,Family,Genus,Species,service_name,account_id"
	s.db = db
	s.user = user
	s.tables = tables
	s.column = column
	s.value = value
	s.operator = op
	s.common = com
	s.infant = inf
	s.na = []string{"NA", "NA", "NA", "NA", "NA", "NA", "NA"}
	return s
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
	tmp := strarray.NewSet()
	for _, i := range table {
		if i[idx] != "-1" {
			tmp.Add(i[idx])
		}
	}
	return tmp.ToSlice()
}

func (s *searcher) getIDs() {
	// Gets ids from target table and get patient records
	s.ids = getColumn(0, s.db.GetRows(s.tables[0], s.column, s.value, "ID"))
	s.res = dbupload.ToMap(s.db.GetRows("Patient", "ID", strings.Join(s.ids, ","), "*"))
}

func (s *searcher) setIDs() {
	// Sets IDs from s.res (ID must be in first column)
	for k := range s.res {
		s.ids = append(s.ids, k)
	}
}

func (s *searcher) setTaxaIDs() {
	// Stores taxa ids from patient results
	s.taxaids = getColumn(3, s.toSlice())
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
			s.res[k] = append(v, s.na[:2]...)
		}
	}
}

func (s *searcher) appendTaxonomy() {
	// Appends raxonomy to s.res
	taxa := s.getTaxonomy(s.taxaids, true)
	for k, v := range s.res {
		// Apppend taxonomy to records
		fmt.Println(k, v)
		taxonomy, ex := taxa[v[3]]
		if ex == true {
			s.res[k] = append(v, taxonomy...)
		} else {
			s.res[k] = append(v, s.na...)
		}
	}
}

func (s *searcher) appendDiagnosis() {
	// Appends data from tumor and tumor relation tables
	d := dbupload.ToMap(s.db.GetRows("Diagnosis", "ID", strings.Join(s.ids, ","), "*"))
	t := dbupload.ToMap(s.db.GetRows("Tumor", "ID", strings.Join(s.ids, ","), "*"))
	for k, v := range s.res {
		// Concatenate tables
		diag, ex := d[k]
		if ex == true {
			s.res[k] = append(v, diag...)
		} else {
			s.res[k] = append(v, s.na[:5]...)
		}
		tumor, exists := t[k]
		if exists == true {
			s.res[k] = append(v, tumor...)
		} else {
			s.res[k] = append(v, s.na[:5]...)
		}
	}
}

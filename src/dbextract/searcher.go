// This script contains methods for searching tumor tables

package dbextract

import (
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
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
	res      [][]string
	ids      []string
	taxaids  []string
	header   string
	na       []string
}

func newSearcher(db *dbIO.DBIO, tables []string, user, column, op, value string, com, inf bool) *searcher {
	// Assigns starting values to searcher
	s := new(searcher)
	// Add default header
	s.header = "ID,Sex,Age,Castrated,taxa_id,source_id,Date,Comments,"
	s.header = s.header + "Masspresent,Hyperplasia,Necropsy,Metastasis,primary_tumor,Malignant,Type,Location,"
	s.header = s.header + "Kingdom,Phylum,Class,Order,Family,Genus,Species,service_name,account_id"
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

func (s *searcher) getIDs() {
	// Gets ids from target table and get patient records
	ids := s.db.GetRows(s.tables[0], s.column, s.value, "ID")
	for _, i := range ids {
		// Convert to single slice
		s.ids = append(s.ids, i[0])
	}
	s.res = s.db.GetRows("Patient", "ID", strings.Join(s.ids, ","), "*")
}

func (s *searcher) setIDs() {
	// Sets IDs from s.res (ID must be in first column)
	for _, i := range s.res {
		s.ids = append(s.ids, i[0])
	}
}

func (s *searcher) setTaxaIDs() {
	// Stores taxa ids from patient results
	for _, i := range s.res {
		s.taxaids = append(s.taxaids, i[4])
	}
}

func (s *searcher) filterInfantRecords() {
	// Removes infant records from search results
	// In summary.go
	ages := getMinAges(s.db, s.taxaids)
	// Filter results
	for idx, i := range s.res {
		if len(i) >= 5 {
			min, ex := ages[i[4]]
			if ex == true {
				age, err := strconv.ParseFloat(i[2], 64)
				if err == nil && age <= min {
					// Remove infant record
					if idx < len(s.res) - 1 {
						s.res = append(s.res[:idx], s.res[idx+1:]...)
					} else {
						// Drop last element
						s.res = s.res[:idx]
					}
				}
			}
		}
	}
	// Update ids
	s.ids = nil
	s.taxaids = nil
	s.setIDs()
	s.setTaxaIDs()
}

func (s *searcher) appendSource() {
	// Appends data from source table to res
	m := dbupload.ToMap(s.db.GetRows("Source", "ID", strings.Join(s.ids, ","), "*"))
	for idx, i := range s.res {
		row, ex := m[i[0]]
		if ex == true {
			s.res[idx] = append(i, row...)
		} else {
			s.res[idx] = append(i, s.na[:2]...)
		}
	}
}

func (s *searcher) appendTaxonomy() {
	// Appends raxonomy to s.res
	taxa := s.getTaxonomy(s.taxaids, true)
	for idx, i := range s.res {
		// Apppend taxonomy to records
		taxonomy, ex := taxa[i[4]]
		if ex == true {
			s.res[idx] = append(i, taxonomy...)
		} else {
			s.res[idx] = append(i, s.na...)
		}
	}
}

func (s *searcher) appendDiagnosis() {
	// Appends data from tumor and tumor relation tables
	d := dbupload.ToMap(s.db.GetRows("Diagnosis", "ID", strings.Join(s.ids, ","), "*"))
	t := dbupload.ToMap(s.db.GetRows("Tumor", "ID", strings.Join(s.ids, ","), "*"))
	for idx, i := range s.res {
		// Concatenate tables
		id := i[0]
		diag, ex := d[id]
		if ex == true {
			i = append(i, diag...)
		} else {
			i = append(i, s.na[:5]...)
		}
		tumor, exists := t[id]
		if exists == true {
			i = append(i, tumor...)
		} else {
			i = append(i, s.na[:5]...)
		}
		s.res[idx] = i
	}
}

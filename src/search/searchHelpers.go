// Contains functions for filtering and appending results

package search

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
	"strconv"
	"strings"
)

func GetMinAges(db *dbIO.DBIO, taxaids []string) map[string]float64 {
	// Returns map of minumum ages by taxa id
	var table map[string]string
	ages := make(map[string]float64)
	if len(taxaids) >= 1 {
		table = codbutils.EntryMap(db.GetRows("Life_history", "taxa_id", strings.Join(taxaids, ","), "Infancy,taxa_id"))
	} else {
		table = codbutils.EntryMap(db.GetColumns("Life_history", []string{"Infancy", "taxa_id"}))
	}
	// Convert string ages to float
	for k, v := range table {
		a, err := strconv.ParseFloat(v, 64)
		if err == nil {
			ages[k] = a
		}
	}
	return ages
}

func TumorMap(db *dbIO.DBIO) map[string][]string {
	// Returns map of all tumor entries per ID ni 2d slice
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

func (s *searcher) filterInfantRecords() {
	// Removes infant records from search results
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
	t := TumorMap(s.db)
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

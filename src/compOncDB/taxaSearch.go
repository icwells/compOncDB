// This script will search for given records from the comparative oncology database

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/icwells/go-tools/strarray"
	"os"
	"strings"
)

func getTumorRecords(ch chan []string, db *sql.DB, id string, tumor map[string][]string, primary bool) {
	// Returns tumor information for given id
	var loc, typ, mal, prim []string
	rows := dbIO.GetRows(db, "Tumor_relation", "ID", id, "*")
	for _, i := range rows {
		if primary == false || primary == true && i[2] == "1" {
			j, ex := tumor[i[1]]
			if ex == true {
				if i[2] == "1" {
					// Prepend primary tumor
					prim = append([]string{i[2]}, prim...)
					mal = append([]string{i[3]}, mal...)
					typ = append([]string{j[1]}, typ...)
					loc = append([]string{j[2]}, loc...)
				} else {
					// Append tumor type and location
					prim = append(prim, i[2])
					mal = append(mal, i[3])
					typ = append(typ, j[1])
					loc = append(loc, j[2])
				}
			}
		}
	}
	diag := []string{strings.Join(prim, ";"), strings.Join(mal, ";"), strings.Join(typ, ";"), strings.Join(loc, ";")}
	ch <- diag
}

func getTumor(db *sql.DB, ids []string, primary bool) map[string][]string {
	// Returns map of tumor data from patient ids
	ch := make(chan []string)
	// {id: [types], [locations]}
	rec := make(map[string][]string)
	tumor := toMap(dbIO.GetTable(db, "Tumor"))
	for _, id := range ids {
		// Get records for each patient concurrently
		go getTumorRecords(ch, db, id, tumor, primary)
		ret := <-ch
		if len(ret) >= 1 {
			rec[id] = ret
		}
	}
	return rec
}

func getMetastasis(ch chan []string, db *sql.DB, id string, diag map[string][]string, mass bool) {
	// Returns diagnosis and metastasis data
	var d []string
	row, ex := diag[id]
	if ex == true {
		if mass == false || mass == true && row[1] == "1" {
			d = row[1:]
		}
	}
	ch <- d
}

func getDiagosis(db *sql.DB, ids []string, mass bool) map[string][]string {
	// Returns metastatis info from patient ids
	ch := make(chan []string)
	diagnoses := make(map[string][]string)
	diag := toMap(dbIO.GetTable(db, "Diagnosis"))
	for _, id := range ids {
		// Get records for each patient concurrently
		go getMetastasis(ch, db, id, diag, mass)
		ret := <-ch
		if len(ret) >= 1 {
			diagnoses[id] = ret
		}
	}
	return diagnoses
}

func (s *searcher) getRecords(ids []string, mass, primary bool) map[string][]string {
	// Gets diagnosis and metastasis data and formats values
	fmt.Println("\tExtracting diagnosis information...")
	diagnoses := make(map[string][]string)
	meta := getDiagosis(s.db, ids, mass)
	tumor := getTumor(s.db, ids, primary)
	for _, i := range ids {
		// Join multiple entires for same record
		temp := append(meta[i], tumor[i]...)
		diagnoses[i] = temp
	}
	return diagnoses
}

func (s *searcher) getTaxonomy() map[string][]string {
	// Returns taxonomy as map with taxa id as key
	taxa := make(map[string][]string)
	fmt.Println("\tExtracting taxonomy information...")
	table := dbIO.GetTableMap(s.db, "Taxonomy")
	for _, id := range s.ids {
		if strarray.InMapSli(table, id) == true {
			// Exclude source and species (in patient table)
			taxa[id] = table[id][:6]
		}
	}
	return taxa
}

func (s *searcher) getTaxa() (map[string][][]string, []string) {
	// Extracts patient data using taxa ids
	patients := make(map[string][][]string)
	var uid []string
	table := dbIO.GetTable(s.db, "Patient")
	for _, i := range table {
		for _, id := range s.ids {
			if id == i[4] {
				patients[id] = append(patients[id], i)
				uid = append(uid, i[0])
				break
			}
		}
	}
	return patients, uid
}

func (s *searcher) getTaxaIDs(names []string) {
	// Returns taxa id from species name
	var table [][]string
	if s.common == true {
		// Get taxonomy ids from common name list
		table = dbIO.SearchColumnText(s.db, "Common", "Name", names)
	} else {
		// Get ids from taxonomy table
		table = dbIO.SearchColumnText(s.db, "Taxonomy", s.column, names)
	}
	for _, row := range table {
		s.ids = append(s.ids, row[0])
	}
}

func (s *searcher) checkLevel(level string) {
	// Makes sure a valid taxonomic level is given
	found := false
	if s.common == true {
		// Overwrite to species for common name comparison
		s.column = "Species"
	} else {
		levels := []string{"Kingdom", "Phylum", "Class", "Orders", "Family", "Genus", "Species"}
		// Convert level to title case
		level = strings.Title(level)
		for _, i := range levels {
			if level == i {
				found = true
				break
			}
		}
		if found == false {
			fmt.Println("\n\t[Error] Please enter a valid taxonomic level. Exiting.\n")
			os.Exit(11)
		}
		s.column = level
	}
}

func searchTaxonomicLevels(db *sql.DB, col map[string]string, names []string) ([][]string, string) {
	// Extracts data using species names
	s := newSearcher(db, col, []string{"Taxonomy"})
	s.header = "ID,Sex,Age,Castrated,taxa_id,source_id,Species,Date,Comments,"
	s.header = s.header + "Masspresent,Necropsy,Metastasis,primary_tumor,Malignant,Type,Location,Kingdom,Phylum,Class,Orders,Family,Genus"
	fmt.Println("\tExtracting patient information...")
	s.checkLevel(*level)
	s.getTaxaIDs(names)
	patients, uid := s.getTaxa()
	// Leaving primary tumor and mass present switches false for now
	records := s.getRecords(uid, false, false)
	taxonomy := s.getTaxonomy()
	for _, id := range s.ids {
		for _, i := range patients[id] {
			_, ex := records[i[0]]
			if ex == true {
				rec := append(i, records[i[0]]...)
				rec = append(rec, taxonomy[id]...)
				s.res = append(s.res, rec)
			}
		}
	}
	return s.res, s.header
}

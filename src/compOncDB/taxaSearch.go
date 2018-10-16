// This script will search for given records from the comparative oncology database

package main

import (
	"bytes"
	"database/sql"
	"dbIO"
	"fmt"
	"os"
	"strings"
)

func (s *searcher) getRecords() map[string][]string {
	// Gets diagnosis and metastasis data and formats values
	fmt.Println("\tExtracting diagnosis information...")
	d := toMap(dbIO.GetRows(s.db, "Diagnosis", "ID", strings.Join(s.ids, ","), "*"))
	tumor := s.getTumor()
	for k, v := range d {
		// Join multiple entires for same record
		d[k] = append(v, tumor[k]...)
	}
	return d
}

func (s *searcher) getTaxa() map[string][][]string {
	// Extracts patient data using taxa ids
	patients := make(map[string][][]string)
	table := dbIO.GetRows(s.db, "Patient", "taxa_id", strings.Join(s.taxaids, ","), "*")
	for _, i := range table {
		id := i[4]
		patients[id] = append(patients[id], i)
		s.ids = append(s.ids, i[0])
	}
	return patients
}

func (s *searcher) getTaxonomy(names []string, ids bool) map[string][]string {
	// Stores taxa ids from species name and returns taxonomy
	ret := make(map[string][]string)
	var table [][]string
	if s.common == true {
		// Get taxonomy ids from common name list
		c := dbIO.SearchColumnText(s.db, "Common", "Name", names)
		buffer := bytes.NewBufferString(c[0][0])
		for _, i := range c[1:] {
			buffer.WriteByte(',')
			buffer.WriteString(i[0])
		}
		// Get taxonomy entries
		table = dbIO.GetRows(s.db, "Taxonomy", "taxa_id", buffer.String(), "*")
	} else if ids == false {
		// Get matching taxonomies
		table = dbIO.GetRows(s.db, "Taxonomy", s.column, strings.Join(names, ","), "*")
	} else {
		table = dbIO.GetRows(s.db, "Taxonomy", "taxa_id", strings.Join(names, ","), "*")
	}
	for _, row := range table {
		if ids == false {
			// Append taxa id and return map of taxonomy entries
			s.taxaids = append(s.taxaids, row[0])
		}
		// Exclude taxa id, source, and species (in patient table)
		ret[row[0]] = row[1:7]
	}
	return ret
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
	var records map[string][]string
	s := newSearcher(db, col, []string{"Taxonomy"})
	s.header = "ID,Sex,Age,Castrated,taxa_id,source_id,Species,Date,Comments,"
	s.header = s.header + "Masspresent,Necropsy,Metastasis,primary_tumor,Malignant,Type,Location,Kingdom,Phylum,Class,Orders,Family,Genus"
	fmt.Println("\tExtracting patient information...")
	s.checkLevel(*level)
	taxonomy := s.getTaxonomy(names, false)
	patients := s.getTaxa()
	if *count == false {
		// Skip if not needed since this is the most time consuming step
		records = s.getRecords()
	}
	for _, id := range s.taxaids {
		for _, i := range patients[id] {
			_, ex := records[i[0]]
			var rec []string
			if ex == true {
				rec = append(i, records[i[0]]...)
			} else {
				rec = i
			}
			rec = append(rec, taxonomy[id]...)
			s.res = append(s.res, rec)
		}
	}
	return s.res, s.header
}

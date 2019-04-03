// This script will search for given records from the comparative oncology database

package dbextract

import (
	"bytes"
	"fmt"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/strarray"
	"os"
	"strings"
)

func toTitle(names []string) []string {
	// Converts all input names to title case
	var ret []string
	for _, i := range names {
		ret = append(ret, strings.TitleCase(i))
	}
	return ret
}

func (s *searcher) getTaxa() {
	// Extracts patient data using taxa ids
	s.res = dbupload.ToMap(s.db.GetRows("Patient", "taxa_id", strings.Join(s.taxaids, ","), "*"))
	s.setIDs()
}

func (s *searcher) getTaxonomy() {
	// Stores taxonomy (ids must be set first)
	s.taxa = dbupload.ToMap(s.db.GetRows("Taxonomy", "taxa_id", strings.Join(s.taxaids, ","), "taxa_id,Kingdom,Phylum,Class,Orders,Family,Genus,Species"))
}

func (s *searcher) setTaxonomy(names []string) bool {
	// Stores taxa ids and taxonomies from species name
	var ret bool
	if s.common == true {
		// Get taxonomy ids from common name list
		c := s.db.GetRows("Common", "Name", strings.Join(toTitle(names), ","), "taxa_id")
		if len(c) >= 1 {
			// Colect taxa IDs
			buffer := bytes.NewBufferString(c[0][0])
			for _, i := range c[1:] {
				buffer.WriteByte(',')
				buffer.WriteString(i[0])
			}
			// Get taxonomy entries
			s.taxa = dbupload.ToMap(s.db.GetRows("Taxonomy", "taxa_id", buffer.String(), "taxa_id,Kingdom,Phylum,Class,Orders,Family,Genus,Species"))
		}
	} else {
		// Get matching taxonomies
		s.taxa = dbupload.ToMap(s.db.GetRows("Taxonomy", s.column, strings.Join(names, ","), "taxa_id,Kingdom,Phylum,Class,Orders,Family,Genus,Species"))
	}
	if len(s.taxa) > 0 {
		ret = true
		for k := range s.taxa {
			s.taxaids = append(s.taxaids, k)
		}
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
			fmt.Print("\n\t[Error] Please enter a valid taxonomic level. Exiting.\n\n")
			os.Exit(11)
		}
		s.column = level
	}
}

func SearchTaxonomicLevels(db *dbIO.DBIO, names []string, user, level string, count, com, inf bool) ([][]string, string) {
	// Extracts data using species names
	s := newSearcher(db, []string{"Taxonomy"}, user, level, "=", "", com, inf)
	s.checkLevel(level)
	fmt.Printf("\tExtracting patient information from %s...\n", s.column)
	pass := s.setTaxonomy(names)
	if pass == true {
		s.getTaxa()
		if len(s.res) >= 1 {
			if s.infant == false {
				s.filterInfantRecords()
			}
			if count == false {
				// Skip if not needed since this is the most time consuming step
				s.appendDiagnosis()
				s.appendTaxonomy()
				s.appendSource()
			}
		}
		//fmt.Println(s.res)
	}
	return s.toSlice(), s.header
}

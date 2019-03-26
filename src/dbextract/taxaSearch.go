// This script will search for given records from the comparative oncology database

package dbextract

import (
	"bytes"
	"fmt"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"os"
	"strings"
)

func toTitle(names []string) []string {
	// Converts all input names to title case
	var ret []string
	for _, i := range names {
		ret = append(ret, dbupload.TitleCase(i))
	}
	return ret
}

func (s *searcher) getTaxa() {
	// Extracts patient data using taxa ids
	s.res = s.db.GetRows("Patient", "taxa_id", strings.Join(s.taxaids, ","), "*")
	s.setIDs()
}

func (s *searcher) getTaxonomy(names []string, ids bool) map[string][]string {
	// Stores taxa ids from species name and returns taxonomy
	ret := make(map[string][]string)
	var table [][]string
	if s.common == true {
		// Get taxonomy ids from common name list
		c := s.db.GetRows("Common", "Name", strings.Join(toTitle(names), ","), "*")
		if len(c) >= 1 {
			// Colect taxa IDs
			buffer := bytes.NewBufferString(c[0][0])
			for _, i := range c[1:] {
				buffer.WriteByte(',')
				buffer.WriteString(i[0])
			}
			// Get taxonomy entries
			table = s.db.GetRows("Taxonomy", "taxa_id", buffer.String(), "*")
		}
	} else if ids == false {
		// Get matching taxonomies
		table = s.db.GetRows("Taxonomy", s.column, strings.Join(names, ","), "*")
	} else {
		table = s.db.GetRows("Taxonomy", "taxa_id", strings.Join(names, ","), "*")
	}
	for _, row := range table {
		if ids == false {
			// Append taxa id and return map of taxonomy entries
			s.taxaids = append(s.taxaids, row[0])
		}
		// Exclude taxa id and source
		ret[row[0]] = row[1:8]
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
	taxonomy := s.getTaxonomy(names, false)
	if len(taxonomy) >= 1 {
		s.getTaxa()
		fmt.Println(s.res)
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
	}
	return s.res, s.header
}

// This script will summarize and upload the taxonomy table for the comparative oncology database

package dbupload

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"log"
	"strconv"
	"strings"
)

func getTaxon(genus, species string) string {
	// Returns lowest taxon present
	var ret string
	species = strings.TrimSpace(species)
	genus = strings.TrimSpace(genus)
	if len(species) > 1 && strings.ToUpper(species) != "NA" {
		ret = speciesCaps(species)
	} else if len(genus) > 1 && strings.ToUpper(genus) != "NA" {
		ret = strings.Title(genus)
	}
	return ret
}

func getSource(sources []string) string {
	// Get first returned source
	source := "NA"
	for _, i := range sources {
		if i != "NA" && len(i) >= 5 {
			// Assumes at least "http:"
			source = i
			break
		}
	}
	return source
}

func speciesCaps(species string) string {
	// Returns properly capitalized species name
	var ret string
	s := strings.Split(strings.ToLower(species), " ")
	if len(s) > 1 {
		// Save with genus capitalized and species in lower case
		ret = strings.Title(s[0]) + " " + s[1]
		if len(s) > 2 && ret == "Canis lupus" && strings.TrimSpace(s[2]) == "familiaris" {
			ret += " " + s[2]
		}
	} else {
		ret = "NA"
	}
	return ret
}

func checkCaps(taxonomy []string) []string {
	// Returns slice with proper capitization
	l := len(taxonomy) - 1
	for idx, i := range taxonomy {
		if idx < l {
			taxonomy[idx] = strings.Title(i)
		} else {
			taxonomy[idx] = speciesCaps(i)
		}
	}
	return taxonomy
}

func GetTaxaIDs(db *dbIO.DBIO, commonNames bool) map[string]string {
	// Returns map of taxa ids corresponding to common names, binomial names, or genus
	ret := make(map[string]string)
	taxa := db.GetColumns("Taxonomy", []string{"taxa_id", "Genus", "Species"})
	for _, i := range taxa {
		key := getTaxon(i[1], i[2])
		if len(key) >= 1 {
			ret[key] = i[0]
		}
	}
	if commonNames == true {
		for _, i := range db.GetTable("Common") {
			ret[i[1]] = i[0]
		}
	}
	return ret
}

//----------------------------------------------------------------------------

type taxa struct {
	col         map[string]int
	common      map[string][][]string
	commonNames bool
	count       int
	db          *dbIO.DBIO
	logger      *log.Logger
	neu         map[string][]string
	taxaids     map[string]string
}

func newTaxa(db *dbIO.DBIO, common bool) *taxa {
	// Returns new taxonomy struct
	t := new(taxa)
	t.db = db
	t.common = make(map[string][][]string)
	t.commonNames = common
	t.count = t.db.GetMax("Taxonomy", "taxa_id") + 1
	t.logger = codbutils.GetLogger()
	t.taxaids = GetTaxaIDs(t.db, t.commonNames)
	t.neu = make(map[string][]string)
	return t
}

func (t *taxa) uploadTable() {
	// Uploads table to database
	var com [][]string
	for _, v := range t.common {
		// Convert common names map to slice
		com = append(com, v...)
	}
	if len(t.neu) > 0 {
		vals, l := dbIO.FormatMap(t.neu)
		t.db.UpdateDB("Taxonomy", vals, l)
	}
	if len(com) > 0 {
		vals, l := dbIO.FormatSlice(com)
		t.db.UpdateDB("Common", vals, l)
	}
}

func (t *taxa) getTaxaID(s string) string {
	// Returns taxa id for given scientific name
	if _, ex := t.neu[s]; ex == true {
		return t.neu[s][0]
	} else if _, ex := t.taxaids[s]; ex == true {
		return t.taxaids[s]
	} else {
		return "NA"
	}
}

func (t *taxa) setCommon(spl []string, c, s string, cur bool) {
	// Adds common name entry to map
	if _, ex := t.neu[c]; ex == false {
		// Make sure entry in common name column is not a scientific name
		if _, ex := t.taxaids[c]; ex == false {
			// Add unique common name entries to slice
			curator := "NA"
			if cur == true {
				// Store curator name
				curator = spl[t.col["Name"]]
			}
			id := t.getTaxaID(s)
			if id != "NA" {
				row, ex := t.common[s]
				if ex == true && strarray.InSliceSli(row, c, 0) == false {
					// Add to existing species record
					t.common[s] = append(t.common[s], []string{id, c, curator})
				} else {
					t.common[s] = append(t.common[s], []string{id, c, curator})
				}
			}
		}
	}
}

func (t *taxa) setTaxon(spl []string, s string, l int, cur bool) {
	// Adds taxonomy to map
	if _, ex := t.taxaids[s]; ex == false {
		// Skip entries which are already in db
		if _, ex := t.neu[s]; ex == false {
			// Add unique taxonomies
			id := strconv.Itoa(t.count)
			taxonomy := append([]string{id}, checkCaps(spl[t.col["Kingdom"]:t.col["Species"]+1])...)
			taxonomy = append(taxonomy, spl[t.col["SearchTerm"]])
			if cur == true {
				taxonomy = append(taxonomy, getSource(spl[t.col["Species"]+1:t.col["Name"]]))
			} else {
				taxonomy = append(taxonomy, getSource(spl[t.col["Species"]+1:l]))
			}
			t.neu[s] = taxonomy
			t.count++
		}
	}
}

func (t *taxa) extractTaxa(infile string) {
	// Extracts taxonomy from input file
	var rows [][]string
	cur := true
	t.logger.Printf("Extracting taxa from %s\n", infile)
	rows, t.col = iotools.ReadFile(infile, true)
	l := len(t.col)
	if _, ex := t.col["Name"]; ex == false {
		// Set dummy currator column
		t.col["Name"] = l + 1
		cur = false
	}
	for _, i := range rows {
		c := strarray.TitleCase(i[t.col["SearchTerm"]])
		s := getTaxon(i[t.col["Genus"]], i[t.col["Species"]])
		t.setTaxon(i, s, l, cur)
		if t.commonNames == true {
			t.setCommon(i, c, s, cur)
		}
	}
}

func LoadTaxa(db *dbIO.DBIO, infile string, commonNames bool) {
	// Loads unique entries into comparative oncology taxonomy table
	t := newTaxa(db, commonNames)
	t.extractTaxa(infile)
	t.uploadTable()
}

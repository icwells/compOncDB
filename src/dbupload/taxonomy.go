// This script will summarize and upload the taxonomy
//table for the comparative oncology database

package dbupload

import (
	"bufio"
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"strconv"
	"strings"
)

func uploadTable(db *dbIO.DBIO, taxa map[string][]string, common map[string][][]string, count int) {
	// Uploads table to database
	var com [][]string
	for k, v := range taxa {
		// Add unique taxa ID
		count++
		c := strconv.Itoa(count)
		taxa[k] = append([]string{c}, v...)
		row, ex := common[k]
		if ex == true {
			// Join common names to taxa id in paired entries
			for _, n := range row {
				com = append(com, []string{c, n[0], n[1]})
			}
		}
	}
	if len(taxa) > 0 {
		vals, l := dbIO.FormatMap(taxa)
		db.UpdateDB("Taxonomy", vals, l)
	}
	if len(com) > 0 {
		vals, l := dbIO.FormatSlice(com)
		db.UpdateDB("Common", vals, l)
	}
}

func getTaxon(family, genus, species string) string {
	// Returns lowest taxon present
	var ret string
	if species != "NA" {
		s := strings.Split(strings.ToLower(species), " ")
		if len(s) > 1 {
			// Save with genus capitalized and species in lower case
			ret = strings.Title(s[0]) + " " + s[1]
		} else {
			ret = strings.Title(species)
		}
	} else if genus != "NA" {
		ret = strings.Title(species)
	} else if family != "NA" {
		ret = strings.Title(family)
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

func extractTaxa(infile string, taxaids map[string]string, commonNames bool) (map[string][]string, map[string][][]string) {
	// Extracts taxonomy from input file
	var col map[string]int
	var l int
	first := true
	cur := true
	taxa := make(map[string][]string)
	common := make(map[string][][]string)
	fmt.Printf("\n\tExtracting taxa from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := strings.TrimSpace(string(input.Text()))
		spl := strings.Split(line, ",")
		if first == false {
			c := strings.Title(spl[col["SearchTerm"]])
			s := getTaxon(spl[col["Family"]], spl[col["Genus"]], spl[col["Species"]])
			if _, ex := taxaids[s]; ex == false {
				// Skip entries which are already in db
				if _, ex := taxa[s]; ex == false {
					// Add unique taxonomies
					taxonomy := spl[col["Kingdom"] : col["Species"]+1]
					if cur == true {
						taxonomy = append(taxonomy, getSource(spl[col["Species"]+1:col["Name"]]))
					} else {
						taxonomy = append(taxonomy, getSource(spl[col["Species"]+1:l]))
					}
					taxa[s] = taxonomy
				}
			}
			if commonNames == true {
				if _, ex := taxaids[c]; ex == false {
					// Add unique common name entries to slice
					curator := "NA"
					if cur == true {
						// Store curator name
						curator = spl[col["Name"]]
					}
					row, ex := common[s]
					if ex == true {
						// Add to existing species record
						if strarray.InSliceSli(row, c, 0) == false {
							common[s] = append(common[s], []string{c, curator})
						}
					} else {
						common[s] = append(common[s], []string{c, curator})
					}
				}
			}
		} else {
			col = getColumns(spl)
			l = len(spl)
			if _, ex := col["Name"]; ex == false {
				col["Name"] = l + 1
				cur = false
			}
			first = false
		}
	}
	return taxa, common
}

func getTaxaIDs(db *dbIO.DBIO, commonNames bool) map[string]string {
	// Returns map of taxa ids corresponding to common names, binomial names, genus, or family (whichever is most descriptive)
	ret := make(map[string]string)
	taxa := db.GetColumns("Taxonomy", []string{"taxa_id", "Family", "Genus", "Species"})
	for _, i := range taxa {
		key := getTaxon(i[1], i[2], i[3])
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

func LoadTaxa(db *dbIO.DBIO, infile string, commonNames bool) {
	// Loads unique entries into comparative oncology taxonomy table
	var taxa map[string][]string
	var common map[string][][]string
	m := db.GetMax("Taxonomy", "taxa_id")
	taxaids := getTaxaIDs(db, commonNames)
	taxa, common = extractTaxa(infile, taxaids, commonNames)
	uploadTable(db, taxa, common, m)
}

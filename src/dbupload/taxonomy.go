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

func extractTaxa(infile string, species, com []string, commonNames bool) (map[string][]string, map[string][][]string) {
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
			s := spl[col["Species"]]
			if strarray.InSliceStr(species, s) == false {
				// Skip entries which are already in db
				if _, ex := taxa[s]; ex == false {
					// Add unique taxonomies
					taxonomy := spl[col["Kingdom"] : col["Species"]+1]
					// Get first returned source
					sources := spl[col["Species"]+1 : col["Name"]]
					source := "NA"
					for _, i := range sources {
						if i != "NA" && len(i) >= 5 {
							// Assumes at least "http:"
							source = i
							break
						}
					}
					taxonomy = append(taxonomy, source)
					taxa[s] = taxonomy
				}
			}
			if commonNames == true && strarray.InSliceStr(com, c) == false {
				// Add unique common name entries to slice
				curator := "NA"
				if cur == true {
					// Store curator name
					curator = spl[col["Name"]]
				}
				row, ex := common[s]
				if ex == true {
					if strarray.InSliceSli(row, c, 0) == false {
						common[s] = append(common[s], []string{c, curator})
					}
				} else {
					common[s] = append(common[s], []string{c, curator})
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

func LoadTaxa(db *dbIO.DBIO, infile string, commonNames bool) {
	// Loads unique entries into comparative oncology taxonomy table
	var taxa map[string][]string
	var common map[string][][]string
	m := db.GetMax("Taxonomy", "taxa_id")
	species := db.GetColumnText("Taxonomy", "Species")
	com := db.GetColumnText("Common", "Name")
	taxa, common = extractTaxa(infile, species, com, commonNames)
	uploadTable(db, taxa, common, m)
}

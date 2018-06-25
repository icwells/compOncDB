// This script will summarize and upload the taxonomy
//table for the comparative oncology database

package main

import (
	"bufio"
	"database/sql"
	"dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"strconv"
	"strings"
)

func uploadTable(db *sql.DB, col map[string]string, taxa, common map[string][]string, count int) {
	// Uploads table to database
	var com [][]string
	for k, v := range taxa {
		// Add unique taxa ID
		taxa[k] = append([]string{strconv.Itoa(count)}, v...)
		if strarray.InMapSli(common, k) == true {
			// Join common names to taxa id in paired entries
			for _, n := range common[k] {
				com = append(com, []string{string(count), n})
			}
		}
		count++
	}
	tvals := dbIO.FormatMap(taxa)
	dbIO.UpdateDB(db, "Taxonomy", col["Taxonomy"], tvals)
	cvals := dbIO.FormatSlice(com)
	dbIO.UpdateDB(db, "Common", col["Common"], cvals)
}

func extractTaxa(infile string, species, com []string) (map[string][]string, map[string][]string) {
	// Extracts taxonomy from input file
	first := true
	taxa := make(map[string][]string)
	common := make(map[string][]string)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			spl := strings.Split(line, ",")
			c := spl[1]
			s := spl[8]
			if strarray.InSliceStr(species, s) == false {
				// Skip entries which are already in db
				if strarray.InMapSli(taxa, s) == false {
					// Add unique taxonomies
					taxonomy := spl[2:9]
					// Get first returned source
					sources := spl[9:]
					for _, i := range sources {
						if i != "NA" && len(i) >= 5 {
							// Assumes at least "http:"
							taxonomy = append(taxonomy, i)
							break
						}
					}
					taxa[s] = taxonomy
				}
			}
			if strarray.InSliceStr(com, c) == false {
				// Add unique common name entries to slice
				if strarray.InMapSli(common, s) == true {
					if strarray.InSliceStr(common[s], c) == false {
						common[s] = append(common[s], c)
					}
				} else {
					common[s] = append(common[s], c)
				}
			}
		} else {
			first = false
		}
	}
	return taxa, common
}

func LoadTaxa(db *sql.DB, col map[string]string, infile string) {
	// Loads unique entries into comparative oncology taxonomy table
	var taxa, common map[string][]string
	m := dbIO.GetMax(db, "Taxonomy", "taxa_id")
	species := dbIO.GetColumnText(db, "Taxonomy", "Species")
	com := dbIO.GetColumnText(db, "Common", "Name")
	taxa, common = extractTaxa(infile, species, com)
	uploadTable(db, col, taxa, common, m)
}

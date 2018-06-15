// This script will summarize and upload the taxonomy
//table for the comparative oncology database

package main

import (
	"bufio"
	"database/sql"
	"dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"strings"
)

func uploadTable(db *sql.DB, col map[string]string, taxa, common map[string][]string) {
	// Uploads table to database
	var com [][]string
	count = 1
	l := len(data)
	for k, v := range taxa {
		// Add unique taxa ID
		taxa[k] = append(string(count), v)
		if strarray.InMapSli(common, k) == true {
			// Join common names to taxa id in paired entries
			for _, n := range common[k] {
				com = append(com, []string{string(count), n})
			}
		}
	}
	tvals := dbIO.FormatMap(taxa)
	UpdateDB(db, "Taxonomy", columns["Taxonomy"], tvals)
	cvals := dbIO.FormatSlice(com)
	UpdateDB(db, "Common", columns["Common"], cvals)
}

func extractTaxa(taxa, common map[string][]string, infile string) (map[string][]string, map[string][]string) {
	// Extracts taxonomy from input file
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			spl := strings.Split(line, ",")
			c := spl[1]
			s := spl[8]
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
			// Add unique common name entries to slice
			if strarray.InMapSli(common, s) == true {
				if strarray.InSliceStr(common[s], c) == false {
					common[s] = append(common[s], c)
				}
			} else {
				var common [s][]string
				common[s] = append(common[s], c)
			}
		} else {
			first = false
		}
	}
	return taxa, common
}

func LoadTaxa(db *BD, col map[string]string, nwzp, zeps, msu string) {
	// Loads unique entries into comparative oncology taxaonomy table
	var taxa, common map[string][]string
	taxa, common = extractTaxa(taxa, common, nwzp)
	taxa, common = extractTaxa(taxa, common, zeps)
	taxa, common = extractTaxa(taxa, common, msu)
	uploadTables(db, col, taxa, common)
}

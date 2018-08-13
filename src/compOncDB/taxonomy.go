// This script will summarize and upload the taxonomy
//table for the comparative oncology database

package main

import (
	"bufio"
	"database/sql"
	"dbIO"
	"fmt"
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
		count++
		c := strconv.Itoa(count)
		taxa[k] = append([]string{c}, v...)
		if strarray.InMapSli(common, k) == true {
			// Join common names to taxa id in paired entries
			for _, n := range common[k] {
				com = append(com, []string{c, n})
			}
		}
	}
	if len(taxa) > 0 {
		vals, l := dbIO.FormatMap(taxa)
		dbIO.UpdateDB(db, "Taxonomy", col["Taxonomy"], vals, l)
	}
	if len(com) > 0 {
		vals, l := dbIO.FormatSlice(com)
		dbIO.UpdateDB(db, "Common", col["Common"], vals, l)
	}
}

func extractTaxa(infile string, species, com []string, commonNames bool) (map[string][]string, map[string][]string) {
	// Extracts taxonomy from input file
	first := true
	taxa := make(map[string][]string)
	common := make(map[string][]string)
	fmt.Printf("\n\tExtracting taxa from %s\n", infile)
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
			if commonNames == true && strarray.InSliceStr(com, c) == false {
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

func LoadTaxa(db *sql.DB, col map[string]string, infile string, commonNames bool) {
	// Loads unique entries into comparative oncology taxonomy table
	var taxa, common map[string][]string
	m := dbIO.GetMax(db, "Taxonomy", "taxa_id")
	species := dbIO.GetColumnText(db, "Taxonomy", "Species")
	com := dbIO.GetColumnText(db, "Common", "Name")
	taxa, common = extractTaxa(infile, species, com, commonNames)
	uploadTable(db, col, taxa, common, m)
}

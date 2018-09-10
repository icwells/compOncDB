// This script will summarize and upload the life history
//table for the comparative oncology database

package main

import (
	"bufio"
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"strings"
)

func uploadTraits(db *sql.DB, col map[string]string, traits [][]string) {
	// Uploads table to database
	if len(traits) > 0 {
		vals, l := dbIO.FormatSlice(traits)
		dbIO.UpdateDB(db, "Life_history", col["Life_history"], vals, l)
	}
}

func fmtEntry(tid string, row []string) []string {
	// Returns row formatted for upload with NAs replaced with 0.0
	entry := []string{tid}
	for _, i := range row[1:14] {
		if i == "NA" {
			entry = append(entry, "0.0")
		} else {
			entry = append(entry, i)
		}
	}
	return entry
}

func extractTraits(infile string, ids []string, species map[string]string) [][]string {
	// Extracts taxonomy from input file
	first := true
	var traits [][]string
	fmt.Printf("\n\tExtracting life history data from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			line = strings.Trim(line, "\n\t ")
			spl := strings.Split(line, ",")
			s := strings.Trim(spl[0], "\t ")
			if strarray.InMapStr(species, s) == true {
				// Get taxa id from species name
				tid := species[s]
				if strarray.InSliceStr(ids, tid) == false {
					// Skip entries which are already in db
					traits = append(traits, fmtEntry(tid, spl))
				}
			} else {
				fmt.Printf("\t[Warning] %s not in taxonomy database. Skipping.\n", s)
			}
		} else {
			first = false
		}
	}
	return traits
}

func loadLifeHistory(db *sql.DB, col map[string]string, infile string) {
	// Loads unique entries into comparative oncology taxonomy table
	species := entryMap(dbIO.GetColumns(db, "Taxonomy", []string{"taxa_id", "Species"}))
	ids := dbIO.GetColumnText(db, "Life_history", "taxa_id")
	traits := extractTraits(infile, ids, species)
	uploadTraits(db, col, traits)
}

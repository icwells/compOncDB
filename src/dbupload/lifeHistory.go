// This script will summarize and upload the life history
//table for the comparative oncology database

package dbupload

import (
	"bufio"
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"strings"
)

func uploadTraits(db *dbIO.DBIO, traits [][]string) {
	// Uploads table to database
	if len(traits) > 0 {
		vals, l := dbIO.FormatSlice(traits)
		db.UpdateDB("Life_history", vals, l)
	}
}

func fmtEntry(tid string, row []string) []string {
	// Returns row formatted for upload with NAs replaced with -1.0
	entry := []string{tid}
	for _, i := range row[1:14] {
		if i == "NA" {
			entry = append(entry, "-1.0")
		} else {
			entry = append(entry, i)
		}
	}
	return entry
}

func extractTraits(infile string, ids []string, species map[string]string) [][]string {
	// Extracts taxonomy from input file
	missed := 0
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
				missed++
			}
		} else {
			first = false
		}
	}
	if missed > 0 {
		fmt.Printf("\t[Warning] %d records not in taxonomy database.\n", missed)
	}
	return traits
}

func LoadLifeHistory(db *dbIO.DBIO, infile string) {
	// Loads unique entries into comparative oncology taxonomy table
	species := EntryMap(db.GetColumns("Taxonomy", []string{"taxa_id", "Species"}))
	ids := db.GetColumnText("Life_history", "taxa_id")
	traits := extractTraits(infile, ids, species)
	uploadTraits(db, traits)
}

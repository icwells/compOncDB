// This script will summarize and upload the taxonomy 
 //table for the comparative oncology database

package preupload

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"os"
	"strconv"
	"strings"
)

func uploadTable(db *DB, common map[string][string]) {
	// Uploads table to database
	

}

func extractTaxa(taxa, common map[string][string], infile string) (map[string][string], map[string][string]) {
	// Extracts taxonomy from input file
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			spl := strings.Split(line, ",")
			/*if strarray.InSliceStr[] == true {
				taxa[spl[]] = strings.Join(spl[2:7], ",")
				// Get most applicable source

			} */
		} else {
			first = false
		}
	}
	return taxa, common
}

func LoadTaxa(db *BD, nwzp, zeps, msu string) {
	// Loads unique entries into comparative oncology taxaonomy table
	var taxa, common map[string][string]
	taxa, common := extractTaxa(taxa, common, nwzp)
	taxa, common = extractTaxa(taxa, common, zeps)
	taxa, common = extractTaxa(taxa, common, msu)
	uploadTable(db, taxa)
	uploadTable(db, common)
}

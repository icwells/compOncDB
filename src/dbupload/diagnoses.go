// This script will upload unique tumor and metastasis data to the comparative oncology database

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

func uploadDiagnosis(db *dbIO.DBIO, tumor map[string][]string, t int) {
	// Uploads unique tumor and metastasis entries with random ID number
	var tmr [][]string
	// Convert tumor map to slice
	for k, v := range tumor {
		if k != "NA" {
			for _, i := range v {
				// Add unique taxa ID and append if tumor type is present
				t++
				c := strconv.Itoa(t)
				tmr = append(tmr, []string{c, k, i})
			}
		}
	}
	if len(tmr) > 0 {
		vals, l := dbIO.FormatSlice(tmr)
		db.UpdateDB("Tumor", vals, l)
	}
}

func tumorPairs(typ, loc string) [][]string {
	// Returns slice of pairs of type, location
	var ret [][]string
	types := strings.Split(typ, ";")
	locations := strings.Split(loc, ";")
	for idx, i := range types {
		if idx < len(locations) {
			ret = append(ret, []string{strings.TrimSpace(i), strings.TrimSpace(locations[idx])})
		}
	}
	return ret
}

func extractDiagnosis(infile string, tmr map[string]map[string]string) map[string][]string {
	// Extracts accounts from input file
	first := true
	tumor := make(map[string][]string)
	fmt.Printf("\n\tExtracting diagnosis data from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		s := strings.Split(line, ",")
		if first == false && len(s) == 17 {
			// Iterate through type, location pairs individually
			pairs := tumorPairs(s[10], s[11])
			for _, i := range pairs {
				_, intmr := tmr[i[0]]
				if intmr == false || intmr == true && strarray.InMapStr(tmr[i[0]], i[1]) == false {
					// Skip entries present in database or already in map
					if strarray.InMapSli(tumor, i[0]) == true && strarray.InSliceStr(tumor[i[0]], i[1]) == false {
						// Add new location info
						tumor[i[0]] = append(tumor[i[0]], i[1])
					} else {
						// Add new list
						tumor[i[0]] = []string{i[1]}
					}
				}
			}
		} else {
			first = false
		}
	}
	return tumor
}

func LoadDiagnoses(db *dbIO.DBIO, infile string) {
	// Loads unique entries into comparative oncology metastatis, tumor, and account tables
	t := db.GetMax("Tumor", "tumor_id")
	tmr := mapOfMaps(db.GetTable("Tumor"))
	tumor := extractDiagnosis(infile, tmr)
	uploadDiagnosis(db, tumor, t)
}

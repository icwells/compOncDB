// This script will upload unique tumor and metastasis data to the comparative oncology database

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

func uploadDiagnosis(db *sql.DB, col map[string]string, tumor map[string][]string, meta []string, t, m int) {
	// Uploads unique tumor and metastasis entries with random ID number
	var mts, tmr [][]string
	// Convert tumor map to slice
	for k, v := range tumor {
		for _, i := range v {
			// Add unique taxa ID
			t++
			c := strconv.Itoa(t)
			tmr = append(tmr, []string{c, k, i})
		}
	}
	if len(tmr) > 0 {
		vals, l := dbIO.FormatSlice(tmr)
		dbIO.UpdateDB(db, "Tumor", col["Tumor"], vals, l)
	}
	// Add ids to metastasis data
	for _, i := range meta {
		m++
		c := strconv.Itoa(m)
		mts = append(mts, []string{c, i})
	}
	if len(mts) > 0 {
		vals, l := dbIO.FormatSlice(mts)
		dbIO.UpdateDB(db, "Metastasis", col["Metastasis"], vals, l)
	}
}

func extractDiagnosis(infile string, tmr map[string]map[string]string, mts []string) (map[string][]string, []string) {
	// Extracts accounts from input file
	first := true
	tumor := make(map[string][]string)
	var meta []string
	fmt.Printf("\n\tExtracting accounts from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		s := strings.Split(line, ",")
		if first == false && len(s) == 15 {
			// Determine if entry is unique
			if strarray.InSliceStr(mts, s[8]) == false {
				meta = append(meta, s[8])
			}
			intmr := strarray.InMapMapStr(tmr, s[9])
			if intmr == false || intmr == true && strarray.InMapStr(tmr[s[9]], s[10]) == false {
				if strarray.InMapSli(tumor, s[9]) == true {
					tumor[s[9]] = append(tumor[s[9]], s[10])
				} else {
					tumor[s[9]] = []string{s[10]}
				}
			}
		} else {
			first = false
		}
	}
	return tumor, meta
}

func LoadDiagnoses(db *sql.DB, col map[string]string, infile string) {
	// Loads unique entries into comparative oncology metastatis, tumor, and account tables
	t := dbIO.GetMax(db, "Tumor", "tumor_id")
	m := dbIO.GetMax(db, "Metastasis", "metastasis_id")
	tmr := mapOfMaps(dbIO.GetTable(db, "Tumor"))
	mts := dbIO.GetColumnText(db, "Metastasis", "location")
	tumor, meta := extractDiagnosis(infile, tmr, mts)
	uploadDiagnosis(db, col, tumor, meta, t, m)
}

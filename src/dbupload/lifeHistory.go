// This script will summarize and upload the life history table for the comparative oncology database

package dbupload

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"strconv"
	"strings"
)

// Average proportion of weaning/max_longevity (determined in ageOfInfancy script)
var PROP = 0.032964

func uploadTraits(db *dbIO.DBIO, traits [][]string) {
	// Uploads table to database
	if len(traits) > 0 {
		vals, l := dbIO.FormatSlice(traits)
		db.UpdateDB("Life_history", vals, l)
	}
}

func getAvgMaturity(male, female string) string {
	// Returns age of infancy from average of male and female maturity
	var ret float64
	m, er := strconv.ParseFloat(male, 64)
	f, err := strconv.ParseFloat(female, 64)
	if err != nil && er != nil {
		ret = 1.0
	} else if er != nil && err == nil {
		ret = f
	} else if err != nil {
		ret = m
	//} else {
		//ret = (((f + m) / 2) * 0.1)
	}
	return strconv.FormatFloat(ret, 'f', -1, 64)
}

func calculateInfancy(weaning, male, female, longevity string) string {
	// Returns age for infancy column
	var ret string
	w, err := strconv.ParseFloat(weaning, 64)
	if err == nil && w >= 0.0 {
		// Assign weaning age
		ret = weaning
	} else {
		ret = getAvgMaturity(male, female)
	}
	if ret == "-1" && longevity != "NA" {
		l, _ := strconv.ParseFloat(longevity, 64)
		ret = strconv.FormatFloat(PROP * longevity, 'f', -1, 64)
	}
	return ret
}

func fmtNA(val string) string {
	// Converts NAs to -1.0
	if val == "NA" {
		val = "-1.0"
	}
	return val
}

func fmtEntry(col map[string]int, l int, tid string, row []string) []string {
	// Returns row formatted for upload with NAs replaced with -1.0; skips source column; adds age of infancy
	entry := []string{tid}
	for _, i := range row[col["FemaleMaturity"] : col["Weaning"]+1] {
		entry = append(entry, fmtNA(i))
	}
	entry = append(entry, calculateInfancy(row[col["Weaning"]], row[col["MaleMaturity"]], row[col["FemaleMaturity"]], row[col["MaximumLongevity"]]))
	for _, i := range row[col["Litter"]:col["Source"]] {
		entry = append(entry, fmtNA(i))
	}
	return entry
}

func getColumnIndeces(head []string) map[string]int {
	// Returns map of column indeces by name
	ret := make(map[string]int)
	for idx, i := range head {
		if strings.Contains(i, "/") == true {
			i = i[:strings.Index(i, "/")]
		}
		if strings.Contains(i, "(") == true {
			i = i[:strings.Index(i, "(")]
		}
		ret[i] = idx
	}
	return ret
}

func extractTraits(infile string, ids []string, species map[string]string) [][]string {
	// Extracts taxonomy from input file
	var l int
	var traits [][]string
	var col map[string]int
	missed := 0
	first := true
	logger := codbutils.GetLogger()
	logger.Printf("Extracting life history data from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := iotools.GetScanner(f)
	for input.Scan() {
		line := strings.TrimSpace(string(input.Text()))
		spl := strings.Split(line, ",")
		if first == false {
			s := strings.TrimSpace(spl[0])
			// Get taxa id from species name
			tid, ex := species[s]
			if ex == true {
				if strarray.InSliceStr(ids, tid) == false {
					// Skip entries which are already in db
					traits = append(traits, fmtEntry(col, l, tid, spl))
				}
			} else {
				missed++
			}
		} else {
			l = len(spl)
			col = getColumnIndeces(spl)
			first = false
		}
	}
	if missed > 0 {
		logger.Printf("[Warning] %d records not in taxonomy database.\n", missed)
	}
	return traits
}

func LoadLifeHistory(db *dbIO.DBIO, infile string) {
	// Loads unique entries into comparative oncology taxonomy table
	species := codbutils.EntryMap(db.GetColumns("Taxonomy", []string{"taxa_id", "Species"}))
	ids := db.GetColumnText("Life_history", "taxa_id")
	traits := extractTraits(infile, ids, species)
	uploadTraits(db, traits)
}

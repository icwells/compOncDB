// This script will upload patient data to the comparative oncology database

package main

import (
	"bufio"
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"math"
	"strconv"
	"strings"
)

type Entries struct {
	p	[][]string
	d	[][]string
	t	[][]string
	s	[][]string
}

func (e *Entries) update(p, d, t, s []string) {
	// Appends new entries to appropriate slice
	e.p = append(e.p, p)
	e.d = append(e.d, d)
	e.t = append(e.t, t)
	e.s = append(e.s, s)
}

func getDenominator(l, row int) int {
	// Returns denominator for subsetting upload slice
	p := float64(l * row)
	max := 200000.0
	return int(math.Floor(p / max))
}

func uploadPatients(db *sql.DB, table string, col map[string]string, list [][]string) {
	// Uploads patient entries to db
	l := len(list)
	den := getDenominator(l, len(list[0]))
	if den <= 1 {
		// Upload slice at once
		vals, l := dbIO.FormatSlice(list)
		dbIO.UpdateDB(db, table, col[table], vals, l)
	} else {
		// Upload in chunks
		var set [][][]string
		idx := l / den
		ind := 0
		for i := 0; i < den; i++ {
			if ind+idx > l {
				// Get last less than idx rows
				idx = l - ind + 1
			}
			sub := list[ind : ind+idx]
			set = append(set, sub)
			ind = ind + idx
		}
		for _, i := range set {
			vals, ln := dbIO.FormatSlice(i)
			dbIO.UpdateDB(db, table, col[table], vals, ln)
		}
	}
}

func extractPatients(infile string, count int, tumor, acc map[string]map[string]string, meta, species map[string]string) Entries {
	// Assigns patient data to appropriate slices with unique entry IDs
	missed := 0
	first := true
	start := count
	var entries Entries
	fmt.Printf("\n\tExtracting patient data from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			pass := false
			spl := strings.Split(line, ",")
			if strings.ToUpper(spl[4]) != "NA" {
				if len(spl) == 17 && strarray.InMapStr(species, spl[4]) == true && strarray.InMapMapStr(acc, spl[15]) == true {
					// Skip entries without valid species and source data
					if strarray.InMapStr(acc[spl[15]], spl[16]) == true {
						var d, t []string
						count++
						id := strconv.Itoa(count)
						if strings.Contains(spl[3], "NA") == true {
							// Make sure source ID is not NA
							spl[3] = "-1"
						} else if len(spl[3]) > 6 {
							// Make sure age does not exceed decimal precision
							spl[3] = spl[3][:7]
						}
						// ID, Sex, Age, Castrated, taxa_id, source_id, Species, Date, Comments
						p := []string{id, spl[0], spl[1], spl[2], species[spl[4]], spl[3], spl[4], spl[5], spl[6]}
						// ID, service, account_id
						s := []string{id, spl[14], acc[spl[15]][spl[16]]}
						// Diagnosis entry
						if strarray.InMapStr(meta, spl[9]) == true {
							// ID, masspresent, necropsy, metastasis_id
							d = []string{id, spl[7], spl[8], meta[spl[9]]}
						} else {
							d = []string{id, spl[7], spl[8], "-1"}
						}
						// Assign ID to all tumor, location pairs tumorPairs (in diagnoses.go)
						pairs := tumorPairs(spl[10], spl[11])
						for _, i := range pairs {
							if strarray.InMapMapStr(tumor, i[0]) == true && strarray.InMapStr(tumor[i[0]], i[1]) == true {
								// ID, tumor_id, primary_tumor, malignant
								t = []string{id, tumor[i[0]][i[1]], spl[12], spl[13]}
							} else {
								t = []string{id, "-1", spl[12], spl[13]}
							}
						}
						entries.update(p, d, t, s)
						pass = true
					}
				}
				if pass == false {
					missed++
				}
			}
		} else {
			first = false
		}
	}
	fmt.Printf("\tExtracted %d records.\n", count-start)
	if missed > 0 {
		fmt.Printf("\t[Warning] Count not find taxa ID or source ID for %d records.\n", missed)
	}
	return entries
}

func loadPatients(db *sql.DB, col map[string]string, infile string) {
	// Loads unique patient info to appropriate tables
	m := dbIO.GetMax(db, "Patient", "ID")
	tumor := mapOfMaps(dbIO.GetTable(db, "Tumor"))
	acc := mapOfMaps(dbIO.GetTable(db, "Accounts"))
	meta := entryMap(dbIO.GetTable(db, "Metastasis"))
	species := entryMap(dbIO.GetColumns(db, "Taxonomy", []string{"taxa_id", "Species"}))
	// Get entry slices and upload to db
	entries := extractPatients(infile, m, tumor, acc, meta, species)
	uploadPatients(db, "Patient", col, entries.p)
	uploadPatients(db, "Diagnosis", col, entries.d)
	uploadPatients(db, "Tumor_relation", col, entries.t)
	uploadPatients(db, "Source", col, entries.s)
}

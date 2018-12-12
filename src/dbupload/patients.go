// This script will upload patient data to the comparative oncology database

package dbupload

import (
	"bufio"
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"math"
	"strconv"
	"strings"
)

type Entries struct {
	p [][]string
	d [][]string
	t [][]string
	s [][]string
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

func uploadPatients(db *dbIO.DBIO, table string, list [][]string) {
	// Uploads patient entries to db
	l := len(list)
	den := getDenominator(l, len(list[0]))
	if den <= 1 {
		// Upload slice at once
		vals, l := dbIO.FormatSlice(list)
		db.UpdateDB(table, vals, l)
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
			db.UpdateDB(table, vals, ln)
		}
	}
}

func extractPatients(infile string, count int, tumor, acc map[string]map[string]string, species map[string]string) Entries {
	// Assigns patient data to appropriate slices with unique entry IDs
	var col map[string]int
	var l int
	var entries Entries
	missed := 0
	first := true
	start := count
	fmt.Printf("\n\tExtracting patient data from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		spl := strings.Split(string(input.Text()), ",")
		if first == false {
			pass := false
			if strings.ToUpper(spl[4]) != "NA" {
				sp, exists := species[spl[col["Species"]]]
				ac, ex := acc[spl[col["Account"]]]
				if len(spl) == l && exists == true && ex == true {
					// Skip entries without valid species and source data
					aid, e := ac[spl[col["Submitter"]]]
					if e == true {
						var t []string
						count++
						id := strconv.Itoa(count)
						if strings.Contains(spl[col["ID"]], "NA") == true {
							// Make sure source ID is an integer
							spl[col["ID"]] = "-1"
						} else if len(spl[col["Age"]]) > 6 {
							// Make sure age does not exceed decimal precision
							spl[col["Age"]] = spl[col["Age"]][:7]
						}
						// ID, Sex, Age, Castrated, taxa_id, source_id, Species, Date, Comments
						p := []string{id, spl[col["Sex"]], spl[col["Age"]], spl[col["Castrated"]], sp, spl[col["ID"]], spl[col["Species"]], spl[col["Date"]], spl[col["Comments"]]}
						// ID, service, account_id
						s := []string{id, spl[col["Service"]], aid}
						// Diagnosis entry: ID, masspresent, hyperplasia, necropsy, metastasis_id
						d := []string{id, spl[col["MassPresent"]], spl[col["Hyperplasia"]], spl[col["Necropsy"]], spl[col["Metastasis"]]}
						// Assign ID to all tumor, location pairs tumorPairs (in diagnoses.go)
						pairs := tumorPairs(spl[col["Type"]], spl[col["Location"]])
						for _, i := range pairs {
							row, intmr := tumor[i[0]]
							r, inloc := row[i[1]]
							if intmr == true && inloc == true {
								// ID, tumor_id, primary_tumor, malignant
								t = []string{id, r, spl[col["Primary"]], spl[col["Malignant"]]}
							} else {
								t = []string{id, "-1", spl[col["Primary"]], spl[col["Malignant"]]}
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
			col = getColumns(spl)
			l = len(spl)
			first = false
		}
	}
	fmt.Printf("\tExtracted %d records.\n", count-start)
	if missed > 0 {
		fmt.Printf("\t[Warning] Count not find taxa ID or source ID for %d records.\n", missed)
	}
	return entries
}

func LoadPatients(db *dbIO.DBIO, infile string) {
	// Loads unique patient info to appropriate tables
	m := db.GetMax("Patient", "ID")
	tumor := MapOfMaps(db.GetTable("Tumor"))
	acc := MapOfMaps(db.GetTable("Accounts"))
	species := EntryMap(db.GetColumns("Taxonomy", []string{"taxa_id", "Species"}))
	// Get entry slices and upload to db
	entries := extractPatients(infile, m, tumor, acc, species)
	uploadPatients(db, "Patient", entries.p)
	uploadPatients(db, "Diagnosis", entries.d)
	uploadPatients(db, "Tumor_relation", entries.t)
	uploadPatients(db, "Source", entries.s)
	// Recacluate species totals
	SpeciesTotals(db)
}

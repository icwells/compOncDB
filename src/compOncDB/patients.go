// This script will upload patient data to the comparative oncology database

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

type Entries struct {
	p [][]string
	d [][]string
	t [][]string
	s [][]string
}

func (e *Entries) update(p, d, t, s []string) {
	// Appends new entries to appriate slice
	e.p = append(e.p, p)
	e.d = append(e.d, d)
	e.t = append(e.t, t)
	e.s = append(e.s, s)
}

func uploadPatients(db *sql.DB, table string, col map[string]string, list [][]string) {
	// Uploads patient entries to db
	fmt.Printf("\tUploading %s to database\n", table)
	vals, l := dbIO.FormatSlice(list)
	dbIO.UpdateDB(db, table, col[table], vals, l)
}

func extractPatients(infile string, count int, tumor, acc map[string]map[string]string, meta, species map[string]string) Entries {
	// Assigns patient data to appropriate slices with unique entry IDs
	first := true
	var entries Entries
	fmt.Printf("\n\tExtracting accounts from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			pass := false
			spl := strings.Split(line, ",")
			if len(spl) == 15 && strarray.InMapStr(species, spl[4]) == true && strarray.InMapMapStr(acc, spl[13]) == true {
				// Skip entries without valid species and source data
				if strarray.InMapStr(acc[spl[13]], spl[14]) == true {
					count++
					id := strconv.Itoa(count)
					var d, t []string
					// ID, Sex, Age, Castrated, taxa_id, source_id, Species, Date, Comments
					p := []string{id, spl[0], spl[1], spl[2], species[spl[4]], spl[3], spl[4], spl[5], spl[6]}
					// ID, service, account_id
					s := []string{id, spl[12], acc[spl[13]][spl[14]]}
					// Diagnosis entry
					if strarray.InMapStr(meta, spl[8]) == true {
						// ID, masspresent, metastasis_id
						d = []string{id, spl[7], meta[spl[8]]}
					} else {
						d = []string{id, spl[7], `\N`}
					}
					if strarray.InMapMapStr(tumor, spl[9]) == true {
						// Tumor relation entry
						if strarray.InMapStr(tumor[spl[9]], spl[10]) == true {
							// ID, tumor_id, primary_tumor
							t = []string{id, tumor[spl[9]][spl[10]], spl[11]}
						} else {
							t = []string{id, `\N`, spl[11]}
						}
					} else {
						t = []string{id, `\N`, spl[11]}
					}
					entries.update(p, d, t, s)
					pass = true
				}
			}
			if pass == false {
				fmt.Printf("\t[Error] Count not find taxa ID or source ID for %s.\n", spl[4])
			}
		} else {
			first = false
		}
	}
	return entries
}

func LoadPatients(db *sql.DB, col map[string]string, infile string) {
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
	uploadPatients(db, "Tumor_Relation", col, entries.t)
	uploadPatients(db, "Source", col, entries.s)
}

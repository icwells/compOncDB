// This script will search for given records from the comparative oncology database

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

func getTumorRecords(ch chan [][]string, db *mysql.DB, id string, tumor [][]string, primary bool) {
	// Returns tumor information for given id
	var loc []string
	var typ []string	
	rows := GetRows(db, "Tumor_Relation", "ID", id)
	for _, i := range rows {
		if primary == false || primary == true && i[2] == "1" {
			for _, j := range tumor {
				if i[1] == j[0] {
					if i[2] == "1" {
						// Prepend primary tumor
						typ = append(j[1], typ)
						loc = append(j[2], loc)
					} else {
						// Append tumor type and location
						typ = append(typ, j[1])
						loc = append(loc, j[2])
					}
					break
				}
			}
		}
	}
	diag := [][]string{typ, loc}
	ch <- diag
}

func getTumor(db *mysql.DB, ids []string, primary bool) map[string][]string {
	// Returns map of tumor data from patient ids
	ch := make(chan [][]string)
	rec := make(map[string][]string)
	tumor := dbIO.GetTable(db, "Tumor")
	for _, id := range ids {
		// Get records for each patient concurrently
		go getTumorRecords(ch, db, id, tumor, primary)
		ret := <-ch
		if len(ret) >= 1 {
			rec[id] = ret
		}
	}
	return rec
}

func getMetastasis(ch chan []string, db *mysql.DB, id string, meta [][]string, mass bool) {
	// Returns diagnosis and metastasis data
	var diag []string
	rows := GetRows(db, "Diagnosis", "ID", id)
	for _, i := range rows {
		if mass == false || mass == true && i[1] == "1" {
			for _, j := range meta {
				if i[2] == j[0] {
					// Append metastasis location
					diag = append(diag, j[1])
					break
				}
			}
		}
	}
	ch <- diag
}

func getDiagosis(db *mysql.DB, ids []string, mass bool) map[string][]string {
	// Returns metastatis info from patient ids
	ch := make(chan []string)
	diagnoses := make(map[string][]string)
	meta := dbIO.GetTable(db, "Metastasis")
	for _, id := range ids {
		// Get records for each patient concurrently
		go getMetastasis(ch, db, id, meta, mass)
		ret := <-ch
		if len(ret) >= 1 {
			diagnoses[id] = ret
		}
	}
	return diagnoses
}

func getRecords(db *mysql.DB, ids []string, mass, primary bool) map[string][]string {
	// Gets diagnosis and metastasis data and formats values
	fmt.Println("\tExtracting diagnosis information...")
	diagnoses := make(map[string][]string)
	meta := getDiagosis(db, ids, mass)
	tumor := getTumor(db, ids, primary)
	for _, i := range ids {
		temp := []string{strings.Join(tumor[i][0], ";"), strings.Join(tumor[i][1]), strings.Join(meta[i], ";")}
		diagnoses[i] = temp
	}
	return diagnoses
}

func getTaxonomy(db *mysql.DB, ids []string, source bool) map[string][]string {
	// Returns taxonomy as map with taxa id as key
	taxa := make(map[string][]string)
	fmt.Println("\tExtracting taxonomy information...")
	table := dbIO.GetTableMap(db, "Taxonomy")
	for _, id := range ids {
		if strarray.InMapSli(table, id) == true {
			if source == true {
				taxa[id] = table[id]
			} else {
				// Exclude source
				taxa[id] = table[id][:7]
			}
		}
	}
	return taxa
}

func getPatients(db *mysql.DB, ids []string) (map[string][]string, map[string]string) {
	// Returns map of target patient data (without id numbers) and map of taxa ids
	patients := make(map[string][]string)
	tids := make(map[string]string)
	table := dbIO.GetTableMap(db, "Patient")
	for _, id := range ids {
		if strarray.InMapSli(table, id) == true {
			patients[id] = table[id][:3]
			patients[id] = append(patients[id], table[id][5:]...)
			tids[id] = table[id][3]
		}
	}
	return patients, tids
}

func searchPatients(db *mysql.DB, col map[string]string, ids []string, outfile, header string) {
	// Extracs patient data using IDs
	var records [][]string
	var taxaids []string
	fmt.Println("\tExtracting patient information...")
	patients, tid := getPatients(db, ids)
	for _, i := range tid {
		taxaids = append(taxaids, i)
	}
	// Leaving primary tumor and mass present switches false for now
	records := getRecords(db, ids, false, false)
	taxonomy := getTaxonomy(db, taxaids, false)
	for _, i := range ids {
		rec := append(patients[i], records[i]...)
		rec = append(rec, taxonomy[tid[i]]...)
		records = append(records, rec)
	}
	iotools.WriteToCSV(outfile, header, records)
}

func getTaxa(db *mysql.DB, ids []string) map[string][]string {
	// Extracts patient data using taxa ids
	patients := make(map[string][]string)
	table := dbIO.GetTable(db, "Patient")
	for _, i := range table {
		for _, id := range ids {
			if id == i[4] {
				patients[id] = table[id][:3]
				patients[id] = append(patients[id], table[id][5:]...)
				break
			}
		}
	}
	return patients
}

func getTaxaIDs(db *mysql.DB, names []string, common bool) []string {
	// Returns taxa id from species name
	var ids []string
	var table [][]string
	if common == true {
		// Get taxonomy ids from common name list
		table = dbIO.GetTable(db, "Common")
	} else {
		// Get ids from taxonomy table
		table = dbIO.GetColumns(db, "Taxonomy", []string{"taxa_id", "Species"})
	}
	for _, n := range names {
		for _, i := range table {
			if n == i[1] {
				ids = append(ids, i[0])
				break
			}
		}
	}
	return ids
}

func searchSpecies(db *mysql.DB, col map[string]string, names []string, outfile, header string, common bool) {
	// Extracts data using species names
	var records [][]string
	var tid []string
	fmt.Println("\tExtracting patient information...")
	ids := getTaxaIDs(db, names, common)
	pat := getTaxa(db, ids)
	// Leaving primary tumor and mass present switches false for now
	rec := getRecords(db, ids, false, false)
	taxonomy := getTaxonomy(db, ids, false)
	for _, i := range ids {
		rec := append(patients[i], records[i]...)
		rec = append(rec, taxonomy[tid[i]]...)
		records = append(records, rec)
	}
	iotools.WriteToCSV(outfile, header, records)
}

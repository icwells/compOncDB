// This script will search for given records from the comparative oncology database

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/icwells/go-tools/strarray"
	"os"
	"strings"
)

func getTumorRecords(ch chan []string, db *sql.DB, id string, tumor [][]string, primary bool) {
	// Returns tumor information for given id
	var loc, typ, mal, prim []string
	rows := dbIO.GetRows(db, "Tumor_relation", "ID", id, "*")
	for _, i := range rows {
		if primary == false || primary == true && i[2] == "1" {
			for _, j := range tumor {
				if i[1] == j[0] {
					if i[2] == "1" {
						// Prepend primary tumor
						prim = append([]string{i[2]}, prim...)
						mal = append([]string{i[3]}, mal...)
						typ = append([]string{j[1]}, typ...)
						loc = append([]string{j[2]}, loc...)
					} else {
						// Append tumor type and location
						prim = append(prim, i[2])
						mal = append(mal, i[3])
						typ = append(typ, j[1])
						loc = append(loc, j[2])
					}
					break
				}
			}
		}
	}
	diag := []string{strings.Join(prim, ";"), strings.Join(mal, ";"), strings.Join(typ, ";"), strings.Join(loc, ";")}
	ch <- diag
}

func getTumor(db *sql.DB, ids []string, primary bool) map[string][]string {
	// Returns map of tumor data from patient ids
	ch := make(chan []string)
	// {id: [types], [locations]}
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

func getMetastasis(ch chan []string, db *sql.DB, id string, meta [][]string, mass bool) {
	// Returns diagnosis and metastasis data
	var diag, loc []string
	rows := dbIO.GetRows(db, "Diagnosis", "ID", id, "*")
	for _, i := range rows {
		if mass == false || mass == true && i[1] == "1" {
			diag = i[1:3]
			for _, j := range meta {
				if i[2] == j[0] {
					// Append metastasis location
					loc = append(loc, j[1])
					break
				}
			}
		}
	}
	if len(diag) > 1 {
		diag = append(diag, strings.Join(loc, ";"))
	}
	ch <- diag
}

func getDiagosis(db *sql.DB, ids []string, mass bool) map[string][]string {
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

func getRecords(db *sql.DB, ids []string, mass, primary bool) map[string][]string {
	// Gets diagnosis and metastasis data and formats values
	fmt.Println("\tExtracting diagnosis information...")
	diagnoses := make(map[string][]string)
	meta := getDiagosis(db, ids, mass)
	tumor := getTumor(db, ids, primary)
	for _, i := range ids {
		// Join multiple entires for same record
		temp := append(meta[i], tumor[i]...)
		diagnoses[i] = temp
	}
	return diagnoses
}

func getTaxonomy(db *sql.DB, ids []string, source bool) map[string][]string {
	// Returns taxonomy as map with taxa id as key
	taxa := make(map[string][]string)
	fmt.Println("\tExtracting taxonomy information...")
	table := dbIO.GetTableMap(db, "Taxonomy")
	for _, id := range ids {
		if strarray.InMapSli(table, id) == true {
			if source == true {
				// Keep source column
				taxa[id] = table[id][:6]
				taxa[id] = append(taxa[id], table[id][7])
			} else {
				// Exclude source and species (in patient table)
				taxa[id] = table[id][:6]
			}
		}
	}
	return taxa
}

func getTaxa(db *sql.DB, ids []string) (map[string][][]string, []string) {
	// Extracts patient data using taxa ids
	patients := make(map[string][][]string)
	var uid []string
	table := dbIO.GetTable(db, "Patient")
	for _, i := range table {
		for _, id := range ids {
			if id == i[4] {
				var rec []string
				// Skip source and taxonomy ids
				rec = i[:4]
				rec = append(rec, i[6:]...)
				patients[id] = append(patients[id], rec)
				uid = append(uid, i[0])
				break
			}
		}
	}
	return patients, uid
}

func getTaxaIDs(db *sql.DB, names []string, level string, common bool) []string {
	// Returns taxa id from species name
	var ids []string
	var table [][]string
	if common == true {
		// Get taxonomy ids from common name list
		table = dbIO.GetTable(db, "Common")
	} else {
		// Get ids from taxonomy table
		table = dbIO.GetColumns(db, "Taxonomy", []string{"taxa_id", level})
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

func checkLevel(level string, common bool) string {
	// Makes sure a valid taxonomic level is given
	found := false
	if common == true {
		// Overwrite to species for common name comparison
		level = "Species"
	} else {
		levels := []string{"Kingdom", "Phylum", "Class", "Orders", "Family", "Genus", "Species"}
		// Convert level to title case
		level = strings.Title(level)
		for _, i := range levels {
			if level == i {
				found = true
				break
			}
		}
		if found == false {
			fmt.Println("\n\t[Error] Please enter a valid taxonomic level. Exiting.\n")
			os.Exit(11)
		}
	}
	return level
}

func searchTaxonomicLevels(db *sql.DB, col map[string]string, level string, names []string, common bool) [][]string {
	// Extracts data using species names
	var ret [][]string
	fmt.Println("\tExtracting patient information...")
	level = checkLevel(level, common)
	ids := getTaxaIDs(db, names, level, common)
	patients, uid := getTaxa(db, ids)
	// Leaving primary tumor and mass present switches false for now
	records := getRecords(db, uid, false, false)
	taxonomy := getTaxonomy(db, ids, false)
	for _, id := range ids {
		for _, i := range patients[id] {
			_, ex := records[i[0]] 
			if ex == true {
				rec := append(i, records[i[0]]...)
				rec = append(rec, taxonomy[id]...)
				ret = append(ret, rec)
			}
		}
	}
	return ret
}

/*func searchTaxaRank(db *sql.DB, col map[string]string, target, column, outfile, header string) {
	// Extracts ids from taxonomy table where target is present in column
	var ids []string
	tids := dbIO.GetRows(db, "Taxonomy", column, target, "taxa_id")
	for _, i := range tids {
		for _, j := range i {
			// Convert ot slice of strings	
			ids = append(ids, j)
		}
	}
	patients := getTaxa(db, ids)
	
}*/

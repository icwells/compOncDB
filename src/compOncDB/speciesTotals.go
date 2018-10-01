// This script will calculate and store total occurances for each species

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	"strconv"
)

func (r *Record) toSlice() []string {
	// Returns string slice of total, adult, and cancer incidences
	var ret []string
	ret = append(ret, strconv.Itoa(r.total))
	ret = append(ret, strconv.Itoa(r.adult))
	ret = append(ret, strconv.Itoa(r.cancer))
	return ret
}

func uploadTotals(db *sql.DB, col map[string]string, records map[string]*Record) {
	// Converts map to slice and uploads to table
	var totals [][]string
	fmt.Println("\tUploading new species totals...")
	for k, v := range records {
		// Taxa id, total, adult, cancer
		t := append([]string{k}, v.toSlice()...)
		totals = append(totals, t)
	}
	vals, l := dbIO.FormatSlice(totals)
	dbIO.UpdateDB(db, "Totals", col["Totals"], vals, l)
}

func getTotals(db *sql.DB, records map[string]*Record) map[string]*Record {
	// Returns struct with number of total, adult, and adult cancer occurances by species
	diag := entryMap(dbIO.GetColumns(db, "Diagnosis", []string{"Masspresent", "ID"}))
	rows := dbIO.GetColumns(db, "Patient", []string{"taxa_id", "Age", "ID"})
	for _, i := range rows {
		_, exists := records[i[0]]
		if exists == true {
			// Increment total
			records[i[0]].total++
			age, _ := strconv.ParseFloat(i[1], 64)
			if age > records[i[0]].infant {
				// Increment adult if age is greater than age of infancy
				records[i[0]].adult++
			}
			d, e := diag[i[2]]
			if e == true {
				if d == "1" {
					// Increment cancer count if masspresent == true
					records[i[0]].cancer++
				}
			}
		}
	}
	return records
}

func speciesTotals(db *sql.DB, col map[string]string) {
	// Recalculates occurances for each species
	dbIO.TruncateTable(db, "Totals")
	fmt.Println("\tCalculating total occurances by species...")
	records := getTargetSpecies(db, 0)
	records = getAgeOfInfancy(db, records)
	records = getTotals(db, records)
	uploadTotals(db, col, records)
}

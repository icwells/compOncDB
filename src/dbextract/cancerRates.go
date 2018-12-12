// This scrpt will calculate cancer rates for species with  at least a given number of entries

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
)

func formatRates(records map[string]*dbupload.Record) [][]string {
	// Calculates rates and formats for printing
	var ret [][]string
	for _, v := range records {
		ret = append(ret, v.CalculateRates())
	}
	return ret
}

func getSpeciesNames(db *dbIO.DBIO, records map[string]*dbupload.Record) map[string]*dbupload.Record {
	// Adds species names to structs
	species := dbupload.EntryMap(db.GetRows("Taxonomy", "taxa_id", dbupload.getRecKeys(records), "Species,taxa_id"))
	for k, v := range species {
		if dbupload.InMapRec(records, k) == true {
			records[k].Species = v
		}
	}
	return records
}

func getTargetSpecies(db *dbIO.DBIO, min int) map[string]*dbupload.Record {
	// Returns map of empty species records with >= min occurances
	records := make(map[string]*dbupload.Record)
	target := db.GetRowsMin("Totals", "Adult", "*", min)
	for _, i := range target {
		var rec dbupload.Record
		rec.SetRecord(i)
		records[i[0]] = &rec
	}
	return records
}

func GetCancerRates(db *dbIO.DBIO, min int, nec bool) [][]string {
	// Returns slice of string slices of cancer rates and related info
	var ret [][]string
	fmt.Printf("\n\tCalculating rates for species with at least %d entries...\n", min)
	records := getTargetSpecies(db, min)
	if len(records) > 0 {
		records = getSpeciesNames(db, records)
		if len(records) > 0 {
			ret = formatRates(records)
		}
	}
	return ret
}

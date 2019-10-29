// This scrpt will calculate cancer rates for species with  at least a given number of entries

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/strarray"
)

func cancerRateHeader() []string {
	// Returns header for hancer rate file
	return []string{"Kingdom", "Phylum", "Class", "Orders", "Family", "Genus", "ScientificName",
		"TotalRecords", "CancerRecords", "CancerRate", "AverageAge(months)", "AvgAgeCancer(months)",
		"Male", "Female", "MaleCancer", "FemaleCancer"}
}

func formatRates(records map[string]*dbupload.Record) *dataframe.Dataframe {
	// Calculates rates and formats for printing
	ret := dataframe.NewDataFrame(-1)
	ret.SetHeader(cancerRateHeader())
	for _, v := range records {
		if len(v.Species) > 0 {
			err := ret.AddRow(v.CalculateRates())
			if err != nil {

			}
		}
	}
	return ret
}

func getSpeciesNames(db *dbIO.DBIO, records map[string]*dbupload.Record) map[string]*dbupload.Record {
	// Adds taxonomies to structs
	species := dbupload.ToMap(db.GetRows("Taxonomy", "taxa_id", dbupload.GetRecKeys(records), "*"))
	for k, v := range species {
		if dbupload.InMapRec(records, k) == true {
			records[k].Taxonomy = v[:6]
			records[k].Species = v[6]
		}
	}
	return records
}

func filterRecords(taxaids []string, records map[string]*dbupload.Record) map[string]*dbupload.Record {
	// Deletes records not in taxaids
	if len(records) > 0 && len(taxaids) > 0 {
		for k := range records {
			if !strarray.InSliceStr(taxaids, k) {
				delete(records, k)
			}
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
		if len(i[0]) > 0 {
			records[i[0]] = &rec
		}
	}
	return records
}

func minRecords(records map[string]*dbupload.Record, min int) map[string]*dbupload.Record {
	// Removes records with < min entries
	for k := range records {
		if records[k].Adult < min {
			delete(records, k)
		} else {
			// Calculate average ages
			records[k].CalculateAvgAges()
		}
	}
	return records
}

func getNecropsySpecies(db *dbIO.DBIO, min int) map[string]*dbupload.Record {
	// Returns species with at least min necropsy records
	records := make(map[string]*dbupload.Record)
	fmt.Println("\tCounting necropsy records...")
	records = dbupload.GetAllSpecies(db)
	records = dbupload.GetAgeOfInfancy(db, records)
	records = dbupload.GetTotals(db, records, true)
	return minRecords(records, min)
}

func GetCancerRates(db *dbIO.DBIO, min int, nec bool, eval []codbutils.Evaluation) *dataframe.Dataframe {
	// Returns slice of string slices of cancer rates and related info
	var ret *dataframe.Dataframe
	var records map[string]*dbupload.Record
	fmt.Printf("\n\tCalculating rates for species with at least %d entries...\n", min)
	if nec == false {
		records = getTargetSpecies(db, min)
	} else {
		records = getNecropsySpecies(db, min)
	}
	if len(eval) > 0 {
		s := newSearcher(db, false)
		s.assignSearch(eval)
		records = filterRecords(s.taxaids, records)
	}
	if len(records) > 0 {
		records = getSpeciesNames(db, records)
		if len(records) > 0 {
			ret = formatRates(records)
		}
	}
	fmt.Printf("\tFound %d records with at least %d entries.\n", ret.Length(), min)
	return ret
}

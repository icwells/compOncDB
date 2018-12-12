// This scrpt will calculate cancer rates for species with  at least a given number of entries

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"strconv"
)

func (r *dbupload.Record) calculateRates() []string {
	// Returns string slice of rates
	//"ScientificName,AdultRecords,CancerRecords,CancerRate,AverageAge(months),AvgAgeCancer(months),Male:Female\n"
	ret := []string{r.species}
	ret = append(ret, strconv.Itoa(r.adult))
	ret = append(ret, strconv.Itoa(r.cancer))
	// Calculate rates
	rate := float64(r.cancer) / float64(r.adult)
	// Append rates to slice and return
	ret = append(ret, strconv.FormatFloat(rate, 'f', 2, 64))
	ret = append(ret, strconv.FormatFloat(r.age, 'f', 2, 64))
	ret = append(ret, strconv.FormatFloat(r.cancerage, 'f', 2, 64))
	ret = append(ret, strconv.Itoa(r.male))
	ret = append(ret, strconv.Itoa(r.female))
	ret = append(ret, strconv.Itoa(r.malecancer))
	ret = append(ret, strconv.Itoa(r.femalecancer))
	return ret
}

func (r *dbupload.Record) setRecord(row []string) {
	// Reads values from Totals table entry
	r.total, _ = strconv.Atoi(row[1])
	r.age, _ = strconv.ParseFloat(row[2], 64)
	r.adult, _ = strconv.Atoi(row[3])
	r.male, _ = strconv.Atoi(row[4])
	r.female, _ = strconv.Atoi(row[5])
	r.cancer, _ = strconv.Atoi((row[6]))
	r.cancerage, _ = strconv.ParseFloat(row[7], 64)
}

//----------------------------------------------------------------------------

func formatRates(records map[string]*dbupload.Record) [][]string {
	// Calculates rates and formats for printing
	var ret [][]string
	for _, v := range records {
		ret = append(ret, v.calculateRates())
	}
	return ret
}

func getSpeciesNames(db *dbIO.DBIO, records map[string]*dbupload.Record) map[string]*dbupload.Record {
	// Adds species names to structs
	species := dbupload.entryMap(db.GetRows("Taxonomy", "taxa_id", dbupload.getRecKeys(records), "Species,taxa_id"))
	for k, v := range species {
		if dbupload.inMapRec(records, k) == true {
			records[k].species = v
		}
	}
	return records
}

func getTargetSpecies(db *dbIO.DBIO, min int) map[string]*dbupload.Record {
	// Returns map of empty species records with >= min occurances
	records := make(map[string]*dbupload.Record)
	target := db.GetRowsMin("Totals", "Adult", "*", min)
	for _, i := range target {
		var rec Record
		rec.setRecord(i)
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

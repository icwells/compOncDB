// This script will calculate and store total occurances for each species

package dbupload

import (
	"fmt"
	"github.com/icwells/dbIO"
	"strconv"
)

func uploadTotals(db *dbIO.DBIO, records map[string]*Record) {
	// Converts map to slice and uploads to table
	var totals [][]string
	fmt.Println("\tUploading new species totals...")
	for k, v := range records {
		// Taxa id, total, adult, cancer
		totals = append(totals, v.ToSlice(k))
	}
	vals, l := dbIO.FormatSlice(totals)
	db.UpdateDB("Totals", vals, l)
}

func addDenominators(db *dbIO.DBIO, records map[string]*Record) map[string]*Record {
	// Adds fixed values from denominators table
	d := ToMap(db.GetTable("Denominators"))
	for k, v := range d {
		_, ex := records[k]
		if ex == true {
			t, err := strconv.Atoi(v[0])
			if err == nil {
				records[k].Total += t
				records[k].Adult += t
			}
		}
	}
	return records
}

func getTotals(db *dbIO.DBIO, records map[string]*Record) map[string]*Record {
	// Returns struct with number of total, adult, and adult cancer occurances by species
	diag := EntryMap(db.GetColumns("Diagnosis", []string{"Masspresent", "ID"}))
	rows := db.GetColumns("Patient", []string{"taxa_id", "Age", "ID", "Sex"})
	for _, i := range rows {
		_, exists := records[i[0]]
		if exists == true {
			// Increment total
			records[i[0]].Total++
			age, err := strconv.ParseFloat(i[1], 64)
			if err == nil && age > records[i[0]].Infant {
				// Increment adult if age is greater than age of infancy
				records[i[0]].Adult++
				records[i[0]].Age = records[i[0]].Age + age
				if i[3] == "male" {
					records[i[0]].Male++
				} else if i[3] == "female" {
					records[i[0]].Female++
				}
				d, e := diag[i[2]]
				if e == true {
					if d == "1" {
						// Increment cancer count and age if masspresent == true
						records[i[0]].Cancer++
						records[i[0]].Cancerage = records[i[0]].Cancerage + age
						if i[3] == "male" {
							records[i[0]].Malecancer++
						} else if i[3] == "female" {
							records[i[0]].Femalecancer++
						}
					}
				}
			}
		}
	}
	return records
}

func getAgeOfInfancy(db *dbIO.DBIO, records map[string]*Record) map[string]*Record {
	// Updates structs with min age for each species
	// Get appropriate ages for each taxon
	ages := db.GetRows("Life_history", "taxa_id", GetRecKeys(records), "taxa_id,female_maturity,male_maturity,Weaning")
	for _, i := range ages {
		// Assign ages to structs
		if InMapRec(records, i[0]) == true {
			w, _ := strconv.ParseFloat(i[3], 64)
			f, _ := strconv.ParseFloat(i[1], 64)
			m, _ := strconv.ParseFloat(i[2], 64)
			if w > 0.0 {
				// Assign weaning age
				records[i[0]].Infant = w
			} else if f > 0.0 && m > 0.0 {
				// Assign 10% of average age of maturity
				records[i[0]].Infant = (((f + m) / 2) * 0.1)
			} else {
				// Default to 1 month
				records[i[0]].Infant = 1.0
			}
		}
	}
	return records
}

func getAllSpecies(db *dbIO.DBIO) map[string]*Record {
	// Returns map of empty species records with >= min occurances
	records := make(map[string]*Record)
	unique := db.GetColumnText("Taxonomy", "taxa_id")
	for _, v := range unique {
		var rec Record
		records[v] = &rec
	}
	return records
}

func SpeciesTotals(db *dbIO.DBIO) {
	// Recalculates occurances for each species
	db.TruncateTable("Totals")
	fmt.Println("\tCalculating total occurances by species...")
	records := getAllSpecies(db)
	records = getAgeOfInfancy(db, records)
	records = getTotals(db, records)
	records = addDenominators(db, records)
	uploadTotals(db, records)
}

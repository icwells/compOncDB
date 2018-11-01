// This script will calculate and store total occurances for each species

package main

import (
	"fmt"
	"github.com/icwells/dbIO"
	"strconv"
)

type Record struct {
	species   		string
	infant    		float64
	total     		int
	age       		float64
	male      		int
	female    		int
	cancer    		int
	cancerage 		float64
	adult	  		int
	malecancer		int
	femalecancer	int
}

func avgAge(n float64, d int) string {
	// Returns string of n/d
	var ret string
	if n > 0.0 && d > 0 {
		age := n / float64(d)
		ret = strconv.FormatFloat(age, 'f', -1, 64)
	} else {
		ret = "-1"
	}
	return ret
}

func (r *Record) String() string {
	// Returns formatted string of record attributes
	ret := fmt.Sprintf("\nSpecies: %s\n", r.species)
	ret += fmt.Sprintf("Total: %d\n", r.total)
	ret += fmt.Sprintf("Cancer Records: %d", r.cancer)
	return ret
}

func (r *Record) getAvgAge() string {
	// Returns string of avg age
	return avgAge(r.age, r.adult)
}

func (r *Record) getCancerAge() string {
	// Returns string of average cancer record age
	return avgAge(r.cancerage, r.cancer)
}

func (r *Record) toSlice(id string) []string {
	// Returns string slice of values for upload to table
	var ret []string
	ret = append(ret, id)
	ret = append(ret, strconv.Itoa(r.total))
	ret = append(ret, r.getAvgAge())
	ret = append(ret, strconv.Itoa(r.adult))
	ret = append(ret, strconv.Itoa(r.male))
	ret = append(ret, strconv.Itoa(r.female))
	ret = append(ret, strconv.Itoa(r.cancer))
	ret = append(ret, r.getCancerAge())
	ret = append(ret, strconv.Itoa(r.malecancer))
	ret = append(ret, strconv.Itoa(r.femalecancer))
	return ret
}

func uploadTotals(db *dbIO.DBIO, records map[string]*Record) {
	// Converts map to slice and uploads to table
	var totals [][]string
	fmt.Println("\tUploading new species totals...")
	for k, v := range records {
		// Taxa id, total, adult, cancer
		totals = append(totals, v.toSlice(k))
	}
	vals, l := dbIO.FormatSlice(totals)
	db.UpdateDB("Totals", vals, l)
}

func addDenominators(db *dbIO.DBIO, records map[string]*Record) map[string]*Record {
	// Adds fixed values from denominators table
	d := toMap(db.GetTable("Denominators"))
	for k, v := range d {
		_, ex := records[k]
		if ex == true {
			t, err := strconv.Atoi(v[0])
			if err == nil {
				records[k].total += t
				records[k].adult += t
			}
		}
	}
	return records
}

func getTotals(db *dbIO.DBIO, records map[string]*Record) map[string]*Record {
	// Returns struct with number of total, adult, and adult cancer occurances by species
	diag := entryMap(db.GetColumns("Diagnosis", []string{"Masspresent", "ID"}))
	rows := db.GetColumns("Patient", []string{"taxa_id", "Age", "ID", "Sex"})
	for _, i := range rows {
		_, exists := records[i[0]]
		if exists == true {
			// Increment total
			records[i[0]].total++
			age, err := strconv.ParseFloat(i[1], 64)
			if err == nil && age > records[i[0]].infant {
				// Increment adult if age is greater than age of infancy
				records[i[0]].adult++
				records[i[0]].age = records[i[0]].age + age
				if i[3] == "male" {
					records[i[0]].male++
				} else if i[3] == "female" {
					records[i[0]].female++
				}
				d, e := diag[i[2]]
				if e == true {
					if d == "1" {
						// Increment cancer count and age if masspresent == true
						records[i[0]].cancer++
						records[i[0]].cancerage = records[i[0]].cancerage + age
						if i[3] == "male" {
							records[i[0]].malecancer++
						} else if i[3] == "female" {
							records[i[0]].femalecancer++
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
	ages := db.GetRows("Life_history", "taxa_id", getRecKeys(records), "taxa_id,female_maturity,male_maturity,Weaning")
	for _, i := range ages {
		// Assign ages to structs
		if inMapRec(records, i[0]) == true {
			w, _ := strconv.ParseFloat(i[3], 64)
			f, _ := strconv.ParseFloat(i[1], 64)
			m, _ := strconv.ParseFloat(i[2], 64)
			if w > 0.0 {
				// Assign weaning age
				records[i[0]].infant = w
			} else if f > 0.0 && m > 0.0 {
				// Assign 10% of average age of maturity
				records[i[0]].infant = (((f + m) / 2) * 0.1)
			} else {
				// Default to 1 month
				records[i[0]].infant = 1.0
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

func speciesTotals(db *dbIO.DBIO) {
	// Recalculates occurances for each species
	db.TruncateTable("Totals")
	fmt.Println("\tCalculating total occurances by species...")
	records := getAllSpecies(db)
	records = getAgeOfInfancy(db, records)
	records = getTotals(db, records)
	records = addDenominators(db, records)
	uploadTotals(db, records)
}

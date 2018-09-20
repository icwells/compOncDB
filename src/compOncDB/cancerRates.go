// This scrpt will calculate cancer rates for species with  at least a given number of entries

package main

import (
	"bytes"
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/icwells/go-tools/strarray"
	"strconv"
)

type Entry struct {
	// For holding values until diagnosis is determined
	age    float64
	male   bool
	female bool
}

type Record struct {
	species   string
	infant    float64
	total     int
	age       float64
	male      int
	female    int
	entries   map[string]*Entry
	cancer    int
	cancerage float64
}

func (r *Record) String() string {
	// Returns formatted string of record attributes
	ret := fmt.Sprintf("\nSpecies: %s\n", r.species)
	ret += fmt.Sprintf("Total: %d\n", r.total)
	ret += fmt.Sprintf("Cancer Records: %d", r.cancer)
	return ret
}

func (r *Record) calculateRates() []string {
	// Returns string slice of rates
	//"ScientificName,TotalRecords,CancerRecords,CancerRate,AverageAge(months),AvgAgeCancer(months),Male:Female\n"
	ret := []string{r.species}
	ret = append(ret, strconv.Itoa(r.total))
	ret = append(ret, strconv.Itoa(r.cancer))
	// Calculate rates
	rate := float64(r.cancer) / float64(r.total)
	avgage := r.age / float64(r.total)
	cage := r.cancerage / float64(r.cancer)
	ratio := float64(r.male) / float64(r.female)
	// Append rates to slice and return
	ret = append(ret, strconv.FormatFloat(rate, 'f', 2, 64))
	ret = append(ret, strconv.FormatFloat(avgage, 'f', 2, 64))
	ret = append(ret, strconv.FormatFloat(cage, 'f', 2, 64))
	ret = append(ret, strconv.FormatFloat(ratio, 'f', 2, 64))
	return ret
}

func inMapRec(m map[string]*Record, s string) bool {
	// Return true if s is a key in m
	_, ret := m[s]
	return ret
}

func getRecKeys(records map[string]*Record) string {
	// Returns string of taxa_ids
	first := true
	buffer := bytes.NewBufferString("")
	for k, _ := range records {
		if first == false {
			// Write name with preceding comma
			buffer.WriteByte(',')
			buffer.WriteString(k)
		} else {
			buffer.WriteString(k)
			first = false
		}
	}
	return buffer.String()
}

//----------------------------------------------------------------------------

func formatRates(records map[string]*Record) [][]string {
	// Calculates rates and formats for printing
	var ret [][]string
	for _, v := range records {
		ret = append(ret, v.calculateRates())
	}
	return ret
}

func getSpeciesNames(db *sql.DB, records map[string]*Record) map[string]*Record {
	// Adds species names to structs
	species := entryMap(dbIO.GetRows(db, "Taxonomy", "taxa_id", getRecKeys(records), "Species,taxa_id"))
	for k, v := range species {
		if inMapRec(records, k) == true {
			records[k].species = v
		}
	}
	return records
}

func getSpeciesDiagnoses(db *sql.DB, records map[string]*Record, nec bool) map[string]*Record {
	// Adds diagnosis info
	diag := toMap(dbIO.GetTable(db, "Diagnosis"))
	for _, val := range records {
		for k, v := range val.entries {
			if strarray.InMapSli(diag, k) == true {
				if nec == true {
					v, ex := diag[k]
					if ex == false || len(v) < 2 || v[1] != "1" {
						// Delete non-necropsy records from the map
						delete(val.entries, k)
					}
				} else {
					// Add values to species total
					val.total++
					val.age += v.age
					if v.male == true {
						val.male++
					} else if v.female == true {
						val.female++
					}
					if diag[k][0] == "1" {
						val.cancer++
						val.cancerage += v.age
					}
				}
			}
		}
	}
	for k, v := range records {
		// Remove empty records
		if v.total == 0 {
			delete(records, k)
		}
	}
	return records
}

func getSpeciesSummaries(db *sql.DB, records map[string]*Record, min int) map[string]*Record {
	// Updates structs with total age, number of males/females, and patient IDs; deletes entries with fewer than min adult records
	fmt.Println("\tGetting records...")
	patients := dbIO.GetColumns(db, "Patient", []string{"taxa_id", "Age", "Sex" ,"ID"})
	for _, i := range patients {
		if inMapRec(records, i[0]) == true {
			age, _ := strconv.ParseFloat(i[1], 64)
			if age >= records[i[0]].infant {
				var e Entry
				// Store age by id to calculate average age of cancer records later
				e.age = age
				if i[2] == "male" {
					e.male = true
				} else if i[2] == "female" {
					e.female = true
				}
				records[i[0]].entries[i[3]] = &e
			}
		}
	}
	for k, v := range records {
		// Make sure each record still exceeds the minimum
		if len(v.entries) < min {
			delete(records, k)
		}
	}
	return records
}

func getAgeOfInfancy(db *sql.DB, records map[string]*Record) map[string]*Record {
	// Updates structs with min age for each species
	// Get appropriate ages for each taxon
	ages := dbIO.GetRows(db, "Life_history", "taxa_id", getRecKeys(records), "taxa_id,female_maturity,male_maturity,Weaning")
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

func getTargetSpecies(db *sql.DB, min int) map[string]*Record {
	// Returns map of empty species records with >= min occurances
	records := make(map[string]*Record)
	unique := dbIO.GetNumOccurances(db, "Patient", "taxa_id")
	for k, v := range unique {
		if v >= min {
			var rec Record
			rec.entries = make(map[string]*Entry)
			records[k] = &rec
		}
	}
	return records
}

func getCancerRates(db *sql.DB, col map[string]string, min int, nec bool) [][]string {
	// Returns slice of string slices of cancer rates and related info
	var ret [][]string
	fmt.Printf("\n\tCalculating rates for species with at least %d entries...\n", min)
	records := getTargetSpecies(db, min)
	records = getAgeOfInfancy(db, records)
	if len(records) > 0 {
		records = getSpeciesSummaries(db, records, min)
		if len(records) > 0 {
			records = getSpeciesDiagnoses(db, records, nec)
			if len(records) > 0 {
				records = getSpeciesNames(db, records)
				if len(records) > 0 {
					ret = formatRates(records)
				}
			}
		}
	}
	return ret
}

// This scrpt will calculate cancer rates for species with  at least a given number of entries

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/strarray"
	"strings"
)

type cancerRates struct {
	db      *dbIO.DBIO
	min     int
	nec     bool
	lh      bool
	header  []string
	nas     []string
	records map[string]*dbupload.Record
	rates   *dataframe.Dataframe
}

func newCancerRates(db *dbIO.DBIO, min int, nec, lh bool) *cancerRates {
	// Returns initialized cancerRates struct
	c := new(cancerRates)
	c.db = db
	c.min = min
	c.nec = nec
	c.lh = lh
	c.header = []string{"taxa_id", "Kingdom", "Phylum", "Class", "Orders", "Family", "Genus", "ScientificName",
		"TotalRecords", "CancerRecords", "CancerRate", "AverageAge(months)", "AvgAgeCancer(months)",
		"Male", "Female", "MaleCancer", "FemaleCancer"}
	if c.lh {
		// Omit taxa_id column
		tail := strings.Split(c.db.Columns["Life_history"], ",")[1:]
		c.header = append(c.header, tail...)
		for i := 0; i < len(tail); i++ {
			c.nas = append(c.nas, "NA")
		}
	}
	c.records = make(map[string]*dbupload.Record)
	c.rates, _ = dataframe.NewDataFrame(-1)
	c.rates.SetHeader(c.header)
	return c
}

func (c *cancerRates) formatRates() {
	// Adds taxonomies to structs, calculates rates, and formats for printing
	if len(c.records) > 0 {
		species := dbupload.ToMap(c.db.GetRows("Taxonomy", "taxa_id", dbupload.GetRecKeys(c.records), "*"))
		lifehist := dbupload.ToMap(c.db.GetRows("Life_history", "taxa_id", dbupload.GetRecKeys(c.records), "*"))
		for k, v := range c.records {
			// Add taxonomy
			if val, ex := species[k]; ex {
				v.Taxonomy = val[:6]
				v.Species = val[6]
			}
			if c.lh {
				// Add life history
				if val, ex := lifehist[k]; ex {
					v.Lifehistory = val
				} else {
					v.Lifehistory = c.nas
				}
			}
			if len(v.Species) > 0 {
				// Calculate cancer rates
				r := v.CalculateRates(k)
				for idx, i := range r {
					// Replace -1 with NA
					if strings.Split(i, ".")[0] == "-1" {
						r[idx] = "NA"
					}
				}
				// Add to dataframe
				err := c.rates.AddRow(r)
				if err != nil {
					fmt.Printf("\t[Error] Adding row to dataframe: %v\n", err)
				}
			}
		}
	}
}

func (c *cancerRates) filterRecords(taxaids []string) {
	// Deletes records not in taxaids
	if len(c.records) > 0 && len(taxaids) > 0 {
		for k := range c.records {
			if !strarray.InSliceStr(taxaids, k) {
				delete(c.records, k)
			}
		}
	}
}

func (c *cancerRates) getNecropsySpecies() {
	// Returns species with at least min necropsy records
	fmt.Println("\tCounting necropsy records...")
	c.records = dbupload.GetAllSpecies(c.db)
	c.records = dbupload.GetAgeOfInfancy(c.db, c.records)
	c.records = dbupload.GetTotals(c.db, c.records, true)
	for k := range c.records {
		if c.records[k].Adult < c.min {
			delete(c.records, k)
		} else {
			// Calculate average ages
			c.records[k].CalculateAvgAges()
		}
	}
}

func (c *cancerRates) getTargetSpecies() {
	// Returns map of empty species records with >= min occurances
	if c.nec {
		c.getNecropsySpecies()
	} else {
		target := c.db.GetRowsMin("Totals", "Adult", "*", c.min)
		for _, i := range target {
			var rec dbupload.Record
			rec.SetRecord(i)
			if len(i[0]) > 0 {
				c.records[i[0]] = &rec
			}
		}
	}
}

func GetCancerRates(db *dbIO.DBIO, min int, nec, lh bool, eval []codbutils.Evaluation) *dataframe.Dataframe {
	// Returns slice of string slices of cancer rates and related info
	c := newCancerRates(db, min, nec, lh)
	fmt.Printf("\n\tCalculating rates for species with at least %d entries...\n", c.min)
	c.getTargetSpecies()
	if len(eval) > 0 {
		s := newSearcher(c.db, false)
		s.assignSearch(eval)
		c.filterRecords(s.taxaids)
	}
	c.formatRates()
	fmt.Printf("\tFound %d records with at least %d entries.\n", c.rates.Length(), c.min)
	return c.rates
}

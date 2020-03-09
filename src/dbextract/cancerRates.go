// This scrpt will calculate cancer rates for species with  at least a given number of entries

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"strings"
)

type cancerRates struct {
	db     *dbIO.DBIO
	df     *dataframe.Dataframe
	header []string
	lh     bool
	min    int
	//nec     bool
	nas     []string
	rates   *dataframe.Dataframe
	records map[string]*Record
}

func newCancerRates(db *dbIO.DBIO, min int, lh bool) *cancerRates {
	// Returns initialized cancerRates struct
	c := new(cancerRates)
	c.db = db
	c.min = min
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
		for k, v := range c.records {
			if v.Total >= c.min {
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

func (c *cancerRates) countRecords() {
	// Returns struct with number of total, adult, and adult cancer occurances by species
	var d []string
	var e bool
	for idx := range c.df.Rows {
		var err error
		tid, _ := df.GetCell(idx, "taxa_id")
		if _, ex := c.records[tid]; ex {
			// Increment total
			c.records[tid].Total++
			age, err := c.df.GetCellFloat(idx, "Age")
			if err == nil && age > c.records[tid].Infant {
				// Increment adult if age is greater than age of infancy
				c.records[tid].Adult++
				c.records[tid].Age = c.records[tid].Age + age
				sex, er := c.df.GetCell(idx, "Sex")
				if er == nil {
					if sex == "male" {
						c.records[tid].Male++
					} else if sex == "female" {
						c.records[tid].Female++
					}
					if mp, e := c.df.GetCell(idx, "Masspresent"); e == nil {
						if mp == "1" {
							// Increment cancer count and age if masspresent == true
							c.records[tid].Cancer++
							c.records[tid].Cancerage = c.records[tid].Cancerage + age
							if sex == "male" {
								c.records[tid].Malecancer++
							} else if sex == "female" {
								c.records[tid].Femalecancer++
							}
						}
					}
				}
			}
		}
	}
}

func (c *cancerRates) appendLifeHistory() {
	// Determines age of infancy and adds life history if needed
	lifehist := dbupload.ToMap(c.db.GetRows("Life_history", "taxa_id", getRecKeys(c.records), "*"))
	for k, v := range c.records {
		if lh, ex := lifehist[k]; ex {
			v.Lifehistory = lh
			v.Infant, _ = strconv.ParseFloat(lh[4], 64)
		} else {
			v.Lifehistory = c.nas
		}
	}
}

func (c *cancerRates) addDenominators() {
	// Adds fixed values from denominators table
	for k, v := range ToMap(c.db.GetTable("Denominators")) {
		if _, ex := c.records[k]; ex {
			if t, err := strconv.Atoi(v[0]); err == nil {
				c.records[k].Total += t
				c.records[k].Adult += t
			}
		}
	}
}

func (c *cancerRates) setTaxonomy(idx int) []string {
	// Stores taxonomy for given record
	var ret []string
	for _, k := range c.db.Columns["Taxonomy"] {
		if k != "Source" {
			val, _ := df.GetCell(idx, k)
			ret = append(ret, val)
		}
	}
	return ret
}

func (c *cancerRates) getTargetSpecies() {
	// Stores map of empty species records with >= min occurances
	for idx := range df.Rows {
		tid, _ := df.GetCell(idx, "taxa_id")
		if _, ex := c.records[tid]; !ex {
			c.records[tid] = newRecord(c.setTaxonomy(idx))
		}

	}
	c.addDenominators()
	c.appendLifeHistory()
}

func (c *cancerRates) setDataframe(eval []codbutils.Evaluation) {
	// Gets dataframe of matching records
	if len(eval) == 0 {
		// Set evaluation to return everything
		eval = codbutils.SetOperations(c.db.Columns, "ID > 0")
	}
	c.df, _ = SearchColumns(c.db, "", eval, false, false)
}

func GetCancerRates(db *dbIO.DBIO, min int, lh bool, eval []codbutils.Evaluation) *dataframe.Dataframe {
	// Returns slice of string slices of cancer rates and related info
	c := newCancerRates(db, min, lh)
	fmt.Printf("\n\tCalculating rates for species with at least %d entries...\n", c.min)
	c.setDataframe(eval)
	c.getTargetSpecies()
	c.countRecords()
	c.formatRates()
	fmt.Printf("\tFound %d records with at least %d entries.\n", c.rates.Length(), c.min)
	return c.rates
}

// Contains functions for populating records map

package dbextract

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"strconv"
	"strings"
)

func (c *cancerRates) getTotals() {
	// Adds all sub-records to field total
	id := "total"
	for key, val := range c.records {
		for k, v := range val {
			if k != id {
				c.records[key][id].Add(v)
			}
		}
	}
}

func (c *cancerRates) countNeoplasia(idx int, field, id, sex string, age float64) {
	// Counts cancer related data
	if mp, err := c.df.GetCell(idx, "Masspresent"); err == nil {
		if mp == "1" {
			// Increment cancer count and age if masspresent == true
			c.records[field][id].Cancer++
			c.records[field][id].Cancerage = c.records[field][id].Cancerage + age
			if sex == "male" {
				c.records[field][id].Malecancer++
			} else if sex == "female" {
				c.records[field][id].Femalecancer++
			}
			if mal, er := c.df.GetCell(idx, "Malignant"); er == nil {
				if mal == "1" {
					c.records[field][id].Malignant++
				} else if mal == "0" {
					c.records[field][id].Benign++
				}
			}
		}
	}
}

func (c *cancerRates) countRecords() {
	// Counts the number of total, adult, and adult cancer occurances by species
	var field, id string
	for idx := range c.df.Rows {
		if c.species {
			field = TID
			id, _ = c.df.GetCell(idx, c.key)
		} else {
			field, _ = c.df.GetCell(idx, c.key)
			id, _ = c.df.GetCell(idx, c.sec)
		}
		if _, ex := c.records[field][id]; ex {
			// Increment total
			c.records[field][id].Total++
			if nec, _ := c.df.GetCell(idx, "Necropsy"); nec == "1" {
				c.records[field][id].Necropsy++
			}
			age, err := c.df.GetCellFloat(idx, "Age")
			if err == nil {
				// Increment adult if age is greater than age of infancy
				c.records[field][id].Age = c.records[field][id].Age + age
				sex, er := c.df.GetCell(idx, "Sex")
				if er == nil {
					if sex == "male" {
						c.records[field][id].Male++
					} else if sex == "female" {
						c.records[field][id].Female++
					}
					c.countNeoplasia(idx, field, id, sex, age)
				}
			}
		}
	}
}

func (c *cancerRates) setTumorRecords() {
	// Stores map of empty species records with >= min occurances
	for idx := range c.df.Rows {
		field, _ := c.df.GetCell(idx, c.key)
		id, _ := c.df.GetCell(idx, c.sec)
		if field != "NA" && id != "NA" {
			if _, ex := c.records[id]; !ex {
				c.records[field] = make(map[string]*Record)
				c.records[field]["total"] = NewRecord()
			}
			if _, ex := c.records[field][id]; !ex {
				c.records[field][id] = NewRecord()
			}
		}
	}
}

//----------------------------------------------------------------------------

func (c *cancerRates) appendLifeHistory() {
	// Determines age of infancy and adds life history if needed
	lifehist := codbutils.ToMap(c.db.GetRows("Life_history", TID, getRecKeys(c.records[TID]), "*"))
	for k, v := range c.records[TID] {
		if lh, ex := lifehist[k]; ex {
			v.Lifehistory = lh
		} else {
			v.Lifehistory = c.nas
		}
	}
}

func (c *cancerRates) addDenominators() {
	// Adds fixed values from denominators table
	for k, v := range codbutils.ToMap(c.db.GetTable("Denominators")) {
		if _, ex := c.records[TID][k]; ex {
			if t, err := strconv.Atoi(v[0]); err == nil {
				c.records[TID][k].Total += t
			}
		}
	}
}

func (c *cancerRates) setTaxonomy(idx int) []string {
	// Stores taxonomy for given record
	var ret []string
	for _, k := range strings.Split(c.db.Columns["Taxonomy"], ",") {
		if k != "Source" && k != TID {
			val, _ := c.df.GetCell(idx, k)
			ret = append(ret, val)
		}
	}
	return ret
}

func (c *cancerRates) setTaxaRecords() {
	// Sets records by taxa_ id
	c.records[TID] = make(map[string]*Record)
	c.records[TID]["total"] = NewRecord()
	// Store blank taxonomy to preserve spacing
	c.records[TID]["total"].setTaxonomy([]string{"", "", "", "", "", "", ""})
	for idx := range c.df.Rows {
		id, _ := c.df.GetCell(idx, c.key)
		if _, ex := c.records[TID][id]; !ex {
			c.records[TID][id] = NewRecord()
			c.records[TID][id].setTaxonomy(c.setTaxonomy(idx))
		}
	}
	c.addDenominators()
	if c.lh {
		c.appendLifeHistory()
	}
}

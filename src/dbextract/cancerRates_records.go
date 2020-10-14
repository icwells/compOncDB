// Contains functions for populating records map

package dbextract

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"strconv"
	"strings"
)

func (c *cancerRates) getTotals() {
	// Adds all sub-records to field total
	id := c.total
	for key, val := range c.records {
		for k, v := range val {
			if k != id {
				c.records[key][id].Add(v)
			}
		}
	}
}

func (c *cancerRates) countNeoplasia(idx int, field, id, service, sex string, age float64) {
	// Counts cancer related data
	if mp, err := c.df.GetCell(idx, "Masspresent"); err == nil {
		if mp == "1" {
			// Increment cancer count and age if masspresent == true
			c.records[id][field].allcancer++
			mal, _ := c.df.GetCell(idx, "Malignant")
			if mal == "1" {
				c.records[id][field].maltotal++
			} else if mal == "0" {
				c.records[id][field].bentotal++
			}
			if service != "MSU" {
				c.records[id][field].Cancer++
				c.records[id][field].Cancerage += age
				if sex == "male" {
					c.records[id][field].Malecancer++
				} else if sex == "female" {
					c.records[id][field].Femalecancer++
				}
				if mal == "1" {
					c.records[id][field].Malignant++
				} else if mal == "0" {
					c.records[id][field].Benign++
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
			id, _ = c.df.GetCell(idx, c.key)
			field = c.total
		} else {
			id, _ = c.df.GetCell(idx, c.key)
			field, _ = c.df.GetCell(idx, c.sec)
		}
		if _, ex := c.records[id][field]; ex {
			// Increment total
			c.records[id][field].grandtotal++
			age, _ := c.df.GetCellFloat(idx, "Age")
			service, _ := c.df.GetCell(idx, "service_name")
			sex, _ := c.df.GetCell(idx, "Sex")
			source, _ := c.df.GetCell(idx, "account_id")
			c.records[id][field].Sources.Add(source)
			c.countNeoplasia(idx, field, id, service, sex, age)
			if service != "MSU" {
				c.records[id][field].Total++
				if nec, _ := c.df.GetCell(idx, "Necropsy"); nec == "1" {
					c.records[id][field].Necropsy++
				}
				// Increment adult if age is greater than age of infancy
				c.records[id][field].Age = c.records[id][field].Age + age
				if sex == "male" {
					c.records[id][field].Male++
				} else if sex == "female" {
					c.records[id][field].Female++
				}
			}
		}
	}
}

//----------------------------------------------------------------------------

func (c *cancerRates) appendLifeHistory() {
	// Determines age of infancy and adds life history if needed
	lifehist := codbutils.ToMap(c.db.GetRows("Life_history", TID, strings.Join(c.tids, ","), "*"))
	for k, v := range c.records {
		if lh, ex := lifehist[k]; ex {
			v[c.total].Lifehistory = lh
		} else {
			v[c.total].Lifehistory = c.nas
		}
	}
}

func (c *cancerRates) addDenominators() {
	// Adds fixed values from denominators table
	for k, v := range codbutils.ToMap(c.db.GetTable("Denominators")) {
		if _, ex := c.records[k]; ex {
			if t, err := strconv.Atoi(v[0]); err == nil {
				c.records[k][c.total].Total += t
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

func (c *cancerRates) setRecords() {
	// Stores map of empty species records with >= min occurances
	blank := make([]string, len(c.nas))
	for idx := range c.df.Rows {
		if id, err := c.df.GetCell(idx, c.key); err == nil {
			if _, ex := c.records[id]; !ex {
				c.records[id] = make(map[string]*Record)
				c.records[id][c.total] = NewRecord()
				c.records[id][c.total].setTaxonomy(c.setTaxonomy(idx))
				c.tids = append(c.tids, id)
			}
			if !c.species {
				field, _ := c.df.GetCell(idx, c.sec)
				if _, ex := c.records[id][field]; !ex {
					c.records[id][field] = NewRecord()
					c.records[id][field].setTaxonomy(blank[:7])
					c.records[id][field].Lifehistory = blank
				}
			}
		}
	}
	c.addDenominators()
	if c.lh {
		c.appendLifeHistory()
	}
}

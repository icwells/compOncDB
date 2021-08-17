// Includes calculation functions for cancerrates

package cancerrates

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"strings"
)

func (c *cancerRates) formatRates() {
	// Calculates rates, and formats for printing
	for k, v := range c.Records {
		if v.total.total >= c.min {
			c.ids.Add(k)
			for _, i := range v.ToSlice(c.keep) {
				if len(i) > 0 {
					// Add to dataframe
					err := c.rates.AddRow(i)
					if err != nil {
						c.logger.Printf("Adding row to dataframe: %v\n", err)
						break
					} else {
						c.species++
					}
				}
			}
		}
	}
}

func (c *cancerRates) CountRecords() {
	// Counts Patient records
	source := codbutils.ToMap(c.db.GetColumns("Source", []string{"ID", "service_name", "account_id", "Approved", "Aza", "Zoo", "Institute"}))
	diagnosis := codbutils.ToMap(c.db.GetColumns("Diagnosis", []string{"ID", "Masspresent", "Necropsy"}))
	tumor := search.TumorMap(c.db)
	for _, i := range c.db.GetRows("Patient", TID, strings.Join(c.tids, ","), "ID,Sex,Age,Infant,Wild,"+TID) {
		var location string
		s := c.Records[i[5]]
		id := i[0]
		acc := source[id]
		diag := diagnosis[id]
		if c.checkSettings(i[3], i[4], acc[0], acc[2], acc[3], acc[4], acc[5], diag[1]) {
			allrecords := c.checkService(acc[0], "")
			if c.checkService(acc[0], diag[0]) {
				// Add non-cancer values (skips non-cancer msu records)
				s.addNonCancer(allrecords, i[2], i[1], diag[1], acc[0], acc[1])
			}
			if diag[0] == "1" {
				if v, ex := tumor[id]; ex {
					// Add tumor values and add tissue denominator
					s.addCancer(allrecords, i[2], i[1], diag[1], v[1], v[3], acc[0], acc[1])
					location = v[3]
				} else {
					// Add values where masspresent is known, but further diagnosis data is missing
					s.addCancer(allrecords, i[2], i[1], diag[1], "-1", "", acc[0], acc[1])
				}
			}
			if allrecords {
				s.addDenominator(diag[0], location)
			}
		}
	}
}

func (c *cancerRates) GetTaxa(eval string) {
	// Gets records map
	var taxa map[string][]string
	if eval != "" && eval != "nil" {
		var e codbutils.Evaluation
		e.SetOperation(eval)
		taxa = codbutils.ToMap(c.db.GetRows("Taxonomy", e.Column, e.Value, strings.Join(c.header[:8], ",")))
	} else {
		taxa = codbutils.ToMap(c.db.GetColumns("Taxonomy", c.header[:8]))
	}
	for k, v := range taxa {
		c.tids = append(c.tids, k)
		c.Records[k] = newSpecies(k, c.location, v)
	}
	if c.lh {
		c.addLifeHistory()
	}
	//c.addDenominators()
}

func GetCancerRates(db *dbIO.DBIO, min, nec int, inf, lh, wild, keepall bool, zoo, eval, location string) *dataframe.Dataframe {
	// Returns dataframe of cancer rates
	c := NewCancerRates(db, min, nec, inf, lh, wild, keepall, zoo, location)
	c.logger.Printf("Calculating rates for species with at least %d entries...\n", c.min)
	c.GetTaxa(eval)
	c.CountRecords()
	c.formatRates()
	c.logger.Printf("Found %d species with at least %d entries.\n", c.species, c.min)
	c.setMetaData(eval)
	return c.rates
}

func GetRatesAndRecords(db *dbIO.DBIO, min, nec int, inf, lh, keepall bool, zoo, eval, location string) (*dataframe.Dataframe, *dataframe.Dataframe) {
	// Returns dataframe of cancer rates and pathology reports used to caclulate them
	c := NewCancerRates(db, min, nec, inf, lh, false, keepall, zoo, location)
	c.logger.Printf("Calculating rates for species with at least %d entries...\n", c.min)
	c.GetTaxa(eval)
	c.CountRecords()
	c.formatRates()
	c.logger.Printf("Found %d species with at least %d entries.\n", c.species, c.min)
	c.setMetaData(eval)
	ids := fmt.Sprintf("taxa_id ^ (%s)", strings.Join(c.ids.ToStringSlice(), ";"))
	res, msg := search.SearchRecords(db, c.logger, ids, false, true)
	c.logger.Println(msg)
	return c.rates, res
}

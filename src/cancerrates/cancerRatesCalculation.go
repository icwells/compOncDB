// Contains methods for calculating cancer rates

package cancerrates

import (
	"fmt"
	"github.com/icwells/go-tools/dataframe"
	"strings"
)

func (c *cancerRates) checkService(service, masspresent string) bool {
	// Returns true if record should be counted (skips non-cancer msu, national zoo, and zeps records)
	var ret bool
	if masspresent == "1" {
		ret = true
	} else if SERVICES.AllRecords(service) {
		ret = true
	}
	return ret
}

func (c *cancerRates) formatRates() {
	// Calculates rates, and formats for printing
	for _, v := range c.Records {
		if v.total.total >= c.min {
			for _, i := range v.ToSlice(c.keep, c.age, c.sex) {
				if len(i) > 0 {
					// Add to dataframe
					if err := c.rates.AddRow(i); err != nil {
						fmt.Println(c.header)
						fmt.Println(i)
						c.logger.Fatalf("Adding row to dataframe: %v\n", err)
						break
					} else {
						c.species++
					}
				}
			}
		} else {
			// Remove records from search results that don't meet the minimum
			for _, i := range v.ids.ToStringSlice() {
				c.search.DeleteRow(i)
			}
		}
	}
}

func (c *cancerRates) getSpecies(s *dataframe.Series, tid string) *Species {
	// Initializes records entry, stores taxonomy and life history, and returns species entry
	if _, ex := c.Records[tid]; !ex {
		var cols, taxa []string
		if c.taxa {
			// Store complete taxonomy
			cols = H.Taxonomy[1 : len(H.Taxonomy)-1]
		} else {
			// Store species and common name
			cols = H.Taxonomy[7 : len(H.Taxonomy)-1]
		}
		for _, i := range cols {
			v, _ := s.GetCell(i)
			taxa = append(taxa, v)
		}
		// Initialize new species entry
		c.Records[tid] = newSpecies(tid, c.location, taxa)
		if c.lh {
			// Store life history
			c.Records[tid].lifehistory = append(c.Records[tid].lifehistory, DASH)
			for _, i := range H.Life_history[1:] {
				if strings.Contains(i, "(") {
					i = i[:strings.Index(i, "(")]
				}
				v, _ := s.GetCell(i)
				c.Records[tid].lifehistory = append(c.Records[tid].lifehistory, v)
			}
		}
	}
	return c.Records[tid]
}

func (c *cancerRates) CountRecords() {
	// Counts Patient records
	for i := range c.search.Iterate() {
		sex, _ := i.GetCell("Sex")
		age, _ := i.GetCell("age_months")
		tid, _ := i.GetCell(TID)
		service, _ := i.GetCell("service_name")
		aid, _ := i.GetCell("account_id")
		mass, _ := i.GetCell("Masspresent")
		nec, _ := i.GetCell("Necropsy")
		mal, _ := i.GetCell("Malignant")
		loc, _ := i.GetCell(c.lcol)
		s := c.getSpecies(i, tid)
		allrecords := c.checkService(service, "")
		if c.checkService(service, mass) {
			// Add non-cancer values (skips records from services without denominators)
			s.addNonCancer(allrecords, age, sex, nec, service, aid, i.Name)
		}
		if mass == "1" {
			// Add tumor values and add tissue denominator
			s.addCancer(allrecords, age, sex, nec, mal, loc, service, aid)
		}
		if allrecords {
			s.addDenominator(service, loc)
		}
	}
}

func (c *cancerRates) GetCancerRates(eval string) (*dataframe.Dataframe, *dataframe.Dataframe) {
	// Returns dataframe of cancer rates
	c.setDataFrame()
	c.logger.Printf("Calculating rates for species with at least %d entries...\n", c.min)
	c.SetSearch(eval)
	c.CountRecords()
	c.formatRates()
	c.logger.Printf("Found %d species with at least %d entries.\n", c.species, c.min)
	return c.rates, c.search
}

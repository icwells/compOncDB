// This script will calculate cancer rates for species with  at least a given number of entries

package cancerrates

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/simpleset"
	"log"
	"strconv"
	"strings"
)

var (
	SERVICES = codbutils.NewServices()
	TID      = "taxa_id"
)

func checkService(service, masspresent string) bool {
	// Returns true if record should be counted (skips non-cancer msu and national zoo records)
	var ret bool
	if masspresent == "1" || SERVICES.HasDenominators(service) {
		ret = true
	}
	return ret
}

type cancerRates struct {
	approval *simpleset.Set
	approved bool
	aza      bool
	db       *dbIO.DBIO
	header   []string
	infant   bool
	lh       bool
	location string
	logger   *log.Logger
	min      int
	nas      []string
	nec      int
	rates    *dataframe.Dataframe
	Records  map[string]*Species
	species  int
	tids     []string
	total    string
	zoo      bool
}

func NewCancerRates(db *dbIO.DBIO, min, nec int, inf, lh, appr, aza, zoo bool, location string) *cancerRates {
	// Returns initialized cancerRates struct
	idx := 0
	c := new(cancerRates)
	c.approved = appr
	c.aza = aza
	c.db = db
	c.infant = inf
	c.location = location
	c.lh = lh
	c.logger = codbutils.GetLogger()
	c.min = min
	c.nec = nec
	c.setHeader()
	if location != "" {
		// Don't store by index when repeated taxa_ids are present
		idx = -1
	}
	c.rates, _ = dataframe.NewDataFrame(idx)
	c.rates.SetHeader(c.header)
	c.Records = make(map[string]*Species)
	c.total = "total"
	c.zoo = zoo
	return c
}

func (c *cancerRates) setHeader() {
	// Stores target column name
	c.header = codbutils.CancerRateHeader()
	tail := strings.Split(c.db.Columns["Life_history"], ",")[1:]
	for i := 0; i < len(tail); i++ {
		c.nas = append(c.nas, "NA")
	}
	if c.lh {
		c.header = append(c.header, tail...)
	}
}

func (c *cancerRates) formatRates() {
	// Calculates rates, and formats for printing
	for _, v := range c.Records {
		if v.total.total >= c.min {
			var err error
			for _, i := range v.ToSlice() {
				if len(i) > 0 {
					// Add to dataframe
					err = c.rates.AddRow(i)
					if err != nil {
						c.logger.Printf("Adding row to dataframe: %v\n", err)
						break
					}
				}
			}
			if err == nil {
				c.species++
			}
		}
	}
}

func (c *cancerRates) checkNecropsy(service, nec string) bool {
	// Returns if records should be processed
	var ret bool
	if c.nec == 0 {
		ret = true
	} else if c.nec == 1 && nec == "1" {
		ret = SERVICES.AllRecords(service)
	} else if c.nec == -1 && nec != "1" {
		ret = SERVICES.AllRecords(service)
	}
	return ret
}

func (c *cancerRates) checkSource(approved, aza, zoo string) bool {
	// Compares source information to filtering settings
	ret := true
	if c.approved && approved != "1" {
		ret = false
	} else if c.aza && aza != "1" {
		ret = false
	} else if c.zoo && zoo != "1" {
		ret = false
	}
	return ret
}

func (c *cancerRates) CountRecords() {
	// Counts Patient records
	source := codbutils.ToMap(c.db.GetColumns("Source", []string{"ID", "service_name", "account_id", "Approved", "Aza", "Zoo"}))
	diagnosis := codbutils.ToMap(c.db.GetColumns("Diagnosis", []string{"ID", "Masspresent", "Necropsy"}))
	tumor := search.TumorMap(c.db)
	for _, i := range c.db.GetRows("Patient", TID, strings.Join(c.tids, ","), "ID,Sex,Age,Infant,"+TID) {
		var location string
		s := c.Records[i[4]]
		id := i[0]
		acc := source[id]
		if c.checkSource(acc[2], acc[3], acc[4]) {
			// Ignore infant records if infant flag not set
			if c.infant || i[3] != "1" {
				if diag, ex := diagnosis[id]; ex {
					// Compare record against necropsy settings
					if c.checkNecropsy(acc[0], diag[1]) {
						if checkService(acc[0], diag[0]) {
							// Add non-cancer values (skips non-cancer msu records)
							s.addNonCancer(i[2], i[1], diag[1], acc[0], acc[1])
						}
						if diag[0] == "1" {
							if v, ex := tumor[id]; ex {
								// Add tumor values and add tissue denominator
								s.addCancer(i[2], i[1], diag[1], v[1], v[3], acc[0], acc[1])
								location = v[3]
							} else {
								// Add values where masspresent is known, but further diagnosis data is missing
								s.addCancer(i[2], i[1], diag[1], "-1", "", acc[0], acc[1])
							}
						}
						if checkService(acc[0], "") {
							s.addDenominator(diag[0], location)
						}
					}
				}
			}
		}
	}
}

func (c *cancerRates) addDenominators() {
	// Adds fixed values from denominators table
	if c.nec == 0 {
		for k, v := range codbutils.ToMap(c.db.GetRows("Denominators", TID, strings.Join(c.tids, ","), "*")) {
			if _, ex := c.Records[k]; ex {
				if t, err := strconv.Atoi(v[0]); err == nil {
					c.Records[k].total.addTotal(t)
				}
			}
		}
	}
}

func (c *cancerRates) addLifeHistory() {
	// Add life history data
	lifehist := codbutils.ToMap(c.db.GetRows("Life_history", TID, strings.Join(c.tids, ","), "*"))
	for k, v := range c.Records {
		if lh, ex := lifehist[k]; ex {
			v.lifehistory = lh
		} else {
			v.lifehistory = c.nas
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
	c.addDenominators()
}

func GetCancerRates(db *dbIO.DBIO, min, nec int, inf, lh, appr, aza, zoo bool, eval, location string) *dataframe.Dataframe {
	// Returns dataframe of cancer rates
	c := NewCancerRates(db, min, nec, inf, lh, appr, aza, zoo, location)
	c.logger.Printf("Calculating rates for species with at least %d entries...\n", c.min)
	c.GetTaxa(eval)
	c.CountRecords()
	c.formatRates()
	c.logger.Printf("Found %d species with at least %d entries.\n", c.species, c.min)
	return c.rates
}

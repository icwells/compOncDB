// This script will calculate cancer rates for species with  at least a given number of entries

package cancerrates

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"log"
	"strconv"
	"strings"
)

var TID = "taxa_id"

type cancerRates struct {
	db       *dbIO.DBIO
	header   []string
	infant   bool
	lh       bool
	location string
	logger   *log.Logger
	min      int
	nas      []string
	nec      bool
	rates    *dataframe.Dataframe
	records  map[string]*species
	species  int
	tids     []string
	total    string
}

func newCancerRates(db *dbIO.DBIO, min int, nec, inf, lh bool, location string) *cancerRates {
	// Returns initialized cancerRates struct
	c := new(cancerRates)
	c.db = db
	c.infant = inf
	c.location = location
	c.lh = lh
	c.logger = codbutils.GetLogger()
	c.min = min
	c.nec = nec
	c.setHeader()
	c.rates, _ = dataframe.NewDataFrame(-1)
	c.rates.SetHeader(c.header)
	c.records = make(map[string]*species)
	c.total = "total"
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
	for _, v := range c.records {
		if v.total.total >= c.min {
			var err error
			for _, i := range v.toSlice() {
				// Add to dataframe
				err = c.rates.AddRow(i)
				if err != nil {
					c.logger.Printf("Adding row to dataframe: %v\n", err)
					break
				}
			}
			if err == nil {
				c.species++
			}
		}
	}
}

func (c *cancerRates) countRecords() {
	// Counts Patient records
	source := codbutils.ToMap(c.db.GetColumns("Source", []string{"ID", "service_name", "account_id"}))
	diagnosis := codbutils.ToMap(c.db.GetColumns("Diagnosis", []string{"ID", "Masspresent", "Necropsy"}))
	tumor := codbutils.ToMap(c.db.GetColumns("Tumor", []string{"ID", "Malignant", "Location"}))
	for _, i := range c.db.GetRows("Patient", TID, strings.Join(c.tids, ","), "ID,Sex,Age,"+TID) {
		s := c.records[i[3]]
		id := i[0]
		if f, err := strconv.ParseFloat(i[2], 64); err == nil {
			// Ignore infant records if infant flag not set
			if c.infant || f >= s.infancy {
				diag := diagnosis[id]
				// Subset necropsy records if nec == true
				if !c.nec || diag[1] == "1" {
					acc := source[id]
					// Add non-cancer values
					s.addNonCancer(f, i[1], diag[1], acc[0], acc[1])
					if diag[0] == "1" {
						if v, ex := tumor[id]; ex {
							// Add tumor values
							s.addCancer(f, i[1], diag[1], v[0], v[1], acc[0], acc[1])
						}
					}
				}
			}
		}
	}
}

func (c *cancerRates) addDenominators() {
	// Adds fixed values from denominators table
	for k, v := range codbutils.ToMap(c.db.GetRows("Denominators", TID, strings.Join(c.tids, ","), "*")) {
		if _, ex := c.records[k]; ex {
			if t, err := strconv.Atoi(v[0]); err == nil {
				c.records[k].addDenominator(t)
			}
		}
	}
}

func (c *cancerRates) addLifeHistory() {
	// Add life history data
	lifehist := codbutils.ToMap(c.db.GetRows("Life_history", TID, strings.Join(c.tids, ","), "*"))
	for k, v := range c.records {
		if lh, ex := lifehist[k]; ex {
			v.lifehistory = lh
		} else {
			v.lifehistory = c.nas
		}
	}
}

func (c *cancerRates) addInfancy() {
	// Adds age of infancy to records
	for k, v := range dbextract.GetMinAges(c.db, c.tids) {
		if r, ex := c.records[k]; ex {
			r.infancy = v
		}
	}
}

func (c *cancerRates) getTaxa(eval string) {
	// Gets records map
	var taxa map[string][]string
	if eval != "" {
		var e codbutils.Evaluation
		e.SetOperation(eval)
		taxa = codbutils.ToMap(c.db.GetRows("Taxonomy", e.Column, e.Value, strings.Join(c.header[:8], ",")))
	} else {
		taxa = codbutils.ToMap(c.db.GetColumns("Taxonomy", c.header[:8]))
	}
	for k, v := range taxa {
		c.tids = append(c.tids, k)
		c.records[k] = newSpecies(k, c.location, v)
	}
	if !c.infant {
		c.addInfancy()
	}
	if c.lh {
		c.addLifeHistory()
	}
	c.addDenominators()
}

func GetCancerRates(db *dbIO.DBIO, min int, nec, inf, lh bool, eval, location string) *dataframe.Dataframe {
	// Returns dataframe of cancer rates
	c := newCancerRates(db, min, nec, inf, lh, location)
	c.logger.Printf("Calculating rates for species with at least %d entries...\n", c.min)
	c.getTaxa(eval)
	c.countRecords()
	c.formatRates()
	c.logger.Printf("Found %d species with at least %d entries.\n", c.species, c.min)
	return c.rates
}

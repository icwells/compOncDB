// This script will calculate cancer rates for species with  at least a given number of entries

package cancerrates

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/simpleset"
	"log"
	"strings"
)

var (
	H        = codbutils.NewHeaders()
	SERVICES = codbutils.NewServices()
	TID      = "taxa_id"
)

type cancerRates struct {
	approval *simpleset.Set
	db       *dbIO.DBIO
	header   []string
	infant   bool
	keep     bool
	lcol     string
	lh       bool
	location string
	logger   *log.Logger
	min      int
	nas      []string
	nec      int
	rates    *dataframe.Dataframe
	Records  map[string]*Species
	search   *dataframe.Dataframe
	species  int
	tids     []string
	total    string
	wild     bool
	zoo      string
}

func NewCancerRates(db *dbIO.DBIO, min, nec int, inf, lh, wild, keepall bool, zoo, tissue, location string) *cancerRates {
	// Returns initialized cancerRates struct
	idx := 0
	c := new(cancerRates)
	c.db = db
	c.infant = inf
	c.keep = keepall
	if tissue != "" {
		c.lcol = "Tissue"
		c.location = tissue
	} else {
		c.lcol = "Location"
		c.location = location
	}
	c.lh = lh
	c.logger = codbutils.GetLogger()
	c.min = min
	c.nec = nec
	c.setHeader()
	if c.location != "" {
		// Don't store by index when repeated taxa_ids are present
		idx = -1
	}
	c.rates, _ = dataframe.NewDataFrame(idx)
	c.rates.SetHeader(c.header)
	c.Records = make(map[string]*Species)
	c.total = "total"
	c.wild = wild
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

func (c *cancerRates) ChangeLocation(loc string) {
	// Prepares struct to analyze a new location
	c.location = loc
	c.Records = make(map[string]*Species)
}

func (c *cancerRates) setMetaData(eval string) {
	// Stores search options as string
	var m []string
	var nec string
	switch c.nec {
	case 1:
		nec = "Necropsy"
	case 0:
		nec = "All"
	case -1:
		nec = "NonNecropsy"
	}
	m = append(m, codbutils.GetTimeStamp())
	if eval != "" && eval != "nil" {
		m = append(m, eval)
	}
	m = append(m, fmt.Sprintf("%s=%s", c.lcol, c.location))
	m = append(m, fmt.Sprintf("min=%d", c.min))
	m = append(m, fmt.Sprintf("necropsyStatus=%s", nec))
	m = append(m, fmt.Sprintf("SourceType=%s", c.zoo))
	m = append(m, fmt.Sprintf("KeepInfantRecords=%v", c.infant))
	m = append(m, fmt.Sprintf("KeepWildRecords=%v", c.wild))
	c.rates.SetMetaData(strings.Join(m, ","))
}

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

func (c *cancerRates) SetSearch(eval string) {
	// Sets dataframe using filtering options
	var msg string
	// Store metadata before adding to string
	c.setMetaData(eval)
	switch c.nec {
	case 1:
		eval += ",Necropsy=1"
	case -1:
		eval += ",Necropsy=0"
	}
	if !c.infant {
		eval += ",Infant!=1"
	} else {
		eval += ",Infant=1"
	}
	if c.wild {
		eval += ",Wild=1"
	} else {
		eval += ",Wild!=1"
	}
	switch c.zoo {
	case "approved":
		eval += ",Approved=1"
	case "aza":
		eval += ",Aza=1"
	case "zoo":
		eval += ",Zoo=1"
	}
	if eval[0] == ',' {
		// Remove initial comma
		eval = eval[1:]
	}
	c.search, msg = search.SearchRecords(c.db, c.logger, eval, c.infant, c.lh)
	c.logger.Print(msg)
}

func (c *cancerRates) formatRates() {
	// Calculates rates, and formats for printing
	for _, v := range c.Records {
		if v.total.total >= c.min {
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
		} else {
			// Remove records from search results that don't meet the minimum
			for _, i := range v.ids.ToStringSlice() {
				c.search.DeleteRow(i)
			}
		}
	}
}

func (c *cancerRates) getSpecies(k, tid string) *Species {
	// Initializes records entry, stores taxonomy and life history, and returns species entry
	if _, ex := c.Records[tid]; !ex {
		var taxa []string
		// Store taxonomy
		for _, i := range H.Taxonomy[1 : len(H.Taxonomy)-1] {
			v, _ := c.search.GetCell(k, i)
			taxa = append(taxa, v)
		}
		// Initialize new species entry
		c.Records[tid] = newSpecies(tid, c.location, taxa)
		if c.lh {
			// Store life history
			for _, i := range H.Life_history[1:] {
				if strings.Contains(i, "(") {
					i = i[:strings.Index(i, "(")]
				}
				v, _ := c.search.GetCell(k, i)
				if v[0] == '%' {
					v = "-1"
				}
				c.Records[tid].lifehistory = append(c.Records[tid].lifehistory, v)
			}
		}
	}
	return c.Records[tid]
}

func (c *cancerRates) CountRecords() {
	// Counts Patient records
	for k := range c.search.Index {
		sex, _ := c.search.GetCell(k, "Sex")
		age, _ := c.search.GetCell(k, "age_months")
		tid, _ := c.search.GetCell(k, TID)
		service, _ := c.search.GetCell(k, "service_name")
		aid, _ := c.search.GetCell(k, "account_id")
		mass, _ := c.search.GetCell(k, "Masspresent")
		nec, _ := c.search.GetCell(k, "Necropsy")
		mal, _ := c.search.GetCell(k, "Malignant")
		loc, _ := c.search.GetCell(k, c.lcol)
		s := c.getSpecies(k, tid)
		allrecords := c.checkService(service, "")
		if c.checkService(service, mass) {
			// Add non-cancer values (skips records from services without denominators)
			s.addNonCancer(allrecords, age, sex, nec, service, aid, k)
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

func GetCancerRates(db *dbIO.DBIO, min, nec int, inf, lh, wild, keepall bool, zoo, eval, tissue, location string) (*dataframe.Dataframe, *dataframe.Dataframe) {
	// Returns dataframe of cancer rates
	c := NewCancerRates(db, min, nec, inf, lh, wild, keepall, zoo, tissue, location)
	c.logger.Printf("Calculating rates for species with at least %d entries...\n", c.min)
	c.SetSearch(eval)
	c.CountRecords()
	c.formatRates()
	c.logger.Printf("Found %d species with at least %d entries.\n", c.species, c.min)
	return c.rates, c.search
}

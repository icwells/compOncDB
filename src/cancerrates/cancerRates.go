// Defines cnacer rate struct and getting/setting methods

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
	DASH     = "-"
	H        = codbutils.NewHeaders()
	SERVICES = codbutils.NewServices()
	TID      = "taxa_id"
)

type cancerRates struct {
	age      bool
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
	sex      bool
	species  int
	taxa     bool
	tids     []string
	total    string
	wild     bool
	zoo      string
}

func NewCancerRates(db *dbIO.DBIO, min int, keepall bool, tissue, location string) *cancerRates {
	// Returns initialized cancerRates struct
	c := new(cancerRates)
	c.db = db
	c.keep = keepall
	c.taxa = true
	if tissue != "" {
		c.lcol = "Tissue"
		c.location = tissue
	} else {
		c.lcol = "Location"
		c.location = location
	}
	c.logger = codbutils.GetLogger()
	c.min = min
	c.Records = make(map[string]*Species)
	c.total = "total"
	return c
}

func (c *cancerRates) SearchSettings(nec int, inf, wild bool, zoo string) {
	// Stores settings for searching database
	c.nec = nec
	c.infant = inf
	c.wild = wild
	c.zoo = zoo
}

func (c *cancerRates) OutputSettings(age, lifehistory, sex, taxonomy bool) {
	// Stores output file settings
	c.age = age
	c.lh = lifehistory
	c.sex = sex
	c.taxa = taxonomy
}

func (c *cancerRates) setHeader() {
	// Stores target column name
	var loc bool
	if c.location != "" {
		loc = true
	}
	c.header = codbutils.CancerRateHeader(c.age, c.lh, loc, c.sex, c.taxa)
}

func (c *cancerRates) setDataFrame() {
	// Initializes header and dataframe
	idx := 0
	if c.location != "" {
		// Don't store by index when repeated taxa_ids are present
		idx = -1
	}
	c.setHeader()
	c.rates, _ = dataframe.NewDataFrame(idx)
	c.rates.SetHeader(c.header)
}

//----------------------------------------------------------------------------

func (c *cancerRates) ChangeLocation(loc string, typ bool) {
	// Prepares struct to analyze a new location
	if typ {
		c.lcol = "Type"
	}
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
	if c.location != "" {
		m = append(m, fmt.Sprintf("%s=%s", c.lcol, c.location))
	}
	m = append(m, fmt.Sprintf("min=%d", c.min))
	m = append(m, fmt.Sprintf("necropsyStatus=%s", nec))
	m = append(m, fmt.Sprintf("SourceType=%s", c.zoo))
	m = append(m, fmt.Sprintf("KeepInfantRecords=%v", c.infant))
	m = append(m, fmt.Sprintf("KeepWildRecords=%v", c.wild))
	c.rates.SetMetaData(strings.Join(m, ","))
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
	c.logger.Println(strings.TrimSpace(msg))
}

// This script will calculate cancer rates for species with  at least a given number of entries

package cancerrates

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/simpleset"
	"log"
	"strings"
)

var (
	SERVICES = codbutils.NewServices()
	TID      = "taxa_id"
)

type cancerRates struct {
	approval *simpleset.Set
	db       *dbIO.DBIO
	header   []string
	ids      *simpleset.Set
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
	wild     bool
	zoo      string
}

func NewCancerRates(db *dbIO.DBIO, min, nec int, inf, lh, wild bool, zoo, location string) *cancerRates {
	// Returns initialized cancerRates struct
	idx := 0
	c := new(cancerRates)
	c.db = db
	c.ids = simpleset.NewStringSet()
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

func (c *cancerRates) checkNecropsy(service, nec string) bool {
	// Returns true if records should be processed
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

func (c *cancerRates) checkSource(approved, aza, zoo, inst string) bool {
	// Compares source information to filtering settings
	var ret bool
	switch c.zoo {
	case "all":
		ret = true
	case "approved":
		if approved == "1" {
			ret = true
		}
	case "aza":
		if aza == "1" {
			ret = true
		}
	case "noprivate":
		if zoo == "1" || inst == "1" {
			ret = true
		}
	case "zoo":
		if zoo == "1" {
			ret = true
		}
	}
	return ret
}

func (c *cancerRates) checkSettings(infant, wild, service, approved, aza, zoo, inst, nec string) bool {
	// Returns true if record should be analyzed
	var ret bool
	if c.checkSource(approved, aza, zoo, inst) && c.checkNecropsy(service, nec) {
		if c.infant || infant != "1" {
			if c.wild && wild == "1" {
				ret = true
			} else if !c.wild && wild != "1" {
				ret = true
			}
		}
	}
	return ret
}

/*func (c *cancerRates) addDenominators() {
	// Adds fixed values from denominators table
	if c.nec == 0 && c.zoo == "all" {
		for k, v := range codbutils.ToMap(c.db.GetRows("Denominators", TID, strings.Join(c.tids, ","), "*")) {
			if _, ex := c.Records[k]; ex {
				if t, err := strconv.Atoi(v[0]); err == nil {
					c.Records[k].total.addTotal(t)
				}
			}
		}
	}
}*/

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
	if c.location != "" {
		m = append(m, fmt.Sprintf("tissue=%s", c.location))
	}
	if eval != "" && eval != "nil" {
		m = append(m, eval)
	}
	m = append(m, fmt.Sprintf("min=%d", c.min))
	m = append(m, fmt.Sprintf("necropsyStatus=%s", nec))
	m = append(m, fmt.Sprintf("SourceType=%s", c.zoo))
	m = append(m, fmt.Sprintf("KeepInfantRecords=%v", c.infant))
	c.rates.SetMetaData(strings.Join(m, ","))
}

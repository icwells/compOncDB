// This script will calculate cancer rates for species with  at least a given number of entries

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"strings"
)

var TID = "taxa_id"

type cancerRates struct {
	db      *dbIO.DBIO
	df      *dataframe.Dataframe
	header  []string
	infant  bool
	key     string
	lh      bool
	min     int
	nas     []string
	rates   *dataframe.Dataframe
	records map[string]map[string]*Record
	sec     string
	species bool
	total   string
}

func newCancerRates(db *dbIO.DBIO, min int, lh, inf bool, typ string) *cancerRates {
	// Returns initialized cancerRates struct
	c := new(cancerRates)
	c.setKey(typ)
	c.db = db
	c.min = min
	c.infant = inf
	c.header = codbutils.CancerRateHeader(c.sec)
	c.lh = lh
	if c.lh {
		// Omit taxa_id column
		tail := strings.Split(c.db.Columns["Life_history"], ",")[1:]
		c.header = append(c.header, tail...)
		for i := 0; i < len(tail); i++ {
			c.nas = append(c.nas, "NA")
		}
	}
	c.records = make(map[string]map[string]*Record)
	c.rates, _ = dataframe.NewDataFrame(-1)
	c.rates.SetHeader(c.header)
	c.total = "total"
	return c
}

func (c *cancerRates) setKey(t string) {
	// Stores target column name
	c.key = TID
	switch strings.ToLower(t) {
	case "location":
		c.sec = "Location"
	case "type":
		c.sec = "Type"
	default:
		c.species = true
	}
}

func (c *cancerRates) calculateRates(v *Record, name string) {
	// Calclates cancer rates and adds to dataframe
	if v.Total >= c.min {
		// Calculate cancer rates
		r := v.CalculateRates(name, c.lh)
		// Add to dataframe
		err := c.rates.AddRow(r)
		if err != nil {
			fmt.Printf("\t[Error] Adding row to dataframe: %v\n", err)
		}
	}
}

func (c *cancerRates) formatRates() {
	// Calculates rates, and formats for printing
	if len(c.records) > 0 {
		for _, val := range c.records {
			if !c.species {
				c.calculateRates(val[c.total], c.total)
				for k, v := range val {
					if k != c.total {
						c.calculateRates(v, k)
					}
				}
			} else {
				// Omit location/type column
				c.calculateRates(val[c.total], "")
			}
		}
	}
}

func (c *cancerRates) setDataframe(eval [][]codbutils.Evaluation, nec bool) {
	// Gets dataframe of matching records
	if nec {
		e := codbutils.SetOperations(c.db.Columns, "Necropsy = 1")
		for idx := range eval {
			eval[idx] = append(eval[idx], e[0][0])
		}
	} else if len(eval) == 0 {
		// Set evaluation to return everything
		eval = codbutils.SetOperations(c.db.Columns, "ID > 0")
	}
	c.df, _ = SearchColumns(c.db, "", eval, c.infant)
}

func GetCancerRates(db *dbIO.DBIO, min int, nec, inf, lh bool, eval [][]codbutils.Evaluation, typ string) *dataframe.Dataframe {
	// Returns slice of string slices of cancer rates and related info
	c := newCancerRates(db, min, lh, inf, typ)
	fmt.Printf("\n\tCalculating rates for species with at least %d entries...\n", c.min)
	c.setDataframe(eval, nec)
	c.setRecords()
	c.countRecords()
	if !c.species {
		c.getTotals()
	}
	c.formatRates()
	fmt.Printf("\tFound %d records with at least %d entries.\n", c.rates.Length(), c.min)
	return c.rates
}

func SearchCancerRates(db *dbIO.DBIO, min int, nec, inf, lh bool, typ, eval, infile string) *dataframe.Dataframe {
	// Wraps call to GetCancerRates
	var e [][]codbutils.Evaluation
	if eval != "nil" {
		e = codbutils.SetOperations(db.Columns, eval)
	} else if infile != "nil" {
		e = codbutils.OperationsFromFile(db.Columns, infile)
	}
	return GetCancerRates(db, min, nec, inf, lh, e, typ)
}

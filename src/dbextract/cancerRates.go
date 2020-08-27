// This scrpt will calculate cancer rates for species with  at least a given number of entries

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"strconv"
	"strings"
)

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
	records map[string]*Record
}

func newCancerRates(db *dbIO.DBIO, min int, lh, inf bool, typ string) *cancerRates {
	// Returns initialized cancerRates struct
	c := new(cancerRates)
	c.setKey(typ)
	c.db = db
	c.min = min
	if c.key == "species" {
		c.lh = lh
	}
	c.infant = inf
	c.header = codbutils.CancerRateHeader(c.key)
	if c.lh {
		// Omit taxa_id column
		tail := strings.Split(c.db.Columns["Life_history"], ",")[1:]
		c.header = append(c.header, tail...)
		for i := 0; i < len(tail); i++ {
			c.nas = append(c.nas, "NA")
		}
	}
	c.records = make(map[string]*Record)
	c.rates, _ = dataframe.NewDataFrame(0)
	c.rates.SetHeader(c.header)
	return c
}

func (c *cancerRates) setKey(t string) {
	// Stores target column name
	switch strings.ToLower(t) {
	case "location":
		c.key = "Location"
	case "type":
		c.key = "Type"
	default:
		c.key = "taxa_id"
	}
}

func (c *cancerRates) formatRates() {
	// Adds taxonomies to structs, calculates rates, and formats for printing
	if len(c.records) > 0 {
		for k, v := range c.records {
			if v.Total >= c.min {
				// Calculate cancer rates
				r := v.CalculateRates(k, c.lh)
				// Add to dataframe
				err := c.rates.AddRow(r)
				if err != nil {
					fmt.Printf("\t[Error] Adding row to dataframe: %v\n", err)
				}
			}
		}
	}
}

func (c *cancerRates) countNeoplasia(idx int, tid, sex string, age float64) {
	// Counts cancer related data
	if mp, err := c.df.GetCell(idx, "Masspresent"); err == nil {
		if mp == "1" {
			// Increment cancer count and age if masspresent == true
			c.records[tid].Cancer++
			c.records[tid].Cancerage = c.records[tid].Cancerage + age
			if sex == "male" {
				c.records[tid].Malecancer++
			} else if sex == "female" {
				c.records[tid].Femalecancer++
			}
			if mal, er := c.df.GetCell(idx, "Malignant"); er == nil {
				if mal == "1" {
					c.records[tid].Malignant++
				}
			}
		}
	}
}

func (c *cancerRates) countRecords() {
	// Counts the number of total, adult, and adult cancer occurances by species
	for idx := range c.df.Rows {
		tid, _ := c.df.GetCell(idx, "taxa_id")
		if _, ex := c.records[tid]; ex {
			// Increment total
			c.records[tid].Total++
			if nec, _ := c.df.GetCell(idx, "Necropsy"); nec == "1" {
				c.records[tid].Necropsy++
			}
			age, err := c.df.GetCellFloat(idx, "Age")
			if err == nil {
				// Increment adult if age is greater than age of infancy
				c.records[tid].Age = c.records[tid].Age + age
				sex, er := c.df.GetCell(idx, "Sex")
				if er == nil {
					if sex == "male" {
						c.records[tid].Male++
					} else if sex == "female" {
						c.records[tid].Female++
					}
					c.countNeoplasia(idx, tid, sex, age)
				}
			}
		}
	}
}

func (c *cancerRates) appendLifeHistory() {
	// Determines age of infancy and adds life history if needed
	lifehist := codbutils.ToMap(c.db.GetRows("Life_history", "taxa_id", getRecKeys(c.records), "*"))
	for k, v := range c.records {
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
		if _, ex := c.records[k]; ex {
			if t, err := strconv.Atoi(v[0]); err == nil {
				c.records[k].Total += t
			}
		}
	}
}

func (c *cancerRates) setTaxonomy(idx int) []string {
	// Stores taxonomy for given record
	var ret []string
	for _, k := range strings.Split(c.db.Columns["Taxonomy"], ",") {
		if k != "Source" && k != "taxa_id" {
			val, _ := c.df.GetCell(idx, k)
			ret = append(ret, val)
		}
	}
	return ret
}

func (c *cancerRates) setRecords() {
	// Stores map of empty species records with >= min occurances
	for idx := range c.df.Rows {
		id, err := c.df.GetCell(idx, c.key)
		fmt.Println(c.key, id, err)
		if id != "NA" && id != "-1" {
			if _, ex := c.records[id]; !ex {
				c.records[id] = NewRecord()
				if c.key == "taxa_id" {
					c.records[id].setTaxonomy(c.setTaxonomy(idx))
				}
			}
		}
	}
	if c.key == "taxa_id" {
		c.addDenominators()
		if c.lh {
			c.appendLifeHistory()
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

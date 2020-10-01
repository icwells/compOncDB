// Converts bmr data in life history table from watts to mLO2/hr

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	//"strconv"
	"time"
)

var (
	infile = kingpin.Flag("infile", "Path to input file (if using).").Short('i').Default("nil").String()
	user   = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()
)

type converter struct {
	db     *dbIO.DBIO
	factor float64
	table  [][]string
	taxa   map[string]string
}

func newConverter() *converter {
	// Returns initialized converter struct
	c := new(converter)
	c.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false))
	c.factor = 1.0 / 1.68875
	//c.table = c.db.EvaluateRows("Life_history", "metabolic_rate", ">", "-1", "taxa_id,metabolic_rate")
	c.taxa = codbutils.EntryMap(c.db.GetColumns("Taxonomy", []string{"taxa_id", "Species"}))
	return c
}

func (c *converter) upload() {
	// Updates life history table with converted values
	fmt.Println("\tUpdating Life History table...")
	for _, i := range c.table {
		c.db.UpdateRow("Life_history", "metabolic_rate", i[1], "taxa_id", "=", i[0])
	}
}

/*func (c *converter) convertBMR() {
	// Converts Watts to mLO2/hr
	var count int
	for _, i := range c.table {
		if w, err := strconv.ParseFloat(i[1], 64); err == nil {
			ml := w * c.factor
			i[1] = strconv.FormatFloat(ml, 'f', -1, 64)
			count++
		}
	}
	fmt.Printf("\tSuccessfully formatted %d of %d entries.\n", count, len(c.table))
}*/

func (c *converter) getBMR() {
	// Stores BMR data from pantheria file
	rows, header := iotools.YieldFile(*infile, true)
	for i := range rows {
		species := fmt.Sprintf("%s %s", i[header["MSW05_Genus"]], i[header["MSW05_Species"]])
		bmr := i[header["18-1_BasalMetRate_mLO2hr"]]
		if bmr != "-999" {
			if tid, ex := c.taxa[species]; ex {
				c.table = append(c.table, []string{tid, bmr})
			}
		}
	}
	fmt.Printf("\tFound %d BMR entries.\n", len(c.table))
}

func main() {
	kingpin.Parse()
	c := newConverter()
	fmt.Println("\n\tConverting Watts to mLO2/hr...")
	c.getBMR()
	//c.convertBMR()
	c.upload()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(c.db.Starttime))
}

// Converts anage interbirth interval in life history table from days to months

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"strconv"
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
	c.factor = 30.42
	//c.table = c.db.EvaluateRows("Life_history", "metabolic_rate", ">", "-1", "taxa_id,metabolic_rate")
	c.taxa = codbutils.EntryMap(c.db.GetColumns("Taxonomy", []string{"taxa_id", "Species"}))
	return c
}

func (c *converter) upload() {
	// Updates life history table with converted values
	fmt.Println("\tUpdating Life History table...")
	for _, i := range c.table {
		c.db.UpdateRow("Life_history", "interbirth_interval", i[1], "taxa_id", "=", i[0])
	}
}

func (c *converter) convertToMonths(days string) string {
	// Converts days to months
	var ret string
	if d, err := strconv.ParseFloat(days, 64); err == nil && d > 0 {
		ret = strconv.FormatFloat(d/c.factor, 'f', 2, 64)
	}
	return ret
}

func (c *converter) getIBI() {
	// Stores BMR data from anage file
	rows, header := iotools.YieldFile(*infile, true)
	for i := range rows {
		if len(i) > header["Inter-litter/Interbirth interval"] {
			species := fmt.Sprintf("%s %s", i[header["Genus"]], i[header["Species"]])
			if tid, ex := c.taxa[species]; ex {
				if days := c.convertToMonths(i[header["Inter-litter/Interbirth interval"]]); days != "" {
					c.table = append(c.table, []string{tid, days})
				}
			}
		}
	}
	fmt.Printf("\tFound %d IBI entries.\n", len(c.table))
}

func main() {
	kingpin.Parse()
	c := newConverter()
	fmt.Println("\n\tConverting inter-birth interval days to months...")
	c.getIBI()
	c.upload()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(c.db.Starttime))
}

// Summarizes the number of species with a minimum number of records with age information

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"strconv"
	"time"
)

var (
	column  = kingpin.Flag("column", "Name of column to summarize.").Short('c').Default("age_months").String()
	outfile = kingpin.Flag("outfile", "Name of output file.").Short('o').Required().String()
	user    = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type colSummary struct {
	header  string
	logger  *log.Logger
	min     map[int]int
	species map[string][]int
	steps   []int
	table   *dataframe.Dataframe
	total   map[int]int
}

func newColSummary() *colSummary {
	// Returns initialized struct
	c := new(colSummary)
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	c.header = "Min,TotalSpecies,SpeciesWithValues"
	c.logger = codbutils.GetLogger()
	c.min = make(map[int]int)
	c.species = make(map[string][]int)
	c.steps = []int{10, 15, 20, 25, 30, 40, 45, 50}
	c.table, _ = search.SearchRecords(db, c.logger, "Approved=1", false, false)
	c.total = make(map[int]int)
	for _, i := range c.steps {
		c.min[i] = 0
		c.total[i] = 0
	}
	return c
}

func (c *colSummary) write() {
	// Writes results to file
	c.logger.Println("Writing results to file...")
	var res [][]string
	for k, v := range c.min {
		res = append(res, []string{strconv.Itoa(k), strconv.Itoa(c.total[k]), strconv.Itoa(v)})
	}
	iotools.WriteToCSV(*outfile, c.header, res)
}

func (c *colSummary) getTotals() {
	// Counts number of species at each step
	c.logger.Println("Calculating species minimums...")
	for _, v := range c.species {
		for _, i := range c.steps {
			if v[0] >= i {
				c.total[i]++
				if v[1] >= 1 {
					c.min[i]++
				}
			}
		}
	}
}

func (c *colSummary) setSpecies() {
	// Creates entries for species
	c.logger.Println("Setting species slice...")
	for _, i := range c.table.Index {
		sp, _ := c.table.GetCell(i, "Species")
		if _, ex := c.species[sp]; !ex {
			c.species[sp] = []int{0, 0}
		}
		val, err := c.table.GetCell(i, *column)
		if err != nil {
			panic(err)
		}
		c.species[sp][0]++
		if val != "NA" && val != "-1" && val != "-1.00" && val != "" {
			c.species[sp][1]++
		}
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	c := newColSummary()
	c.setSpecies()
	c.getTotals()
	c.write()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

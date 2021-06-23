// Returns cancer rates for wild and captive records for each species

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/cancerrates"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"time"
)

var (
	outfile = kingpin.Flag("outfile", "Name of output file.").Short('o').Required().String()
	user    = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Required().String()
)

type species struct {
	captive, lifehist, taxa, wild []string
}

func (s *species) row() []string {
	// Returns concatenated string
	row := append(s.taxa, s.captive...)
	row = append(row, s.wild...)
	row = append(row, s.lifehist...)
	return row
}

//----------------------------------------------------------------------------

type captiveVsWild struct {
	approved string
	db       *dbIO.DBIO
	header   []string
	lh       bool
	lifehist []string
	metadata []string
	min      int
	nec      int
	output   [][]string
	records  map[string]*species
	target   []string
	taxa     []string
}

func newCaptiveVsWild() *captiveVsWild {
	c := new(captiveVsWild)
	c.approved = "all"
	c.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	c.lh = true
	c.lifehist = []string{"female_maturity(months)", "male_maturity(months)", "Gestation(months)", "Weaning(months)", "Infancy(months)", "litter_size", "litters_year",
		"interbirth_interval", "birth_weight(g)", "weaning_weight(g)", "adult_weight(g)", "growth_rate(1/days)", "max_longevity(months)", "metabolic_rate(mLO2/hr)"}
	c.min = 1
	c.nec = 0
	c.records = make(map[string]*species)
	c.target = []string{"RecordsWithDenominators", "NeoplasiaWithDenominators", "NeoplasiaPrevalence", "Malignant", "MalignancyPrevalence", "PropMalignant", "Necropsies"}
	c.taxa = []string{"taxa_id", "Kingdom", "Phylum", "Class", "Orders", "Family", "Genus", "Species"}
	c.setWild()
	c.setCaptive()
	return c
}

func (c *captiveVsWild) setHeader() string {
	// Sets outfile header
	var ret []string
	ret = append(ret, c.taxa...)
	prefix := []string{"Captive", "Wild"}
	prevalence := []string{"Records", "Neoplasia", "NeoplasiaPrevalence", "Malignant", "MalignancyPrevalence", "PropMalignant", "Necropsies"}
	for _, p := range prefix {
		for _, i := range prevalence {
			ret = append(ret, p+i)
		}
	}
	ret = append(ret, c.lifehist...)
	return strings.Join(ret, ",")
}

func (c *captiveVsWild) setWild() {
	// Pulls cancer rates for wild records
	var count int
	r := cancerrates.GetCancerRates(c.db, c.min, c.nec, true, c.lh, true, false, c.approved, "", "")
	for idx := range r.Index {
		count++
		c.records[idx] = new(species)
		sp := c.records[idx]
		for _, i := range strings.Split(c.db.Columns["Life_history"], ",")[1:] {
			v, _ := r.GetCell(idx, i)
			sp.lifehist = append(sp.lifehist, v)
		}
		for _, i := range c.target {
			v, _ := r.GetCell(idx, i)
			sp.wild = append(sp.wild, v)
		}
		// Append taxa_id
		sp.taxa = append(sp.taxa, idx)
		for _, i := range c.taxa[1:] {
			v, _ := r.GetCell(idx, i)
			sp.taxa = append(sp.taxa, v)
		}
	}
	fmt.Printf("\tFound %d wild species.\n", count)
}

func (c *captiveVsWild) setCaptive() {
	// Stores cancer rates for captive records if wild records were also found for that species
	var count int
	r := cancerrates.GetCancerRates(c.db, c.min, c.nec, false, c.lh, true, false, c.approved, "", "")
	for idx := range r.Index {
		if sp, ex := c.records[idx]; ex {
			for _, i := range c.target {
				v, _ := r.GetCell(idx, i)
				sp.captive = append(sp.captive, v)
			}
			c.output = append(c.output, sp.row())
			count++
		}
	}
	fmt.Printf("\tFound %d captive species.\n", count)
}

func main() {
	kingpin.Parse()
	c := newCaptiveVsWild()
	iotools.WriteToCSV(*outfile, c.setHeader(), c.output)
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(c.db.Starttime))
}

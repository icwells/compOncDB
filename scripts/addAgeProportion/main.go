// Appends proportion between age of diagnosis and max longevity to search results

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	infant   = kingpin.Flag("infant", "Include infant records in results (excluded by default).").Default("false").Bool()
	min      = kingpin.Flag("min", "Minimum number of entries required for calculations.").Short('m').Default("1").Int()
	nec      = kingpin.Flag("necropsy", "2: Extract only necropsy records, 1: extract all records by default, 0: extract non-necropsy records.").Default("2").Int()
	outfile  = kingpin.Flag("outfile", "Name of output file (writes to stdout if not given).").Short('o').Required().String()
	source   = kingpin.Flag("source", "Zoo/institute records to calculate prevalence with; all: use all records, approved: used zoos approved for publication, aza: use only AZA member zoos, zoo: use only zoos.").Short('z').Default("approved").String()
	user     = kingpin.Flag("user", "MySQL username.").Short('u').Default("").String()
	wild     = kingpin.Flag("wild", "Return results for wild records only (returns non-wild only by default).").Default("false").Bool()
)

type ageProportion struct {
	db			*dbIO.DBIO
	logger		*log.Logger
	longevity	map[string]float64
	records		*dataframe.Dataframe
}

func newAgeProportion() *ageProportion {
	// Returns initialized struct
	a := new(ageProportion)
	a.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	a.logger = codbutils.GetLogger()
	a.longevity = make(map[string]float64)
	a.setLongevity()
	return a
}

func (a *ageProportion) setLongevity() {
	// Formats longevity table
	a.logger.Println("Getting species longevity data...")
	for _, i := range a.db.GetColumns("Life_history", []string{"taxa_id", "max_longevity"}) {
		if val, err := strconv.ParseFloat(i[1], 64); err == nil {
			if val > 0.0 {
				a.longevity[i[0]] = val
			}
		}
	}
}

func (a *ageProportion) getProportion(tid string, age float64) string {
	// Returns age/max longevity as string
	ret := "NA"
	if age > 0.0 {
		if max, ex := a.longevity[tid]; ex {
			ret = strconv.FormatFloat(age/max, 'f', 4, 64)
		}
	}
	return ret
}

func (a *ageProportion) addProportion() {
	// Calculates age of diagnosis to max longevity proportion for all records
	var count int
	a.logger.Println("Calculating proportion of age of diagnosis over max longevity...")
	for i := range a.records.Iterate() {
		if age, err := i.GetCellFloat("age_months"); err == nil {
			tid, _ := i.GetCell("taxa_id")
			prop := a.getProportion(tid, age)
			if prop != "NA" {
				i.UpdateCell("proportion_longevity", prop)
				count++
			}
		} else {
			panic(err)
		}
	}
	a.logger.Printf("Calculated longevity proportion for %d records.", count)
}

func (a *ageProportion) filterMinSpecies() {
	// Removes Records from species with too few entries
	var rm []string
	counts := make(map[string]int)
	l := a.records.Length()
	for i := range a.records.Iterate() {
		tid, _ := i.GetCell("taxa_id")
		if _, ex := counts[tid]; !ex {
			counts[tid] = 0
		}
		counts[tid]++
	}
	for i := range a.records.Iterate() {
		tid, _ := i.GetCell("taxa_id")
		if counts[tid] < *min {
			rm = append(rm, i.Name)
		}
	}
	for _, i := range rm {
		a.records.DeleteRow(i)
	}
	a.logger.Printf("Removed %d records from species with fewer than %d records.", l - a.records.Length(), *min)
}

func (a *ageProportion) getSearchResults() {
	// Sets dataframe using filtering options
	var eval, msg string
	*nec--
	switch *nec {
	case 1:
		eval += ",Necropsy=1"
	case -1:
		eval += ",Necropsy=0"
	}
	if *infant {
		eval += ",Infant=1"
	} else {
		eval += ",Infant!=1"
	}
	if *wild {
		eval += ",Wild=1"
	} else {
		eval += ",Wild!=1"
	}
	switch *source {
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
	a.records, msg = search.SearchRecords(a.db, a.logger, eval, *infant, false)
	a.logger.Println(strings.TrimSpace(msg))
	a.filterMinSpecies()
	a.records.AddColumn("proportion_longevity", "NA")
}

func main() {
	start := time.Now()
	kingpin.Parse()
	a := newAgeProportion()
	a.getSearchResults()
	a.addProportion()
	a.records.ToCSV(*outfile)
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

// Returns cancer rates for given combination of types

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/cancerrates"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"gopkg.in/alecthomas/kingpin.v2"
	"sort"
	"strings"
	"time"
)

var (
	eval     = kingpin.Flag("eval", "Evaluation argument for taxonic level such that level=taxon (i.e. genus=canis).").Short('e').Default("").String()
	min      = kingpin.Flag("min", "Minimum number of records required for cancer rates.").Default("1").Int()
	necropsy = kingpin.Flag("necropsy", "2: extract only necropsy records, 0: extract only non-necropsy records.").Short('n').Default("1").Int()
	outfile  = kingpin.Flag("outfile", "Name of output file.").Short('o').Required().String()
	target   = kingpin.Flag("target", "Comma seperated list of tumor types to extact").Required().String()
	user     = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Required().String()
)

type record struct {
	other     *cancerrates.Species
	otherset  bool
	target    *cancerrates.Species
	targetset bool
}

func newRecord(s *cancerrates.Species) *record {
	// Returns new record struct
	r := new(record)
	return r
}

func (r *record) addTarget(s *cancerrates.Species) {
	// Adds s.tissue to target
	if !r.targetset {
		r.target = s
		r.target.Location = strings.Replace(*target, ",", ";", -1)
		r.targetset = true
	} else {
		r.target.AddTissue(s)
	}
}

func (r *record) addOther(s *cancerrates.Species) {
	// Adds s.tissue to other
	if !r.otherset {
		r.other = s
		r.other.Location = "other"
		r.otherset = true
	} else {
		r.other.AddTissue(s)
	}
}

func (r *record) format() [][]string {
	// Returns records as string slice
	var ret [][]string
	if target := r.target.ToSlice(false); len(target) > 0 {
		ret = append(ret, target[0])
		if r.targetset && len(target) > 1 {
			ret = append(ret, target[1])
		}
	}
	if r.otherset {
		other := r.other.ToSlice(false)
		if len(ret) == 0 && len(other) > 0 {
			ret = append(ret, other[0])
		}
		if len(other) >= 1 {
			ret = append(ret, other[1])
		}
	}
	return ret
}

type aggregator struct {
	db      *dbIO.DBIO
	records []*record
	taxa    map[string]*record
	target  []string
	types   []string
}

func newAggregator() *aggregator {
	*necropsy--
	a := new(aggregator)
	a.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	a.taxa = make(map[string]*record)
	a.setTypes()
	return a
}

func (a *aggregator) setTypes() {
	// Stores type and target slices
	for _, i := range strings.Split(*target, ",") {
		a.target = append(a.target, strings.TrimSpace(i))
	}
	for _, i := range a.db.Execute("SELECT DISTINCT(Type) FROM Tumor;") {
		if !strarray.InSliceStr(a.types, i[0]) {
			a.types = append(a.types, i[0])
		}
	}
}

func (a *aggregator) setTissues() {
	// Gets cancer rates for every tissue
	fmt.Println("\n\tCalculating cancer rates...")
	c := cancerrates.NewCancerRates(a.db, *min, *necropsy, false, true, false, false, "approved", "", "")
	c.SetSearch(*eval)
	for idx, list := range [][]string{a.target, a.types} {
		for _, i := range list {
			c.ChangeLocation(i, true)
			fmt.Printf("\tCalculating rates for %s...\n", i)
			c.CountRecords()
			for k, v := range c.Records {
				if v.Grandtotal > 0 {
					if _, ex := a.taxa[k]; !ex {
						a.taxa[k] = newRecord(v)
					}
					if idx == 0 {
						a.taxa[k].addTarget(v)
					} else {
						a.taxa[k].addOther(v)
					}
				}
			}
		}
	}
}

func (a *aggregator) Len() int {
	return len(a.records)
}

func (a *aggregator) Less(i, j int) bool {
	return a.records[i].target.Grandtotal > a.records[j].target.Grandtotal
}

func (a *aggregator) Swap(i, j int) {
	a.records[i], a.records[j] = a.records[j], a.records[i]
}

func (a *aggregator) sort() {
	// Sorts records slice by number of records
	fmt.Println("\tSorting results...")
	for _, v := range a.taxa {
		if v.target.Grandtotal >= *min {
			a.records = append(a.records, v)
		}
	}
	sort.Sort(a)
}

func (a *aggregator) printRecords() {
	// Writes records to file
	var res [][]string
	fmt.Println("\tFormatting results...")
	header := append(codbutils.CancerRateHeader(), strings.Split(a.db.Columns["Life_history"], ",")[1:]...)
	for _, v := range a.records {
		if row := v.format(); len(row) > 0 {
			res = append(res, row...)
		}
	}
	iotools.WriteToCSV(*outfile, strings.Join(header, ","), res)
}

func main() {
	kingpin.Parse()
	a := newAggregator()
	a.setTissues()
	a.sort()
	a.printRecords()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(a.db.Starttime))
}

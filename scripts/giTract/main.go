// Returns cancer rates for gi tract and other tissues

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/cancerrates"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
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
	password = kingpin.Flag("password", "Password (for testing or scripting).").Default("").String()
	repro    = kingpin.Flag("repro", "Extract reproductive tissues instead of gi tract.").Default("false").Bool()
	user     = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Required().String()
)

type record struct {
	gi       *cancerrates.Species
	giset    bool
	other    *cancerrates.Species
	otherset bool
}

func newRecord(s *cancerrates.Species) *record {
	// Returns new record struct
	r := new(record)
	return r
}

func (r *record) addGI(s *cancerrates.Species) {
	// Adds s.tissue to gi tract
	if !r.giset {
		r.gi = s
		if *repro {
			r.gi.Location = "reproductive"
		} else {
			r.gi.Location = "gi tract"
		}
		r.giset = true
	} else {
		r.gi.AddTissue(s)
	}
}

func (r *record) addOther(s *cancerrates.Species) {
	// Adds s.tissue to gi tract
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
	if gi := r.gi.ToSlice(false); len(gi) > 0 {
		ret = append(ret, gi[0])
		if r.giset && len(gi) > 1 {
			ret = append(ret, gi[1])
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

type gimerger struct {
	approved string
	db       *dbIO.DBIO
	gi       []string
	records  []*record
	repro    []string
	taxa     map[string]*record
	tissues  []string
}

func newGImerger() *gimerger {
	*necropsy--
	g := new(gimerger)
	g.approved = "approved"
	g.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)
	g.gi = []string{"liver", "bile duct", "gall bladder", "stomach", "small intestine", "colon", "esophagus", "oral", "duodenum"}
	g.repro = []string{"testis", "prostate", "ovary", "vulva", "uterus"}
	g.taxa = make(map[string]*record)
	g.tissues = []string{"abdomen", "fibrous", "myxomatous tissue", "fat", "notochord", "smooth muscle", "striated muscle", "peripheral nerve sheath", "blood", "cartilage", "synovium", "bone", "bone marrow", "lymph nodes", "spleen", "mast cell", "dendritic cell", "pigment cell", "skin", "hair follicle", "gland", "mammary", "glial cell", "meninges", "nerve cell", "pnet", "neuroepithelial", "spinal cord", "brain", "pituitary gland", "parathyroid gland", "thyroid", "adrenal medulla ", "adrenal cortex", "pancreas", "carotid body", "neuroendocrine", "kidney", "bladder", "oviduct", "iris", "pupil", "larynx", "trachea", "lung", "nose", "transitional epithelium", "mesothelium", "heart", "widespread"}
	if *repro {
		g.tissues = append(g.tissues, g.gi...)
		g.gi = g.repro
	} else {
		g.tissues = append(g.tissues, g.repro...)
	}
	return g
}

func (g *gimerger) setTissues() {
	// Gets cancer rates for every tissue
	fmt.Println("\n\tCalculating cancer rates...")
	c := cancerrates.NewCancerRates(g.db, *min, *necropsy, false, true, false, false, g.approved, "")
	c.SetSearch(*eval)
	for idx, list := range [][]string{g.gi, g.tissues} {
		for _, i := range list {
			c.ChangeLocation(i)
			fmt.Printf("\tCalculating rates for %s...\n", i)
			c.CountRecords()
			for k, v := range c.Records {
				if v.Grandtotal > 0 {
					if _, ex := g.taxa[k]; !ex {
						g.taxa[k] = newRecord(v)
					}
					if idx == 0 {
						g.taxa[k].addGI(v)
					} else {
						g.taxa[k].addOther(v)
					}
				}
			}
		}
	}
}

func (g *gimerger) Len() int {
	return len(g.records)
}

func (g *gimerger) Less(i, j int) bool {
	return g.records[i].gi.Grandtotal > g.records[j].gi.Grandtotal
}

func (g *gimerger) Swap(i, j int) {
	g.records[i], g.records[j] = g.records[j], g.records[i]
}

func (g *gimerger) sort() {
	// Sorts records slice by number of records
	fmt.Println("\tSorting results...")
	for _, v := range g.taxa {
		if v.gi.Grandtotal >= *min {
			g.records = append(g.records, v)
		}
	}
	sort.Sort(g)
}

func (g *gimerger) printRecords() {
	// Writes records to file
	var res [][]string
	fmt.Println("\tFormatting results...")
	header := append(codbutils.CancerRateHeader(), strings.Split(g.db.Columns["Life_history"], ",")[1:]...)
	for _, v := range g.records {
		if row := v.format(); len(row) > 0 {
			res = append(res, row...)
		}
	}
	iotools.WriteToCSV(*outfile, strings.Join(header, ","), res)
}

func main() {
	kingpin.Parse()
	g := newGImerger()
	g.setTissues()
	g.sort()
	g.printRecords()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(g.db.Starttime))
}

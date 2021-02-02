// Returns cancer rates for gi tract and other tissues

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
	outfile = kingpin.Flag("outfile", "Name of output file (writes to stdout if not given).").Short('o').Required().String()
	user    = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Required().String()
)

type record struct {
	gi       *cancerrates.Species
	giset    bool
	other    *cancerrates.Species
	otherset bool
	total    int
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
		r.gi.Location = "gi tract"
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
		r.gi.AddTissue(s)
	}
}

func (r *record) format() [][]string {
	// Returns records as string slice
	ret := r.gi.ToSlice()
	other := r.other.ToSlice()
	return append(ret, other[1])
}

type gimerger struct {
	db      *dbIO.DBIO
	gi      []string
	taxa    map[string]*record
	tissues []string
}

func newGImerger() *gimerger {
	g := new(gimerger)
	g.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	g.gi = []string{"liver", "bile duct", "gall bladder", "stomach", "small intestine", "colon", "esophagus", "oral", "duodenum", "abdomen"}
	g.taxa = make(map[string]*record)
	g.tissues = []string{"fibrous", "myxomatous tissue", "fat", "notochord", "smooth muscle", "striated muscle", "peripheral nerve sheath", "blood", "cartilage", "synovium", "bone", "bone marrow",
		"lymph nodes", "spleen", "mast cell", "dendritic cell", "pigment cell", "skin", "hair follicle", "gland", "mammary", "glial cell", "meninges", "nerve cell", "pnet", "neuroepithelial", "spinal cord", "brain", "pituitary gland", "parathyroid gland", "thyroid", "adrenal medulla ", "adrenal cortex", "pancreas", "carotid body", "neuroendocrine", "testis", "prostate", "ovary", "vulva", "uterus", "kidney", "bladder", "oviduct", "iris", "pupil", "larynx", "trachea", "lung", "nose", "transitional epithelium", "mesothelium", "heart", "widespread"}
	return g
}

func (g *gimerger) setTissues() {
	// Gets cancer rates for every tissue
	fmt.Println("\n\tCalculating cancer rates...")
	for idx, list := range [][]string{g.gi, g.tissues} {
		for _, i := range list {
			fmt.Printf("\tCalculating rates for %s...\n", i)
			c := cancerrates.NewCancerRates(g.db, 1, false, false, true, false, i)
			c.GetTaxa("")
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

func (g *gimerger) printRecords() {
	// Writes records to file
	var res [][]string
	fmt.Println("\tFormatting results...")
	header := append(codbutils.CancerRateHeader(), strings.Split(g.db.Columns["Life_history"], ",")[1:]...)
	for _, v := range g.taxa {
		res = append(res, v.format()...)
	}
	iotools.WriteToCSV(*outfile, strings.Join(header, ","), res)
}

func main() {
	kingpin.Parse()
	g := newGImerger()
	g.setTissues()
	g.printRecords()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(g.db.Starttime))
}

// Assigns records with no species to genus

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	infile  = kingpin.Flag("infile", "Path to taxonomy file.").Short('i').Required().String()
	outfile = kingpin.Flag("outfile", "Name of output file (writes to stdout if not given).").Short('o').Default("nil").String()
	user    = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()
)

type record struct {
	family, genus, id, name, taxaid string
}

func newRecord(id, name string) *record {
	// Returns initialized record
	r := new(record)
	r.id = id
	r.name = name
	return r
}

func (r *record) set() bool {
	// Return true if ids are filled
	if r.id != "" && r.taxaid != "" {
		return true
	}
	return false
}

func (r *record) slice() []string {
	// Returns tring slcie of output
	var ret []string
	ret = append(ret, r.id)
	ret = append(ret, r.taxaid)
	ret = append(ret, r.name)
	ret = append(ret, r.genus)
	return ret
}

type genera struct {
	db      *dbIO.DBIO
	genus   *dataframe.Dataframe
	header  string
	id      string
	key     string
	records []*record
	taxa    map[string][]string
}

func newGenera(db *dbIO.DBIO) *genera {
	// Returns iniitialzed struct
	g := new(genera)
	g.db = db
	g.genus, _ = dataframe.FromFile(*infile, 0)
	g.header = "ID,taxa_id,source_name,genus"
	g.id = "69"
	g.key = "NA"
	g.taxa = dbupload.ToMap(g.db.GetRows("Taxonomy", "Species", g.key, "taxa_id,Family,Genus"))
	g.setRecords()
	return g
}

func (g *genera) setRecords() {
	// Matches source names to genus and taxa_id
	fmt.Println("\n\tReading source names...")
	for _, i := range g.db.GetRows("Patient", "taxa_id", g.id, "ID,source_name") {
		r := newRecord(i[0], i[1])
		if _, ex := g.genus.Index[r.name]; ex {
			var err error
			r.genus, err = g.genus.GetCell(r.name, "Genus")
			if err == nil {
				r.family, _ = g.genus.GetCell(r.name, "Family")
				g.records = append(g.records, r)
			}
		}
	}
}

func (g *genera) getGenera() {
	// Gets taxa_id for source name and genus match
	fmt.Println("\tMerging name and genera..")
	for _, i := range g.records {
		for k, v := range g.taxa {
			if v[0] == i.family && v[1] == i.genus {
				i.taxaid = k
				break
			}
		}
	}
}

func (g *genera) writefile() {
	// Writes records to file
	var rec [][]string
	fmt.Println("\tWriting to file...")
	for _, i := range g.records {
		if i.set() {
			rec = append(rec, i.slice())
		}
	}
	iotools.WriteToCSV(*outfile, g.header, rec)
}

func main() {
	kingpin.Parse()
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false))
	g := newGenera(db)
	g.getGenera()
	g.writefile()
}

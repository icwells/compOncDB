// Summarizes neoplasia types by species

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/simpleset"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"strings"
	"time"
)

var (
	min     = kingpin.Flag("min", "Minimum number of entries required for calculations.").Short('m').Default("1").Int()
	outfile = kingpin.Flag("outfile", "Name of output file.").Short('o').Required().String()
	user    = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type typesSummary struct {
	benign    *simpleset.Set
	clades    []string
	db        *dbIO.DBIO
	header    []string
	logger    *log.Logger
	malignant *simpleset.Set
	records   *dataframe.Dataframe
	taxa      map[string]*species
	taxonomy  []string
}

func newTypesSummary() *typesSummary {
	// Returns initialized struct
	h := codbutils.NewHeaders()
	t := new(typesSummary)
	t.benign = simpleset.NewStringSet()
	t.clades = []string{"Amphibia", "Aves", "Mammalia", "Reptilia"}
	t.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	t.logger = codbutils.GetLogger()
	t.malignant = simpleset.NewStringSet()
	t.taxa = make(map[string]*species)
	t.taxonomy = h.Taxonomy[:len(h.Taxonomy)-1]
	t.header = append(t.taxonomy, []string{"TotalRecords", "NeoplasiaRecords", "NeoplasiaRate", "Malignant", "MalignantRate", "Benign", "BenignRate"}...)
	t.setRecords()
	return t
}

func (t *typesSummary) setRecords() {
	// Stores records from target clades
	t.logger.Println("Storing target records...")
	t.records, _ = dataframe.NewDataFrame(0)
	for idx, i := range t.clades {
		records, msg := search.SearchRecords(t.db, t.logger, fmt.Sprintf("Necropsy=1,Approved=1,Class=%s", i), false, false)
		t.logger.Println(msg)
		if idx == 0 {
			t.records = records
		} else {
			if err := t.records.Extend(records); err != nil {
				panic(err)
			}
		}
	}
}

func (t *typesSummary) write() {
	// Writes results to file
	var rows [][]string
	t.logger.Println("Writing results to file...")
	benign := t.benign.ToStringSlice()
	malignant := t.malignant.ToStringSlice()
	for _, i := range malignant {
		t.header = append(t.header, fmt.Sprintf("m-%s", i))
	}
	for _, i := range benign {
		t.header = append(t.header, fmt.Sprintf("b-%s", i))
	}
	for _, v := range t.taxa {
		rows = append(rows, v.toSlice(malignant, benign))
	}
	iotools.WriteToCSV(*outfile, strings.Join(t.header, ","), rows)
}

func (t *typesSummary) setTypes() {
	// Counts number of specific tumor types
	t.logger.Println("Counting neoplasia types...")
	for i := range t.records.Iterate() {
		tid, _ := i.GetCell("taxa_id")
		if v, ex := t.taxa[tid]; ex {
			if mp, _ := i.GetCellInt("Masspresent"); mp == 1 {
				mal, _ := i.GetCellInt("Malignant")
				types, _ := i.GetCell("Type")
				v.addNeoplasia(mal)
				for _, typ := range strings.Split(types, ";") {
					if mal == 1 {
						t.malignant.Add(typ)
					} else {
						t.benign.Add(typ)
					}
					v.addType(mal, typ)
				}
			}
		}
	}
}

func (t *typesSummary) setSpecies() {
	// Creates entries for species
	t.logger.Println("Setting species map...")
	for i := range t.records.Iterate() {
		tid, _ := i.GetCell("taxa_id")
		if _, ex := t.taxa[tid]; !ex {
			taxonomy := []string{tid}
			for _, t := range t.taxonomy[1:] {
				v, _ := i.GetCell(t)
				taxonomy = append(taxonomy, v)
			}
			t.taxa[tid] = newSpecies(taxonomy)
		}
		t.taxa[tid].total++
	}
	// Remove species without enough records
	if *min > 1 {
		for k, v := range t.taxa {
			if v.total < *min {
				delete(t.taxa, k)
			}
		}
	}
	t.logger.Printf("Found %d species with at least %d records.", len(t.taxa), *min)
}

func main() {
	start := time.Now()
	kingpin.Parse()
	t := newTypesSummary()
	t.setSpecies()
	t.setTypes()
	t.write()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

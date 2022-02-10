// Summarizes neoplasia types by species

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/simpleset"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"path"
	"strconv"
	"time"
)

var (
	min     = kingpin.Flag("min", "Minimum number of entries required for calculations.").Short('m').Default("1").Int()
	outfile = kingpin.Flag("outfile", "Name of output file.").Short('o').Required().String()
	user    = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type species struct {
	benign		map[string]int
	btotal		int
	cancer		int
	id			string
	malignant	map[string]int
	mtotal		int
	name		string
	species		string
	total		int
}

func newSpecies(species, name, id string) *species {
	// Returns initialized struct
	s := new(species)
	s.benign = make(map[string]int)
	s.diagnoses = codbutils.Getutils()
	s.id = id
	s.malignant = make(map[string]int)
	s.name = name
	s.species = species
	return s
}	

func (s *species) add(malignant int, typ string) {
	// Adds record to appropriate map
	if malignant == 1 {
		s.mtotal++
		if _, ex := s.malignant[typ]; !ex {
			s.malignant[typ]++
		}
	} else {
		s.btotal++
		if _, ex := s.benign[typ]; !ex {
			s.benign[typ]++
		}
	}
}

func (s *species) toSlice(malignant, benign []string) []string {
	// Returns values as string slice
	var ret []string
	ret = append(ret, s.id)

	return ret
}

type typesSummary struct {
	benign		*simpleset.Set
	clades		[]string
	db			*dbIO.DBIO
	header  	[]string
	logger  	*log.Logger
	malignant	*simpleset.Set
	records 	*dataframe.Dataframe
	taxa		map[string]*species
}

func newTypesSummary() *typesSummary {
	// Returns initialized struct
	t := new(typesSummary)
	t.clades = []string{"Amphibia", "Aves", "Mammalia", "Reptilia"}
	t.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")

	t.header = []string{"taxa_id", "Species", "Common"}

	t.logger = codbutils.GetLogger()
	t.min = *min
	t.taxa = make(map[string]*species)
	t.types = simpleset.NewStringSet()
	t.setRecords()
	return t
}

func (t *typesSummary) setRecords() {
	// Stores records from target clades
	t.logger.Println("Storing target records...")
	var res [][]string
	t.records, _ := dataframe.NewDataFrame(0)
	for _, i := range t.clades {
		records, err := search.SearchRecords(t.db, t.logger, fmt.Sprintf("Approved=1,Class=%s", i), false, false)
		if err != nil {
			t.logger.Fatal(err)
		}
		t.records.Extend(records)
	}
}

func (t *typesSummary) write() {
	// Writes results to file
	t.logger.Println("Writing results to file...")
	benign := t.benign.ToStringSlice()
	malignant := t.malignant.ToStringSlice()
	t.header = append(t.header, malignant...)

	t.header = append(t.header, benign...)
	
	iotools.WriteToCSV(*outfile, strings.Join(t.header), res)
}

func (t *typesSummary) setTypes() {
	// Counts number of specific tumor types
	t.logger.Println("Counting neoplasia types...")
	for _, i := range t.records.Iterate {
		sp, _ := i.GetCell("Species")
		if v, ex := t.taxa[sp]; ex {
			if mp, _ := i.GetCellInt("Masspresent"); i == 1 {
				mal, _ := i.GetCellInt("Malignant")
				typ, _ := i.GetCell("Type")
				if mal == 1 {
					t.malignant.Add(typ)
				} else {
					t.benign.Add(typ)
				}
				v.add(mal, typ)
			}
		}
	}
}

func (t *typesSummary) setSpecies() {
	// Creates entries for species
	t.logger.Println("Setting species map...")
	for _, i := range t.records.Iterate {
		name, _ := i.GetCell("Common")
		sp, _ := i.GetCell("Species")
		tid, _ := i.GetCell("taxa_id")
		if _, ex := t.taxa[tid]; !ex {
			t.taxa[tid] = newSpecies(sp, name, tid)
		}
		t.taxa[tid].total++
	}
	// Remove species without enough records
	if *min > 1 {
		for k, v := range t.taxa {
			if v.total < * min {
				delete(t.taxa, k)
			}
		}
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	t := newTypesSummary()
	t.setSpecies()
	t.getTotals()
	t.write()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

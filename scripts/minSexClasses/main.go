// Calculates prevalence for species in target classes

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/cancerrates"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/simpleset"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"path"
	"time"
)

var (
	min      = kingpin.Flag("min", "Minimum number of each male and female required for cancer rates.").Default("10").Int()
	necropsy = kingpin.Flag("necropsy", "2: extract only necropsy records, 0: extract only non-necropsy records.").Short('n').Default("2").Int()
	outdir   = kingpin.Flag("outdir", "Name of output directory.").Short('o').Required().String()
	user     = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Required().String()
)

type minSex struct {
	approved string
	classes  []string
	db       *dbIO.DBIO
	id       string
	logger   *log.Logger
	min      int
	outdir   string
	taxonomy []string
	tissue   string
}

func newMinSex() *minSex {
	*necropsy--
	m := new(minSex)
	m.approved = "approved"
	m.classes = []string{"Amphibia", "Aves", "Mammalia", "Reptilia"}
	m.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	m.id = "taxa_id"
	m.logger = codbutils.GetLogger()
	m.min = *min
	m.outdir, _ = iotools.FormatPath(*outdir, true)
	m.taxonomy = []string{"Kingdom", "Phylum", "Class", "Orders", "Family", "Genus", "Species", "common_name"}
	m.tissue = "Reproductive"
	return m
}

func (m *minSex) getOutfile(class, typ string) string {
	// Formats name for output file
	return path.Join(m.outdir, fmt.Sprintf("%s.%s.%s.csv", class, typ, codbutils.GetTimeStamp()))
}

func (m *minSex) insertTaxonomy(row *dataframe.Series, taxonomy []string) []string {
	// Inserts taxonomy and converts to slice
	for idx, i := range m.taxonomy {
		row.UpdateCell(i, taxonomy[idx])
	}
	return row.ToSlice()
}

func (m *minSex) splitRates(rates *dataframe.Dataframe) (*dataframe.Dataframe, *dataframe.Dataframe) {
	// Sorts reproductive and other rows into individual dataframes
	taxa := make(map[string][]string)
	repro := rates.Clone()
	other := rates.Clone()
	// Store taxonmies to fill in tissue specific columns
	for i := range rates.Iterate() {
		if loc, _ := i.GetCell("Location"); loc == "all" {
			id, _ := i.GetCell(m.id)
			taxa[id] = []string{}
			for _, t := range m.taxonomy {
				v, _ := i.GetCell(t)
				taxa[id] = append(taxa[id], v)
			}
		}
	}
	for i := range rates.Iterate() {
		id, _ := i.GetCell(m.id)
		loc, _ := i.GetCell("Location")
		row := m.insertTaxonomy(i, taxa[id])
		if loc == m.tissue {
			repro.AddRow(row)
		} else if loc == "Other" {
			other.AddRow(row)
		}
	}
	return repro, other
}

func (m *minSex) filterMinSex(rates *dataframe.Dataframe) *dataframe.Dataframe {
	// Removes rows from species with too few male/female records
	ret := rates.Clone()
	ids := simpleset.NewStringSet()
	for i := range rates.Iterate() {
		id, _ := i.GetCell(m.id)
		if loc, _ := i.GetCell("Location"); loc == "all" {
			female, _ := i.GetCellInt("Female")
			male, _ := i.GetCellInt("Male")
			if male < m.min && female < m.min {
				ids.Add(id)
				ret.AddRow(i.ToSlice())
			}
		} else if ex, _ := ids.InSet(id); ex {
			ret.AddRow(i.ToSlice())
		}
	}
	return ret
}

func (m *minSex) getRates() {
	// Gets cancer rates for every tissue
	for _, c := range m.classes {
		m.logger.Printf("Calculating cancer rates for %s...", c)
		eval := fmt.Sprintf("Class=%s", c)
		rates, _ := cancerrates.GetCancerRates(m.db, *min*2, *necropsy, false, false, false, true, m.approved, eval, m.tissue, "")
		rates = m.filterMinSex(rates)
		m.logger.Printf("Found %d species with greater than %d male and female records.", int(float64(rates.Length()/3.0)), m.min)
		rates.ToCSV(m.getOutfile(c, "All"))
		repro, nonrepro := m.splitRates(rates)
		repro.ToCSV(m.getOutfile(c, "Reproductive"))
		nonrepro.ToCSV(m.getOutfile(c, "Nonreproductive"))
	}
}

func main() {
	kingpin.Parse()
	m := newMinSex()
	m.getRates()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(m.db.Starttime))
}

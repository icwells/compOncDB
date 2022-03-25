// Returns cancer rates with metastasis and metastasis rate

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/cancerrates"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"strconv"
	"time"
)

var (
	min      = kingpin.Flag("min", "Minimum number of records required for cancer rates.").Default("10").Int()
	necropsy = kingpin.Flag("necropsy", "2: extract only necropsy records, 0: extract only non-necropsy records.").Short('n').Default("2").Int()
	outfile  = kingpin.Flag("outfile", "Name of output file.").Short('o').Required().String()
	source   = kingpin.Flag("source", "Zoo/institute records to calculate prevalence with; all: use all records, approved: used zoos approved for publication, aza: use only AZA member zoos, zoo: use only zoos.").Short('z').Default("approved").String()
	user     = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type metastasisRate struct {
	db         *dbIO.DBIO
	logger     *log.Logger
	mal        string
	met        string
	metastasis map[string]int
	rate       string
	rates      *dataframe.Dataframe
	records    *dataframe.Dataframe
	tid        string
}

func newMetastasisRate() *metastasisRate {
	*necropsy--
	m := new(metastasisRate)
	m.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	m.logger = codbutils.GetLogger()
	m.logger.Println("Initializing struct...")
	m.mal = "Malignant"
	m.met = "Metastasis"
	m.metastasis = make(map[string]int)
	m.rate = "MetastasisRate"
	m.rates, m.records = cancerrates.GetCancerRates(m.db, *min, *necropsy, false, false, false, false, *source, "", "", "")
	m.tid = "taxa_id"
	m.rates.AddColumn(m.met, "")
	m.rates.AddColumn(m.rate, "")
	m.setMetastasis()
	return m
}

func (m *metastasisRate) setMetastasis() {
	// Stores the number of metastases per species
	m.logger.Println("Counting metastases...")
	for i := range m.records.Iterate() {
		tid, _ := i.GetCell(m.tid)
		if _, ex := m.metastasis[tid]; !ex {
			m.metastasis[tid] = 0
		}
		mal, _ := i.GetCellInt(m.mal)
		met, _ := i.GetCellInt(m.met)
		if met == 1 && mal == 1 {
			m.metastasis[tid]++
		}
	}
}

func (m *metastasisRate) addMetastasis() {
	// Adds mets count to species and calculates rate
	var rm []string
	m.logger.Println("Calculating rates...")
	for i := range m.rates.Iterate() {
		if mal, _ := i.GetCellInt(m.mal); mal >= *min {
			met := m.metastasis[i.Name]
			m.rates.UpdateCell(i.Index, m.met, strconv.Itoa(met))
			rate := strconv.FormatFloat(float64(met)/float64(mal), 'f', 4, 64)
			m.rates.UpdateCell(i.Index, m.rate, rate)
		} else {
			rm = append(rm, i.Name)
		}
	}
	for _, i := range rm {
		m.rates.DeleteRow(i)
	}
	m.logger.Printf("Found %d species with at least %d malignancies.", m.rates.Length(), *min)
}

func main() {
	start := time.Now()
	kingpin.Parse()
	m := newMetastasisRate()
	m.addMetastasis()
	m.rates.ToCSV(*outfile)
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

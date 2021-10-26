// Reasigns age of infancy for species without weaining/maturity information

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"gopkg.in/alecthomas/kingpin.v2"
	"strconv"
	"strings"
	"time"
)

var user = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Required().String()

type infancy struct {
	db     *dbIO.DBIO
	infant map[string]float64
	prop   float64
}

func newInfancy() *infancy {
	// Initializes struct
	i := new(infancy)
	i.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	fmt.Println("\n\tInitializing struct...")
	i.infant = make(map[string]float64)
	return i
}

func (i *infancy) getAge(v string) float64 {
	// Cancerts string to float if possible
	if val, err := strconv.ParseFloat(v, 64); err == nil {
		return val
	}
	return -1.0
}

func (i *infancy) getTaxaIds() string {
	// Returns taxa_ids formatted for sql search
	var ret []string
	for k := range i.infant {
		ret = append(ret, k)
	}
	return strings.Join(ret, ",")
}

func (i *infancy) updatePatient() {
	// Updates infant flag if species infancy value has changed
	var add, remove int
	fmt.Println("\tUpdating Patient table...")
	inf := dbupload.NewInfancy(i.db)
	for _, v := range i.db.GetRows("Patient", "taxa_id", i.getTaxaIds(), "ID,Age,taxa_id,Infant,Comments") {
		if val := inf.SetInfant(v[2], v[1], v[4]); val != v[3] {
			i.db.UpdateRow("Patient", "Infant", val, "ID", "=", v[0])
			if val == "1" {
				add++
			} else if val == "0" {
				remove++
			}
		}
	}
	fmt.Printf("\tUpdated %d infant and %d non-infant records.\n", add, remove)
}

func (i *infancy) uploadInfancy() {
	// Uploads new infancy records
	var count int
	fmt.Println("\tUpdating Life_history table...")
	for k, v := range i.infant {
		count++
		i.db.UpdateRow("Life_history", "Infancy", strconv.FormatFloat(v, 'f', -1, 64), "taxa_id", "=", k)
		fmt.Printf("\tUpdated %d of %d records.\r", count, len(i.infant))
	}
}

func (i *infancy) setInfancy() {
	// Sets new infancy value for approriate taxa
	var count int
	fmt.Println("\tCalculating infancy for species missing maturity info...")
	for _, v := range i.db.EvaluateRows("Life_history", "Weaning", "<=", "0.0", "taxa_id,female_maturity,male_maturity,max_longevity") {
		if v[1] == "-1" && v[2] == "-1" {
			if l := i.getAge(v[3]); l > 0.0 {
				i.infant[v[0]] = l * i.prop
				count++
			}
		}
	}
	fmt.Printf("\tCalculated infancy for %d species.\n", count)
}

func (i *infancy) setProportion() {
	// Sets average proportion of weaning age/max_longevity
	var count int
	var val float64
	fmt.Println("\tCalculating average proportion...")
	for _, v := range i.db.EvaluateRows("Life_history", "Weaning", ">", "0.0", "Weaning,max_longevity") {
		w := i.getAge(v[0])
		l := i.getAge(v[1])
		if w > 0.0 && l > 0.0 {
			val += w / l
			count++
		}
	}
	i.prop = val / float64(count)
	fmt.Printf("\tCalculated %f from %d species.\n", i.prop, count)
}

func main() {
	start := time.Now()
	kingpin.Parse()
	i := newInfancy()
	i.setProportion()
	i.setInfancy()
	i.uploadInfancy()
	i.updatePatient()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

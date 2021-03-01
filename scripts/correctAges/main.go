// Corrects regex floating point error when assigning ages

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/diagnoses"
	"github.com/icwells/dbIO"
	"gopkg.in/alecthomas/kingpin.v2"
	"strconv"
	"time"
)

var (
	user    = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Required().String()
	outfile = kingpin.Flag("outfile", "Name of output file (writes to stdout if not given).").Short('o').Default("nil").String()
)

// ageCorrection is a struct for identifying incorrect ages and determining correct ones.
type ageCorrection struct {
	db      *dbIO.DBIO
	header  string
	matcher diagnoses.Matcher
	outfile string
	results [][]string
	table   [][]string
}

// newAgeCorrection connects to database and returns initialized struct.
func newAgeCorrection() *ageCorrection {
	a := new(ageCorrection)
	a.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false))
	a.header = "ID,Age,Comments"
	a.matcher = diagnoses.NewMatcher(codbutils.GetLogger())
	a.outfile = *outfile
	fmt.Println("\n\tExtracting table from database...")
	a.table = a.db.EvaluateRows("Patient", "Age", "!=", "NA", a.header)
	return a
}

// write prints output to screen/file.
func (a *ageCorrection) write() {
	codbutils.WriteResults(a.outfile, a.header, a.results)
}

// correctAges extracts correct ages from comments and stores in results slice.
func (a *ageCorrection) correctAges() {
	fmt.Println("\tCorrecting erroneous ages...")
	for _, i := range a.table {
		if i[2] != "NA" {
			if newage := a.matcher.GetAge(i[2]); newage != "-1" {
				na, _ := strconv.ParseFloat(newage, 64)
				age := strconv.FormatFloat(na, 'f', 2, 64)
				if age != i[1] {
					a.results = append(a.results, append(i, age))
				}
			}
		}
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	a := newAgeCorrection()
	a.correctAges()
	a.write()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

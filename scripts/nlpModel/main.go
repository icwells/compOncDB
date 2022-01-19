// Extracts and preformats data for use as nlpModel input

package main

import (
	"fmt"
	"log"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/diagnoses"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"time"
    "unicode"
)

var (
	outfile = kingpin.Flag("outfile", "Name of output file.").Short('o').Required().String()
	user    = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type formatter struct {
	columns	[]string
	db		*dbIO.DBIO
	match	*diagnoses.Matcher
	records	[][]string
}

func newFormatter() {
	// Connects to db and returns initialized struct
	f := new(formatter)
	f.columns = []string{"ID", "Comments", "Masspresent", "Hyperplasia", "Necropsy", "Metastasis", "primary_tumor", "Type", "Location", "service_name"}
	f.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)
	f.match = diagnosis.NewMatcher(codbutils.GetLogger())
	f.records = f.db.EvaluateRows("Records", "Comments", "!=", "NA", strings.Join(f.columns, ","))
	return f
}

func (f *formatter) inferSentences(val string) string {
	// Reinserts periods that were removed during data cleanup
	val = strings.Replace(val, ";", ".", -1)
	s := strings.Split(val, " ")
	if len(s) > 1 {
		for idx, i := range s[:len(s) - 1]:
			if v := s[idx + 1]; len(v) > 1 {
				if unicode.IsUpper(v[0]) && !unicode.IsUpper(v[1]):
					val[idx] += '.'
				}
			}
		}
	}
	return strings.Join(s, " ")
}

func (f *formatter) checkMass() {


	tumorType, tissue, location, malignant = f.match.GetTumor(line, rec.sex, cancer)
}

func (f *formatter) subsetTumors() {

}

func (f *formatter) formatRows() {
	// Preformats data for nlp modeling
	for _, i := range f.records {

	}
}

func main() {
	start := time.Now()
	f := newFormatter()

	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

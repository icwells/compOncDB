// Extracts and preformats data for use as nlpModel input

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/diagnoses"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/strarray"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"strings"
	"time"
	"unicode"
)

var (
	outfile = kingpin.Flag("outfile", "Name of output file.").Short('o').Required().String()
	user    = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type formatter struct {
	columns []string
	logger  *log.Logger
	match   diagnoses.Matcher
	records *dataframe.Dataframe
}

func newFormatter() *formatter {
	// Connects to db and returns initialized struct
	var msg string
	f := new(formatter)
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	f.columns = []string{"Sex", "Comments", "Masspresent", "Metastasis", "primary_tumor", "Type", "Location"}
	f.logger = codbutils.GetLogger()
	f.match = diagnoses.NewMatcher(f.logger)
	f.logger.Println("Extracting records from database...")
	f.records, msg = search.SearchRecords(db, f.logger, "Comments!=NA", true, false)
	f.logger.Println(msg)
	for k := range f.records.Header {
		if !strarray.InSliceStr(f.columns, k) {
			f.records.DeleteColumn(k)
		}
	}
	return f
}

func (f *formatter) write() {
	// Writes records to outfile
	f.logger.Println("Writing to file...")
	f.records.SetMetaData("")
	f.records.DeleteColumn("Sex")
	f.records.ToCSV(*outfile)
}

func (f *formatter) inferSentences(val string) string {
	// Reinserts periods that were removed during data cleanup
	val = strings.Replace(val, ";", ".", -1)
	s := strings.Split(val, " ")
	if len(s) > 1 {
		for idx, i := range s[:len(s)-1] {
			if v := s[idx+1]; len(v) > 1 {
				if unicode.IsUpper(rune(v[0])) && !unicode.IsUpper(rune(v[1])) {
					i += "."
				}
			}
		}
	}
	return strings.Join(s, " ")
}

func (f *formatter) formatRows() {
	// Preformats data for nlp modeling
	f.logger.Println("Formatting comments...")
	for idx := range f.records.Index {
		comments, _ := f.records.GetCell(idx, "Comments")
		sex, _ := f.records.GetCell(idx, "Sex")
		mp, _ := f.records.GetCell(idx, "Masspresent")
		//service, _ := f.records.GetCell(idx, "service_name")
		typ, _ := f.records.GetCell(idx, "Type")
		loc, _ := f.records.GetCell(idx, "Location")
		if strings.Contains(typ, ";") || strings.Contains(loc, ";") {
			f.records.DeleteRow(idx)
		} else {
			tumor, _, _, _ := f.match.GetTumor(comments, sex, true)
			if tumor == "NA" && mp == "1" {
				// Masspresent only equals 1 if it is identifiable from the comments
				f.records.UpdateCell(idx, "Masspresent", tumor)
			}
			f.records.UpdateCell(idx, "Comments", f.inferSentences(comments))
		}
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	f := newFormatter()
	f.formatRows()
	f.write()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

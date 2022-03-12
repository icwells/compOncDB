// Extracts and preformats data for use as nlpModel input

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/diagnoses"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/iotools"
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
	results [][]string
}

func newFormatter() *formatter {
	// Connects to db and returns initialized struct
	var msg string
	f := new(formatter)
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	f.columns = []string{"Sex", "Comments", "Masspresent", "Hyperplasia", "Type", "Location"}
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
	iotools.WriteToCSV(*outfile, strings.Join(f.columns[1:], ","), f.results)
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
	var count int
	f.logger.Println("Formatting comments...")
	for idx := range f.records.Index {
		comments, _ := f.records.GetCell(idx, "Comments")
		sex, _ := f.records.GetCell(idx, "Sex")
		comments = f.inferSentences(comments)
		for _, i := range f.match.SplitOnTumors(comments, sex) {
			f.results = append(f.results, i)
		}
		count++
	}
	f.logger.Printf("Formatted %d of %d records.", len(f.results), count)
}

func main() {
	start := time.Now()
	kingpin.Parse()
	f := newFormatter()
	f.formatRows()
	f.write()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

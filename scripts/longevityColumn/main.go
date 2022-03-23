// Appends longevity column to input file

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	column = kingpin.Flag("column", "Name of column to append.").Short('c').Default("max_longevity").String()
	infile = kingpin.Flag("infile", "Name of input file (writes output to same file).").Short('i').Required().String()
	user   = kingpin.Flag("user", "MySQL username.").Short('u').Default("").String()
)

type longevity struct {
	column    string
	db        *dbIO.DBIO
	id        string
	logger    *log.Logger
	longevity map[string]string
	records   *dataframe.Dataframe
}

func newLongevity() *longevity {
	// Returns initialized struct
	l := new(longevity)
	l.column = *column
	l.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	l.id = "ID"
	l.logger = codbutils.GetLogger()
	l.logger.Println("Initializing struct...")
	l.longevity = make(map[string]string)
	l.records, _ = dataframe.FromFile(*infile, 0)
	l.records.AddColumn(l.column, "-1")
	return l
}

func (l *longevity) setLongevity() {
	// Formats longevity table
	l.logger.Println("Appending max longevity data...")
	var b strings.Builder
	first := true
	for k := range l.records.Index {
		if !first {
			b.WriteByte(',')
		}
		b.WriteString(k)
		first = false
	}
	for _, i := range l.db.GetRows("Records", l.id, b.String(), fmt.Sprintf("%s,%s", l.id, l.column)) {
		val := i[1]
		if _, err := strconv.Atoi(val); err != nil {
			val = "-1"
		}
		l.records.UpdateCell(i[0], l.column, val)
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	l := newLongevity()
	l.setLongevity()
	l.records.ToCSV(*infile)
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

// Adds values to Wild column in patient table

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"time"
)

var user = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()

type wildColumn struct {
	db      *dbIO.DBIO
	records map[string]string
	patient map[string][]string
	terms   []string
}

func newWildColumn() *wildColumn {
	// Return new struct
	w := new(wildColumn)
	w.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	fmt.Println("\n\tInitializing struct...")
	w.records = make(map[string]string)
	w.patient = codbutils.ToMap(w.db.GetColumns("Patient", []string{"ID", "Comments", "Wild"}))
	w.terms = []string{"wild caught", "free rang"}
	return w
}

func (w *wildColumn) setWildColumn() {
	// Determines if records are infant records
	var count int
	fmt.Println("\tIdentifying wild caught records...")
	for k, v := range w.patient {
		val := "0"
		comment := strings.ToLower(v[0])
		for _, i := range w.terms {
			if strings.Contains(comment, i) {
				val = "1"
				break
			}
		}
		if val != v[1] {
			count++
			w.records[k] = val
		}
	}
	fmt.Printf("\tFound %d wild records.\n", count)
}

func (w *wildColumn) upload() {
	// Updates life history table with converted values
	var count int
	fmt.Println("\tUpdating Patient table...")
	for k, v := range w.records {
		count++
		w.db.UpdateRow("Patient", "Wild", v, "ID", "=", k)
		fmt.Printf("\tUpdated %d of %d records.\r", count, len(w.records))
	}
	fmt.Println()
}

func main() {
	start := time.Now()
	kingpin.Parse()
	w := newWildColumn()
	w.setWildColumn()
	w.upload()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

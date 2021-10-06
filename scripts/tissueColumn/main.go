// Adds values to tissue column in tumor table

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/diagnoses"
	"github.com/icwells/dbIO"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"time"
)

var user = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()

type tissueColumn struct {
	db      *dbIO.DBIO
	logger  *log.Logger
	tissues map[string]string
}

func newTissueColumn() *tissueColumn {
	// Return new struct
	t := new(tissueColumn)
	t.logger = codbutils.GetLogger()
	t.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	t.logger.Println("\n\tInitializing struct...")
	m := diagnoses.NewMatcher(t.logger)
	t.tissues = m.GetTissues()
	return t
}

func (t *tissueColumn) update() {
	// Updates life history table with converted values
	var count int
	fmt.Println("\tUpdating Tumor table...")
	for k, v := range t.tissues {
		count++
		t.db.UpdateRow("Tumor", "Tissue", v, "Location", "=", k)
		t.logger.Printf("\tUpdated %d of %d terms.\r", count, len(t.tissues))
	}
	fmt.Println()
}

func main() {
	start := time.Now()
	kingpin.Parse()
	t := newTissueColumn()
	t.update()
	t.logger.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

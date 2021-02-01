// Adds values to infant column in patient table

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
)

var user = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()

type infantColumn struct {
	ages    *dbupload.Infancy
	db      *dbIO.DBIO
	infant  [][]string
	patient map[string][]string
}

func newinfantColumn() *infantColumn {
	// Return new struct
	i := new(infantColumn)
	i.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	fmt.Println("\n\tInitializing struct...")
	i.ages = dbupload.NewInfancy(i.db)
	i.patient = codbutils.ToMap(i.db.GetColumns("Patient", []string{"ID", "taxa_id", "Age", "Infant", "Comments"}))
	return i
}

func (i *infantColumn) setInfancy() {
	// Determines if records are infant records
	for k, v := range i.patient {
		val := i.ages.SetInfant(v[0], v[1], v[3])
		if val != v[2] {
			i.infant = append(i.infant, []string{k, val})
		}
	}
}

func (i *infantColumn) upload() {
	// Updates life history table with converted values
	fmt.Println("\tUpdating Patient table...")
	for idx, j := range i.infant {
		i.db.UpdateRow("Patient", "Infant", j[1], "ID", "=", j[0])
		fmt.Printf("\tUpdated %d of %d records.\r", idx+1, len(i.infant))
	}
	fmt.Println()
}

func main() {
	start := time.Now()
	kingpin.Parse()
	i := newinfantColumn()
	i.setInfancy()
	i.upload()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

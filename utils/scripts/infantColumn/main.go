// Adds values to infant column in patient table

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"gopkg.in/alecthomas/kingpin.v2"
	"strconv"
	"strings"
	"time"
)

var user = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()

type infantColumn struct {
	ages    map[string]float64
	db      *dbIO.DBIO
	infant  [][]string
	patient map[string][]string
}

func newinfantColumn() *infantColumn {
	// Return new struct
	var l []string
	i := new(infantColumn)
	i.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false))
	fmt.Println("\n\tInitializing struct...")
	i.ages = codbutils.GetMinAges(i.db, l)
	i.patient = codbutils.ToMap(i.db.GetColumns("Patient", []string{"ID", "taxa_id", "Age", "Comments"}))
	return i
}

func (i *infantColumn) setInfancy() {
	// Determines if records are infant records
	for k, v := range i.patient {
		val := "-1"
		if min, ex := i.ages[v[0]]; ex {
			if age, err := strconv.ParseFloat(v[1], 64); err == nil {
				if age >= 0 {
					if age <= min {
						val = "1"
					} else {
						val = "0"
					}
				}
			}
		}
		if val == "-1" {
			// Check comments for keywords
			comments := strings.ToLower(v[2])
			for idx, j := range []string{"infant", "fetus", "juvenile", "immature", "adult", "mature"} {
				if strings.Contains(comments, j) {
					if idx <= 3 {
						val = "1"
					} else {
						val = "0"
					}
					break
				}
			}
		}
		i.infant = append(i.infant, []string{k, val})
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

// Uploads verified common names for taxonomy records

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
)

var (
	infile = kingpin.Arg("infile", "Path to input file.").Required().String()
	user   = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type commonNames struct {
	common   map[string]string
	db       *dbIO.DBIO
	infile   string
	names    [][]string
	taxa     map[string]string
	taxonomy map[string]string
}

func newCommonNames() *commonNames {
	// Returns struct
	c := new(commonNames)
	c.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	c.common = codbutils.EntryMap(c.db.GetColumns("Common", []string{"taxa_id", "Name"}))
	c.infile = *infile
	c.taxa = make(map[string]string)
	c.taxonomy = codbutils.EntryMap(c.db.GetColumns("Taxonomy", []string{"taxa_id", "Species"}))
	return c
}

func (c *commonNames) update() {
	// Updates database with new common names
	fmt.Println("\tUpdating taxonomy table...")
	for k, v := range c.taxa {
		c.db.UpdateRow("Taxonomy", "common_name", v, "taxa_id", "=", k)
	}
	fmt.Println("\tUpdating common names table...")
	for _, i := range c.names {
		c.db.UpdateRow("Taxonomy", "common_name", i[1], "taxa_id", "=", i[0])
	}
}

func (c *commonNames) readInfile() {
	// Reads common names from input file
	fmt.Println("\n\tReading common names from input file...")
	reader, header := iotools.YieldFile(c.infile, true)
	for i := range reader {
		species := i[header["Species"]]
		common := i[header["Common"]]
		if species != "" && species != "NA" && common != "" && common != "NA" {
			if id, ex := c.taxonomy[species]; ex {
				if _, exists := c.taxa[id]; !exists {
					c.taxa[id] = common
					if _, e := c.common[common]; !e {
						// Add missing common names
						c.names = append(c.names, []string{id, common, "NA"})
					}
				}
			}
		}
	}
	fmt.Printf("\tFound %d verified common names.\n", len(c.taxa))
	fmt.Printf("\tFound %d novel common names.\n", len(c.names))
}

func main() {
	kingpin.Parse()
	c := newCommonNames()
	c.readInfile()
	c.update()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(c.db.Starttime))
}

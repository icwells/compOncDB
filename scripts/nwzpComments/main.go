// Appends nwzp diagnosis comments to existing comment column

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"time"
)

var (
	infile = kingpin.Flag("infile", "Path to input file.").Short('i').Required().String()
	user   = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type comments struct {
	db    *dbIO.DBIO
	ids   map[string][]string
	table [][]string
}

func newComments() *comments {
	// Returns initialized converter struct
	c := new(comments)
	c.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false))
	c.setIDs()
	return c
}

func (c *comments) setIDs() {
	// Stores source_id to ID map
	var ids []string
	for _, i := range c.db.GetRows("Source", "service_name", "NWZP", "ID") {
		ids = append(ids, i[0])
	}
	c.ids = codbutils.ToMap(c.db.GetRows("Patient", "ID", strings.Join(ids, ","), "source_id,source_name,ID"))
	fmt.Printf("\tFound %d NWZP IDs.", len(c.ids))
}

func (c *comments) upload() {
	// Updates life history table with converted values
	l := 5000
	var start, end int
	fmt.Println("\tUpdating Life History table...")
	u := dbextract.NewUpdater(c.db)
	u.SetHeader("ID,Comments")
	for end < len(c.table) {
		end += l
		if end > len(c.table) {
			end = len(c.table)
		}
		for _, i := range c.table[start:end] {
			u.EvaluateRow(i)
		}
		u.UpdateTables()
		u.ClearTables()
		fmt.Printf("\tUploaded %d of %d comments.\n", end, len(c.table))
		start = end
	}
}

func (c *comments) getComments() {
	// Stores comment by source_id
	first := true
	d := "\t"
	f := iotools.OpenFile(*infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		if !first {
			i := strings.Split(strings.TrimSpace(string(scanner.Text())), d)
			if len(i) >= 12 {
				uid := i[0]
				if v, ex := c.ids[uid]; ex && strings.TrimSpace(i[6]) == v[0] {
					diag := strings.TrimSpace(i[9])
					comment := strings.TrimSpace(i[len(i)-1])
					if diag != "" && comment != "" {
						if diag[len(diag)-1] != '.' {
							diag += "."
						}
						diag += " " + comment
						diag = strings.Replace(diag, ",", "", -1)
						diag = strings.Replace(diag, "\"", "", -1)
						diag = strings.Replace(diag, "'", "", -1)
						diag = strings.Replace(diag, "  ", " ", -1)
						c.table = append(c.table, []string{v[1], diag})
					}
				}
			}
		} else {
			first = false
		}
	}
	fmt.Printf("\tFound %d NWZP comments.\n", len(c.table))
}

func main() {
	kingpin.Parse()
	c := newComments()
	fmt.Println("\n\tReading NWZP comments...")
	c.getComments()
	c.upload()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(c.db.Starttime))
}

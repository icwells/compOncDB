// Assigns years to nwzp records in the database

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/simpleset"
	"gopkg.in/alecthomas/kingpin.v2"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	infile  = kingpin.Flag("infile", "Path to NWZP input file.").Short('i').Required().String()
	outfile = kingpin.Flag("outfile", "Name of output file (writes to stdout if not given).").Short('o').Default("nil").String()
	user    = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()
)

type nwzpyears struct {
	db      *dbIO.DBIO
	digit   *regexp.Regexp
	ids     map[string]string
	records [][]string
	years   map[string]string
}

func newNwzpyears(db *dbIO.DBIO) *nwzpyears {
	// Returns new date instance
	n := new(nwzpyears)
	n.db = db
	n.digit = regexp.MustCompile(`([0-9]*[.])?[0-9]+`)
	n.ids = make(map[string]string)
	n.years = make(map[string]string)
	return n
}

func (n *nwzpyears) setYears() {
	// Merges nwzp ids with years
	fmt.Println("\tMerging NWZP IDs with years...")
	for k, v := range n.years {
		if id, ex := n.ids[k]; ex {
			n.records = append(n.records, []string{id, v, k})
		}
	}
}

func (n *nwzpyears) formatYear(val string) string {
	//Stores year in 4 digit format
	var ret string
	if val != "" {
		year := n.digit.FindString(val)
		if len(year) == 2 {
			if y, _ := strconv.Atoi(year); y > 50 {
				ret = "19" + year
			} else {
				ret = "20" + year
			}
		} else if len(year) == 4 {
			ret = year
		}
	}
	return ret
}

func (n *nwzpyears) nwzpIDs() {
	// Stores ids of nwzp records
	fmt.Println("\tRetrieving NWZP IDs...")
	ids := simpleset.NewStringSet()
	for _, i := range n.db.GetRows("Source", "service_name", "NWZP", n.db.Columns["Source"]) {
		ids.Add(i[0])
	}
	for _, i := range n.db.GetRows("Patient", "ID", strings.Join(ids.ToStringSlice(), ","), "ID,source_id") {
		n.ids[i[1]] = i[0]
	}
}

func (n *nwzpyears) readInfile() {
	// Stores formatted years by uid
	fmt.Printf("\n\tReading %s...\n", iotools.GetFileName(*infile))
	s, h := iotools.ReadFile(*infile, true)
	for _, i := range s {
		uid := strings.TrimSpace(i[h["UID"]])
		year := n.formatYear(strings.TrimSpace(i[h["CASE"]]))
		if year != "" {
			n.years[uid] = year
		}
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	n := newNwzpyears(codbutils.ConnectToDatabase(codbutils.SetConfiguration("config.txt", *user, false)))
	n.readInfile()
	n.nwzpIDs()
	//n.setYears()
	//codbutils.WriteResults(*outfile, "ID,Year", n.records)
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

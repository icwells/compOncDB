// Assigns years to nwzp records in the database

package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
)

var (
	infile  = kingpin.Flag("infile", "Path to NWZP input file.").Short('i').Required().String()
	outfile = kingpin.Flag("outfile", "Name of output file (writes to stdout if not given).").Short('o').Default("nil").String()
	user    = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()
)

type rec struct {
	id		string
	date	string
	year	string
}

func newRec(id, date, year string) *rec {
	// Returns new date instance
	r := new(rec)
	r.id = id
	r.date = date
	r.year = year
	return r
}

func (d *rec) equals(v *rec) bool {
	// Returns true if record ids and dates are equal
	if d.id == v.id && d.date == r.date {
		return true
	} else {
		return false
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration("config.txt", *user, false))

	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

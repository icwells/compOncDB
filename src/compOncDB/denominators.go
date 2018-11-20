// This script contains functions for updating/deleting values from the database

package main

import (
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"os"
	"strconv"
	"strings"
)

type denominators struct {
	db		*dbIO.DBIO
	infile	string
	delim	string
	species int
	cancer	int
	total	int
	rec		map[string]int
}

func newDenominators(db *dbIO.DBIO, infile string) denominators {
	// Returns initialized struct with existing data from table
	var d denominators
	d.db = db
	d.infile = infile
	d.species = -1
	d.cancer = -1
	d.total = -1
	d.rec = make(map[string]int)
	table := db.GetTable("Denominators")
	for _, i := range table {
		n, err := strconv.Atoi(i[1])
		if err == nil {
			d.rec[i[0]] = n
		}
	}
	return d
}

func (d *denominators) parseHeader(line string) {
	// Gets delimiter and target column numbers
	d.delim = iotools.GetDelim(line)
	s := strings.Split(line, d.delim)
	for idx, i := range s {
		i = strings.TrimSpace(i)
		if i == "Species" || i == "CommonNames" {
			d.species = idx
		} else if i == "Cancer" || i == "Tumor count" {
			d.cancer = idx
		} else if i == "Total" || i == "Total accessions" {
			d.total = idx
		}
	}
	if d.species < 0 || d.cancer < 0 || d.total < 0 {
		fmt.Print("\n\t[Error] Cannot determine column numbers. Exiting.\n\n")
		os.Exit(110)
	}
}

func (d *denominators) getNonCancer(s []string) int {
	// Returns number of non-cancer occurances
	ret := 0
	t, err := strconv.Atoi(s[d.total])
	c, er := strconv.Atoi(s[d.cancer])
	if err == nil && er == nil {
		ret = t - c
	}
	return ret
}

func (d *denominators) upload() {
	// Converts map to slice and uploads to table
	var den [][]string
	for k, v := range d.rec {
		// Taxa id, total
		r := []string{k, strconv.Itoa(v)}
		den = append(den, r)
	}
	vals, l := dbIO.FormatSlice(den)
	if len(vals) >= 1 {
		d.db.UpdateDB("Denominators", vals, l)
	}
}

func (d *denominators) readDenominators() {
	// Reads data from input file
	first := true
	com := entryMap(d.db.GetTable("Common"))
	f := iotools.OpenFile(d.infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		if first == false {
			s := strings.Split(line, d.delim)
			// Get taxa id
			c := strings.Title(strings.TrimSpace(s[d.species]))
			id, ex := com[c]
			if ex == true {
				x := d.getNonCancer(s)
				_, e := d.rec[id]
				if e == true {
					// Update record
					d.rec[id] += x
				} else {
					// Create new record
					d.rec[id] = x
				}
			}
		} else {
			d.parseHeader(line)
			first = false
		}
	}
}

func loadNonCancerTotals(db *dbIO.DBIO, infile string) {
	// Loads denominator 
	fmt.Println("\n\tUploading to denominators table...")
	d := newDenominators(db, infile)
	db.TruncateTable("Denominators")
	d.readDenominators()
	d.upload()
}

// Updates malignant code for NWZP records in upload file

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"time"
)

var (
	app     = kingpin.New("malignancyUpdater", "Updates malignant code for NWZP records in upload file.")
	infile  = kingpin.Flag("infile", "Path to input upload file.").Required().Short('i').String()
	update  = kingpin.Flag("update", "Path to file with updated malignant codes.").Required().Short('u').String()
	outfile = kingpin.Flag("outfile", "Path to outpur upload file.").Required().Short('o').String()
)

type updater struct {
	infile  string
	update  string
	outfile string
	codes   map[string]string
	head    map[string]int
	rows    [][]string
}

func newUpdater() *updater {
	// Initializes struct
	var u updater
	u.infile = *infile
	u.update = *update
	u.outfile = *outfile
	u.codes = make(map[string]string)
	u.rows, u.head = iotools.ReadFile(u.infile, true)
	return &u
}

func (u *updater) headerString() string {
	// Converts header to string
	ret := make([]string, len(u.head))
	for k, v := range u.head {
		ret[v] = k
	}
	return strings.Join(ret, ",")
}

func (u *updater) updateMalignant() {
	// Updates malignant codes in upload file
	fmt.Println("\tUpdating malignancy codes...")
	count := 0
	for _, i := range u.rows {
		if len(i) >= u.head["Service"] && i[u.head["Service"]] == "NWZP" {
			id := strings.TrimSpace(i[u.head["ID"]])
			code, ex := u.codes[id]
			if ex {
				i[u.head["Malignant"]] = code
				count++
			}
		}
	}
	fmt.Printf("\tUpdated %d records.\n", count)
	iotools.WriteToCSV(u.outfile, u.headerString(), u.rows)
}

func (u *updater) setCodes() {
	// Reads malignant codes into map
	fmt.Println("\tReading updated malignancy codes...")
	rows, h := iotools.ReadFile(u.update, true)
	for _, i := range rows {
		if len(i) >= h["Malignant"] {
			id := strings.TrimSpace(i[h["ID"]])
			code := strings.TrimSpace(i[h["Malignant"]])
			switch code {
			case "N":
				u.codes[id] = "0"
			case "Y":
				u.codes[id] = "1"
			case "NA":
				u.codes[id] = "-1"
			}
		}
	}
	fmt.Printf("\tFound %d updated codes.\n", len(u.codes))
}

func main() {
	start := time.Now()
	kingpin.Parse()
	fmt.Println("\n\tUpdating malignant codes in updload file...")
	u := newUpdater()
	u.setCodes()
	u.updateMalignant()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

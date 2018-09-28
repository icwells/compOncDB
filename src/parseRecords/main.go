// This script will parse and organize records for upload to the comparative oncology database

package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
	"time"
)

var (
	// Kingpin arguments
	app     = kingpin.New("parseRecords", "This script will parse and organize records for upload to the comparative oncology database.")
	infile  = kingpin.Flag("infile", "Path to input file.").Short('i').Required().String()
	outfile = kingpin.Flag("outfile", "Path to output file.").Short('o').Required().String()
	service = kingpin.Flag("service", "Service database name.").Short('s').Required().String()

	extract = kingpin.Command("extract", "Extract diagnosis data from infile.")
	dict    = extract.Flag("dict", "Path to dictionary of cancer terms.").Short('d').Default("cancerdict.tsv").String()

	merge = kingpin.Command("merge", "Merges taxonomy and diagnosis info with infile.")
	taxa  = merge.Flag("taxa", "Path to kestrel output.").Short('t').Default("nil").String()
	diag  = merge.Flag("diagnoses", "Path to diagnosis data.").Short('d').Default("nil").String()
)

func printFatal(msg string, code int) {
	// Prints error and exits
	fmt.Printf("\n\t[Error] %s. Exiting. \n\n", msg)
	os.Exit(code)
}

func getDelim(header string) string {
	// Returns delimiter
	var d string
	found := false
	for _, i := range []string{"\t", ",", " "} {
		if strings.Contains(header, i) == true {
			d = i
			found = true
		}
	}
	if found == false {
		printFatal("Cannot determine delimeter", 10)
	}
	return d
}

func mergeRecords(ent entries) {
	// Merges data into upload file
	fmt.Println("\n\tMerging records...")
	if *taxa != "nil" {
		ent.getTaxonomy(*taxa)
	}
	if *diag != "nil" {
		ent.getDiagnosis(*diag)
	}
	ent.sortRecords(*infile, *outfile)
}

func main() {
	start := time.Now()
	switch kingpin.Parse() {
	case extract.FullCommand():
		ent := newEntries(*service)
		ent.extractDiagnosis(*dict, *infile, *outfile)
	case merge.FullCommand():
		ent := newEntries(*service)
		mergeRecords(ent)
	}
	fmt.Printf("\tFinished. Run time: %s\n\n", time.Since(start))
}
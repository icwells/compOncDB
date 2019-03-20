// This script will parse and organize records for upload to the comparative oncology database

package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	// Kingpin arguments
	app     = kingpin.New("parseRecords", "This script will parse and organize records for upload to the comparative oncology database.")
	service = kingpin.Flag("service", "Service database name.").Short('s').Required().String()
	infile  = kingpin.Flag("infile", "Path to input file.").Short('i').Required().String()
	outfile = kingpin.Flag("outfile", "Path to output file.").Short('o').Required().String()
	taxa    = kingpin.Flag("taxa", "Path to kestrel output.").Short('t').Required().String()
)

func printFatal(msg string, code int) {
	// Prints error and exits
	fmt.Printf("\n\t[Error] %s. Exiting. \n\n", msg)
	os.Exit(code)
}

func main() {
	start := time.Now()
	kingpin.Parse()
	fmt.Print("\n\tProcessing input records...\n")
	ent := newEntries(*service)
	ent.getTaxonomy(*taxa)
	ent.sortRecords(*infile, *outfile)
	fmt.Printf("\tFinished. Run time: %s\n\n", time.Since(start))
}

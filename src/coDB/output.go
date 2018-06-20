// This script contains functions for writing database info to file

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"os"
	"strings"
)

func printCSV(outfile, header []string, table [][]string) {
	// Writes list of lists to csv
	out := iotools.CreateFile(outfile)
	defer out.Close()
	_, err := out.WriteString(strings.Join(header, ",") + "\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] Writing header line: %v\n", err)
	}
	for _, i := range table {
		_, err = out.WriteString(strings.Join(i, ",") + "\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "[Error] Writing table data: %v\n", err)
		}
	}
}

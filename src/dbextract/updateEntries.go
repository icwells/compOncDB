// Updates taxonmy entries in place

package dbextract

import (
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"strings"
)

func getColumns(row []string) {

}

func getUpdateFile(infile string) map[string][]string {
	// Returns map of data to be updated
	var d string
	ret := make(map[string][]string)
	first := true
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(string(input.Text()))
		if first == false {
			s := strings.Split(line, d)

		} else {
			d = iotools.GetDelim(line)
			
			first = false
		}
	}
}

func UpdateEntries(db *dbIO.DBIO, infile string) {
	// Updates taxonomy entries
	fields := getUpdateFile(infile)
}

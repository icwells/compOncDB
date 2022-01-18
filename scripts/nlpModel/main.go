// Extracts and preformats data for use as nlpModel input

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
	outfile = kingpin.Flag("outfile", "Name of output file.").Short('o').Required().String()
	user    = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type formatter struct {
	db       *dbIO.DBIO
}

func newFormatter() {
	// Connects to db and returns initialized struct
	f := new(formatter)
	f.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)

	return f
}

func main() {
	start := time.Now()

	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

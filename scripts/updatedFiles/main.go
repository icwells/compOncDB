// Produces updated files using commands in config file

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/Songmu/prompter"
	"gopkg.in/alecthomas/kingpin.v2"
	"os/exec"
	"path"
	"sync"
	"time"
)

var (
	app    = kingpin.New("updatedFiles", "Produces updated files using commands in config file.")
	outdir = kingpin.Flag("outdir", "Path to output directory.").Short('o').Required().String()
	user   = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type command struct {
	command		string
	directory	string
}

func newCommand(com, dir string) command {
	// Returns new command
	var c command
	c.command = com
	c.directory = dir
	return c
}

func (c *command) setOutfile(s string) string {
	// Returns formatted output file name
	stamp := codbutils.GetTimeStamp()

}

func (c *command) runCommand(wg *sync.WaitGroup) {
	// Runs given command
	defer wg.Done()
	cmd := 


}

func ping() string {
	// Returns connects to database and returns password
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	return db.Password
}

func setConfig() []command {
	// Returns input file values
	var ret []command
	first := true
	infile := path.Join(iotools.GetGOPATH(), "src/github.com/icwells/compOncDB/scripts/updatedFiles/config.csv")
	reader, header := iotools.YieldFile(infile, true)
	for i := range reader {
		if !first {
			ret = append(ret, newCommand(i[header["Command"]], i[header["Directory"]]))
		} else {
			first = false
		}
	}
	return ret
}

func main() {
	start := time.Now()
	kingpin.Parse()
	password := ping()
	fmt.Println("\n\tIssuing commands...")
	var wg sync.WaitGroup
	for _, i := range setConfig() {
		wg.Add(1)
		go i.runCommand(&wg)
	}
	fmt.Println("\tWaiting for results...")
	wg.Wait()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

// Produces updated files using commands in config file

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	app    = kingpin.New("updatedFiles", "Produces updated files using commands in config file.")
	outdir = kingpin.Flag("outdir", "Path to output directory.").Short('o').Required().String()
	user   = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type command struct {
	command   string
	directory string
	options   string
}

func newCommand(com, dir string) command {
	// Returns new command
	var c command
	s := strings.Split(com, " ")
	c.command = s[0]
	c.directory = dir
	c.options = strings.Join(s[1:], " ")
	return c
}

func (c *command) formatOptions(pw string) {
	// Returns formatted output file name
	stamp := codbutils.GetTimeStamp()
	// Add outdir and time stamp to outfile
	cmd := strings.Split(c.options, "-o ")
	tail := strings.Split(cmd[1], " ")
	tail[0] = "-o " + path.Join(*outdir, strings.Replace(tail[0], ".csv", stamp+".csv", 1))
	// Add username and password
	cmd = append([]string{cmd[0]}, "-u "+*user)
	cmd = append([]string{cmd[0]}, "--password "+pw)
	cmd = append(cmd, tail...)
	c.options = strings.Join(cmd, " ")
}

func (c *command) runCommand(wg *sync.WaitGroup, pw string) {
	// Runs given command
	defer wg.Done()
	c.formatOptions(pw)
	cmd := exec.Command(c.command, c.options)
	cmd.Dir = c.directory
	if err := cmd.Run(); err != nil {
		fmt.Printf("\tCommand failed. %v\n", err)
	}
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
		go i.runCommand(&wg, password)
	}
	fmt.Println("\tWaiting for results...")
	wg.Wait()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

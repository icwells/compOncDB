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
	options   []string
	outfile   string
	password  string
	program   string
	user      string
}

func newCommand(com, dir, pw string) *command {
	// Returns new command
	idx := 2
	c := new(command)
	s := strings.Split(com, " ")
	if s[1] == "run" {
		c.command = c.setOption(s[1], s[2])
	} else {
		idx++
		c.command = s[1]
	}
	if dir != "" {
		c.directory = dir
	}
	c.password = "--password " + pw
	c.program = s[0]
	c.user =  "-u " + *user
	c.setOptions(s[idx:])
	return c
}

func (c *command) setOption(x, y string) string {
	// Formats individual command
	return fmt.Sprintf("%s %s", x, y)
}

func (c *command) setOutfile(v string) {
	// Formats output file name with outdir and time stamp
	stamp := codbutils.GetTimeStamp()
	c.outfile = "-o " + path.Join(*outdir, strings.Replace(v, "csv", stamp+".csv", 1))
}

func (c *command) setOptions(s []string) {
	// Stores options in struct
	for idx, i := range s {
		if i[0] == '-' && idx < len(s)-1 {
			switch i[1] {
			case 'o':
				c.setOutfile(s[idx+1])
			case '-':
				c.options = append(c.options, i)
			default:
				c.options = append(c.options, c.setOption(i, s[idx+1]))
			}
			idx++
		}
	}
}

func (c *command) formatCommand() *exec.Cmd {
	// Formats command with variable number of options
	var ret *exec.Cmd
	switch len(c.options) {
	case 1:
		ret = exec.Command(c.program, c.command, c.user, c.password, c.options[0], c.outfile)
	case 2:
		ret = exec.Command(c.program, c.command, c.user, c.password, c.options[0], c.options[1], c.outfile)
	case 3:
		ret = exec.Command(c.program, c.command, c.user, c.password, c.options[0], c.options[1], c.options[2], c.outfile)
	case 4:
		ret = exec.Command(c.program, c.command, c.user, c.password, c.options[0], c.options[1], c.options[2], c.options[3], c.outfile)
	//case 5:
	//	ret = exec.Command(c.program, c.command, c.user, c.password, c.options[0], c.options[1], c.options[2], c.options[3], c.options[4], c.outfile)
	}
	return ret
}

func (c *command) runCommand(wg *sync.WaitGroup) {
	// Runs given command
	defer wg.Done()
	cmd := c.formatCommand()
	if c.directory != "" {
		cmd.Dir = c.directory
	}
	fmt.Println(cmd.String())
	if err := cmd.Run(); err != nil {
		fmt.Printf("\tCommand failed. %v\n", err)
	}
}

//----------------------------------------------------------------------------

func ping() string {
	// Returns connects to database and returns password
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	return db.Password
}

func setConfig(pw string) []*command {
	// Returns input file values
	var ret []*command
	infile := path.Join(iotools.GetGOPATH(), "src/github.com/icwells/compOncDB/scripts/updatedFiles/config.csv")
	reader, header := iotools.YieldFile(infile, true)
	for i := range reader {
		ret = append(ret, newCommand(i[header["Command"]], i[header["Directory"]], pw))
	}
	return ret
}

func main() {
	start := time.Now()
	kingpin.Parse()
	password := ping()
	fmt.Println("\n\tIssuing commands...")
	var wg sync.WaitGroup
	for _, i := range setConfig(password) {
		wg.Add(1)
		go i.runCommand(&wg)
	}
	fmt.Println("\tWaiting for results...")
	wg.Wait()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

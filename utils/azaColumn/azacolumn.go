// Inetifies AZA status for existing accounts

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/clusteraccounts"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/simpleset"
	"gopkg.in/alecthomas/kingpin.v2"
	"path"
	"path/filepath"
	"strings"
)

var infile = kingpin.Arg("infile", "Path to input file.").Required().String()

type azaidentifier struct {
	infile   string
	outfile  string
	header   string
	queries  *simpleset.Set
	zoos     map[string][]string
	accounts map[string]string
}

func newAZAidentifier() *azaidentifier {
	// Returns struct
	var a azaidentifier
	a.infile = *infile
	a.outfile = path.Join(filepath.Dir(a.infile), "azaAccounts.csv")
	a.header = "acount_id,AZAzoo,Name,Zoo,AZA,Inst\n"
	a.queries = simpleset.NewStringSet()
	a.accounts = make(map[string]string)
	return &a
}

func (a *azaidentifier) mergeAccounts() {
	// Writes merged data to file
	fmt.Println("\tWriting accounts to file...")
	out := iotools.CreateFile(a.outfile)
	defer out.Close()
	out.WriteString(a.header)
	for k, v := range a.zoos {
		if id, ex := a.accounts[k]; ex {
			out.WriteString(id + "," + strings.Join(v, ",") + "\n")
		}
	}
}

func (a *azaidentifier) readAccounts() {
	// Reads account data from input file
	var delim string
	first := true
	fmt.Println("\n\tReading accounts from input file...")
	f := iotools.OpenFile(a.infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Text()))
		if first == false {
			s := strings.Split(line, delim)
			a.queries.Add(s[2])
			a.accounts[s[2]] = s[0]
		} else {
			delim = iotools.GetDelim(line)
			first = false
		}
	}
}

func main() {
	kingpin.Parse()
	id := newAZAidentifier()
	id.readAccounts()
	a := clusteraccounts.NewAccounts("")
	a.Queries = id.queries
	id.zoos = a.IdentifyAZA()
	id.mergeAccounts()
}

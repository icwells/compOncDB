// Uses spell checking and fuzzy matching to condense submitter names

package clusteraccounts

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"github.com/trustmaster/go-aspell"
	"os"
	"strings"
)

type Accounts struct {
	speller         aspell.Speller
	Queries, corpus strarray.Set
	terms           []*term
}

func NewAccounts(infile string) *Accounts {
	// Returns pointer to initialized struct
	a := new(Accounts)
	var err error
	a.speller, err = aspell.NewSpeller(map[string]string{"lang": "en_US"})
	if err != nil {
		fmt.Printf("\n\t[Error] Cannot initialize speller. Exiting.\n%v", err)
		os.Exit(500)
	}
	a.Queries = strarray.NewSet()
	a.corpus = strarray.NewSet()
	if infile != "" {
		a.readAccounts(infile)
	}
	return a
}

func (a *Accounts) getAccounts() map[string][]string {
	// Returns map of original term: corrected term
	counter := strarray.NewSet()
	total := strarray.NewSet()
	ret := make(map[string][]string)
	for _, i := range a.terms {
		ret[i.query] = i.toSlice()
		counter.Add(i.name)
		total.Add(i.query)
	}
	fmt.Printf("\tFormatted %d names from a total of %d account entries.\n", counter.Length(), total.Length())
	return ret
}

func (a *Accounts) getIndeces(row []string) int {
	// Returns indeces for submitter column
	ret := -1
	for idx, i := range row {
		i = strings.TrimSpace(i)
		i = strings.Replace(i, " ", "", -1)
		if i == "Client" || i == "Owner" || i == "InstitutionID" {
			ret = idx
		}
	}
	return ret
}

func (a *Accounts) readAccounts(infile string) {
	// Reads account data from input file
	var delim string
	var sub int
	first := true
	fmt.Println("\tReading accounts from input file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		if first == false {
			s := strings.Split(line, delim)
			a.Queries.Add(s[sub])
		} else {
			delim = iotools.GetDelim(line)
			sub = a.getIndeces(strings.Split(line, delim))
			first = false
			if sub == -1 {
				// Skip if column is not present
				break
			}
		}
	}
}

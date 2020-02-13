// Uses spell checking and fuzzy matching to condense submitter names

package clusteraccounts

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/simpleset"
	"github.com/trustmaster/go-aspell"
	"os"
	"path"
	"strings"
)

type Accounts struct {
	speller         aspell.Speller
	Queries, corpus *simpleset.Set
	terms           []*term
	zoos            []string
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
	a.Queries = simpleset.NewStringSet()
	a.corpus = simpleset.NewStringSet()
	a.zoos = codbutils.ReadList(path.Join(codbutils.Getutils(), "AZA_Zoos.csv"), 0)
	for idx, i := range a.zoos {
		a.zoos[idx] = strings.ToLower(i)
	}
	if infile != "" {
		a.readAccounts(infile)
	}
	return a
}

func (a *Accounts) getAccounts() map[string][]string {
	// Returns map of original term: corrected term
	counter := simpleset.NewStringSet()
	total := simpleset.NewStringSet()
	ret := make(map[string][]string)
	for _, i := range a.terms {
		a.azaStatus(i)
		ret[i.query] = i.toSlice()
		counter.Add(i.name)
		total.Add(i.query)
	}
	fmt.Printf("\tFormatted %d names from a total of %d account entries.\n", counter.Length(), total.Length())
	return ret
}

func (a *Accounts) getIndeces(row []string) int {
	// Returns indeces for submitter column
	for idx, i := range row {
		i = strings.TrimSpace(i)
		i = strings.Replace(i, " ", "", -1)
		if i == "Client" || i == "Owner" || i == "InstitutionID" || i == "submitter_name" {
			return idx
		}
	}
	return -1
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

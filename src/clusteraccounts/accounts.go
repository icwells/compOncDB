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
	speller              aspell.Speller
	set, corpus          strarray.Set
	clusters, submitters map[string][]*term
	terms                []*term
	scores               map[string]int
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
	a.set = strarray.NewSet()
	a.corpus = strarray.NewSet()
	a.clusters = make(map[string][]*term)
	a.submitters = make(map[string][]*term)
	if infile != "" {
		a.readAccounts(infile)
	}
	return a
}

func (a *Accounts) getAccounts() map[string][]string {
	// Returns map of original term: corrected term
	var count int
	ret := make(map[string][]string)
	for _, i := range a.terms {
		count++
		ret[i.query] = i.toSlice()
	}
	fmt.Printf("\tFormatted %d account entries.\n", count)
	return ret
}

func (a *Accounts) getIndeces(row []string) (int, int) {
	// Returns indeces for account and submitter columns
	acc, sub := -1, -1
	for idx, i := range row {
		i = strings.TrimSpace(i)
		i = strings.Replace(i, " ", "", -1)
		if i == "Account" {
			acc = idx
		} else if i == "Client" || i == "Owner" || i == "InstitutionID" {
			sub = idx
		}
	}
	return acc, sub
}

func (a *Accounts) parseRow(acc, sub string) {
	// Stores formatted terms and adds correcly spelled words to corpus
	name := a.checkAbbreviations(sub)
	t := newTerm(sub, name, acc)
	a.terms = append(a.terms, t)
	for _, i := range strings.Split(name, " ") {
		if a.speller.Check(i) {
			// Add to corpus of correctly spelled words
			a.corpus.Add(i)
		}
	}
}

func (a *Accounts) readAccounts(infile string) {
	// Reads account data from input file
	var delim string
	var acc, sub int
	first := true
	fmt.Println("\tReading accounts from input file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		if first == false {
			var account string
			s := strings.Split(line, delim)
			if acc != -1 {
				account = s[acc]
			}
			a.parseRow(account, s[sub])
		} else {
			delim = iotools.GetDelim(line)
			acc, sub = a.getIndeces(strings.Split(line, delim))
			first = false
			if sub == -1 {
				// Skip if column is not present
				break
			}
		}
	}
}

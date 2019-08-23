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
	ratio                            float64
	speller                          aspell.Speller
	set                              strarray.Set
	corpus                           strarray.Set
	submitters, pool, queries, terms map[string][]string
	scores                           map[string]int
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
	a.submitters = make(map[string][]string)
	a.queries = make(map[string][]string)
	a.terms = make(map[string][]string)
	a.ratio = 0.05
	if infile != "" {
		a.readAccounts(infile)
	}
	return a
}

func (a *Accounts) setAccountType(term string) (string, string) {
	// Returns 1/0 for zoo/institute columns
	zoo := "0"
	inst := "0"
	term = strings.ToLower(term)
	if strings.Contains(term, "zoo") || strings.Contains(term, "aquarium") || strings.Contains(term, "museum") {
		zoo = "1"
	} else if strings.Contains(term, "center") || strings.Contains(term, "institute") || strings.Contains(term, "service") || strings.Contains(term, "research") {
		inst = "1"
	}
	return zoo, inst
}

func (a *Accounts) getAccounts() map[string][]string {
	// Returns map of original term: corrected term
	var count, total int
	ret := make(map[string][]string)
	for key, val := range a.terms {
		count++
		zoo, inst := a.setAccountType(key)
		for _, i := range val {
			for _, v := range a.queries[i] {
				if _, ex := ret[v]; ex == false {
					// Store original term and consensus term
					total++
					ret[v] = []string{key, zoo, inst}
				}
			}
		}
	}
	fmt.Printf("\tFormatted %d terms from %d total account entries.\n", count, total)
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
			s := strings.Split(line, delim)
			if acc != -1 && s[acc] != "NA" {
				// Store in map by account id
				a.submitters[s[acc]] = append(a.submitters[s[acc]], s[sub])
			} else {
				// Store submitter only
				a.set.Add(s[sub])
			}
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

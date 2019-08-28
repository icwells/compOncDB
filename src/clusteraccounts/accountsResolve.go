// Contains methods for accounts struct

package clusteraccounts

import (
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
	"strings"
)

func (a *Accounts) fuzzymatch(s string, t []string) string {
	// Returns target match/self
	ret := s
	// Maximum of one substitution per word
	max := strings.Count(s, " ") + 1
	matches := fuzzy.RankFind(s, t)
	if len(matches) >= 1 {
		sort.Sort(matches)
		if matches[0].Distance <= max {
			ret = matches[0].Target
		}
	}
	return ret
}

func (a *Accounts) correctSpellings() {
	// Compares spelling of each word to corpus
	found := make(map[string]string)
	corp := a.corpus.ToSlice()
	for _, t := range a.terms {
		var words []string
		for _, i := range strings.Split(t.name, " ") {
			if a.corpus.InSet(i) {
				// Skip words in corpus as they are correctly spelled
				words = append(words, i)
			} else {
				if _, ex := found[i]; ex == true {
					// Use previously identified match
					words = append(words, found[i])
				} else {
					match := a.fuzzymatch(i, corp)
					words = append(words, match)
					found[i] = match
				}
			}
		}
		// Update name and determine source type
		t.name = strings.Join(words, " ")
		t.setType()
	}
}

func (a *Accounts) ResolveAccounts() map[string][]string {
	// Resolves differences in account names
	fmt.Println("\tFormatting account names...")
	for _, i := range a.Queries.ToSlice() {
		name := a.checkAbbreviations(i)
		a.terms = append(a.terms, newTerm(i, name))
	}
	a.correctSpellings()
	return a.getAccounts()
}

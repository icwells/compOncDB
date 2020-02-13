// Contains methods for accounts struct

package clusteraccounts

import (
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
	"strings"
)

func (a *Accounts) fuzzyrank(s string, t []string) (fuzzy.Rank, bool) {
	// Returns highest match
	var ret fuzzy.Rank
	var ex bool
	matches := fuzzy.RankFindFold(s, t)
	if len(matches) >= 1 {
		sort.Sort(matches)
		ret = matches[0]
		ex = true
	}
	return ret, ex
}

func (a *Accounts) fuzzymatch(s string, t []string) string {
	// Returns target match/self
	// Maximum of one substitution per word
	max := strings.Count(s, " ") + 1
	match, ex := a.fuzzyrank(s, t)
	if ex && match.Distance <= max {
		return match.Target
	}
	return s
}

func (a *Accounts) azaStatus(t *term) *term {
	// Sets AZA member status for zoos/institutes
	max := strings.Count(t.name, " ") + 1
	name := strings.ToLower(t.name)
	if match, ex := a.fuzzyrank(name, a.zoos); ex {
		t.match = match.Target
		if match.Distance <= max {
			t.aza = 1
		} else if max > 1 {
			// Only apply if there are multiple words
			if strings.Contains(match.Target, name) || strings.Contains(name, match.Target) {
				t.aza = 1
			}
		}
	}
	return t
}

func (a *Accounts) IdentifyAZA() map[string][]string {
	// Returns slice for identifying aza status for existing records
	ret := make(map[string][]string)
	for _, i := range a.Queries.ToStringSlice() {
		t := newTerm(i, i)
		t = a.azaStatus(t)
		ret[t.query] = append([]string{t.match}, t.toSlice()...)
	}
	return ret
}

func (a *Accounts) correctSpellings() {
	// Compares spelling of each word to corpus
	found := make(map[string]string)
	corp := a.corpus.ToStringSlice()
	for _, t := range a.terms {
		var words []string
		for _, i := range strings.Split(t.name, " ") {
			if ex, _ := a.corpus.InSet(i); ex {
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
	ch := make(chan string)
	for _, i := range a.Queries.ToStringSlice() {
		go a.checkAbbreviations(ch, i)
		name := <-ch
		a.terms = append(a.terms, newTerm(i, name))
	}
	a.correctSpellings()
	return a.getAccounts()
}

// Contains methods for accounts struct

package clusteraccounts

import (
	"github.com/icwells/go-tools/strarray"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
	"strings"
)

func (a *Accounts) fuzzymatch(s string, t []string) string {
	// Returns target match/self
	target := s
	matches := fuzzy.RankFind(s, t)
	if len(matches) > 1 {
		sort.Sort(matches)
		if float64(matches[0].Distance)/float64(len(s)) < a.ratio {
			target = matches[0].Target
		}
	}
	return target
}

func (a *Accounts) clusterNames(key string) int {
	// Clusters names from clusters[key] to nearest match
	count := 0
	bins := make(map[int][]*term)
	// Seperate based on number of words
	for _, i := range a.clusters[key] {
		bins[i.length] = append(bins[i.length], i)
	}
	for _, v := range bins {
		// Score against terms of equal length
		found := make(map[string]string)
		s := strarray.NewSet()
		var names []string
		for _, i := range v {
			names = append(names, i.name)
		}
		for idx, i := range names {
			target := strarray.DeleteSliceIndex(names, idx)
			match := a.fuzzymatch(i, target)
			// Catch all matches
			s.Add(match)
		}
		names = s.ToSlice()
		for _, i := range v {
			match, ex := found[i.name]
			if ex == false {
				match = a.fuzzymatch(i.name, names)
				found[i.name] = match
			}
			if match != i.name {
				// Store corrected match to closest fuzzy match/self
				i.name = match
				count++
			}
		}
	}
	return count
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
	a.correctSpellings()
	for _, i := range a.terms {
		if i.id == "" {
			// Store terms with no ids by source type
			k := i.getType()
			a.clusters[k] = append(a.clusters[k], i)
		} else {
			a.clusters[i.id] = append(a.clusters[i.id], i)
		}
	}
	for k := range a.clusters {
		count := len(a.clusters[k])
		for count > 0 {
			// Cluster until no names change
			count = a.clusterNames(k)
		}
	}
	return a.getAccounts()
}

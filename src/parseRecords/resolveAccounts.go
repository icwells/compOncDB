// Contains methods for accounts struct

package main

import (
	"github.com/icwells/go-tools/strarray"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
	"strings"
)

func (a *accounts) clearPool() {
	// Clears pool map
	for k := range a.pool {
		delete(a.pool, k)
	}
}

func (a *accounts) fuzzymatch(s string, t []string) string {
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

func (a *accounts) checkSpelling(v string) int {
	// Returns number of correctly spelled words
	ret := 0
	for _, i := range strings.Split(v, " ") {
		if a.speller.Check(i) {
			ret++
		}
	}
	return ret
}

func (a *accounts) setScores(v []string) {
	// Scores keys in pool
	a.scores = make(map[string]int)
	for _, i := range v {
		// Greater number of properly spelled words = greater likelihood of being correct
		a.scores[i] = a.checkSpelling(i)
	}
	for _, i := range v {
		// Index closest match or self
		target := a.fuzzymatch(i, strarray.DeleteSliceValue(v, i))
		// Add total number of queries
		l := len(a.queries[target])
		a.scores[target] += l
	}
}

func (a *accounts) setTerms() {
	// Determines best candidate for map key
	var max int
	for key, val := range a.pool {
		a.setScores(val)
		for k, v := range a.scores {
			// Determine consensus key
			if v > max {
				key = k
			}
		}
		for _, i := range val {
			// Append to cluster by key
			a.terms[key] = append(a.terms[key], i)
		}
	}
}

func (a *accounts) clusterNames() {
	// Clusters set values into pool
	a.pool = make(map[string][]string)
	bins := make(map[int][]string)
	// Seperate based on number of words
	for k := range a.queries {
		l := strings.Count(k, " ") + 1
		bins[l] = append(bins[l], k)
	}
	for _, v := range bins {
		// Score against terms of equal length
		s := strarray.NewSet()
		for idx, i := range v {
			target := strarray.DeleteSliceIndex(v, idx)
			match := a.fuzzymatch(i, target)
			// Catch all matches
			s.Add(match)
		}
		terms := s.ToSlice()
		for _, i := range v {
			match := a.fuzzymatch(i, terms)
			// Store corrected match to closest fuzzy match/self
			a.pool[match] = append(a.pool[match], i)
		}
	}
}

func (a *accounts) setQueries(s []string, pool bool) {
	// Pools corrected terms from slice
	for _, i := range s {
		// Get unique corrected terms
		term := a.checkAbbreviations(i)
		a.queries[term] = append(a.queries[term], i)
		if pool == true {
			// Add to pool if skipping clustering
			a.pool[term] = append(a.pool[term], i)
		}
	}
}

func (a *accounts) resolveAccounts() map[string]string {
	// Resolves differneces in account names
	if a.set.Length() >= 1 {
		a.setQueries(a.set.ToSlice(), false)
		a.clusterNames()
		a.setTerms()
	}
	if len(a.submitters) >= 1 {
		for _, v := range a.submitters {
			a.pool = make(map[string][]string)
			// Get keys for each account ID
			a.setQueries(v, true)
			a.setTerms()
		}
	}
	return a.getAccounts()
}

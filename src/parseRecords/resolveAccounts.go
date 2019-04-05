// Contains methods for accounts struct

package 

import (
	"github.com/icwells/go-tools/strarray"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
	"strings"
)
func (a *accounts) clearMaps() {
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

func (a *accounts) setScores() {
	// Scores keys in pool
	keys := a.mapKeys(a.pool)
	for k := range a.pool {
		// Greater number of properly spelled words = greater likelihood of being correct
		a.scores[k] = a.checkSpelling(k)
	}
	for k := range a.pool {
		// Index closest match or self
		target := a.fuzzymatch(k, strarray.DeleteSliceValue(keys, k))
		// Add total number of queries
		l := len(a.queries[target])
		a.scores[target] += l
	}
}

func (a *accounts) setTerms() {
	// Determines best candidate for map key
	var key string
	var max int
	a.setScores()
	for k, v := range a.scores {
		// Determine consensus key
		if v > max {
			key = k
		}
	}
	for _, i := range a.pool {
		// Append to cluster by key
		a.terms[key] = append(a.terms, i)
	}
}

func (a *accounts) clusterNames() {
	// Clusters set values into pool
	bins := make(map[int][]string)
	// Seperate based on number of words
	for k := range a.queries {	
		l := strings.Count(k, " ") + 1
		bins[l] = append(bins[l], k)
	}
	for k, v := range bins {
		// Score against terms of equal length
		s := strarray.NewSet()
		for idx, i := range v {
			target := strarray.DeleteSliceIndex(v, idx)
			match := a.fuzzyMatch(i, target)
			// Catch all matches
			s.Add(match)
		}
		terms := s.ToSlice()
		for idx, i := range v {
			match := a.fuzzyMatch(i, terms)
			// Store corrected match to closest fuzzy match/self
			a.pool[match] = append(a.pool[match], i)
		}
	}
}

func (a *accounts) setQueries(s []string, pool) {
	// Pools corrected terms from slice
	for _, i := range s {
		// Get unique corrected terms
		term = a.checkAbbreviations(i)
		a.queries[term] = append(a.queries[term], i)
		if pool == true {
			// Add to pool if skipping clustering
			a.pool[k] = append(a.pool[k], i)
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
			// Get keys for each account ID
			a.setQueries(v, true)
			a.setTerms()
			a.clearPool()
		}
	}
	return a.getAccounts()
}

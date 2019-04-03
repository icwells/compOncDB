// Uses spell checking and fuzzy matching to condense submitter names

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/trustmaster/go-aspell"
	"os"
	"sort"
	"strings"
)

type accounts struct {
	ratio		float64
	speller		aspell.Speller
	set			strarray.Set
	submitters	map[string][]string
	pool		map[string][]string
	scores		map[string]int
	clusters	map[string][]string
	keys		[]string
}

func newAccounts(infile string) *accounts {
	// Returns pointer to initialized struct
	var a accounts
	a.speller, err = aspell.NewSpeller(map[string]string{"lang": "en_US",})
	if err != nil {
		fmt.Printf("\n\t[Error] Cannot initialize speller. Exiting.\n%v", err)
		os.Exit(500)
	}
	a.set = strarray.NewSet()
	a.submitters = make(map[string][]string)
	a.clusters = make(map[string][]string)
	a.ratio = 0.1
	return &a
}

func (a *accounts) getAccounts() map[string]string {
	// Returns map of original term: corrected term
	ret := make(map[string]string)
	for k := range a.clusters {
		for _, i := range a.clusters[k] {
			ret[i] = k
		}
	}
	return ret
}

func (a *accounts) checkAmpersand(val string) string {
	// Replaces ampersand with "and" and corrects spacing
	var rep string
	if strings.Contains(val, "&") == true {
		if strings.Contains(val, " & ") == true {
			rep = "And"
		} else if strings.Contains(val, " &") == true {
			rep = "And "
		} else if strings.Contains(val, "& ") == true {
			rep = " And"
		} else {
			rep = " And "
		}
		val = strings.Replace(val, "&", rep, -1)
	}
	return val
}

func (a *accounts) checkAbbreviations(val string) string {
	//Store submitter/NA
	terms := map[string]string{"Animal Clinic": "A. C.", "Animal Hospital": "A. H.", "Veterinary Clinic": "V. C.", "University": "Univ",
		"Veterinary Hospital": "V. H.", "Veterinary Services": "V. S.", "Pet Vet": "P. V.", "International": "Intl ", "Animal": "Anim "}
	// in records.go
	val = checkString(val)
	if val != "NA" {
		val = a.checkAmpersand(starray.TitleCase(val))
		// Resolve abbreviations
		for k, v := range terms {
			var alt string
			if strings.Contains(v, ".") == false {
				// Add trailing period
				alt = strings.Replace(v, " ", ".", 1)
			} else {
				// Remove space
				alt = strings.Replace(v, " ", "", 1)
			}
			if strings.Contains(val, v) == true {
				val = strings.Replace(val, v, k, 1)
				break
			} else if strings.Contains(val, alt) == true {
				val = strings.Replace(val, alt, k, 1)
				break
			}
		}
	}
	return val
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

func (a *accounts) resetKeys() {
	// Updates key slice with keys from pool
	var keys []string
	for k := range a.pool {
		keys = append(keys, k)
	}
	s.keys = keys
}

func (a *accounts) setPool(s []string) {
	// Pools corrected terms from slice
	pool := make(map[string][]string)
	for _, i := range s {
		// Get unique corrected terms
		term = a.checkAbbreviations(i)
		a.pool[term] = append(a.pool[term], i)
	}
	a.pool = pool
	a.resetKeys()
}

func (a *accounts) fuzzymatch(s string) (string, bool) {
	// Returns target match and whether a match was found
	ret := false
	target := s
	matches := fuzzy.RankFind(s, a.keys)
	if len(matches) > 1 {
		sort.Sort(matches)
		// Skip match to self
		if float64(matches[1].Distance)/float64(len(s)) < a.ratio {
			target = matches[1].Target
			ret = true
		}
	}
	return target, ret
}

func (a *accounts) setScores() {
	// Scores keys in pool
	scores := make(map[string]int)
	for k := range a.pool {
		// Greater number of properly spelled words = greater likelihood of being correct
		scores[k] = a.checkSpelling(k)
		// Index closest match or self
		target, _ := a.fuzzymatch(k)
		scores[target]++
	}
}

func (a *accounts) setClusters() {
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
		a.clusters[key] = append(a.clusters, i)
	}
}

func (a *accounts) clusterNames() {
	// Clusters set values into submitters map
	a.setPool(a.set.ToSlice())
	a.setScores()
	
}

func (a *accounts) resolveAccounts() map[string]string {
	// Resolves differneces in account names
	if a.set.Length() >= 1 {
		a.clusterNames()
	}
	if len(a.submitters) >= 1 {
		for _, v := range a.submitters {
			a.setPool(v)
			a.setClusters()
		}
	}
	return a.getAccounts()
}

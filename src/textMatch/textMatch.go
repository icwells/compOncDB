// This package will search a slice of strings for terms in a given slice and return matched pairs

import (
	"github.com/icwells/go-tools/strarray"
	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

func LevDist(s, t string) int {
	// Returns Levenshtein distance between q and t
	// adapted from https://github.com/jeffsmith82/gofuzzy/blob/master/fuzzy.go
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}
	// Make grid of indeces
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(s)][len(t)]
}

func searchSlice(ch chan []string, matcher *search.Matcher, query string, target [][]string, idx, ind int) {
	// Searches target slice for query
	var ret []string
	for _, i := range target {
		if matcher.EqualString(query, i[idx]) == true {
			ret = []string{query, i[ind]}
			break
		} else if LevDist(query, i) < len( {

		}
	}
	ch <- ret
}

func Search(query []string, target [][]string, idx, ind int) ([][]string, []string) {
	// Searches target for match in target[idx] and returns [query, target[ind]]
	var rows [][]string
	var misses []string
	ch := make(chan []string)
	matcher := search.New(language.English)
	for _, i := range query {
		go searchSlice(ch, mathcer, query[i], target, idx, ind)
		ret := <-ch
		if len(ret) > 1 {
			rows = append(rows, ret)
		} else {
			// Store unmathced queries
			misses = append(misses, i)
	}
	return rows, misses
}

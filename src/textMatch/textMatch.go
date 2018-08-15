// This package will search a slice of strings for terms in a given slice and return matched pairs

package textMatch

import (
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
	// Score all possible matches
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
	// Return minimum distance
	return d[len(s)][len(t)]
}

func ScoreMatch(query, target string) float64 {
	// Calculates percent score
	dist := float64(LevDist(query, target))
	l := float64((len(query) + len(target)) / 2)
	return (1.0 - (dist / l))
}

func matchSlice(ch chan []string, matcher *search.Matcher, query string, target [][]string, idx, ind int, score float64) {
	// Searches target slice for query
	var ret []string
	min := len(query)
	for _, i := range target {
		if matcher.EqualString(query, i[idx]) == true {
			ret = []string{query, i[ind]}
			break
		} else {
			dist := LevDist(query, i[idx])
			val := float64((1 - (dist / (len(query) + len(i[idx])) / 2)))
			if dist < min && val >= score {
				ret = []string{query, i[ind]}
				min = dist
			}
		}
	}
	ch <- ret
}

func SearchSlice(query []string, target [][]string, idx, ind int, min float64) ([][]string, []string) {
	// Searches target for match in target[idx] and returns [query, target[ind]]
	var rows [][]string
	var misses []string
	ch := make(chan []string)
	matcher := search.New(language.English)
	for _, i := range query {
		go matchSlice(ch, matcher, i, target, idx, ind, min)
		ret := <-ch
		if len(ret) > 1 {
			rows = append(rows, ret)
		} else {
			// Store unmathced queries
			misses = append(misses, i)
		}
	}
	return rows, misses
}

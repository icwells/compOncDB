// Helper fucntions for dbupload

package dbupload

import (
	"github.com/icwells/go-tools/strarray"
	"strings"
)

func ToMap(t [][]string) map[string][]string {
	// Converts slice of string slices to map with first element as key
	// If slice is two columns wide, it will append the second item to map entry
	m := make(map[string][]string)
	for _, i := range t {
		_, ex := m[i[0]]
		if ex == false {
			if len(i) == 2 {
				// Create new slice
				m[i[0]] = []string{i[1]}
			} else {
				m[i[0]] = i[1:]
			}
		} else if len(i) == 2 && strarray.InSliceStr(m[i[0]], i[1]) == false {
			// Append new stirng element
			m[i[0]] = append(m[i[0]], i[1])
		}
	}
	return m
}

func MapOfMaps(t [][]string) map[string]map[string]string {
	// Converts table to map of maps for easier searching
	ret := make(map[string]map[string]string)
	for _, row := range t {
		if m, ex := ret[row[1]]; ex == true {
			if _, e := m[row[2]]; e == false {
				// Add to existing map
				ret[row[1]][row[2]] = row[0]
			}
		} else {
			// Make new sub-map
			ret[row[1]] = make(map[string]string)
			ret[row[1]][row[2]] = row[0]
		}
	}
	return ret
}

func EntryMap(t [][]string) map[string]string {
	// Converts pair of columns to map for easier searching
	m := make(map[string]string)
	for _, i := range t {
		if _, ex := m[i[1]]; ex == false {
			m[i[1]] = i[0]
		}
	}
	return m
}

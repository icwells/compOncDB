// Contains functions for convertng splice of string slices to map

package main

import (
	"github.com/icwells/go-tools/strarray"
)

func mapOfMaps(t [][]string) map[string]map[string]string {
	// Converts table to map of maps for easier searching
	ret := make(map[string]map[string]string)
	for _, row := range t {
		if strarray.InMapMapStr(ret, row[1]) == true {
			if strarray.InMapStr(ret[row[1]], row[2]) == false {
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

func entryMap(t [][]string) map[string]string {
	// Converts table to map for easier searching
	m := make(map[string]string)
	for _, i := range t {
		if strarray.InMapStr(m, i[1]) == false {
			m[i[1]] = i[0]
		}
	}
	return m
}

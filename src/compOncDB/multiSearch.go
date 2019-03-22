// Applys addional search filters to search results

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"os"
	"strconv"
	"strings"
)

type evaluation struct {
	column   string
	operator string
	value    string
}

func (e *evaluation) getOperation(eval string) {
	// Splits eval into column, operator, value
	found := false
	operators := []string{"!=", "==", ">=", "<=", "=", ">", "<"}
	for _, i := range operators {
		if strings.Contains(eval, i) == true {
			e.operator = i
			if e.operator == "==" {
				// Convert to single equals sign for sql
				e.operator = "="
			}
			s := strings.Split(eval, i)
			if len(s) == 2 {
				// Only store properly formed queries
				e.column = strings.TrimSpace(s[0])
				e.value = strings.TrimSpace(s[1])
				found = true
			}
			break
		}
	}
	if found == false {
		fmt.Print("\n\t[Error] Please supply a valid evaluation argument. Exiting.\n\n")
		os.Exit(1001)
	}
}

func setOperations(eval string) []evaluation {
	// Returns slice of evaluation targets
	var ret []evaluation
	for _, i := range strings.Split(eval, ",") {
		var e evaluation
		e.getOperation(i)
		ret = append(ret, e)
	}
	if len(ret) == 0 {
		fmt.Print("\n\t[Error] Please supply an evaluation argument. Exiting.\n\n")
		os.Exit(1002)
	}
	return ret
}

func convertValues(v1, v2 string) (float64, float64, bool) {
	// Converts values to float for comparison
	var r2 float64
	var ret bool
	r1, err := strconv.ParseFloat(v1, 64)
	if err == nil {
		r2, err = strconv.ParseFloat(v2, 64)
		if err == nil {
			ret = true
		}
	}
	return r1, r2, ret
}

func (e *evaluation) evaluateLine(h map[string]int, row []string) bool {
	// Applies evaluation to line, return true if it passes
	ret := false
	idx, ex := h[e.column]
	if ex == true {
		switch e.operator {
		case "!=":
			if row[idx] != e.value {
				ret = true
			}
		case "=":
			if row[idx] == e.value {
				ret = true
			}
		case ">=":
			v1, v2, pass := convertValues(row[idx], e.value)
			if pass == true {
				if v1 >= v2 {
					ret = true
				}
			}
		case "<=":
			v1, v2, pass := convertValues(row[idx], e.value)
			if pass == true {
				if v1 <= v2 {
					ret = true
				}
			}
		case ">":
			v1, v2, pass := convertValues(row[idx], e.value)
			if pass == true {
				if v1 > v2 {
					ret = true
				}
			}
		case "<":
			v1, v2, pass := convertValues(row[idx], e.value)
			if pass == true {
				if v1 < v2 {
					ret = true
				}
			}
		}
	} else {
		fmt.Printf("\t[Warning] Column %s not present in header. Skipping.\n", e.column)
	}
	return ret
}

func filterSearchResults(header string, e []evaluation, res [][]string) [][]string {
	// Applies filters to search results
	var ret [][]string
	fmt.Println("\tFiltering search results...")
	h := iotools.GetHeader(strings.Split(header, ","))
	for _, i := range res {
		keep := true
		for idx := range e {
			keep = e[idx].evaluateLine(h, i)
			if keep == false {
				break
			}
		}
		if keep == true {
			ret = append(ret, i)
		}
	}
	fmt.Printf("\tFound %d records that passed additional search parameters.\n", len(ret))
	return ret
}

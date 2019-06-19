// Applys addional search filters to search results

package codbutils

import (
	"fmt"
	"os"
	"strings"
)

type Evaluation struct {
	Column   string
	Operator string
	Value    string
}

func (e *Evaluation) getOperation(eval string) {
	// Splits eval into column, operator, value
	found := false
	operators := []string{"!=", "==", ">=", "<=", "=", ">", "<"}
	for _, i := range operators {
		if strings.Contains(eval, i) == true {
			e.Operator = i
			if e.Operator == "==" {
				// Convert to single equals sign for sql
				e.Operator = "="
			}
			s := strings.Split(eval, i)
			if len(s) == 2 {
				// Only store properly formed queries
				e.Column = strings.TrimSpace(s[0])
				e.Value = strings.TrimSpace(s[1])
				found = true
			}
			break
		}
	}
	if found == false {
		fmt.Printf("\n\t[Error] %s is not a valid evaluation argument. Exiting.\n\n", eval)
		os.Exit(1001)
	}
}

func SetOperations(eval string) []Evaluation {
	// Returns slice of evaluation targets
	var ret []Evaluation
	for _, i := range strings.Split(eval, ",") {
		var e Evaluation
		e.getOperation(i)
		ret = append(ret, e)
	}
	if len(ret) == 0 {
		fmt.Print("\n\t[Error] Please supply an evaluation argument. Exiting.\n\n")
		os.Exit(1002)
	}
	return ret
}

/*func convertValues(v1, v2 string) (float64, float64, bool) {
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

func (e *Evaluation) evaluateLine(h map[string]int, row []string) bool {
	// Applies evaluation to line, return true if it passes
	ret := false
	idx, ex := h[e.Column]
	if ex == true && idx < len(row) {
		switch e.Operator {
		case "!=":
			if row[idx] != e.Value {
				ret = true
			}
		case "=":
			if row[idx] == e.Value {
				ret = true
			}
		case ">=":
			v1, v2, pass := convertValues(row[idx], e.Value)
			if pass == true {
				if v1 >= v2 {
					ret = true
				}
			}
		case "<=":
			v1, v2, pass := convertValues(row[idx], e.Value)
			if pass == true {
				if v1 <= v2 {
					ret = true
				}
			}
		case ">":
			v1, v2, pass := convertValues(row[idx], e.Value)
			if pass == true {
				if v1 > v2 {
					ret = true
				}
			}
		case "<":
			v1, v2, pass := convertValues(row[idx], e.Value)
			if pass == true {
				if v1 < v2 {
					ret = true
				}
			}
		}
	} else {
		fmt.Printf("\t[Warning] Column %s not present in header. Skipping.\n", e.Column)
	}
	return ret
}

func FilterSearchResults(header string, e []Evaluation, res [][]string) [][]string {
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
}*/

// Applys addional search filters to search results

package codbutils

import (
	"fmt"
	"github.com/icwells/go-tools/dataframe"
	"strconv"
	"strings"
)

type Evaluation struct {
	Table    string
	ID       string
	Column   string
	Operator string
	Value    string
}

func (e *Evaluation) String() string {
	// Returns fromatted string
	val := e.Value
	// Wrap strings in quotes
	if _, err := strconv.ParseFloat(val, 64); err != nil {
		val = fmt.Sprintf("'%s'", val)
	}
	return fmt.Sprintf("%s %s %s", e.Column, e.Operator, val)
}

func (e *Evaluation) SetIDType(columns map[string]string) {
	// Sets target id type
	tid := "taxa_id"
	if e.Table != "Patient" && strings.Contains(columns[e.Table], tid) {
		e.ID = tid
	} else {
		e.ID = "ID"
	}
}

func (e *Evaluation) SetTable(columns map[string]string, quit bool) string {
	// Wraps call to GetTable to set table and id type
	var ret string
	if quit == true {
		e.Table = GetTable(columns, e.Column)
	} else {
		e.Table, ret = FindTable(columns, e.Column)
	}
	if e.Table != "" {
		ret = ""
		e.SetIDType(columns)
	}
	return ret
}

func (e *Evaluation) SetOperation(eval string) {
	// Splits eval into column, operator, value
	found := false
	operators := []string{"!=", "==", ">=", "<=", "=", ">", "<", "^"}
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
		GetLogger().Fatalf("%s is not a valid evaluation argument. Exiting.\n", eval)
	}
}

func SetOperations(columns map[string]string, eval string) [][]Evaluation {
	// Returns slice of evaluation targets
	var ret []Evaluation
	for _, i := range strings.Split(eval, ",") {
		if i != "nil" {
			var e Evaluation
			e.SetOperation(i)
			e.SetTable(columns, true)
			ret = append(ret, e)
		}
	}
	if len(ret) == 0 {
		GetLogger().Fatalf("Please supply an evaluation argument. Exiting.\n")
	}
	return [][]Evaluation{ret}
}

func RecordsEvaluations(columns map[string]string, eval string) []Evaluation {
	// Returns slice of evaluations for Records view
	var ret []Evaluation
	records := "Records"
	col := columns[records]
	for _, i := range strings.Split(eval, ",") {
		if i != "nil" {
			var e Evaluation
			e.SetOperation(i)
			if e.Column == "Age" {
				// Update age column name
				e.Column = "age_months"
			}
			if strings.Contains(col, e.Column) {
				e.Table = records
				e.ID = "ID"
				ret = append(ret, e)
			} else {
				GetLogger().Fatalf("Target column %s not in Records table. Exiting.\n", e.Column)
			}
		}
	}
	if len(ret) == 0 {
		GetLogger().Fatal("Please supply an evaluation argument. Exiting.")
	}
	return ret
}


func OperationsFromFile(columns map[string]string, infile string) [][]Evaluation {
	// Reads evaluations from input file
	var ret [][]Evaluation
	df, err := dataframe.FromFile(infile, -1)
	if err != nil {
		GetLogger().Fatalf("Cannot read evaluation file: %v\n", err)
	}
	header := df.GetHeader()
	for _, row := range df.Rows {
		var eval []Evaluation
		for idx, i := range row {
			if len(strings.TrimSpace(i)) > 0 {
				var e Evaluation
				e.SetOperation(fmt.Sprintf("%s = %s", header[idx], i))
				e.SetTable(columns, true)
				eval = append(eval, e)
			}
		}
		if len(eval) > 0 {
			ret = append(ret, eval)
		}
	}
	return ret
}

//----------------------------------------------------------------------------

func tableFromID(col string) string {
	// Returns table for id columns present in multiple tables
	var ret string
	switch col {
	case "id":
		ret = "Patient"
	case "taxa_id":
		ret = "Taxonomy"
	case "account_id":
		ret = "Source"
	case "source_id":
		ret = "Patient"
	}
	return ret
}

func FindTable(tables map[string]string, col string) (string, string) {
	// Returns single table name and error message. Exits if there is an error and quit is true
	var ret string
	msg := fmt.Sprintf("Cannot find table with column %s.", col)
	col = strings.ToLower(col)
	if col == "id" || strings.Contains(col, "_id") {
		ret = tableFromID(col)
	} else {
		if strings.Contains(col, "_") == false {
			col = strings.Title(col)
		}
		// Iterate through available column names
		for k, val := range tables {
			if k != "Records" {
				for _, i := range strings.Split(val, ",") {
					i = strings.TrimSpace(i)
					if col == i {
						ret = k
						break
					}
				}
			}
		}
	}
	return ret, msg
}

func GetTable(tables map[string]string, col string) string {
	// Determines which table column is in, exits if there is an error
	ret, msg := FindTable(tables, col)
	if len(ret) == 0 {
		GetLogger().Fatalf("%s Exiting.\n", msg)
	}
	return ret
}

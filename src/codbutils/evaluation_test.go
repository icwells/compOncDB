// Tests multiSearch functions and methods

package codbutils

import (
	"testing"
)

func TestSetOperations(t *testing.T) {
	// Tests the getOperation function
	matches := []struct {
		input    string
		table    string
		id       string
		column   string
		operator string
		value    string
	}{
		{"Species == Canis lupus", "Taxonomy", "taxa_id", "Species", "=", "Canis lupus"},
		{"Sex=male", "Patient", "ID", "Sex", "=", "male"},
		{"Avgage>=12", "Totals", "taxa_id", "Avgage", ">=", "12"},
		{" Cancer < 5 ", "Totals", "taxa_id", "Cancer", "<", "5"},
	}
	columns := getTableColumns()
	for _, i := range matches {
		var msg string
		evaluations := SetOperations(columns, i.input)
		e := evaluations[0][0]
		if e.Table != i.table {
			t.Errorf("Actual table %s is not equal to expected: %s", e.Table, i.table)
		} else if e.ID != i.id {
			t.Errorf("Actual id type %s is not equal to expected: %s", e.ID, i.id)
		} else if e.Column != i.column {
			t.Errorf("Actual table column %s is not equal to expected: %s", e.Column, i.column)
		} else if e.Operator != i.operator {
			t.Errorf("Actual table operator %s is not equal to expected: %s", e.Operator, i.operator)
		} else if e.Value != i.value {
			t.Errorf("Actual table value %s is not equal to expected: %s", e.Value, i.value)
		}
		if len(msg) > 1 {
			t.Error(msg)
		}
	}
}

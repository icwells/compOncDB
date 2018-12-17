// White box tests on updateEntries methods

package dbextract

import (
	"strings"
	"testing"
)

func getColumns() map[string]string {
	// Returns map of columns names
	return map[string]string{
		"Patient": "ID,Sex,Age,Castrated,taxa_id,source_id,scientific_name,Date,Comments",
		"Diagnosis": "ID,Masspresent,Hyperplasia,Necropsy,Metastasis",
	}
}

func TestFormatHeader(t *testing.T) {
	// Tests updater.formatHeader
	c := getColumns()
	u := newUpdater(c)
	for k, v := range c {
		var testcase string
		if k == "Patient" {
			testcase =  strings.ToUpper(v)
		} else {
			testcase = strings.ToLower(v)
		}
		input := strings.Split(testcase, ",")
		expected := strings.Split(v, ",")
		actual := u.formatHeader(input)
		if len(actual) != len(expected) {
			t.Errorf("Actual header length %d does not equal expected: %d", len(actual), len(expected))
		} else {
			for idx, i := range actual {
				if i != expected[idx] {
					t.Errorf("Actual value %s does not equal expected: %s", i, expected[idx])
				}
			}
		}
	}
}

func getUpdater() updater {
	// Returns initialized updater with columns map set
	u := newUpdater(getColumns())
	u.setColumns([]string{"ID", "Age", "Sex"})
	return u
}

func TestSetColumns(t *testing.T) {
	// Tests updater.setColumns
	u := getUpdater()
	if _, ex := u.columns["Diagnosis"]; ex == true {
		t.Errorf("Diagnosis table was not deleted from columns map.")
	}
	row, ex := u.columns["Patient"]
	if ex == false {
		t.Errorf("Patient table was deleted from columns map.")
	} else {
		for idx, i := range row {
			if idx == 1 && i != 2 {
				t.Errorf("Actual index for sex column %d does not equal 2", i)
			} else if idx == 2 && i != 1 {
				t.Errorf("Actual index for age column %d does not equal 1", i)
			} else if i != -1 {
				t.Errorf("Actual index %d does not equal -1", i)
			}
		}
	}
	if _, ex := u.tables["Diagnosis"]; ex == true {
		t.Errorf("Diagnosis table was initialized in tables map.")
	}
	s, e := u.tables["Patient"]
	if e == false {
		t.Errorf("Patient table was not initialized in tables map.")
	} else {
		if s.table != "Patient" {
			t.Errorf("Incorrect table name %s stored in tableupdate struct", s.table)
		} else if s.target != "ID" {
			t.Errorf("Incorrect target name %s stored in tableupdate struct", s.target)
		}
	}
}

func TestEvaluateRow(t *testing.T) {
	// Tests updater.evaluateRow
	u := getUpdater()
	cases := []struct {
		input		[]string
		expected	[]string
	}{
		{[]string{"1", "12", "male"}, []string{"male", "12", "", "", "", "", "", ""}},
		{[]string{"2", "5", "female"}, []string{"female", "5", "", "", "", "", "", ""}},
		{[]string{"3", "", "female"}, []string{"female", "", "", "", "", "", "", ""}},
		{[]string{"", "5", "female"}, []string{"", "", "", "", "", "", "", ""}},
	}
	for _, i := range cases {
		u.evaluateRow(i.input)
		row, ex := u.tables["Patient"].values[i.input[0]]
		if len(i.input[0]) >= 1 {
			if ex == false {
				t.Errorf("Row %s was not uploaded", i.input[0])
			} else {
				for idx, r := range row {
					if r != i.expected[idx] {
						t.Errorf("Row %s: actual value %s does not equal expected: %s", i.input[0], r, i.expected[idx])
					}
				}
			}
		} else if ex == true{
			t.Errorf("Value with missing ID stored in map.")
		}
	}
}

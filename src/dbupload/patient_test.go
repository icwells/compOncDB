// White box tests for patients script

package dbupload

import (
	"testing"
)

func TestTumorPairs(t *testing.T) {
	input := []string{"carcinoma;sarcoma", "Gastrointestinal;Round Cell", "liver;lymph nodes; teeth"}
	expected := [][]string{
		[]string{"carcinoma", "Gastrointestinal", "liver"},
		[]string{"sarcoma", "Round Cell", "lymph nodes"},
	}
	actual := tumorPairs(input[0], input[1], input[2])
	if len(actual) != len(expected) {
		t.Errorf("Actual length %d does not equal expected: %d", len(actual), len(expected))
	} else {
		for idx, i := range actual {
			if i[0] != expected[idx][0] || i[1] != expected[idx][1] || i[2] != expected[idx][2] {
				t.Errorf("Actual pair %s:%s:%s does not equal expected: %s:%s:%s", i[0], i[1], i[2], expected[idx][0], expected[idx][1], expected[idx][2])
			}
		}
	}
}

func setEntries() *entries {
	// Returns test entry struct
	e := newEntries(nil, false)
	e.count = 0
	e.col = make(map[string]int)
	s := []string{"Sex", "Age", "Castrated", "ID", "Genus", "Species", "Name", "Date", "Year", "Comments", "MassPresent", "Hyperplasia", "Necropsy", "Metastasis", "TumorType", "Tissue", "Location", "Primary", "Malignant", "Service", "Account", "Submitter", "Zoo", "AZA", "Institute"}
	for idx, i := range s {
		e.col[i] = idx
	}
	e.length = len(s)
	e.taxa = map[string]string{
		"Canis latrans": "1",
		"Canis lupus":   "2",
	}
	e.submitter = make(map[string]string)
	e.submitter["XYZ"] = "1"
	return e
}

func getExpected() *entries {
	// Returns pre-filled struct of expected results
	e := newEntries(nil, false)
	e.count = 4
	e.p = [][]string{
		[]string{"1", "male", "-1.00", "-1", "-1", "0", "1", "1", "coyote", "12-Dec", "2001", "Biopsy: NORMAL BLOOD SMEAR"},
		[]string{"2", "NA", "-1.00", "-1", "-1", "0", "1", "2", "coyote", "13-Jan", "2001", "ERYTHROPHAGOCYTOSIS"},
		[]string{"3", "male", "24.00", "-1", "-1", "0", "1", "3", "coyote", "1-Dec", "2001", "Lymphoma lymph nodes 2 year old male"},
		[]string{"4", "NA", "-1.00", "-1", "-1", "0", "1", "4", "coyote", "1-Dec", "2001", "HIPOTOMAS TOXIC HIPOTOPATHY autopsy"},
	}
	e.d = [][]string{
		[]string{"1", "0", "0", "0", "-1"},
		[]string{"2", "0", "0", "-1", "-1"},
		[]string{"3", "1", "0", "-1", "-1"},
		[]string{"4", "0", "0", "1", "-1"},
	}
	e.t = [][]string{
		[]string{"1", "0", "-1", "carcinoma", "Gastrointestinal", "liver"},
		[]string{"1", "0", "-1", "sarcoma", "Epithelial", "skin"},
		[]string{"2", "0", "-1", "NA", "NA", "NA"},
		[]string{"3", "0", "1", "lymphoma", "Round Cell", "lymph nodes"},
		[]string{"4", "0", "-1", "NA", "NA", "NA"},
	}
	e.s = [][]string{
		[]string{"1", "NWZP", "-1", "0", "-1", "-1", "1"},
		[]string{"2", "NWZP", "-1", "0", "-1", "-1", "1"},
		[]string{"3", "NWZP", "-1", "0", "-1", "-1", "1"},
		[]string{"4", "NWZP", "-1", "0", "-1", "-1", "1"},
	}
	return e
}

func getInput() [][]string {
	// Returns input slice for testing
	return [][]string{
		[]string{"male", "-1", "-1", "1", "Canis", "Canis latrans", "coyote", "12-Dec", "2001", "Biopsy: NORMAL BLOOD SMEAR", "0", "0", "0", "-1", "carcinoma;sarcoma", "Gastrointestinal;Epithelial", "liver;skin", "0", "-1", "NWZP", "X520", "XYZ", "-1", "0", "-1"},
		[]string{"NA", "-1", "-1", "2", "Canis", "Canis latrans", "coyote", "13-Jan", "2001", "ERYTHROPHAGOCYTOSIS", "0", "0", "-1", "-1", "NA", "NA", "NA", "0", "-1", "NWZP", "X520", "XYZ", "-1", "0", "-1"},
		[]string{"male", "24", "-1", "3", "Canis", "Canis latrans", "coyote", "1-Dec", "2001", "Lymphoma lymph nodes 2 year old male", "1", "0", "-1", "-1", "lymphoma", "Round Cell", "lymph nodes", "0", "1", "NWZP", "X520", "XYZ", "-1", "0", "-1"},
		[]string{"NA", "-1", "-1", "4", "Canis", "Canis latrans", "coyote", "1-Dec", "2001", "HIPOTOMAS TOXIC HIPOTOPATHY autopsy", "0", "0", "1", "-1", "NA", "NA", "NA", "0", "-1", "NWZP", "X520", "XYZ", "-1", "0", "-1"},
	}
}

func TestExists(t *testing.T) {
	a := setEntries()
	input := getInput()
	for _, i := range input {
		if a.ex.Exists("1", i[a.col["ID"]], i[a.col["Age"]], "1", i[a.col["Date"]]) {
			t.Error("Exists returned true from an empty struct")
		}
	}
	a.ex.Entries["1"] = make(map[string]*Entry)
	for _, i := range input {
		a.ex.Entries["1"][i[a.col["ID"]]] = NewEntry([]string{i[a.col["Age"]], "1", i[a.col["Date"]]})
	}
	for _, i := range input {
		if !a.ex.Exists("1", i[a.col["ID"]], i[a.col["Age"]], "1", i[a.col["Date"]]) {
			t.Error("Exists returned false from a full struct")
		} else if a.ex.Exists("1", i[a.col["ID"]], i[a.col["Age"]], "2", i[a.col["Date"]]) {
			t.Error("Exists returned true for incorrect taxa id")
		}
	}
}

func compareTables(t *testing.T, table string, a, e [][]string) {
	// Comapres actual table to expected
	if len(a) != len(e) {
		t.Errorf("%s: Actual length %v does not equal expected: %d", table, len(a), len(e))
	} else {
		for ind, row := range a {
			for idx, i := range row {
				if i != e[ind][idx] {
					t.Errorf("%s %d: Actual value %s does not equal expected: %s", table, idx, i, e[ind][idx])
					break
				}
			}
		}
	}
}

func TestEvaluateRow(t *testing.T) {
	// Tests evaluate row and all methods called by it
	a := setEntries()
	e := getExpected()
	input := getInput()
	for _, i := range input {
		a.evaluateRow(i)
	}
	compareTables(t, "patient", a.p, e.p)
	compareTables(t, "diangosis", a.d, e.d)
	compareTables(t, "tumor", a.t, e.t)
	compareTables(t, "source", a.s, e.s)
}

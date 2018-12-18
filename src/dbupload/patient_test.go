// White box tests for patients script

package dbupload

import (
	"testing"
)

func TestTumorPairs(t *testing.T) {
	input := []string{"carcinoma;sarcoma", "liver;lymph nodes; teeth"}
	expected := [][]string{
		[]string{"carcinoma", "liver"},
		[]string{"sarcoma", "lymph nodes"},
	}
	actual := tumorPairs(input[0], input[1])
	if len(actual) != len(expected) {
		t.Errorf("Actual length %d does not equal expected: %d", len(actual), len(expected))
	} else {
		for idx, i := range actual {
			if i[0] != expected[idx][0] || i[1] != expected[idx][1] {
				t.Errorf("Actual pair %s:%s does not equal expected: %s:%s", i[0], i[1], expected[idx][0], expected[idx][1])
			}
		}
	}
}

func setEntries() *entries {
	// Returns test entry struct
	e := newEntries(0)
	e.length = 18
	e.col = map[string]int{
		"Sex": 0, "Age": 1, "Castrated": 2, "ID": 3, "Species": 4, "Date": 5, "Comments": 6, "MassPresent": 7, "Hyperplasia": 8,
		"Necropsy": 9, "Metastasis": 10, "Type": 11, "Location": 12, "Primary": 13, "Malignant": 14, "Service": 15, 
		"Account": 16, "Submitter": 17,}
	e.species = map[string]string{
		"Canis latrans": "1",
		"Canis lupus": "2",
	}
	e.accounts["X520"] = make(map[string]string)
	e.accounts["X520"]["XYZ"] = "1"
	return e
}

func getExpected() *entries {
	// Returns pre-filled struct of expected results
	e := newEntries(4)
	e.p = [][]string{
		[]string{"1", "male", "-1", "-1", "1", "1", "12-Dec", "Biopsy: NORMAL BLOOD SMEAR"},
		[]string{"2", "NA", "-1", "-1", "1", "2", "13-Jan", "ERYTHROPHAGOCYTOSIS"},
		[]string{"3", "male", "24", "-1", "1", "3", "1-Dec", "Lymphoma lymph nodes 2 year old male"},
		[]string{"4", "NA", "-1", "-1", "1", "4", "1-Dec", "HIPOTOMAS TOXIC HIPOTOPATHY autopsy"},
	}
	e.d = [][]string{
		[]string{"1", "0", "0", "0", "-1"},
		[]string{"2", "0", "0", "-1", "-1"},
		[]string{"3", "1", "0", "-1", "-1"},
		[]string{"4", "0", "0", "1", "-1"},
	}
	e.t = [][]string{
		[]string{"1", "0", "-1", "carcinoma", "liver"},
		[]string{"1", "0", "-1", "sarcoma", "skin"},
		[]string{"2", "0", "-1", "NA", "NA"},
		[]string{"3", "0", "1", "lymphoma", "lymph nodes"},
		[]string{"4", "0", "-1", "NA", "NA"},
	}
	e.s = [][]string{
		[]string{"1", "NWZP", "1"},
		[]string{"2", "NWZP", "1"},
		[]string{"3", "NWZP", "1"},
		[]string{"4", "NWZP", "1"},
	}
	return e
}

func comapareTables(t *testing.T, table string, a, e [][]string) {
	// Comapres actual table to expected
	if len(a) != len(e) {
		t.Errorf("%s: Actual length %d does not equal expected: %d", table, len(a), len(e))
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
	input := [][]string{
		[]string{"male", "-1", "-1", "1", "Canis latrans", "12-Dec", "Biopsy: NORMAL BLOOD SMEAR", "0", "0", "0", "-1", "carcinoma;sarcoma", "liver;skin", "0", "-1", "NWZP", "X520", "XYZ"},
		[]string{"NA", "-1", "-1", "2", "Canis latrans", "13-Jan", "ERYTHROPHAGOCYTOSIS", "0", "0", "-1", "-1", "NA", "NA", "0", "-1", "NWZP", "X520", "XYZ"},
		[]string{"male", "24", "-1", "3", "Canis latrans", "1-Dec", "Lymphoma lymph nodes 2 year old male", "1", "0", "-1", "-1", "lymphoma", "lymph nodes", "0", "1", "NWZP", "X520", "XYZ"},
		[]string{"NA", "-1", "-1", "4", "Canis latrans", "1-Dec", "HIPOTOMAS TOXIC HIPOTOPATHY autopsy", "0", "0", "1", "-1", "NA", "NA", "0", "-1", "NWZP", "X520", "XYZ"},
	}
	for _, i := range input {
		a.evaluateRow(i)
	}
	comapareTables(t, "patient", a.p, e.p)
	comapareTables(t, "diangosis", a.d, e.d)
	comapareTables(t, "tumor", a.t, e.t)
	comapareTables(t, "source", a.s, e.s)
}

// Performs white box tests on functions in the compOncDB utils script

package main

import (
	"fmt"
	"testing"
)

func fmtMessage(field, a, e string) string {
	// Returns formatted string
	return fmt.Sprintf("Actual %s %s is not equal to expected: %s", field, a, e)
}

func TestToMap(t *testing.T) {
	// Tests toMap function
	expected1 := map[string][]string{
		"1": {"a"},
		"2": {"b"},
		"3": {"c", "d"},
	}
	expected2 := map[string][]string{
		"1": {"a", "a"},
		"2": {"b", "b"},
		"3": {"c", "c"},
	}
	slice := [][]string{
		{"1", "a"},
		{"2", "b"},
		{"3", "c"},
		{"3", "d"},
	}
	actual := toMap(slice)
	for k, v := range actual {
		for idx, i := range v {
			if i != expected1[k][idx] {
				msg := fmtMessage("appended map value", i, expected1[k][idx])
				t.Error(msg)
			}
		}
	}
	for idx, i := range slice {
		// Lengthen inner slice
		slice[idx] = append(i, i[1])
	}
	actual = toMap(slice)
	for k, v := range actual {
		for idx, i := range v {
			if i != expected2[k][idx] {
				msg := fmtMessage("single map value", i, expected2[k][idx])
				t.Error(msg)
			}
		}
	}
}

func TestMapOfMaps(t *testing.T) {
	// Tests mapOfMaps
	expected := make(map[string]map[string]string)
	expected["1"] = make(map[string]string)
	expected["1"]["a"] = "aa"
	expected["2"] = make(map[string]string)
	expected["2"]["b"] = "bb"
	expected["3"] = make(map[string]string)
	expected["3"]["c"] = "cc"
	expected["3"]["d"] = "dd"
	table := [][]string{
		{"aa", "a", "1"},
		{"bb", "b", "2"},
		{"cc", "c", "3"},
		{"dd", "d", "3"},
	}
	actual := mapOfMaps(table)
	for key, val := range actual {
		for k, v := range val {
			if v != expected[key][k] {
				//msg := fmtMessage("map value", v, expected[key][k])
				t.Error( expected[key])
			}
		}
	}
}

/*func TestEntryMap(t *testing.T) {

}

func TestGetOperation(t *testing.T) {

}

func TestGetTable(t *testing.T) {

}*/

// Tests fucntions in utils

package dbupload

import (
	"testing"
)

func TestCalculateInfancy(t *testing.T) {
	// Tests CalculateInfancy in lifeHistory.go
	cases := []struct {
		male, female, weaning, expected string
	}{
		{"274", "274", "2.07100591715976", "2.07100591715976"},
		{"669", "669", "2.03813280736358", "2.03813280736358"},
		{"345", "365", "1.87376725838264", "1.87376725838264"},
		{"NA", "v", "NA", "1"},
	}
	for _, i := range cases {
		actual := calculateInfancy(i.weaning, i.male, i.female)
		if actual != i.expected {
			t.Errorf("Actual infant age %s does not equal expected: %s", actual, i.expected)
		}
	}
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
	actual := ToMap(slice)
	for k, v := range actual {
		for idx, i := range v {
			if i != expected1[k][idx] {
				t.Errorf("Actual appended map value %s does not equal expected: %s", i, expected1[k][idx])
			}
		}
	}
	for idx, i := range slice {
		// Lengthen inner slice
		slice[idx] = append(i, i[1])
	}
	actual = ToMap(slice)
	for k, v := range actual {
		for idx, i := range v {
			if i != expected2[k][idx] {
				t.Errorf("Actual single map value %s does not equal expected: %s", i, expected1[k][idx])
			}
		}
	}
}

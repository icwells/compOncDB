// Performs white box tests on functions in the compOncDB utils script

package codbutils

import (
	"testing"
)

func getTableColumns() map[string]string {
	// Returns map of table columns with not types
	return map[string]string{
		"Patient":        "ID,Sex,Age,Castrated,taxa_id,source_id,scientific_name,Date,Comments",
		"Taxonomy":       "taxa_id,Kingdon,Phylum,Class,Orders,Family,Genus,Species,Source",
		"Common":         "taxa_id,Name",
		"Totals":         "taxa_id,Total,Avgage,Adult,Male,Female,Cancer,Cancerage,Malecancer,Femalecancer",
		"Denominators":   "taxa_id,Noncancer",
		"Source":         "ID,service_name,account_id",
		"Accounts":       "account_id,Account,submitter_name",
		"Diagnosis":      "ID,Masspresent,Necropsy,Metastasis",
		"Tumor":          "tumor_id,Type,Location",
		"Tumor_relation": "ID,tumor_id,primary_tumor,Malignant",
		"Life_history":   "taxa_id,female_maturity,male_maturity,Gestation,Weaning,litter_size,litters_year,interbirth_interval,birth_weight,weaning_weight,adult_weight,growth_rate,max_longevity,metabolic_rate",
	}
}

func TestGetTable(t *testing.T) {
	// Tests getTable function
	col := getTableColumns()
	matches := []struct {
		column string
		table  string
	}{
		{"sex", "Patient"},
		{"taxa_id", "Taxonomy"},
		{"TOTAL", "Totals"},
		{"account_ID", "Source"},
		{"Id", "Patient"},
		{"Primary_Tumor", "Tumor_relation"},
	}
	for _, e := range matches {
		actual := GetTable(col, e.column)
		if actual != e.table {
			t.Errorf("Actual table %s is not equal to expected: %s", actual, e.table)
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

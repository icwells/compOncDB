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

func TestGetOperation(t *testing.T) {
	// Tests the getOperation function
	matches := []struct{
		input		string
		column		string
		operator	string
		value		string
	} {
		{"Species == Canis lupus", "Species", "=", "Canis lupus"},
		{"Sex=male", "Sex", "=", "male"},
		{"Avgage>=12", "Avgage", ">=", "12"},
		{" Cancer < 5 ", "Cancer", "<", "5"},
	}
	for _, i := range matches {
		var msg string
		col, op, val := getOperation(i.input)
		if col != i.column {
			msg = fmtMessage("column", col, i.column)
		} else if op != i.operator {
			msg = fmtMessage("operator", op, i.operator)
		} else if val != i.value {
			msg = fmtMessage("value", val, i.value)
		}
		if len(msg) > 1 {
			t.Error(msg)
		}
	}
}

func getTableColumns() map[string]string {
	// Returns map of table columns with not types
	return map[string]string{
		"Patient": "ID,Sex,Age,Castrated,taxa_id,source_id,scientific_name,Date,Comments",
		"Taxonomy": "taxa_id,Kingdon,Phylum,Class,Orders,Family,Genus,Species,Source",
		"Common": "taxa_id,Name",
		"Totals": "taxa_id,Total,Avgage,Adult,Male,Female,Cancer,Cancerage,Malecancer,Femalecancer",
		"Denominators": "taxa_id,Noncancer",
		"Source": "ID,service_name,account_id",
		"Accounts": "account_id,Account,submitter_name",
		"Diagnosis": "ID,Masspresent,Necropsy,Metastasis",
		"Tumor": "tumor_id,Type,Location",
		"Tumor_relation": "ID,tumor_id,primary_tumor,Malignant",
		"Life_history": "taxa_id,female_maturity,male_maturity,Gestation,Weaning,litter_size,litters_year,interbirth_interval,birth_weight,weaning_weight,adult_weight,growth_rate,max_longevity,metabolic_rate",
	}
}

func TestGetTable(t *testing.T) {
	// Tests getTable function
	col := getTableColumns()
	matches := []struct{
		column	string
		tables	[]string
	} {
		{"sex", []string{"Patient"}},
		{"taxa_id", []string{"Patient", "Taxonomy", "Common", "Totals", "Life_history"}},
		{"TOTAL", []string{"Totals"}},
		{"account_ID", []string{"Source", "Accounts"}},
		{"Id", []string{"Patient", "Source", "Diagnosis", "Tumor_relation"}},
		{"Primary_Tumor", []string{"Tumor_relation"}},
	}
	for _, e := range matches {
		actual := getTable(col, e.column)
		for idx, i := range actual {
			if i != e.tables[idx] {
				msg := fmtMessage("table", i, e.tables[idx])
				t.Error(msg)
			}
		}
	}
}

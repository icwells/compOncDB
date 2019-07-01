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
		table string
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

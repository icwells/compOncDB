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
		e := evaluations[0]
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

/*func TestConvertValues(t *testing.T) {
	cases := []struct {
		v1, v2 string
		f1, f2 float64
		pass   bool
	}{
		{"2", "3", 2.0, 3.0, true},
		{"6.3", "10.1", 6.3, 10.1, true},
		{"16", "NA", 16.0, 0.0, false},
	}
	for _, i := range cases {
		a1, a2, pass := convertValues(i.v1, i.v2)
		if pass != i.pass {
			t.Errorf("Actual pass value %v does not equal expected: %v", pass, i.pass)
		} else if a1 != i.f1 {
			t.Errorf("Actual float %f does not equal expected: %f", a1, a2)
		} else if a2 != i.f2 {
			t.Errorf("Actual float %f does not equal expected: %f", a2, a2)
		}
	}
}

type filtercases struct {
	row      string
	eval     evaluation
	expected bool
}

func newFilterCases() []filtercases {
	// Returns test struct
	var r1, r2, r3, r4 filtercases
	// Store rows as strings to save keep readable
	r1.row = "1306,NA,-1,-1,369,1351,23-Apr,HEMANGIOSARCOMAS,0,0,-1,-1,0,-1,NA,NA,Animalia,Chordata,Mammalia,Carnivora,Canidae,Canis,Canis lupus,NWZP,968"
	r1.eval.getOperation("Castrated==-1")
	r1.expected = true
	r2.row = "1337,male,60,-1,369,1382,27-Apr,VASCULAR ANOMALY,0,0,1,-1,0,-1,male,NA,Animalia,Chordata,Mammalia,Carnivora,Canidae,Canis,Canis lupus,NWZP,912"
	r2.eval.getOperation("Sex!=male")
	r2.expected = false
	r3.row = "1451,female,36,-1,369,1500,24-May,RENAL ADENOCARCINOMA,0,0,1,-1,0,-1,female,NA,Animalia,Chordata,Mammalia,Carnivora,Canidae,Canis,Canis lupus,NWZP,913"
	r3.eval.getOperation("Sex!=male")
	r3.expected = true
	r4.row = "1306,NA,19,-1,369,1351,23-Apr,HEMANGIOSARCOMAS,0,0,-1,-1,0,-1,NA,NA,Animalia,Chordata,Mammalia,Carnivora,Canidae,Canis,Canis lupus,NWZP,968"
	r4.eval.getOperation("Age>12")
	r4.expected = true
	return []filtercases{r1, r2, r3, r4}
}

func TestEvaluateLine(t *testing.T) {
	header := "ID,Sex,Age,Castrated,taxa_id,source_id,Date,Comments,"
	header += "Masspresent,Hyperplasia,Necropsy,Metastasis,primary_tumor,Malignant,Type,Location,"
	header += "Kingdom,Phylum,Class,Order,Family,Genus,Species,service_name,account_id"
	h := iotools.GetHeader(strings.Split(header, ","))
	cases := newFilterCases()
	for idx, i := range cases {
		actual := i.eval.evaluateLine(h, strings.Split(i.row, ","))
		if actual != i.expected {
			t.Errorf("Actual evaluation %v for case %d does not equal expected: %v", actual, idx, i.expected)
		}
	}
}*/

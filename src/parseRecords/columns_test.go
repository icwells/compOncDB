// This script will perform black box tests on parseRecords

package main

import (
	"fmt"
	"testing"
)

func getColumns(headers [][]string) []columns {
	// Returns column structs to match headers slice
	var ret []columns
	for idx := range headers {
		c := newColumns()
		switch (idx) {
			case 0:
				c.id = 0
				c.common = 2
				c.sex = 3
				c.age = 4 
				c.comments = 5
				c.max = 5
			case 1:
				c.id = 0
				c.common = 2
				c.species = 3
				c.age = 10
				c.sex = 11
				c.castrated = 12
				c.location = 13
				c.typ = 14
				c.malignant = 15
				c.primary = 16
				c.metastasis = 17
				c.necropsy = 18
				c.comments = 19
				c.max = 19
			case 2:
				c.submitter = 1
				c.account = 2
				c.id = 0
				c.date = 5
				c.common = 6
				c.patient = 7
				c.code = 8
				c.comments = 9
				c.max = 9
			case 3:
				c.id = 0
				c.date = 1
				c.submitter = 2
				c.patient = 3
				c.common = 5
				c.days = 6
				c.age = 7
				c.sex = 8
				c.comments = 9
				c.max = 9
			case 4:
				c.submitter = 1
				c.id = 0
				c.species = 5
				c.sex = 6
				c.age = 7
				c.code = 8
				c.typ = 9
				c.location = 10
				c.metastasis = 11
				c.necropsy = 13
				c.common = 15
				c.max = 15
		}
		ret = append(ret, c)
	}
	return ret
}

func compareColumns(a, e columns) (bool, string) {
	// Returns false and a fail message if corresponding fields for each struct ar enot equal
	var msg string
	equal := true
	if a.id != e.id {
		equal = false
		msg = fmt.Sprintf("ID columns are not equal.")
	} else if a.species != e.species {
		equal = false
		msg = fmt.Sprintf("Species columns are not equal.")
	} else if a.common != e.common {
		equal = false
		msg = fmt.Sprintf("Common Name columns are not equal.")
	} else if a.age != e.age {
		equal = false
		msg = fmt.Sprintf("Age columns are not equal.")
	} else if a.days != e.days {
		equal = false
		msg = fmt.Sprintf("Days columns are not equal.")
	} else if a.sex != e.sex {
		equal = false
		msg = fmt.Sprintf("Sex columns are not equal.")
	} else if a.castrated != e.castrated {
		equal = false
		msg = fmt.Sprintf("Castrated columns are not equal.")
	} else if a.location != e.location {
		equal = false
		msg = fmt.Sprintf("Location columns are not equal.")
	} else if a.typ != e.typ {
		equal = false
		msg = fmt.Sprintf("Type columns are not equal.")
	} else if a.primary != e.primary {
		equal = false
		msg = fmt.Sprintf("Primary tumor columns are not equal.")
	} else if a.metastasis != e.metastasis {
		equal = false
		msg = fmt.Sprintf("Metastasis columns are not equal.")
	} else if a.malignant != e.malignant {
		equal = false
		msg = fmt.Sprintf("Malignant columns are not equal.")
	} else if a.necropsy != e.necropsy {
		equal = false
		msg = fmt.Sprintf("Necropsy columns are not equal.")
	} else if a.date != e.date {
		equal = false
		msg = fmt.Sprintf("Date columns are not equal.")
	} else if a.comments != e.comments {
		equal = false
		msg = fmt.Sprintf("Comments columns are not equal.")
	} else if a.account != e.account {
		equal = false
		msg = fmt.Sprintf("Accounts columns are not equal.")
	} else if a.submitter != e.submitter {
		equal = false
		msg = fmt.Sprintf("Submitter columns are not equal.")
	} else if a.code != e.code {
		equal = false
		msg = fmt.Sprintf("Code columns are not equal.")
	} else if a.patient != e.patient {
		equal = false
		msg = fmt.Sprintf("Patient columns are not equal.")
	} else if a.max != e.max {
		equal = false
		msg = fmt.Sprintf("Max column values are not equal.")
	}
	return equal, msg
}

func TestColumns(t *testing.T) {
	// Tests coutn NA method
	headers := [][]string{
		{"Access#", "Category", "CommonName", "Sex", "Age", "Diagnosis"},
		{"Access#", "Category", "Breed", "ScientificName", "Kingdom", "Phylum", "Class", "Order", "Family", "Genus", "Age(months)", "Sex", "Castrated", "Location", "Type", "Malignant", "PrimaryTumor", "Metastasis", "Necropsy", "Diagnosis"},
		{"UID", "Client", "Account", "CASE", "ID", "Date_Rcvd", "PT_Name", "Patient", "Code", "Diagnosis", "caseid"},
		{"ID", "Date", "Owner", "Name", "Species", "Breed", "Days", "Age", "Sex", "Description"},
		{"ID", "Institution ID", "Origin ID", "Family", "Genus", "Binomial Scientific", "Sex", "Age", "Cancer Y/N", "Cancer Type", "Tissue", "Metastatic", "Widespread", "Death via Cancer Y/N", "Old/New World", "Common Name"},
	}
	col := getColumns(headers)
	for idx, h := range headers {
		actual := newColumns()
		actual.setColumns(h)
		equal, msg := compareColumns(actual, col[idx])
		if equal == false {
			t.Error(msg)
		}
	}
}

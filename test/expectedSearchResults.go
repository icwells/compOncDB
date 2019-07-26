// Contains expected out for databse search

package main

import (
	"github.com/icwells/compOncDB/src/codbutils"
)

func getCanisResults() map[string][]string {
	// Returns map of results for male canis records
	ret := make(map[string][]string)
	ret["3"] = []string{"male", "24", "-1", "1", "3", "1-Dec", "Lymphoma lymph nodes 2 year old male", "1", "0", "-1", "-1", "0", "1", "lymphoma", "lymph nodes", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis latrans", "NWZP", "0", "0", "1"}
	ret["12"] = []string{"male", "60", "-1", "2", "12", "NA", "NA", "0", "0", "-1", "-1", "0", "-1", "NA", "NA", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus", "NWZP", "1", "0", "2"}
	ret["17"] = []string{"male", "144", "1", "2", "17", "NA", "neutered", "0", "0", "-1", "-1", "0", "-1", "NA", "NA", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus", "NWZP", "1", "0", "2"}
	//ret["19"] = []string{"male", "56", "0", "2", "19", "32688", "NA", "0", "0", "0", "-1", "0", "-1", "NA", "NA", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus", "NWZP", "1", "0", "2"}
	return ret
}

func getCoyoteResults() map[string][]string {
	// Returns map of coyote records
	ret := make(map[string][]string)
	ret["3"] = []string{"male", "24", "-1", "1", "3", "1-Dec", "Lymphoma lymph nodes 2 year old male", "1", "0", "-1", "-1", "0", "1", "lymphoma", "lymph nodes", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis latrans", "NWZP", "0", "0", "1"}
	return ret
}

func getLitterResults() map[string][]string {
	// Returns map of life history results
	ret := make(map[string][]string)
	ret["1"] = []string{"274", "274", "2.07100591715976", "1.9723865877712", "1.9723865877712", "5.72", "1.1", "365", "250", "1517", "13250", "0.0183", "261.6", "19.423"}
	return ret
}

type searchCase struct {
	name     string
	eval     []codbutils.Evaluation
	expected map[string][]string
	table    string
}

func setCase(columns map[string]string, name, eval, table string, exp map[string][]string) searchCase {
	// Returns initilized struct
	var s searchCase
	s.name = name
	s.eval = codbutils.SetOperations(columns, eval)
	s.table = table
	s.expected = exp
	return s
}

func newSearchCases(columns map[string]string) []searchCase {
	// Returns search cases with expected results
	var ret []searchCase
	ret = append(ret, setCase(columns, "fox", "Name = Gray fox", "", make(map[string][]string)))
	ret = append(ret, setCase(columns, "canis", "Genus = Canis, Sex==male ", "", getCanisResults()))
	ret = append(ret, setCase(columns, "coyote", " Name == coyote", "", getCoyoteResults()))
	ret = append(ret, setCase(columns, "litter size", " litter_size>=5", "Life_history", getLitterResults()))
	return ret
}

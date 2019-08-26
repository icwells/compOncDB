// Contains expected out for databse search and update

package main

import (
	"github.com/icwells/compOncDB/src/codbutils"
)

func getDiagnosisUpdate() map[string][]string {
	// Returns map of account data
	ret := make(map[string][]string)
	ret["1"] = []string{"0", "0", "0", "-1"}
	ret["2"] = []string{"0", "0", "-1", "-1"}
	ret["3"] = []string{"1", "1", "-1", "-1"}
	ret["4"] = []string{"0", "0", "1", "-1"}
	ret["5"] = []string{"0", "0", "-1", "-1"}
	ret["6"] = []string{"0", "0", "1", "-1"}
	ret["7"] = []string{"0", "0", "-1", "-1"}
	ret["8"] = []string{"1", "1", "0", "-1"}
	ret["9"] = []string{"0", "0", "-1", "-1"}
	ret["10"] = []string{"1", "0", "0", "-1"}
	ret["11"] = []string{"1", "0", "-1", "1"}
	ret["12"] = []string{"0", "0", "-1", "-1"}
	ret["13"] = []string{"0", "0", "0", "-1"}
	ret["14"] = []string{"0", "0", "-1", "-1"}
	ret["15"] = []string{"0", "0", "-1", "-1"}
	ret["16"] = []string{"1", "0", "-1", "-1"}
	ret["17"] = []string{"0", "0", "1", "-1"}
	ret["18"] = []string{"0", "0", "-1", "-1"}
	ret["19"] = []string{"0", "0", "0", "-1"}
	return ret
}

func getPatientUpdate() map[string][]string {
	// Returns map of account data
	ret := make(map[string][]string)
	ret["1"] = []string{"male", "-1", "-1", "1", "1", "12-Dec", "Biopsy: NORMAL BLOOD SMEAR"}
	ret["2"] = []string{"NA", "-1", "-1", "1", "2", "13-Jan", "ERYTHROPHAGOCYTOSIS"}
	ret["3"] = []string{"male", "20", "-1", "1", "3", "1-Dec", "Lymphoma lymph nodes 2 year old male"}
	ret["4"] = []string{"NA", "-1", "-1", "1", "4", "1-Dec", "HIPOTOMAS TOXIC HIPOTOPATHY autopsy"}
	ret["5"] = []string{"NA", "-1", "-1", "1", "5", "1-Dec", "MICRONED HEPATITIS"}
	ret["6"] = []string{"NA", "-1", "-1", "1", "6", "1-Dec", "NA"}
	ret["7"] = []string{"NA", "-1", "-1", "1", "7", "1-Dec", "ASPERGILLOSIS"}
	ret["8"] = []string{"NA", "-1", "-1", "2", "8", "1-Dec", "Ovarian adenoma"}
	ret["9"] = []string{"NA", "-1", "-1", "2", "9", "1-Dec", "NA"}
	ret["10"] = []string{"female", "0", "-1", "2", "10", "NA", "skin biopsy:  squamous cell carcinoma; in situ"}
	ret["11"] = []string{"female", "156", "-1", "2", "11", "NA", "Uterus:  Endometrial carcinoma with metastatis"}
	ret["12"] = []string{"male", "60", "-1", "2", "12", "NA", "NA"}
	ret["13"] = []string{"female", "126", "1", "2", "13", "NA", "Spayed biopsy"}
	ret["14"] = []string{"female", "0", "-1", "2", "14", "NA", "NA"}
	ret["15"] = []string{"NA", "192", "-1", "2", "15", "NA", "16 month old"}
	ret["16"] = []string{"female", "132", "-1", "2", "16", "NA", "Malignant liver adenocarcinoma"}
	ret["17"] = []string{"male", "150", "1", "2", "17", "NA", "neutered"}
	ret["18"] = []string{"male", "-1", "0", "2", "18", "30463", "NA"}
	ret["19"] = []string{"male", "56", "0", "2", "19", "32688", "NA"}
	return ret
}

func getExpectedUpdates() map[string]map[string][]string {
	// Returns map of updated tables
	ret := make(map[string]map[string][]string)
	ret["Diagnosis"] = getDiagnosisUpdate()
	ret["Patient"] = getPatientUpdate()
	// Make sure tumor table hasn't changed
	ret["Tumor"] = getTumor()
	return ret
}

func getCleaned() map[string]map[string][]string {
	// Returns map of expected content after deletion and cleaning
	ret := make(map[string]map[string][]string)
	ret["Accounts"] = getAccounts()
	ret["Common"] = getCommon()
	ret["Denominators"] = getDenominators()
	ret["Diagnosis"] = getDiagnosisUpdate()
	ret["Life_history"] = getLifeHistory()
	ret["Patient"] = getPatientUpdate()
	ret["Source"] = getSource()
	ret["Taxonomy"] = getTaxonomy()
	//ret["Totals"] = getTotals()
	ret["Tumor"] = getTumor()
	for _, i := range []string{"Patient", "Diagnosis", "Tumor", "Source"} {
		delete(ret[i], "19")
	}
	return ret
}

//----------------------Search------------------------------------------------

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

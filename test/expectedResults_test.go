// Contains expected out for databse search and update

package main

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/dataframe"
)

func getDiagnosisUpdate() *dataframe.Dataframe {
	// Returns updated dataframe of account data
	ret := getDiagnosis()
	ret.UpdateCell("3", "Hyperplasia", "1")
	ret.UpdateCell("8", "Hyperplasia", "1")
	ret.UpdateCell("8", "Necropsy", "0")
	ret.UpdateCell("17", "Hyperplasia", "0")
	ret.UpdateCell("17", "Necropsy", "1")
	ret.UpdateCell("19", "Necropsy", "0")
	return ret
}

func getPatientUpdate() *dataframe.Dataframe {
	// Returns updated dataframe of account data
	ret := getPatient()
	ret.UpdateCell("3", "Age", "20")
	ret.UpdateCell("17", "Age", "150")
	ret.UpdateCell("19", "Age", "56")
	return ret
}

func getExpectedUpdates() map[string]*dataframe.Dataframe {
	// Returns map of updated tables
	ret := make(map[string]*dataframe.Dataframe)
	ret["Diagnosis"] = getDiagnosisUpdate()
	ret["Patient"] = getPatientUpdate()
	// Make sure tumor table hasn't changed
	ret["Tumor"] = getTumor()
	return ret
}

func getCleaned() map[string]*dataframe.Dataframe {
	// Returns map of expected content after deletion and cleaning
	ret := make(map[string]*dataframe.Dataframe)
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
		ret[i].DeleteRow("19")
	}
	return ret
}

func getExpectedRates() *dataframe.Dataframe {
	// Returns dataframe of account data
	var s [][]string
	coyote := []string{"1", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis latrans", "all", "1", "1", "1", "1", "1.00", "1", "1", "1.00", "1.00", "0", "0.00", "0.00", "24", "24", "1", "0", "1", "0", "0", "1"}
	wolf := []string{"2", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus", "all", "6", "6", "2", "2", "0.33", "2", "2", "0.33", "1.00", "0", "0.00", "0.00", "135", "144", "2", "3", "0", "2", "0", "1"}
	// Fox is in denominators table
	fox := []string{"3", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Urocyon", "Urocyon cinereoargenteus", "all", "1", "1", "0", "0", "0.00", "0", "0", "0.00", "0.00", "0", "0.00", "0.00", "0.00", "0", "0", "0", "0", "0", "0", "0"}
	s = append(s, codbutils.CancerRateHeader())
	s = append(s, wolf)
	s = append(s, coyote)
	s = append(s, fox)
	return setDF(-1, s)
}

//----------------------Search------------------------------------------------

func getCanisResults() *dataframe.Dataframe {
	// Returns map of results for male canis records
	s := [][]string{
		codbutils.RecordsHeader(),
		{"3", "male", "24", "0", "-1", "0", "1", "3", "Coyote", "1-Dec", "2011", "Lymphoma lymph nodes 2 year old male", "1", "0", "-1", "-1", "1", "1", "lymphoma", "lymph nodes", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis latrans", "NWZP", "0", "0", "0", "-1", "1"},
		{"12", "male", "60", "0", "-1", "0", "2", "12", "wolf", "NA", "1990", "NA", "0", "0", "-1", "-1", "0", "-1", "NA", "NA", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus", "NWZP", "1", "0", "0", "-1", "2"},
		{"17", "male", "144", "0", "1", "0", "2", "17", "GRAY WOLF", "NA", "2016", "neutered", "0", "0", "-1", "-1", "0", "-1", "NA", "NA", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus", "NWZP", "1", "0", "0", "-1", "2"},
	}
	return setDF(0, s)
}

func getCoyoteResults() *dataframe.Dataframe {
	// Returns map of coyote records
	s := [][]string{
		codbutils.RecordsHeader(),
		{"3", "male", "24", "0", "-1", "0", "1", "3", "Coyote", "1-Dec", "2011", "Lymphoma lymph nodes 2 year old male", "1", "0", "-1", "-1", "1", "1", "lymphoma", "lymph nodes", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis latrans", "NWZP", "0", "0", "0", "-1", "1"},
	}
	return setDF(0, s)
}

func getLitterResults() *dataframe.Dataframe {
	// Returns map of life history results
	s := [][]string{
		codbutils.LifeHistoryTestHeader(),
		{"1", "274", "274", "2.07100591715976", "1.9723865877712", "1.9723865877712", "5.72", "1.1", "365", "250", "1517", "13250", "0.0183", "261.6", "19.423"},
	}
	return setDF(0, s)
}

type searchCase struct {
	name     string
	eval     [][]codbutils.Evaluation
	expected *dataframe.Dataframe
	table    string
}

func setCase(columns map[string]string, name, eval, table string, exp *dataframe.Dataframe) searchCase {
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
	df, _ := dataframe.NewDataFrame(-1)
	ret = append(ret, setCase(columns, "fox", "Name = Gray fox", "", df))
	ret = append(ret, setCase(columns, "canis", "Genus = Canis, Sex==male , Infant = 0", "", getCanisResults()))
	ret = append(ret, setCase(columns, "coyote", " Name == coyote, Infant = 0", "", getCoyoteResults()))
	ret = append(ret, setCase(columns, "litter size", " litter_size>=5", "Life_history", getLitterResults()))
	return ret
}

// Map of expected output from comparative oncology database

package main

func getAccounts() map[string][]string {
	// Returns map of account data
	ret := make(map[string][]string)
	ret["1"] = []string{"X520", "XYZ"}
	ret["2"] = []string{"A16", "Kv Zoo"}
	return ret
}

func getCommon() map[string][]string {
	// Returns map of common names
	ret := make(map[string][]string)
	ret["1"] = []string{"Coyote", "Shake"}
	ret["2"] = []string{"Wolf", "Shake"}
	ret["3"] = []string{"Gray Fox", "Shake"}
	return ret
}

func getDenominators() map[string][]string {
	// Returns map of noncancer denominators
	ret := make(map[string][]string)
	ret["3"] = []string{"1"}
	return ret
}

func getDiagnosis() map[string][]string {
	// Returns map of account data
	ret := make(map[string][]string)
	ret["1"] = []string{"0", "0", "0", "-1"}
	ret["2"] = []string{"0", "0", "-1", "-1"}
	ret["3"] = []string{"1", "0", "-1", "-1"}
	ret["4"] = []string{"0", "0", "1", "-1"}
	ret["5"] = []string{"0", "0", "-1", "-1"}
	ret["6"] = []string{"0", "0", "1", "-1"}
	ret["7"] = []string{"0", "0", "-1", "-1"}
	ret["8"] = []string{"1", "0", "-1", "-1"}
	ret["9"] = []string{"0", "0", "-1", "-1"}
	ret["10"] = []string{"1", "0", "0", "-1"}
	ret["11"] = []string{"1", "0", "-1", "1"}
	ret["12"] = []string{"0", "0", "-1", "-1"}
	ret["13"] = []string{"0", "0", "0", "-1"}
	ret["14"] = []string{"0", "0", "-1", "-1"}
	ret["15"] = []string{"0", "0", "-1", "-1"}
	ret["16"] = []string{"1", "0", "-1", "-1"}
	ret["17"] = []string{"0", "0", "-1", "-1"}
	ret["18"] = []string{"0", "0", "-1", "-1"}
	ret["19"] = []string{"0", "0", "-1", "-1"}
	return ret
}

func getLifeHistory() map[string][]string {
	// Returns map of account data
	ret := make(map[string][]string)
	ret["1"] = []string{"274", "274", "2.07100591715976", "1.9723865877712", "1.9723865877712", "5.72", "1.1", "365", "250", "1517", "13250", "0.0183", "261.6", "19.423"}
	ret["2"] = []string{"669", "669", "2.03813280736358", "1.54503616042078", "1.54503616042078", "4.98", "0.8", "365", "450", "5250", "26625", "0.0177", "247.2", "33100"}
	ret["3"] = []string{"345", "365", "1.87376725838264", "1.80802103879027", "1.80802103879027", "3.71", "1.1", "365", "95", "519.7", "4750", "0.0127", "194.4", "-1"}
	return ret
}

func getPatient() map[string][]string {
	// Returns map of account data
	ret := make(map[string][]string)
	ret["1"] = []string{"male", "-1", "-1", "1", "1", "12-Dec", "Biopsy: NORMAL BLOOD SMEAR"}
	ret["2"] = []string{"NA", "-1", "-1", "1", "2", "13-Jan", "ERYTHROPHAGOCYTOSIS"}
	ret["3"] = []string{"male", "24", "-1", "1", "3", "1-Dec", "Lymphoma lymph nodes 2 year old male"}
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
	ret["17"] = []string{"male", "144", "1", "2", "17", "NA", "neutered"}
	ret["18"] = []string{"male", "-1", "0", "2", "18", "30463", "NA"}
	ret["19"] = []string{"male", "-1", "0", "2", "19", "32688", "NA"}
	return ret
}

func getSource() map[string][]string {
	// Returns map of account data
	ret := make(map[string][]string)
	ret["1"] = []string{"NWZP", "0", "0", "1"}
	ret["2"] = []string{"NWZP", "0", "0", "1"}
	ret["3"] = []string{"NWZP", "0", "0", "1"}
	ret["4"] = []string{"NWZP", "0", "0", "1"}
	ret["5"] = []string{"NWZP", "0", "0", "1"}
	ret["6"] = []string{"NWZP", "0", "0", "1"}
	ret["7"] = []string{"NWZP", "0", "0", "1"}
	ret["8"] = []string{"NWZP", "1", "0", "2"}
	ret["9"] = []string{"NWZP", "1", "0", "2"}
	ret["10"] = []string{"NWZP", "1", "0", "2"}
	ret["11"] = []string{"NWZP", "1", "0", "2"}
	ret["12"] = []string{"NWZP", "1", "0", "2"}
	ret["13"] = []string{"NWZP", "1", "0", "2"}
	ret["14"] = []string{"NWZP", "1", "0", "2"}
	ret["15"] = []string{"NWZP", "1", "0", "2"}
	ret["16"] = []string{"NWZP", "1", "0", "2"}
	ret["17"] = []string{"NWZP", "1", "0", "2"}
	ret["18"] = []string{"NWZP", "1", "0", "2"}
	ret["19"] = []string{"NWZP", "1", "0", "2"}
	return ret
}

func getTaxonomy() map[string][]string {
	// Returns map of account data
	ret := make(map[string][]string)
	ret["1"] = []string{"Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis latrans", `http://eol.org/api/hierarchy\_entries/1.0.xml?id=52440711`}
	ret["2"] = []string{"Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus", `http://eol.org/api/hierarchy\_entries/1.0.xml?id=52624675`}
	ret["3"] = []string{"Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Urocyon", "Urocyon cinereoargenteus", `http://eol.org/api/hierarchy\_entries/1.0.xml?id=52578011`}
	return ret
}

func getTotals() map[string][]string {
	// Returns map of account data
	ret := make(map[string][]string)
	ret["1"] = []string{"7", "24", "1", "1", "0", "1", "24", "1", "0"}
	ret["2"] = []string{"12", "135", "6", "2", "3", "2", "144", "0", "2"}
	ret["3"] = []string{"1", "-1", "1", "0", "0", "0", "-1", "0", "0"}
	return ret
}

func getTumor() map[string][]string {
	// Returns map of account data
	ret := make(map[string][]string)
	ret["1"] = []string{"0", "-1", "NA", "NA"}
	ret["2"] = []string{"0", "-1", "NA", "NA"}
	ret["3"] = []string{"0", "1", "lymphoma", "lymph nodes"}
	ret["4"] = []string{"0", "-1", "NA", "NA"}
	ret["5"] = []string{"0", "-1", "NA", "NA"}
	ret["6"] = []string{"0", "-1", "NA", "NA"}
	ret["7"] = []string{"0", "-1", "NA", "NA"}
	ret["8"] = []string{"0", "0", "adenoma", "ovary"}
	ret["9"] = []string{"0", "-1", "NA", "NA"}
	ret["10"] = []string{"0", "1", "carcinoma", "skin"}
	ret["11"] = []string{"0", "1", "carcinoma", "uterus"}
	ret["12"] = []string{"0", "-1", "NA", "NA"}
	ret["13"] = []string{"0", "-1", "NA", "NA"}
	ret["14"] = []string{"0", "-1", "NA", "NA"}
	ret["15"] = []string{"0", "-1", "NA", "NA"}
	ret["16"] = []string{"0", "1", "adenocarcinoma", "liver"}
	ret["17"] = []string{"0", "-1", "NA", "NA"}
	ret["18"] = []string{"0", "-1", "NA", "NA"}
	ret["19"] = []string{"0", "-1", "NA", "NA"}
	return ret
}

func getExpectedTables() map[string]map[string][]string {
	// Returns map of expected content after upload
	ret := make(map[string]map[string][]string)
	ret["Accounts"] = getAccounts()
	ret["Common"] = getCommon()
	ret["Denominators"] = getDenominators()
	ret["Diagnosis"] = getDiagnosis()
	ret["Life_history"] = getLifeHistory()
	ret["Patient"] = getPatient()
	ret["Source"] = getSource()
	ret["Taxonomy"] = getTaxonomy()
	ret["Totals"] = getTotals()
	ret["Tumor"] = getTumor()
	return ret
}

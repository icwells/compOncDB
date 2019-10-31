// Map of expected output from comparative oncology database

package main

import(
	"github.com/icwells/go-tools/dataframe"
)

func setDF(s [][]string) *dataframe.Dataframe {
	// Initializes dataframe with given data
	ret, _ := dataframe.NewDataFrame(0)
	ret.SetHeader(s[0])
	for _, i := range s[1:] {
		ret.AddRow(i)
	}
	return ret
}

func getAccounts() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		{"account_id", "Account", "submitter_name"},
		{"1", "X520", "XYZ"},
		{"2", "A16", "Kv Zoo"},
	}
	return setDF(s)
}

func getCommon() *dataframe.Dataframe {
	// Returns dataframe of common names
	s := [][]string{
		{"taxa_id", "Name", "Curator"},
		{"1", "Coyote", "Shake"},
		{"2", "Wolf", "Shake"},
		{"3", "Gray Fox", "Shake"},
	}
	return setDF(s)
}

func getDenominators() *dataframe.Dataframe {
	// Returns dataframe of noncancer denominators
	s := [][]string{
		{"taxa_id", "Noncancer"},
		{"3", "1"},
	}
	return setDF(s)
}

func getDiagnosis() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		{"ID", "Masspresent", "Hyperplasia", "Necropsy", "Metastasis"},
		{"1", "0", "0", "0", "-1"},
		{"2", "0", "0", "-1", "-1"},
		{"3", "1", "0", "-1", "-1"},
		{"4", "0", "0", "1", "-1"},
		{"5", "0", "0", "-1", "-1"},
		{"6", "0", "0", "1", "-1"},
		{"7", "0", "0", "-1", "-1"},
		{"8", "1", "0", "-1", "-1"},
		{"9", "0", "0", "-1", "-1"},
		{"10", "1", "0", "0", "-1"},
		{"11", "1", "0", "-1", "1"},
		{"12", "0", "0", "-1", "-1"},
		{"13", "0", "0", "0", "-1"},
		{"14", "0", "0", "-1", "-1"},
		{"15", "0", "0", "-1", "-1"},
		{"16", "1", "0", "-1", "-1"},
		{"17", "0", "0", "-1", "-1"},
		{"18", "0", "0", "-1", "-1"},
		{"19", "0", "0", "-1", "-1"},
	}
	return setDF(s)
}

func getLifeHistory() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
	{"taxa_id", "female_maturity", "male_maturity", "Gestation", "Weaning", "Infancy", "litter_size", "litters_year", "interbirth_interval", "birth_weight", "weaning_weight", "adult_weight", "growth_rate", "max_longevity", "metabolic_rate"},
		{"1", "274", "274", "2.07100591715976", "1.9723865877712", "1.9723865877712", "5.72", "1.1", "365", "250", "1517", "13250", "0.0183", "261.6", "19.423"},
		{"2", "669", "669", "2.03813280736358", "1.54503616042078", "1.54503616042078", "4.98", "0.8", "365", "450", "5250", "26625", "0.0177", "247.2", "33100"},
		{"3", "345", "365", "1.87376725838264", "1.80802103879027", "1.80802103879027", "3.71", "1.1", "365", "95", "519.7", "4750", "0.0127", "194.4", "-1"},
	}
	return setDF(s)
}

func getPatient() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		{"ID", "Sex", "Age", "Castrated", "taxa_id", "source_id", "source_name", "Date", "Comments"},
		{"1", "male", "-1", "-1", "1", "1", "Coyote", "12-Dec", "Biopsy: NORMAL BLOOD SMEAR"},
		{"2", "NA", "-1", "-1", "1", "2", "Coyote", "13-Jan", "ERYTHROPHAGOCYTOSIS"},
		{"3", "male", "24", "-1", "1", "3", "Coyote", "1-Dec", "Lymphoma lymph nodes 2 year old male"},
		{"4", "NA", "-1", "-1", "1", "4", "Coyote", "1-Dec", "HIPOTOMAS TOXIC HIPOTOPATHY autopsy"},
		{"5", "NA", "-1", "-1", "1", "5", "Coyote", "1-Dec", "MICRONED HEPATITIS"},
		{"6", "NA", "-1", "-1", "1", "6", "Coyote", "1-Dec", "NA"},
		{"7", "NA", "-1", "-1", "1", "7", "Coyote", "1-Dec", "ASPERGILLOSIS"},
		{"8", "NA", "-1", "-1", "2", "8", "wolf", "1-Dec", "Ovarian adenoma"},
		{"9", "NA", "-1", "-1", "2", "9", "wolf", "1-Dec", "NA"},
		{"10", "female", "0", "-1", "2", "10", "wolf", "NA", "skin biopsy:  squamous cell carcinoma; in situ"},
		{"11", "female", "156", "-1", "2", "11", "wolf", "NA", "Uterus:  Endometrial carcinoma with metastatis"},
		{"12", "male", "60", "-1", "2", "12", "wolf", "NA", "NA"},
		{"13", "female", "126", "1", "2", "13", "wolf", "NA", "Spayed biopsy"},
		{"14", "female", "0", "-1", "2", "14", "wolf", "NA", "NA"},
		{"15", "NA", "192", "-1", "2", "15", "GRAY WOLF", "NA", "16 month old"},
		{"16", "female", "132", "-1", "2", "16", "GRAY WOLF", "NA", "Malignant liver adenocarcinoma"},
		{"17", "male", "144", "1", "2", "17", "GRAY WOLF", "NA", "neutered"},
		{"18", "male", "-1", "0", "2", "18", "GRAY WOLF", "30463", "NA"},
		{"19", "male", "-1", "0", "2", "19", "GRAY WOLF", "32688", "NA"},
	}
	return setDF(s)
}

func getSource() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		{"ID", "service_name", "Zoo", "Institute", "account_id"},
		{"1", "NWZP", "0", "0", "1"},
		{"2", "NWZP", "0", "0", "1"},
		{"3", "NWZP", "0", "0", "1"},
		{"4", "NWZP", "0", "0", "1"},
		{"5", "NWZP", "0", "0", "1"},
		{"6", "NWZP", "0", "0", "1"},
		{"7", "NWZP", "0", "0", "1"},
		{"8", "NWZP", "1", "0", "2"},
		{"9", "NWZP", "1", "0", "2"},
		{"10", "NWZP", "1", "0", "2"},
		{"11", "NWZP", "1", "0", "2"},
		{"12", "NWZP", "1", "0", "2"},
		{"13", "NWZP", "1", "0", "2"},
		{"14", "NWZP", "1", "0", "2"},
		{"15", "NWZP", "1", "0", "2"},
		{"16", "NWZP", "1", "0", "2"},
		{"17", "NWZP", "1", "0", "2"},
		{"18", "NWZP", "1", "0", "2"},
		{"19", "NWZP", "1", "0", "2"},
	}
	return setDF(s)
}

func getTaxonomy() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		{"taxa_id", "Kingdom", "Phylum", "Class", "Orders", "Family", "Genus", "Species", "Source"},
		{"1", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis latrans", `http://eol.org/api/hierarchy\_entries/1.0.xml?id=52440711`},
		{"2", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus", `http://eol.org/api/hierarchy\_entries/1.0.xml?id=52624675`},
		{"3", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Urocyon", "Urocyon cinereoargenteus", `http://eol.org/api/hierarchy\_entries/1.0.xml?id=52578011`},
	}
	return setDF(s)
}

func getTotals() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		{"taxa_id", "Total", "Avgage", "Adult", "Male", "Female", "Cancer", "Cancerage", "Malecancer", "Femalecancer"},
		{"1", "7", "24", "1", "1", "0", "1", "24", "1", "0"},
		{"2", "12", "135", "6", "2", "3", "2", "144", "0", "2"},
		{"3", "1", "-1", "1", "0", "0", "0", "-1", "0", "0"},
	}
	return setDF(s)
}

func getTumor() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		{"ID", "primary_tumor", "Malignant", "Type", "Location"},
		{"1", "0", "-1", "NA", "NA"},
		{"2", "0", "-1", "NA", "NA"},
		{"3", "0", "1", "lymphoma", "lymph nodes"},
		{"4", "0", "-1", "NA", "NA"},
		{"5", "0", "-1", "NA", "NA"},
		{"6", "0", "-1", "NA", "NA"},
		{"7", "0", "-1", "NA", "NA"},
		{"8", "0", "0", "adenoma", "ovary"},
		{"9", "0", "-1", "NA", "NA"},
		{"10", "0", "1", "carcinoma", "skin"},
		{"11", "0", "1", "carcinoma", "uterus"},
		{"12", "0", "-1", "NA", "NA"},
		{"13", "0", "-1", "NA", "NA"},
		{"14", "0", "-1", "NA", "NA"},
		{"15", "0", "-1", "NA", "NA"},
		{"16", "0", "1", "adenocarcinoma", "liver"},
		{"17", "0", "-1", "NA", "NA"},
		{"18", "0", "-1", "NA", "NA"},
		{"19", "0", "-1", "NA", "NA"},
	}
	return setDF(s)
}

func getExpectedTables() map[string]*dataframe.Dataframe {
	// Returns dataframe of expected content after upload
	ret := make(map[string]*dataframe.Dataframe)
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

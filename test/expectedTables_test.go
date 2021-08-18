// Map of expected output from comparative oncology database

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/dataframe"
)

var H = codbutils.NewHeaders()

func setDF(col int, s [][]string) *dataframe.Dataframe {
	// Initializes dataframe with given data
	ret, _ := dataframe.NewDataFrame(col)
	ret.SetHeader(s[0])
	for _, i := range s[1:] {
		ret.AddRow(i)
	}
	return ret
}

func getAccounts() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		H.Accounts,
		{"1", "X520", "XYZ"},
		{"2", "NA", "Kv Zoo"},
	}
	return setDF(0, s)
}

func getCommon() *dataframe.Dataframe {
	// Returns dataframe of common names
	s := [][]string{
		H.Common,
		{"1", "Coyote", "Shake"},
		{"2", "Wolf", "Shake"},
		{"3", "Gray Fox", "Shake"},
	}
	return setDF(0, s)
}

func getDenominators() *dataframe.Dataframe {
	// Returns dataframe of noncancer denominators
	s := [][]string{
		H.Denominators,
		{"3", "1"},
	}
	return setDF(0, s)
}

func getDiagnosis() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		H.Diagnosis,
		{"1", "0", "0", "0", "-1"},
		{"2", "0", "0", "-1", "-1"},
		{"3", "1", "0", "-1", "-1"},
		{"4", "0", "0", "1", "-1"},
		{"5", "0", "0", "-1", "-1"},
		{"6", "0", "0", "1", "-1"},
		{"7", "0", "0", "-1", "-1"},
		{"8", "1", "0", "-1", "0"},
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
	return setDF(0, s)
}

func getLifeHistory() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		codbutils.LifeHistoryTestHeader(),
		{"1", "274", "274", "2.07100591715976", "1.9723865877712", "1.9723865877712", "5.72", "1.1", "365", "250", "1517", "13250", "0.0183", "261.6", "19.423"},
		{"2", "669", "669", "2.03813280736358", "1.54503616042078", "1.54503616042078", "4.98", "0.8", "365", "450", "5250", "26625", "0.0177", "247.2", "33100"},
		{"3", "345", "365", "1.87376725838264", "1.80802103879027", "1.80802103879027", "3.71", "1.1", "365", "95", "519.7", "4750", "0.0127", "194.4", "-1"},
	}
	return setDF(0, s)
}

func getPatient() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		H.Patient,
		{"1", "male", "-1", "-1", "-1", "0", "1", "1", "Coyote", "12-Dec", "2011", "Biopsy: NORMAL BLOOD SMEAR"},
		{"2", "NA", "-1", "-1", "-1", "0", "1", "2", "Coyote", "13-Jan", "2011", "ERYTHROPHAGOCYTOSIS"},
		{"3", "male", "24", "0", "-1", "0", "1", "3", "Coyote", "1-Dec", "2011", "Lymphoma lymph nodes 2 year old male"},
		{"4", "NA", "-1", "-1", "-1", "0", "1", "4", "Coyote", "1-Dec", "2011", "HIPOTOMAS TOXIC HIPOTOPATHY autopsy"},
		{"5", "NA", "-1", "-1", "-1", "0", "1", "5", "Coyote", "1-Dec", "2011", "MICRONED HEPATITIS"},
		{"6", "NA", "-1", "-1", "-1", "0", "1", "6", "Coyote", "1-Dec", "1999", "NA"},
		{"7", "NA", "-1", "-1", "-1", "0", "1", "7", "Coyote", "1-Dec", "1999", "ASPERGILLOSIS"},
		{"8", "NA", "-1", "-1", "-1", "0", "2", "8", "wolf", "1-Dec", "1999", "Ovarian adenoma"},
		{"9", "NA", "-1", "-1", "-1", "0", "2", "9", "wolf", "1-Dec", "1999", "NA"},
		{"10", "female", "0", "1", "-1", "0", "2", "10", "wolf", "NA", "1990", "skin biopsy:  squamous cell carcinoma; in situ"},
		{"11", "female", "156", "0", "-1", "0", "2", "11", "wolf", "NA", "1990", "Uterus:  Endometrial carcinoma with metastatis"},
		{"12", "male", "60", "0", "-1", "0", "2", "12", "wolf", "NA", "1990", "NA"},
		{"13", "female", "126", "0", "1", "0", "2", "13", "wolf", "NA", "2016", "Spayed biopsy"},
		{"14", "female", "0", "1", "-1", "0", "2", "14", "wolf", "NA", "2016", "NA"},
		{"15", "NA", "192", "0", "-1", "0", "2", "15", "GRAY WOLF", "NA", "2016", "16 month old"},
		{"16", "female", "132", "0", "-1", "0", "2", "16", "GRAY WOLF", "NA", "2016", "Malignant liver adenocarcinoma"},
		{"17", "male", "144", "0", "1", "0", "2", "17", "GRAY WOLF", "NA", "2016", "neutered"},
		{"18", "male", "-1", "-1", "0", "0", "2", "18", "GRAY WOLF", "30463", "2016", "NA"},
		{"19", "male", "-1", "-1", "0", "0", "2", "19", "GRAY WOLF", "32688", "2016", "NA"},
	}
	return setDF(0, s)
}

func getSource() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		H.Source,
		{"1", "NWZP", "0", "0", "0", "-1", "1"},
		{"2", "NWZP", "0", "0", "0", "-1", "1"},
		{"3", "NWZP", "0", "0", "0", "-1", "1"},
		{"4", "NWZP", "0", "0", "0", "-1", "1"},
		{"5", "NWZP", "0", "0", "0", "-1", "1"},
		{"6", "NWZP", "0", "0", "0", "-1", "1"},
		{"7", "NWZP", "0", "0", "0", "-1", "1"},
		{"8", "NWZP", "1", "0", "0", "-1", "2"},
		{"9", "NWZP", "1", "0", "0", "-1", "2"},
		{"10", "NWZP", "1", "0", "0", "-1", "2"},
		{"11", "NWZP", "1", "0", "0", "-1", "2"},
		{"12", "NWZP", "1", "0", "0", "-1", "2"},
		{"13", "NWZP", "1", "0", "0", "-1", "2"},
		{"14", "NWZP", "1", "0", "0", "-1", "2"},
		{"15", "NWZP", "1", "0", "0", "-1", "2"},
		{"16", "NWZP", "1", "0", "0", "-1", "2"},
		{"17", "NWZP", "1", "0", "0", "-1", "2"},
		{"18", "NWZP", "1", "0", "0", "-1", "2"},
		{"19", "NWZP", "1", "0", "0", "-1", "2"},
	}
	return setDF(0, s)
}

func getTaxonomy() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		H.Taxonomy,
		{"1", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis latrans", `http://eol.org/api/hierarchy\_entries/1.0.xml?id=52440711`},
		{"2", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus", `http://eol.org/api/hierarchy\_entries/1.0.xml?id=52624675`},
		{"3", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Urocyon", "Urocyon cinereoargenteus", `http://eol.org/api/hierarchy\_entries/1.0.xml?id=52578011`},
	}
	return setDF(0, s)
}

func getTumor() *dataframe.Dataframe {
	// Returns dataframe of account data
	s := [][]string{
		H.Tumor,
		{"1", "0", "-1", "NA", "NA"},
		{"2", "0", "-1", "NA", "NA"},
		{"3", "1", "1", "lymphoma", "lymph nodes"},
		{"4", "0", "-1", "NA", "NA"},
		{"5", "0", "-1", "NA", "NA"},
		{"6", "0", "-1", "NA", "NA"},
		{"7", "0", "-1", "NA", "NA"},
		{"8", "1", "0", "adenoma", "ovary"},
		{"9", "0", "-1", "NA", "NA"},
		{"10", "1", "1", "carcinoma", "skin"},
		{"11", "0", "1", "carcinoma", "uterus"},
		{"12", "0", "-1", "NA", "NA"},
		{"13", "0", "-1", "NA", "NA"},
		{"14", "0", "-1", "NA", "NA"},
		{"15", "0", "-1", "NA", "NA"},
		{"16", "1", "1", "adenocarcinoma", "liver"},
		{"17", "0", "-1", "NA", "NA"},
		{"18", "0", "-1", "NA", "NA"},
		{"19", "0", "-1", "NA", "NA"},
	}
	return setDF(0, s)
}

func getRecords(m map[string]*dataframe.Dataframe) *dataframe.Dataframe {
	// Returns records view
	ret, _ := dataframe.NewDataFrame(0)
	header := H.Patient
	header[2] = "age_months"
	header = append(header, H.Diagnosis[1:]...)
	header = append(header, H.Tumor[1:]...)
	header = append(header, H.Taxonomy[1:len(H.Taxonomy) - 1]...)
	header = append(header, H.Source[1:]...)
	header = append(header, codbutils.LifeHistoryTestHeader()[1:]...)
	ret.SetHeader(header)
	for idx := range m["Patient"].Index {
		tid, _ := m["Patient"].GetCell(idx, "taxa_id")
		row, _ := m["Patient"].GetRow(idx)
		row = append([]string{idx}, row...)
		diag, _ := m["Diagnosis"].GetRow(idx)
		tum, _ := m["Tumor"].GetRow(idx)
		taxa, _ := m["Taxonomy"].GetRow(tid)
		src, _ := m["Source"].GetRow(idx)
		lh, _ := m["Life_history"].GetRow(tid)
		row = append(row, diag...)
		row = append(row, tum...)
		row = append(row, taxa[:len(taxa) - 1]...)
		row = append(row, src...)
		row = append(row, lh...)
		if err := ret.AddRow(row); err != nil {
			panic(fmt.Sprintf("%s\n%s\n%v", header, row, err))
		}
	}
	return ret
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
	ret["Tumor"] = getTumor()
	ret["Records"] = getRecords(ret)
	return ret
}

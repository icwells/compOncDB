// Stores headers for output files

package codbutils

import (
	"strings"
)

type Headers struct {
	Accounts []string
	Common []string
	Denominators []string
	Diagnosis []string
	Life_history []string
	Patient []string
	Rates []string
	Source []string
	Taxonomy []string
	Tumor []string
}

func NewHeaders() *Headers {
	// Returns initialized struct
	h := new(Headers)
	h.Accounts = []string{"account_id", "Account", "submitter_name"}
	h.Common = []string{"taxa_id", "Name", "Curator"}
	h.Denominators = []string{"taxa_id", "Noncancer"}
	h.Diagnosis = []string{"ID", "Masspresent", "Hyperplasia", "Necropsy", "Metastasis"}
	h.Life_history = []string{"taxa_id", "female_maturity(months)", "male_maturity(months)", "Gestation(months)", "Weaning(months)", "Infancy(months)", "litter_size", "litters_year", 
"interbirth_interval", "birth_weight(g)", "weaning_weight(g)", "adult_weight(g)", "growth_rate(1/days)", "max_longevity(months)", "metabolic_rate(mLO2/hr)"}
	h.Patient = []string{"ID", "Sex", "Age", "Infant", "Castrated", "taxa_id", "source_id", "source_name", "Date", "Year", "Comments"}
	h.Rates = []string{"Location", "TotalRecords", "RecordsWithDenominators", "TotalNeoplasia", "NeoplasiaWithDenominators", "NeoplasiaPrevalence", 
"MalignancyKnown", "Malignant", "MalignancyPrevalence", "PropMalignant", "Benign", "BenignPrevalence", "PropBenign", 
"AverageAge(months)", "AvgAgeNeoplasia(months)", "Male", "Female", "MaleNeoplasia", "FemaleNeoplasia", "Necropsies", "#Sources"}
	h.Source = []string{"ID", "service_name", "Zoo", "Aza", "Institute", "Approved", "account_id"}
	h.Taxonomy = []string{"taxa_id", "Kingdom", "Phylum", "Class", "Orders", "Family", "Genus", "Species", "Source"}
	h.Tumor = []string{"ID", "primary_tumor", "Malignant", "Type", "Location"}
	return h
}

func LifeHistoryTestHeader() []string {
	// Removes units from life history header
	var ret []string
	h := NewHeaders()
	for _, i := range h.Life_history {
		idx := strings.Index(i, "(")
		if idx > 0 {
			ret = append(ret, i[:idx])
		} else {
			ret = append(ret, i)
		}
	}
	return ret
}

func LifeHistorySummaryHeader() []string {
	// Returns header for life history summary
	tail := []string{"%Complete", "Neoplasia", "Malignant", "Total"}
	h := NewHeaders()
	// Remove source column
	ret := h.Taxonomy[:len(h.Taxonomy)-1]
	ret = append(ret, h.Life_history[1:]...)
	return append(ret, tail...)
}

func RecordsHeader() []string {
	var ret []string
	h := NewHeaders()
	ret = append(ret, h.Patient...)
	ret = append(ret, h.Diagnosis[1:]...)
	ret = append(ret, h.Tumor[1:]...)
	ret = append(ret, h.Taxonomy[1:len(h.Taxonomy) - 1]...)
	return append(ret, h.Source[1:]...)
}

func CancerRateHeader() []string {
	// Returns header for cancer rate output
	h := NewHeaders()
	return append(h.Taxonomy[:len(h.Taxonomy) - 1], h.Rates...)
}

func ParseHeader(debug bool) string {
	// Returns header for parse output
	ret := "Sex,Age,Castrated,ID,Genus,Species,Name,Date,Year,Comments,"
	ret += "MassPresent,Hyperplasia,Necropsy,Metastasis,TumorType,Location,Primary,Malignant"
	ret += ",Service,Account,Submitter,Zoo,AZA,Institute"
	if debug == true {
		ret += ",Cancer,Code"
	}
	return ret + "\n"
}

// Stores headers for output files and service names by denominator status

package codbutils

import (
	"github.com/icwells/simpleset"
	"strings"
)

type Services struct {
	allrecords     *simpleset.Set
	nodenominators *simpleset.Set
}

func (s *Services) setServices(l []string) *simpleset.Set {
	// Returns set of input slice
	ret := simpleset.NewStringSet()
	for _, i := range l {
		ret.Add(i)
	}
	return ret
}

func NewServices() *Services {
	// Returns new string
	s := new(Services)
	s.allrecords = s.setServices([]string{"DLC", "LZ", "NWZP", "SDZ"})
	//s.denominators = s.setServices([]string{"ZEPS"})
	s.nodenominators = s.setServices([]string{"MSU", "SNZ", "WZ", "ZEPS"})
	return s
}

func (s *Services) AllRecords(name string) bool {
	// Returns true if name is in allrecords
	ret, _ := s.allrecords.InSet(name)
	return ret
}

func (s *Services) NoDenominators(name string) bool {
	// Returns true if name is in nodenominators
	ret, _ := s.nodenominators.InSet(name)
	return ret
}

//----------------------------------------------------------------------------

type Headers struct {
	AgeSex       []string
	Accounts     []string
	Common       []string
	Denominators []string
	Diagnosis    []string
	Life_history []string
	Location     string
	Malignancy   []string
	Neoplasia    []string
	Patient      []string
	RatesTail    []string
	Source       []string
	Taxonomy     []string
	Tumor        []string
}

func NewHeaders() *Headers {
	// Returns initialized struct
	h := new(Headers)
	h.Accounts = []string{"account_id", "submitter_name"}
	h.Common = []string{"taxa_id", "Name", "Curator"}
	h.Denominators = []string{"taxa_id", "Noncancer"}
	h.Diagnosis = []string{"ID", "Masspresent", "Hyperplasia", "Necropsy", "Metastasis"}
	h.Life_history = []string{"taxa_id", "female_maturity(months)", "male_maturity(months)", "Gestation(months)", "Weaning(months)", "Infancy(months)", "litter_size", "litters_year",
		"interbirth_interval", "birth_weight(g)", "weaning_weight(g)", "adult_weight(g)", "growth_rate(1/days)", "max_longevity(months)", "metabolic_rate(mLO2/hr)"}
	h.Patient = []string{"ID", "Sex", "Age", "Infant", "Castrated", "Wild", "taxa_id", "source_id", "source_name", "Date", "Year", "Comments"}
	h.Source = []string{"ID", "service_name", "Zoo", "Aza", "Institute", "Approved", "account_id"}
	h.Taxonomy = []string{"taxa_id", "Kingdom", "Phylum", "Class", "Orders", "Family", "Genus", "Species", "common_name", "Source"}
	h.Tumor = []string{"ID", "primary_tumor", "Malignant", "Type", "Tissue", "Location"}
	// Neoplasia Prevalence
	h.Location = "Location"
	h.Neoplasia = []string{"RecordsWithDenominators", "NeoplasiaDenominators", "NeoplasiaWithDenominators", "NeoplasiaPrevalence"}
	h.Malignancy = []string{"MalignancyKnown", "Malignant", "MalignancyPrevalence", "PropMalignant", "Benign", "BenignPrevalence", "PropBenign"}
	h.AgeSex = []string{"AverageAge(months)", "AvgAgeNeoplasia(months)", "Male", "MaleNeoplasia", "MaleMalignant", "Female", "FemaleNeoplasia", "FemaleMalignant"}
	h.RatesTail = []string{"RecordsFromAllSources", "NeoplasiaFromAllSources", "Necropsies", "#Sources", "NoTissueInfo"}
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
	ret[2] = "age_months"
	ret = append(ret, h.Diagnosis[1:]...)
	ret = append(ret, h.Tumor[1:]...)
	ret = append(ret, h.Taxonomy[1:len(h.Taxonomy)-1]...)
	return append(ret, h.Source[1:]...)
}

func CancerRateHeader(taxonomy, location, lifehistory bool) []string {
	// Returns header for cancer rate output
	var ret []string
	h := NewHeaders()
	if taxonomy {
		ret = h.Taxonomy[:len(h.Taxonomy)-1]
	} else {
		// Store taxa_id, species, and common name
		ret = []string{h.Taxonomy[0], h.Taxonomy[7], h.Taxonomy[8]}
	}
	if location {
		ret = append(ret, h.Location)
	}
	ret = append(ret, h.Neoplasia...)
	ret = append(ret, h.Malignancy...)
	ret = append(ret, h.AgeSex...)
	ret = append(ret, h.RatesTail...)
	if lifehistory {
		ret = append(ret, h.Life_history[1:]...)
	}
	return ret
}

func ParseHeader(debug bool) string {
	// Returns header for parse output
	ret := "Sex,Age,Castrated,ID,Genus,Species,Name,Date,Year,Comments,"
	ret += "MassPresent,Hyperplasia,Necropsy,Metastasis,TumorType,Tissue,Location,Primary,Malignant"
	ret += ",Service,Account,Submitter,Zoo,AZA,Institute"
	if debug == true {
		ret += ",Cancer,Code"
	}
	return ret + "\n"
}

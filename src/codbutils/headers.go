// Stores headers for output files

package codbutils

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
	h.Life_history = []string{"taxa_id", "female_maturity", "male_maturity", "Gestation", "Weaning", "Infancy", "litter_size", "litters_year", "interbirth_interval", "birth_weight", "weaning_weight", "adult_weight", "growth_rate", "max_longevity", "metabolic_rate"}
	h.Patient = []string{"ID", "Sex", "Age", "Castrated", "taxa_id", "source_id", "source_name", "Date", "Year", "Comments"}
	h.Rates = []string{"TotalRecords", "NeoplasiaRecords", "NeoplasiaRate", "Malignant", "MalignancyRate", "AverageAge(months)", "AvgAgeNeoplasia(months)",
		"Male", "Female", "MaleNeoplasia", "FemaleNeoplasia"}
	h.Source = []string{"ID", "service_name", "Zoo", "Aza", "Institute", "Approved", "account_id"}
	h.Taxonomy = []string{"taxa_id", "Kingdom", "Phylum", "Class", "Orders", "Family", "Genus", "Species", "Source"}
	h.Tumor = []string{"ID", "primary_tumor", "Malignant", "Type", "Location"}
	return h
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
	var ret []string
	h := NewHeaders()
	ret = append(ret, h.Taxonomy[:len(h.Taxonomy) - 1]...)
	return append(ret, h.Rates...)
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

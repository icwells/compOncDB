// Defines matches struct and data for test scripts

package diagnoses

type matches struct {
	Line       string
	Tissue     string
	Location   string
	Typ        string
	Infant     bool
	Age        string
	Sex        string
	Castrated  string
	Malignant  string
	Metastasis string
	Primary    string
	Necropsy   string
}

func NewMatches() []matches {
	// Initializes test matches
	line1 := "spinal neoplasia, biopsy; castration helps to resolve the situation since it is somewhat hormonal dependent, Female, 2.0 years old"
	line2 := "cause of death: single Malignant liver carcinoma; retarded growth has also been reported. 37 month old male"
	line3 := "metastatis lymphoma, infant, 30 days, not castrated, "
	line4 := "spayed female gray fox, "
	line5 := "male wolf with benign intestitial seminoma and abdominal lesion"
	return []matches{
		{line1, "Nervous", "spinal cord", "neoplasia", false, "24", "female", "1", "0", "-1", "0", "0"},
		{line2, "Gastrointestinal", "liver", "carcinoma", false, "37", "male", "-1", "1", "-1", "1", "1"},
		{line3, "Round Cell", "lymph nodes", "lymphoma", true, "1", "NA", "0", "1", "1", "0", "-1"},
		{line4, "NA", "NA", "NA", false, "-1", "female", "1", "-1", "-1", "0", "-1"},
		{line5, "Reproductive", "testis", "seminoma", false, "-1", "male", "-1", "0", "0", "0", "-1"},
	}
}

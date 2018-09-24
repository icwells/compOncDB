// This script will extract diagnosis information from a given input file

func (e *entries) checkAge(line []string) string {
	// Returns age/-1
	
}

func (e *entries) checkSex(line []string) string {
	// Returns male/female/NA
	ret := "NA"
	val := subsetLine(e.Sex, line)
	val = strings.ToUpper(val)
	if val == "M" || val == "Male" {
		ret = "male"
	} else if val == "F" || val == "FEMALE" {
		ret = "female"
	}
	return ret
}

// Contains methods for checking and re-formatting account names

package clusteraccounts

import (
	"github.com/icwells/go-tools/strarray"
	"strings"
)

func checkString(val string) string {
	// Returns NA if string is malformed
	v := strings.ToLower(val)
	if len(val) <= 0 {
		val = "NA"
	} else if v == "na" || v == "n/a" {
		val = "NA"
	}
	return val
}

func (a *Accounts) checkCaps(val string) string {
	// Recapitalizes abbreviations
	s := strings.Split(val, " ")
	for idx, i := range s {
		if len(i) <= 4 && !strings.Contains(i, ".") && !a.speller.Check(i) {
			// Recapitalize words under 5 charaters without periods and aren't present in dictionary
			s[idx] = strings.ToUpper(i)
		}
	}
	return strings.Join(s, " ")
}

func (a *Accounts) checkAmpersand(val string) string {
	// Replaces ampersand with "and" and corrects spacing
	t := "&"
	rep := " And "
	if strings.Contains(val, t) == true {
		for _, i := range []string{" & ", " &", "& "} {
			if strings.Contains(val, i) {
				t = i
				break
			}
		}
		val = strings.Replace(val, t, rep, 1)
	}
	return val
}

func (a *Accounts) checkPeriods(val string) string {
	// Fixes capitalization in terms with two letter abbreviations
	if strings.Contains(val, " ") {
		s := strings.Split(val, " ")
		for idx, i := range s {
			if strings.Count(i, ".") == 1 && len(i) == 2 {
				s[idx] = strings.ToUpper(i)
			} else if strings.Count(i, ".") == 2 && len(i) >= 3 && len(i) <= 5 {
				s[idx] = strings.ToUpper(i)
			}
		}
		val = strings.Join(s, " ")
	}
	return val
}

func (a *Accounts) checkAbbreviations(val string) string {
	//Store submitter/NA
	terms := map[string]string{"Animal Clinic": "A. C.", "Animal Hospital": "A. H.", "Veterinary Clinic": "V. C.", "University": "Univ",
		"Veterinary Hospital": "V. H.", "Veterinary Services": "V. S.", "Pet Vet": "P. V.", "International": "Intl ", "Animal": "Anim "}
	if strings.Contains(val, "?") || strings.Contains(strings.ToLower(val), "not used") {
		val = "NA"
	} else {
		// in records.go
		val = checkString(val)
		if val != "NA" {
			val = a.checkAmpersand(strarray.TitleCase(val))
			val = a.checkPeriods(val)
			// Resolve abbreviations
			for k, v := range terms {
				var alt string
				if strings.Contains(v, ".") == false {
					// Add trailing period
					alt = strings.Replace(v, " ", ".", 1)
				} else {
					// Remove space
					alt = strings.Replace(v, " ", "", 1)
				}
				if strings.Contains(val, v) == true {
					val = strings.Replace(val, v, k, 1)
					break
				} else if strings.Contains(val, alt) == true {
					val = strings.Replace(val, alt, k, 1)
					break
				}
			}
		}
	}
	return a.checkCaps(val)
}

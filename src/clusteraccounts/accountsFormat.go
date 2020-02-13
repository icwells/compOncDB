// Contains methods for checking and re-formatting account names

package clusteraccounts

import (
	"github.com/icwells/go-tools/strarray"
	"strconv"
	"strings"
	"unicode"
)

func trimPunc(val string) string {
	// Removes leading/trailing punctuation marks
	if len(val) > 1 {
		idx := len(val) - 1
		last := rune(val[idx])
		if unicode.IsPunct(last) && last != '.' {
			// Keep trailing periods
			val = val[:idx]
		}
		if unicode.IsPunct(rune(val[0])) {
			val = val[1:]
		}
	}
	return val
}

func checkString(val string) string {
	// Returns NA if string is malformed
	val = strings.TrimSpace(val)
	v := strings.ToLower(val)
	if len(val) <= 2 {
		val = "NA"
	} else if v == "na" || v == "n/a" {
		val = "NA"
	}
	return val
}

func nonsenseWord(v string) bool {
	// Returns flase if less than half of the characters in a value are letters
	count := 0
	if _, err := strconv.Atoi(v); err == nil {
		if len(v) == 4 || len(v) == 2 {
			// Keep possible years
			return false
		}
	} else {
		for _, i := range v {
			if unicode.IsLetter(i) {
				count++
			}
		}
		if float64(count) >= float64(len(v))/2.0 {
			return false
		}
	}
	return true
}

func (a *Accounts) checkSpelling(val string) string {
	// Checks spellings, adds correct words to corpus, removes nonsense, and recapitalizes abbreviations
	var s []string
	for _, i := range strings.Split(val, " ") {
		i = strings.TrimSpace(i)
		if a.speller.Check(i) {
			// Add to corpus of correctly spelled words
			a.corpus.Add(i)
			s = append(s, i)
		} else if !nonsenseWord(i) {
			s = append(s, i)
		}
	}
	ret := strings.Join(s, " ")
	if len(ret) == 0 {
		ret = "NA"
	} else if len(ret) <= 4 && !strings.Contains(ret, ".") && !a.speller.Check(val) {
		// Recapitalize words under 5 charaters without periods and aren't present in dictionary
		ret = strings.ToUpper(ret)
	}
	return ret
}

func (a *Accounts) checkSemiColon(val string) string {
	// Reverses terms written in "noun; adj..." form
	if strings.Count(val, ";") == 1 {
		s := strings.Split(val, ";")
		s[0] = trimPunc(strings.TrimSpace(s[0]))
		s[1] = strings.TrimSpace(s[1])
		if len(s[0]) >= 3 && len(s[1]) > 0 {
			// Only proceed if term might be a word/name and multiple terms are present
			if strings.ToLower(s[1]) != "zoo" && strings.ToLower(s[1]) != "company" {
				s[0], s[1] = s[1], s[0]
			}
		}
		val = s[0] + " " + s[1]
	}
	// Coorect zoo names in this format without semicolons
	s := strings.Split(val, " ")
	if strings.ToLower(s[0]) == "zoo" {
		s = append(s[1:], s[0])
		val = strings.Join(s, " ")
	}
	return val
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

func (a *Accounts) checkAbbreviations(ch chan string, val string) {
	//Store submitter/NA
	terms := map[string]string{"Animal Clinic": "A. C.", "Animal Hospital": "A. H.", "Veterinary Clinic": "V. C.", "University": "Univ",
		"Veterinary Hospital": "V. H.", "Veterinary Services": "V. S.", "Pet Vet": "P. V.", "International": "Intl ", "Animal": "Anim "}
	if strings.Contains(val, "?") || strings.Contains(strings.ToLower(val), "not used") {
		val = "NA"
	} else {
		// in records.go
		val = checkString(val)
		if val != "NA" {
			val = trimPunc(val)
			val = a.checkAmpersand(strarray.TitleCase(val))
			val = a.checkPeriods(val)
			val = a.checkSemiColon(val)
			// Resolve abbreviations
			for k, v := range terms {
				var alt string
				if !strings.Contains(v, ".") {
					// Add trailing period
					alt = strings.Replace(v, " ", ".", 1)
				} else {
					// Remove space
					alt = strings.Replace(v, " ", "", 1)
				}
				if strings.Contains(val, v) {
					val = strings.Replace(val, v, k, 1)
				} else if strings.Contains(val, alt) {
					val = strings.Replace(val, alt, k, 1)
				}
			}
		}
	}
	ch <- a.checkSpelling(val)
}

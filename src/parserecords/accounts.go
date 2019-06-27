// Uses spell checking and fuzzy matching to condense submitter names

package parserecords

import (
	"fmt"
	"github.com/icwells/go-tools/strarray"
	"github.com/trustmaster/go-aspell"
	"os"
	"strings"
)

type accounts struct {
	ratio                            float64
	speller                          aspell.Speller
	set                              strarray.Set
	submitters, pool, queries, terms map[string][]string
	scores                           map[string]int
}

func newAccounts() *accounts {
	// Returns pointer to initialized struct
	var a accounts
	var err error
	a.speller, err = aspell.NewSpeller(map[string]string{"lang": "en_US"})
	if err != nil {
		fmt.Printf("\n\t[Error] Cannot initialize speller. Exiting.\n%v", err)
		os.Exit(500)
	}
	a.set = strarray.NewSet()
	a.submitters = make(map[string][]string)
	a.queries = make(map[string][]string)
	a.terms = make(map[string][]string)
	a.ratio = 0.05
	return &a
}

func (a *accounts) getAccounts() map[string]string {
	// Returns map of original term: corrected term
	var count, total int
	ret := make(map[string]string)
	for key, val := range a.terms {
		count++
		for _, i := range val {
			for _, v := range a.queries[i] {
				if _, ex := ret[v]; ex == false {
					// Store original term and consensus term
					total++
					ret[v] = key
				}
			}
		}
	}
	fmt.Printf("\tFormatted %d terms from %d total account entries.\n", count, total)
	return ret
}

//----------------------------------------------------------------------------

func (a *accounts) checkCaps(val string) string {
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

func (a *accounts) checkAmpersand(val string) string {
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

func (a *accounts) checkPeriods(val string) string {
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

func (a *accounts) checkAbbreviations(val string) string {
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

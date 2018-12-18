// This script defiens a struct for matching diagnosis data with regular expressions

package main

import (
	"regexp"
	"strconv"
	"strings"
)

type matcher struct {
	location   map[string]*regexp.Regexp
	types      map[string]*regexp.Regexp
	malignancy map[string]string
	infant     *regexp.Regexp
	digit      *regexp.Regexp
	age        *regexp.Regexp
	sex        *regexp.Regexp
	castrated  *regexp.Regexp
	malignant  *regexp.Regexp
	metastasis *regexp.Regexp
	primary    *regexp.Regexp
	necropsy   *regexp.Regexp
}

func (m *matcher) setTypes() {
	// Sets type and location maps
	m.location = map[string]*regexp.Regexp{
		"abdomen":         regexp.MustCompile(`(?i)abdomen|abdom.*|omentum|diaphragm`),
		"bile duct":       regexp.MustCompile(`(?i)bile.*|biliary`),
		"bone":            regexp.MustCompile(`(?i)sacrum|bone.*`),
		"brain":           regexp.MustCompile(`(?i)brain`),
		"adrenal":         regexp.MustCompile(`(?i)adrenal`),
		"bladder":         regexp.MustCompile(`(?i)bladder`),
		"breast":          regexp.MustCompile(`(?i)(breast|mammary)`),
		"colon":           regexp.MustCompile(`(?i)colon|rectum`),
		"duodenum":        regexp.MustCompile(`(?i)duodenum`),
		"fat":             regexp.MustCompile(`(?i)fat|adipose.*`),
		"heart":           regexp.MustCompile(`(?i)heart|cardiac|atrial`),
		"kidney":          regexp.MustCompile(`(?i)kidney.*|ureter|renal`),
		"leukemia":        regexp.MustCompile(`(?i)leukemia`),
		"liver":           regexp.MustCompile(`(?i)hepa.*|liver.*|hep.*|billia.*`),
		"lung":            regexp.MustCompile(`(?i)lung.*|pulm.*|mediasti.*|bronchial|alveol.*`),
		"lymph nodes":     regexp.MustCompile(`(?i)lymph|lymph node`),
		"muscle":          regexp.MustCompile(`(?i)muscle|.*structure.*`),
		"nerve":           regexp.MustCompile(`(?i)nerve.*`),
		"other":           regexp.MustCompile(`(?i)gland|basal.*|islet|multifocal|neck|nasal|neuroendo.*`),
		"oral":            regexp.MustCompile(`(?i)oral|tongue|mouth|lip|palate|pharyn.*|laryn.*`),
		"ovary":           regexp.MustCompile(`(?i)ovar.*`),
		"pancreas":        regexp.MustCompile(`(?i)pancreas.*|islet`),
		"seminal vesicle": regexp.MustCompile(`(?i)seminal vesicle`),
		"skin":            regexp.MustCompile(`(?i)skin|eyelid|(sub)?cutan.*|derm.*`),
		"spinal cord":     regexp.MustCompile(`(?i)spinal|spine`),
		"spleen":          regexp.MustCompile(`(?i)spleen`),
		"testis":          regexp.MustCompile(`(?i)testi.*`),
		"thyroid":         regexp.MustCompile(`(?i)thyroid`),
		"uterus":          regexp.MustCompile(`(?i)uter.*`),
		"vulva":           regexp.MustCompile(`(?i)vulva|vagina`),
		"widespread":      regexp.MustCompile(`(?i)widespread|metastatic|(body as a whole)|multiple|disseminated`),
	}
	m.types = map[string]*regexp.Regexp{
		"adenocarcinoma": regexp.MustCompile(`(?i)adenocarcinoma`),
		"adenoma":        regexp.MustCompile(`(?i)adenoma`),
		"carcinoma":      regexp.MustCompile(`(?i)\scarcinoma|TCC`),
		"cyst":           regexp.MustCompile(`(?i)cyst`),
		"epulis":         regexp.MustCompile(`(?i)epuli.*`),
		"hyperplasia":    regexp.MustCompile(`(?i)(meta|dys|hyper)plas(ia|tic)`),
		"lymphoma":       regexp.MustCompile(`(?i)lymphoma|lymphosarcoma`),
		"leukemia":       regexp.MustCompile(`(?i)leukemia`),
		"meningioma":     regexp.MustCompile(`(?i)meningioma`),
		"papilloma":      regexp.MustCompile(`(?i)papilloma`),
		"neoplasia":      regexp.MustCompile(`(?i)neoplasia|neoplasm|tumor`),
		"polyp":          regexp.MustCompile(`(?i)polyp`),
		"sarcoma":        regexp.MustCompile(`(?i)\ssarcoma`),
	}
	m.malignancy = map[string]string{
		"adenocarcinoma": "Y",
		"adenoma":        "N",
		"carcinoma":      "Y",
		"cyst":           "N",
		"epulis":         "N",
		"hyperplasia":    "N",
		"lymphoma":       "Y",
		"leukemia":       "Y",
		"meningioma":     "N",
		"papilloma":      "N",
		"neoplasia":      "NA",
		"polyp":          "N",
		"sarcoma":        "Y",
	}
}

func newMatcher() matcher {
	// Compiles regular expressions
	var m matcher
	m.infant = regexp.MustCompile(`infant|(peri|neo)nat(e|al)|fet(us|al)`)
	m.digit = regexp.MustCompile(`[0-9]+`)
	m.age = regexp.MustCompile(`[0-9]+(-|\s)(day|week|month|year)s?(-|\s)?(old)?`)
	m.sex = regexp.MustCompile(`(fe)?male`)
	m.castrated = regexp.MustCompile(`(not )?(castrat(ed)?|neuter(ed)?|spay(ed)?)`)
	m.malignant = regexp.MustCompile(`(not )?(malignant|benign)`)
	m.metastasis = regexp.MustCompile(`(no )?(metastatis|metastatic|mets)`)
	m.primary = regexp.MustCompile(`primary|single|solitary|source`)
	m.necropsy = regexp.MustCompile(`(autopsy|necropsy|deceased|cause(-|\s)of(-|\s)death|dissection|euthan.*)|(biopsy)`)
	m.setTypes()
	return m
}

//----------------------------------------------------------------------------

func (m *matcher) getMatch(re *regexp.Regexp, line string) string {
	// Returns match/NA
	match := re.FindString(line)
	if len(match) == 0 {
		match = "NA"
	}
	return match
}

func (m *matcher) binaryMatch(re *regexp.Regexp, line, exp string) string {
	// Returns Y/N/NA
	ret := "NA"
	match := re.FindStringSubmatch(line)
	if len(match) >= 2 {
		if len(exp) >= 2 {
			if strings.Contains(match[1], "no") == true {
				// Negating phrase found
				if match[len(match)-1] == exp {
					ret = "Y"
				} else {
					ret = "N"
				}
			} else {
				if match[len(match)-1] == exp {
					// Negating expression found
					ret = "N"
				} else {
					ret = "Y"
				}
			}
		} else {
			if strings.Contains(match[1], "no") == true {
				// Negating phrase found
				ret = "N"
			} else {
				// No negation
				ret = "Y"
			}
		}
	}
	return ret
}

func (m *matcher) getMalignancy(line, t string) string {
	// Attmepts to determine if tumor is malignant or benign
	ret := "NA"
	if t != "NA" {
		res, ex := m.malignancy[t]
		if ex == true {
			ret = res
		}
	} else {
		ret = m.binaryMatch(m.malignant, line, "benign")
	}
	return ret
}

func (m *matcher) getType(line string, cancer bool) string {
	// Returns location from map
	ret := "NA"
	if cancer == true {
		for k, v := range m.types {
			match := m.getMatch(v, line)
			if match != "NA" {
				ret = k
				break
			}
		}
	}
	return ret
}

func (m *matcher) getLocation(line string, cancer bool) string {
	// Returns location from map
	ret := "NA"
	if cancer == true {
		for k, v := range m.location {
			match := m.getMatch(v, line)
			if match != "NA" {
				if match != "widespread" && match != "other" {
					ret = k
					break
				} else if ret == "NA" {
					// Attempt to find more descriptive match
					ret = k
				}
			}
		}
	}
	return ret
}

func (m *matcher) getCastrated(line string) string {
	// Returns castration status
	match := m.binaryMatch(m.castrated, line, "")
	if match == "NA" && strings.Contains(line, "intact") == true {
		match = "N"
	}
	return match
}

func (m *matcher) infantRecords(line string) bool {
	// Returns true if patient is an infant
	match := m.infant.Match([]byte(line))
	return match
}

func (m *matcher) getAge(line string) string {
	// Returns formatted age in months
	var ret string
	match := m.getMatch(m.age, line)
	if match != "NA" {
		age := m.digit.FindString(match)
		if strings.Contains(match, "month") == true {
			// Keep if already in months
			ret = age
		} else {
			// Convert to float, determine units, convert to months
			a, _ := strconv.ParseFloat(age, 64)
			if a > 0 {
				if strings.Contains(match, "year") == true {
					a = a * 12.0
				} else if strings.Contains(match, "week") == true {
					a = a / 4.0
				} else if strings.Contains(match, "day") == true {
					a = a / 30.0
				}
				// Convert back to string
				ret = strconv.FormatFloat(a, 'f', -1, 64)
			} else {
				ret = "0"
			}
		}
	} else {
		ret = match
	}
	return ret
}

// This script defiens a struct for matching diagnosis data with regular expressions

package main

import (
	"regexp"
	"strconv"
	"strings"
)

type matcher struct {
	location   map[string]*regexp.Regexp
	types      map[string]diagnosis
	infant     *regexp.Regexp
	digit      *regexp.Regexp
	age        *regexp.Regexp
	sex        *regexp.Regexp
	castrated  *regexp.Regexp
	malignant  *regexp.Regexp
	benign     *regexp.Regexp
	metastasis *regexp.Regexp
	primary    *regexp.Regexp
	necropsy   *regexp.Regexp
	biopsy	   *regexp.Regexp
}

func newMatcher() matcher {
	// Compiles regular expressions
	var m matcher
	m.infant = regexp.MustCompile(`(?i)infant|(peri|neo)nat(e|al)|fet(us|al)`)
	m.digit = regexp.MustCompile(`[0-9]+`)
	m.age = regexp.MustCompile(`(?i)[0-9]+(-|\s)(day|week|month|year)s?(-|\s)?(old)?`)
	m.sex = regexp.MustCompile(`(?i)(fe)?male`)
	m.castrated = regexp.MustCompile(`(?i)(not )?(castrat(ed)?|neuter(ed)?|spay(ed)?)`)
	m.malignant = regexp.MustCompile(`(?i)(not )?(malignan(t|cy)|invasive)`)
	m.benign = regexp.MustCompile(`(?i)(not )?(benign|encapsulated)`)
	m.metastasis = regexp.MustCompile(`(?i)(no )?(metastatis|metastatic|mets)`)
	m.primary = regexp.MustCompile(`(?i)primary|single|solitary|source`)
	m.necropsy = regexp.MustCompile(`(?i)(autopsy|necropsy|deceased|cause(-|\s)of(-|\s)death|dissection|euthan.*)`)
	m.biopsy = regexp.MustCompile(`(?i)biopsy`)
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

func (m *matcher) binaryMatch(re *regexp.Regexp, line string) string {
	// Returns Y/N/NA
	ret := "NA"
	match := re.FindStringSubmatch(line)
	if len(match) >= 2 {
		if strings.Contains(match[1], "no ") == true || strings.Contains(match[1], "not ") == true {
			// Negating phrase found
			ret = "N"
		} else {
			// No negation
			ret = "Y"
		}
	}
	return ret
}

func (m *matcher) getNecropsy(line string) string {
	// Returns Y/N/NA; also searches for negating expression
	ret := m.binaryMatch(m.necropsy, line)
	if ret == "NA" {
		// Search for biopsy
		inverse := m.getMatch(m.biopsy, line)
		if inverse != "NA" {
			ret = "N"
		}
	}
	return ret
}

func (m *matcher) getMalignancy(line string) string {
	// Returns Y/N for malignant/benign
	ret := m.binaryMatch(m.malignant, line)
	if ret == "NA" {
		ret = m.binaryMatch(m.benign, line)
		// Reverse benign result
		if ret == "Y" {
			ret = "N"
		} else if ret == "N" {
			ret = "Y"
		}
	}
	return ret
}

func (m *matcher) getType(line string, cancer bool) (string, string) {
	// Returns location from map
	typ := "NA"
	mal := "NA"
	if cancer == true {
		for k, v := range m.types {
			match := m.getMatch(v.expression, line)
			if match != "NA" {
				typ = k
				mal = v.malignant
				if strings.Contains(line, " "+typ) == true {
					// Only break for whole-word matches
					break
				}
			}
		}
		if mal == "NA" {
			mal = m.getMalignancy(line)
		}
	}
	return typ, mal
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
	match := m.binaryMatch(m.castrated, line)
	if match == "NA" && strings.Contains(line, "intact") == true {
		match = "N"
	}
	return match
}

func (m *matcher) infantRecords(line string) bool {
	// Returns true if patient is an infant
	match := m.infant.MatchString(line)
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

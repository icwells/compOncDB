// This script defines a struct for matching diagnosis data with regular expressions

package parserecords

import (
	"regexp"
	"strconv"
	"strings"
)

type matcher struct {
	location   map[string]*regexp.Regexp
	types      map[string]map[string]diagnosis
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
	biopsy     *regexp.Regexp
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
	m.necropsy = regexp.MustCompile(`(?i)(autopsy|necropsy|deceased|cause(-|\s)of(-|\s)death|dissect*|euthan.*)`)
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
	// Returns 1/0/-1
	ret := "-1"
	match := re.FindStringSubmatch(line)
	if len(match) >= 2 {
		if strings.Contains(match[1], "no ") == true || strings.Contains(match[1], "not ") == true {
			// Negating phrase found
			ret = "0"
		} else {
			// No negation
			ret = "1"
		}
	}
	return ret
}

func (m *matcher) getNecropsy(line string) string {
	// Returns 1/0/-1; also searches for negating expression
	ret := m.binaryMatch(m.necropsy, line)
	if ret == "-1" {
		// Search for biopsy
		inverse := m.getMatch(m.biopsy, line)
		if inverse != "NA" {
			ret = "0"
		}
	}
	return ret
}

func (m *matcher) getMalignancy(line string) string {
	// Returns 1/0 for malignant/benign
	ret := m.binaryMatch(m.malignant, line)
	if ret == "-1" {
		ret = m.binaryMatch(m.benign, line)
		// Reverse benign result
		if ret == "1" {
			ret = "0"
		} else if ret == "0" {
			ret = "1"
		}
	}
	return ret
}

func (m *matcher) getCastrated(line string) string {
	// Returns castration status
	match := m.binaryMatch(m.castrated, line)
	if match == "-1" && strings.Contains(line, "intact") == true {
		match = "0"
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
	ret := "-1"
	match := m.getMatch(m.age, line)
	if match != "NA" {
		age := m.digit.FindString(match)
		if strings.Contains(match, "month") == true {
			// Keep if already in months
			ret = age
		} else {
			// Convert to float, determine units, convert to months
			a, err := strconv.ParseFloat(age, 64)
			if err == nil && a > 0 {
				if strings.Contains(match, "year") == true {
					a = a * 12.0
				} else if strings.Contains(match, "week") == true {
					a = a / 4.0
				} else if strings.Contains(match, "day") == true {
					a = a / 30.0
				}
				// Convert back to string
				ret = strconv.FormatFloat(a, 'f', -1, 64)
			}
		}
	}
	return ret
}

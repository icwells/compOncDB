// This script defiens a struct for matching diagnosis data with regular expressions

package main

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"regexp"
	"sort"
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

func (m *matcher) getType(line string) (string, string) {
	// Returns location from map
	typ := "NA"
	mal := "-1"
	for k, v := range m.types {
		match := m.getMatch(v.expression, line)
		if match != "NA" {
			typ = k
			mal = v.malignant
			if strings.Contains(line, " "+typ) == true {
				// Only break for whole-word matches (i.e. "adenocarcinoma", but not "carcinoma")
				break
			}
		}
	}
	if mal == "-1" {
		mal = m.getMalignancy(line)
	}
	return typ, mal
}

func rankLocations(ranks fuzzy.Ranks) string {
	// Gets most descriptive match with lowest Levenshtein distance
	ret := "NA"
	sort.Sort(ranks)
	if len(ranks) == 1 {
		// Return lone match
		ret = ranks[0].Source
	} else if ranks[0].Source != "widespread" && ranks[0].Source != "other" {
		ret = ranks[0].Source
	} else {
		// Attempt to find more descriptive match
		for _, i := range ranks {
			ret = i.Source
			if i.Source != "widespread" && i.Source != "other" {
				break
			}
		}
	}
	return ret
}

func (m *matcher) getLocation(line string, cancer bool) string {
	// Combines regexp and fuzzy searching to determine best match for location
	ret := "NA"
	var ranks fuzzy.Ranks
	for k, v := range m.location {
		match := m.getMatch(v, line)
		if match != "NA" {
			// Call RankFind to get get Ranks struct returned
			rank := fuzzy.RankFindFold(k, []string{match})
			if len(rank) >= 1 {
				ranks = append(ranks, rank[0])
			}
		}
	}
	if len(ranks) >= 1 {
		ret = rankLocations(ranks)
	}
	return ret
}

func (m *matcher) getTumor(line string, cancer bool) (string, string, string) {
	// Returns type, location, and malignancy
	typ := "NA"
	loc := "NA"
	mal = "-1"
	if cancer == true {
		typ, mal := m.getType(line)
		
	}
	return typ, loc, mal
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

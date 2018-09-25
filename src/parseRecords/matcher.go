// This script defiens a struct for matching diagnosis data with regular expressions

package main

import (
	"regexp"
	"strings"
)

type matcher struct {
	location	map[string]*Regexp
	types		map[string]*Regexp
	infant		*Regexp
	digit		*Regexp
	age			*Regexp
	sex			*Regexp
	castrated	*Regexp
	malignant	*Regexp
	metastasis	*Regexp
	primary		*Regexp
	necropsy	*Regexp
}

func (m *matcher) setTypes(infile string) {
	// Sets type and location maps
	var d string
	first := true
	m.location = make(map[string]*Regexp)
	m.types = make(map[string]*Regexp)
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Input() {
		line := string(input.Text())
		if first == false {
			line = strings.TrimSpace(line)
			if strings.Count(line, d) >= 2 {
				s := strings.Split(line, d)
				// String: regexp
				if s[0] == "Location" {
					m.location[s[1]] = regexp.Compile(s[2])
				} else if s[0] == "Type" {
					m.types[s[1]] = regexp.Compile(s[2])
				}
			}
		} else {
			d = getDelim(line)
			first = false
		}
	}
}

func newMatcher(dict string) matcher {
	// Compiles regular expressions
	var m matcher
	m.infant = regexp.Compile(`infant|(peri|neo)nat(e|al)|fet(us|al)`)
	m.digit = regexp.Compile(`[0-9]+`)
	m.age = regexp.Compile(`[0-9]+(-|\s)(day|week|month|year)s?(-|\s)(old)?`)
	m.sex = regexp.Compile(`(fe)?male`)
	m.castrated = regexp.Compile(`(not )?(castrat(ed)?|neuter(ed)?|spay(ed)?)`)
	m.malignant = regexp.Compile(`(not )?(malignant|benign)`)
	m.metastasis = regexp.Compile(`(no )?(metastatis|mets)`)
	m.primary = regexp.Compile(`primary|single|solitary|source`)
	m.necropsy = regexp.Compile(`(necropsy|decesed|cause of death)|(biopsy)`)
	m.setTypes(dict)
	return m
}

func (m *matcher) getMatch(re *Regexp, line string) string {
	// Returns match/NA
	match := re.FindString(line)
	if len(match) == 0 {
		match = "NA"
	}
	return match
}

func (m *matcher) binaryMatch(re *Regexp, line, exp string) string {
	// Returns Y/N/NA
	ret := "NA"
	match := re.FindStringSubmatch(line, 1)
	if match != nil {
		if len(exp) >= 1 {
			if match[1] != nil {
				if match[0] != nil {
					// Negation found
					if match[1] == exp {
						ret = "Y"
					} else {
						ret = "N"
					}
				} else {
					if match[1] == exp {
						ret = "N"
					} else {
						ret = "Y"
					}
				}
			}
		} else {
			if match[0] != nil {
				// Negating phrase found
				ret = "N"
			} else if match[1] != nil {
				// No negation
				ret = "Y"
			}
		}
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
				ret = k
				break
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
		age := m.digit.FindString(line)
		if strings.Contains(match, "month") == true {
			// Keep if already in months
			ret = age
		} else {
			// Convert to float, determine units, convert to months
			a := strconv.ParseFloat(age, 64)
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
	return ret
}

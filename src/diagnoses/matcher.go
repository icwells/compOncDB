// This script defines a struct for matching diagnosis data with regular expressions

package diagnoses

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type Matcher struct {
	age        *regexp.Regexp
	benign     *regexp.Regexp
	biopsy     *regexp.Regexp
	castrated  *regexp.Regexp
	digit      *regexp.Regexp
	infant     *regexp.Regexp
	infile     string
	location   map[string]*regexp.Regexp
	malignant  *regexp.Regexp
	Metastasis *regexp.Regexp
	necropsy   *regexp.Regexp
	Primary    *regexp.Regexp
	Sex        *regexp.Regexp
	tissues    map[string]string
	types      map[string]*tumortype
}

func NewMatcher(logger *log.Logger) Matcher {
	// Compiles regular expressions
	var m Matcher
	digit := `([0-9]*[.])?[0-9]+`
	m.infile = path.Join(codbutils.Getutils(), "diagnoses.csv")
	m.infant = regexp.MustCompile(`(?i)infant|(peri|neo)nat(e|al)|fet(us|al)`)
	m.digit = regexp.MustCompile(digit)
	m.age = regexp.MustCompile(digit + `(-|\s)(day|week|month|year)s?(-|\s)?(old|of age)?`)
	m.Sex = regexp.MustCompile(`(?i)(fe)?male`)
	m.castrated = regexp.MustCompile(`(?i)(not )?(castrat(ed)?|neuter(ed)?|spay(ed)?)`)
	m.malignant = regexp.MustCompile(`(?i)(not )?(malignan(t|cy)|invasive)`)
	m.benign = regexp.MustCompile(`(?i)(not )?(benign|encapsulated)`)
	m.Metastasis = regexp.MustCompile(`(?i)(no )?(metastati(s|c)|mets|disseminated|distant|stage(\s)?(three|3|four|4))`)
	m.Primary = regexp.MustCompile(`(?i)primary|single|solitary|source`)
	m.necropsy = regexp.MustCompile(`(?i)(autopsy|necropsy|deceased|cause(-|\s)of(-|\s)death|dissect*|euthan.*)`)
	m.biopsy = regexp.MustCompile(`(?i)biopsy`)
	m.setTypes(logger)
	m.addLocations()
	return m
}

//----------------------------------------------------------------------------

func (m *Matcher) addLocations() {
	// Adds additional locations to locations map
	for _, i := range []string{"head", "neck", "leg", "arm", "wing"} {
		// Defer to entry from file if present
		if _, ex := m.location[i]; !ex {
			m.location[i] = m.formatExpression(i)
		}
	}
}

func (m *Matcher) GetMatch(re *regexp.Regexp, line string) string {
	// Returns match/NA
	match := re.FindString(line)
	if len(match) == 0 {
		match = "NA"
	}
	return match
}

func (m *Matcher) BinaryMatch(re *regexp.Regexp, line string) string {
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

func (m *Matcher) GetNecropsy(line string) string {
	// Returns 1/0/-1; also searches for negating expression
	ret := m.BinaryMatch(m.necropsy, line)
	if ret == "-1" {
		// Search for biopsy
		inverse := m.GetMatch(m.biopsy, line)
		if inverse != "NA" {
			ret = "0"
		}
	}
	return ret
}

func (m *Matcher) GetMalignancy(line string) string {
	// Returns 1/0 for malignant/benign
	ret := m.BinaryMatch(m.malignant, line)
	if ret == "-1" {
		ret = m.BinaryMatch(m.benign, line)
		// Reverse benign result
		if ret == "1" {
			ret = "0"
		} else if ret == "0" {
			ret = "1"
		}
	}
	return ret
}

func (m *Matcher) GetCastrated(line string) string {
	// Returns castration status
	match := m.BinaryMatch(m.castrated, line)
	if match == "-1" && strings.Contains(line, "intact") == true {
		match = "0"
	}
	return match
}

func (m *Matcher) InfantRecords(line string) bool {
	// Returns true if patient is an infant
	match := m.infant.MatchString(line)
	return match
}

func (m *Matcher) GetAge(line string) string {
	// Returns formatted age in months
	ret := "-1"
	match := m.GetMatch(m.age, line)
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

// Regular expression dictionaries for the matcher struct

package diagnoses

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/strarray"
	"os"
	"path"
	"regexp"
	"strings"
)

type tumortype struct {
	expression *regexp.Regexp
	benign     bool
	malignant  bool
	locations  strarray.Set
}

func newTumorType(exp *regexp.Regexp) *tumortype {
	// Initializes empty struct
	var t tumortype
	t.expression = exp
	t.locations = strarray.NewSet()
	return &t
}

func (t *tumortype) isBenign() {
	// Sets benign to true
	t.benign = true
}

func (t *tumortype) isMalignant() {
	// Sets malignant to true
	t.malignant = true
}

func (t *tumortype) addLocation(loc string) {
	// Adds loc to location set
	t.locations.Add(loc)
}

//----------------------------------------------------------------------------

func (m *Matcher) formatExpression(e string) *regexp.Regexp {
	// Formats and compiles regular expression
	if strings.Contains(e, " cell") {
		e = strings.Replace(e, " cell", "( cell)?", 1)
	}
	e = strings.Replace(e, " ", `\s`, -1)
	e = fmt.Sprintf("(?i)%s", e)
	return regexp.MustCompile(e)
}

func (m *Matcher) checkType(loc, name, exp string) {
	// Makes new entry in types map if needed and adds location to type
	if exp == "" {
		// Set expression to type name
		exp = name
	}
	if _, ex := m.types[name]; !ex {
		m.types[name] = newTumorType(m.formatExpression(exp))
	}
	m.types[name].addLocation(loc)
}

func (m *Matcher) setTumorType(df *dataframe.Dataframe, loc string, idx int) {
	// Stores relevant information for tumor dignosis
	b, err := df.GetCell(idx, "Benign")
	if err == nil && b != "" {
		exp, _ := df.GetCell(idx, "BenignExpression")
		m.checkType(loc, b, exp)
		m.types[b].isBenign()
	}
	mal, er := df.GetCell(idx, "Malignant")
	if er == nil && mal != "" {
		exp, _ := df.GetCell(idx, "MalignantExpression")
		m.checkType(loc, mal, exp)
		m.types[mal].isMalignant()
	}
}

func (m *Matcher) setLocation(l, exp string) string {
	// Adds new location to map
	l = strings.ToLower(l)
	if strings.Count(l, " ") >= 1 {
		// Remove trailing s from second word
		if l[len(l)-1] == 's' {
			l = l[:len(l)-1]
		}
	}
	if exp == "" {
		exp = l
	}
	m.location[l] = m.formatExpression(exp)
	return l
}

func (m *Matcher) setTypes() {
	// Sets type and location maps from file
	var loc string
	m.location = make(map[string]*regexp.Regexp)
	m.types = make(map[string]*tumortype)
	infile := path.Join(codbutils.Getutils(), "diagnoses.csv")
	df, err := dataframe.DataFrameFromFile(infile, -1)
	if err != nil {
		fmt.Printf("\n\t[Error] Reading diagnoses file: %v\n", err)
		os.Exit(1)
	}
	for idx := range df.Rows {
		l, err := df.GetCell(idx, "Location")
		if err == nil && l != "" {
			exp, err := df.GetCell(idx, "LocationExpression")
			if err != nil {
				exp = ""
			}
			loc = m.setLocation(l, exp)
		}
		m.setTumorType(df, loc, idx)
	}
}

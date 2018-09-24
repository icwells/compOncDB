// This script will extract diagnosis information from a given input file

package main

import (
	"fmt"
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

func (m *matcher) setTypes() {
	// Sets type and location maps
	m.location = make(map[string]*Regexp)
	m.types = make(map[string]*Regexp)
}

func newMatcher() matcher {
	// Compiles regular expressions
	var m matcher
	m.infant = regexp.Compile(`infant|(peri|neo)nat(e|al)|fet(us|al)`)
	m.digit = regexp.Compile(`[0-9]+`)
	m.age = regexp.Compile(`[0-9]+(-|\s)(day|week|month|year)s?(-|\s)(old)?`)
	m.sex = regexp.Compile(`(fe)?male`)
	m.castrated = regexp.Compile(`(not )?(castrat(ed)?|neuter(ed)?|spay(ed)?)`)
	m.malignant = regexp.Compile(`(not )?(malignant|benign)`)
	m.metastasis = regexp.Compile(`no )?(metastatis|mets)`)
	m.primary = regexp.Compile(`primary|single|solitary|source`)
	m.necropsy = regexp.Compile(`(necropsy|decesed|cause of death)|(biopsy)`)
	m.setTypes()
}

func (e *entries) getMatcher() {
	

}

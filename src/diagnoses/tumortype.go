// Defines tumorype struct and methods

package diagnoses

import (
	"github.com/icwells/go-tools/strarray"
	"regexp"
)

type tumortype struct {
	expression *regexp.Regexp
	benign      bool
	malignant   bool
	locations   *strarray.Set
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

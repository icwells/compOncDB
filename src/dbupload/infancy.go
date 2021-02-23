// Declares struct for determining whether a record is an infant or an adult

package dbupload

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"strconv"
	"strings"
)

type Infancy struct {
	ages  map[string]float64
	terms []string
}

func NewInfancy(db *dbIO.DBIO) *Infancy {
	i := new(Infancy)
	i.ages = codbutils.GetMinAges(db, []string{})
	i.terms = []string{"adult", "mature", "infant", "fetus", "juvenile", "immature", "adolescent", "hatchling", "subadult", "neonate", "polyp", "placenta", "newborn", "offspring", "fledgling", "snakelet", "brood", "fry", "fingerling"}
	return i
}

func (i *Infancy) checkAges(id, age string) string {
	// Compairs given age with recorded age of infancy
	ret := "-1"
	if min, ex := i.ages[id]; ex {
		if a, err := strconv.ParseFloat(age, 64); err == nil {
			if a >= 0 {
				if a <= min {
					ret = "1"
				} else {
					ret = "0"
				}
			}
		}
	}
	return ret
}

func (i *Infancy) checkComments(comments string) string {
	// Checks comments for age-related key words
	ret := "-1"
	comments = strings.ToLower(comments)
	for idx, i := range i.terms {
		if strings.Contains(comments, i) {
			if idx <= 1 {
				ret = "0"
			} else {
				ret = "1"
			}
			break
		}
	}
	return ret
}

func (i *Infancy) SetInfant(id, age, comments string) string {
	// Determines if records are infant records
	ret := i.checkAges(id, age)
	if ret == "-1" {
		ret = i.checkComments(comments)
	}
	return ret
}

// Declares struct for determining whether a record is an infant or an adult

package dbupload

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"strconv"
	"strings"
)

type Infancy struct {
	adult  []string
	ages   map[string]float64
	infant []string
}

func NewInfancy(db *dbIO.DBIO) *Infancy {
	i := new(Infancy)
	i.adult = []string{"adult", "mature"}
	i.ages = codbutils.GetMinAges(db, []string{})
	i.infant = []string{"infant", "fetus", "juvenile", "immature", "adolescent", "hatchling", "subadult", "neonate", "polyp"}
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

func (i *Infancy) checkComments(list []string, comments string) string {
	// Checks comments for age-related key words
	ret := "-1"
	comments = strings.ToLower(comments)
	for idx, i := range i.adult {
		if strings.Contains(comments, i) {
			if idx <= 3 {
				ret = "1"
			} else {
				ret = "0"
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
		ret = i.checkComments(i.adult, comments)
		if ret == "-1" {
			ret = i.checkComments(i.infant, comments)
		}
	}
	return ret
}

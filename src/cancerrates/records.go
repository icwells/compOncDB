// Defines record struct and methods

package cancerrates

import (
	"github.com/icwells/simpleset"
	"strconv"
	"strings"
)

func avgAge(n float64, d int) float64 {
	// Returns n/d
	var ret float64
	if n > 0.0 && d > 0 {
		ret = n / float64(d)
	} else {
		ret = -1.0
	}
	return ret
}

type record struct {
	age          float64
	allcancer    int
	benign       int
	bentotal     int
	cancer       int
	cancerage    float64
	female       int
	femalecancer int
	grandtotal   int
	male         int
	malecancer   int
	malignant    int
	maltotal     int
	necropsy     int
	sources      *simpleset.Set
	total        int
}

func newRecord() *record {
	// Initializes new record struct
	r := new(record)
	r.sources = simpleset.NewStringSet()
	return r
}

func (r *record) CalculateAvgages() {
	// Sets average age and average cancer age
	r.age = avgAge(r.age, r.total)
	r.cancerage = avgAge(r.cancerage, r.cancer)
}

func (r *record) formatRate(n, d int) string {
	// Divides n by d and returns formatted string
	var v float64
	if d != 0 {
		v = float64(n) / float64(d)
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func (r *record) setsources() string {
	// Returns number of unique sources
	if r.sources.Length() > 0 {
		return strconv.Itoa(r.sources.Length())
	}
	return "0"
}

func (r *record) calculateRates() []string {
	// Returns string slice of rates
	var ret []string
	r.CalculateAvgages()
	ret = append(ret, strconv.Itoa(r.grandtotal))                   //TotalRecords
	ret = append(ret, strconv.Itoa(r.total))                        //RecordsWithDenominators
	ret = append(ret, strconv.Itoa(r.cancer))                       //NeoplasiaRecords
	ret = append(ret, r.formatRate(r.cancer, r.grandtotal))         //NeoplasiaPrevalence
	ret = append(ret, strconv.Itoa(r.malignant))                    //Malignant
	ret = append(ret, r.formatRate(r.malignant, r.grandtotal))      //MalignancyPrevalence
	ret = append(ret, r.formatRate(r.maltotal, r.allcancer))        //PropMalignant
	ret = append(ret, strconv.Itoa(r.benign))                       //benign
	ret = append(ret, r.formatRate(r.benign, r.grandtotal))         //benignPrevalence
	ret = append(ret, r.formatRate(r.bentotal, r.allcancer))        //Propbenign
	ret = append(ret, strconv.FormatFloat(r.age, 'f', 2, 64))       //AverageAge(months)
	ret = append(ret, strconv.FormatFloat(r.cancerage, 'f', 2, 64)) //AvgAgeNeoplasia(months)
	ret = append(ret, strconv.Itoa(r.male))                         //Male
	ret = append(ret, strconv.Itoa(r.female))                       //Female
	ret = append(ret, strconv.Itoa(r.malecancer))                   //MaleNeoplasia
	ret = append(ret, strconv.Itoa(r.femalecancer))                 //FemaleNeoplasia
	ret = append(ret, strconv.Itoa(r.necropsy))                     //Necropsies
	ret = append(ret, r.setsources())                               //Sources
	for idx, i := range ret {
		// Replace -1 with NA
		if strings.Split(i, ".")[0] == "-1" {
			ret[idx] = "NA"
		}
	}
	return ret
}

func (r *record) cancerMeasures(age float64, sex, mal, service string) {
	// Adds cancer measures
	r.allcancer++
	if service != "MSU" {
		r.cancer++
		r.cancerage++
		if sex == "male" {
			r.malecancer++
		} else if sex == "female" {
			r.femalecancer++
		}
		if mal == "1" {
			r.malignant++
		} else if mal == "0" {
			r.benign++
		}
	}
	// Count all malignant and benign
	if mal == "1" {
		r.maltotal++
	} else if mal == "0" {
		r.bentotal++
	}
}

func (r *record) addTotal(n int) {
	// Adds n to total and grandtotal
	r.grandtotal += n
	r.total += n
}

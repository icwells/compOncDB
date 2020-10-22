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

func NewRecord() *record {
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
	//TotalRecords,RecordsWithDenominators,NeoplasiaRecords,NeoplasiaPrevalence,Malignant,MalignancyPrevalence,PropMalignant,
	//benign,benignPrevalence,Propbenign,AverageAge(months),AvgAgeNeoplasia(months),Male,Female,MaleNeoplasia,FemaleNeoplasia,Necropsies,Sources
	ret = append(ret, strconv.Itoa(r.grandtotal))
	ret = append(ret, strconv.Itoa(r.total))
	ret = append(ret, strconv.Itoa(r.cancer))
	ret = append(ret, r.formatRate(r.cancer, r.grandtotal))
	ret = append(ret, strconv.Itoa(r.malignant))
	ret = append(ret, r.formatRate(r.malignant, r.grandtotal))
	ret = append(ret, r.formatRate(r.maltotal, r.allcancer))
	ret = append(ret, strconv.Itoa(r.benign))
	ret = append(ret, r.formatRate(r.benign, r.grandtotal))
	ret = append(ret, r.formatRate(r.bentotal, r.allcancer))
	ret = append(ret, strconv.FormatFloat(r.age, 'f', 2, 64))
	ret = append(ret, strconv.FormatFloat(r.cancerage, 'f', 2, 64))
	ret = append(ret, strconv.Itoa(r.male))
	ret = append(ret, strconv.Itoa(r.female))
	ret = append(ret, strconv.Itoa(r.malecancer))
	ret = append(ret, strconv.Itoa(r.femalecancer))
	ret = append(ret, strconv.Itoa(r.necropsy))
	ret = append(ret, r.setsources())
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

/*func (r *record) Add(v *record) {
	// Add values from v to r
	r.total += v.total
	r.age += v.age
	r.male += v.male
	r.female += v.female
	r.cancer += v.cancer
	r.cancerage += v.cancerage
	r.malecancer += v.malecancer
	r.femalecancer += v.femalecancer
	r.malignant += v.malignant
	r.benign += v.benign
	r.necropsy += v.necropsy
	r.allcancer += v.allcancer
	r.maltotal += v.maltotal
	r.bentotal += v.bentotal
	for _, i := range v.sources.ToStringSlice() {
		r.sources.Add(i)
	}
}*/

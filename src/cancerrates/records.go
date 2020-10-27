// Defines record struct and methods

package cancerrates

import (
	"github.com/icwells/simpleset"
	"strconv"
	//"strings"
)

func avgAge(n float64, d int) string {
	// Returns n/d
	if n > 0.0 && d > 0 {
		r := n / float64(d)
		return strconv.FormatFloat(r, 'f', 2, 64)
	}
	return "0.00"
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

func (r *record) formatRate(n, d int) string {
	// Divides n by d and returns formatted string
	if d != 0 {
		v := float64(n) / float64(d)
		return strconv.FormatFloat(v, 'f', 2, 64)
	}
	return "0.00"
}

func (r *record) setsources() string {
	// Returns number of unique sources
	if r.sources.Length() > 0 {
		return strconv.Itoa(r.sources.Length())
	}
	return "0"
}

func (r *record) calculateRates(d int) []string {
	// Returns string slice of rates
	var ret []string
	if d < 0 {
		d = r.total
	}
	ret = append(ret, strconv.Itoa(r.grandtotal))            //TotalRecords
	ret = append(ret, strconv.Itoa(r.total))                 //RecordsWithDenominators
	ret = append(ret, strconv.Itoa(r.allcancer))             //TotalNeoplasia
	ret = append(ret, strconv.Itoa(r.cancer))                //NeoplasiaWithDenominators
	ret = append(ret, r.formatRate(r.cancer, d))             //NeoplasiaPrevalence
	ret = append(ret, strconv.Itoa(r.malignant))             //Malignant
	ret = append(ret, r.formatRate(r.malignant, d))          //MalignancyPrevalence
	ret = append(ret, r.formatRate(r.maltotal, r.allcancer)) //PropMalignant
	ret = append(ret, strconv.Itoa(r.benign))                //benign
	ret = append(ret, r.formatRate(r.benign, d))             //benignPrevalence
	ret = append(ret, r.formatRate(r.bentotal, r.allcancer)) //Propbenign
	ret = append(ret, avgAge(r.age, r.total))                //AverageAge(months)
	ret = append(ret, avgAge(r.cancerage, r.cancer))         //AvgAgeNeoplasia(months)
	ret = append(ret, strconv.Itoa(r.male))                  //Male
	ret = append(ret, strconv.Itoa(r.female))                //Female
	ret = append(ret, strconv.Itoa(r.malecancer))            //MaleNeoplasia
	ret = append(ret, strconv.Itoa(r.femalecancer))          //FemaleNeoplasia
	ret = append(ret, strconv.Itoa(r.necropsy))              //Necropsies
	ret = append(ret, r.setsources())                        //Sources
	/*for idx, i := range ret {
		// Replace -1 with NA
		if strings.Split(i, ".")[0] == "-1" {
			ret[idx] = "NA"
		}
	}*/
	return ret
}

func (r *record) cancerMeasures(age float64, sex, mal, service string) {
	// Adds cancer measures
	r.allcancer++
	if mal == "1" {
		r.maltotal++
	} else if mal == "0" {
		r.bentotal++
	}
	if service != "MSU" {
		r.cancer++
		r.cancerage += age
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
}

func (r *record) addTotal(n int) {
	// Adds n to total and grandtotal
	r.grandtotal += n
	r.total += n
}

func (r *record) Copy() *record {
	// Returns deep copy of struct
	c := newRecord()
	c.age = r.age
	c.allcancer = r.allcancer
	c.benign = r.benign
	c.bentotal = r.bentotal
	c.cancer = r.cancer
	c.cancerage = r.cancerage
	c.female = r.female
	c.femalecancer = r.femalecancer
	c.grandtotal = r.grandtotal
	c.male = r.male
	c.malecancer = r.malecancer
	c.malignant = r.malignant
	c.maltotal = r.maltotal
	c.necropsy = r.necropsy
	c.total = r.total
	for _, i := range r.sources.ToStringSlice() {
		c.sources.Add(i)
	}
	return c
}

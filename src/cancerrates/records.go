// Defines record struct and methods

package cancerrates

import (
	"github.com/icwells/simpleset"
	"strconv"
)

func avgAge(n float64, d int) string {
	// Returns n/d
	if d > 0 {
		r := n / float64(d)
		return strconv.FormatFloat(r, 'f', 2, 64)
	}
	return "NA"
}

type Record struct {
	age          float64
	agetotal     int
	allcancer    int
	benign       int
	bentotal     int
	cancer       int
	cancerage    float64
	catotal      int
	female       int
	femalecancer int
	femalemal    int
	grandtotal   int
	male         int
	malecancer   int
	malemal      int
	malignant    int
	maltotal     int
	necropsy     int
	sources      *simpleset.Set
	total        int
}

func newRecord() *Record {
	// Initializes new record struct
	r := new(Record)
	r.sources = simpleset.NewStringSet()
	return r
}

func (r *Record) formatRate(n, d int) string {
	// Divides n by d and returns formatted string
	if d != 0 {
		v := float64(n) / float64(d)
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return "NA"
}

func (r *Record) setsources() string {
	// Returns number of unique sources
	if r.sources.Length() > 0 {
		return strconv.Itoa(r.sources.Length())
	}
	return "0"
}

func (r *Record) calculateRates(d, notissue int, age, sex bool) []string {
	// Returns string slice of rates
	var ret []string
	var nt string
	den := true
	if d < 0 {
		// Remove tissue-specific columns
		den = false
		d = r.total
	} else if notissue >= 0 {
		nt = strconv.Itoa(notissue)
	}
	malknown := r.maltotal + r.bentotal
	ret = append(ret, strconv.Itoa(r.total)) //RecordsWithDenominators
	if den {
		ret = append(ret, strconv.Itoa(d)) //NeoplasiaDenominators
	}
	ret = append(ret, strconv.Itoa(r.cancer))    //NeoplasiaWithDenominators
	ret = append(ret, r.formatRate(r.cancer, d)) //NeoplasiaPrevalence
	ret = append(ret, DASH)
	ret = append(ret, strconv.Itoa(malknown))             //MalignancyKnown
	ret = append(ret, strconv.Itoa(r.malignant))          //Malignant
	ret = append(ret, r.formatRate(r.malignant, d))       //MalignancyPrevalence
	ret = append(ret, r.formatRate(r.maltotal, malknown)) //PropMalignant
	ret = append(ret, strconv.Itoa(r.benign))             //benign
	ret = append(ret, r.formatRate(r.benign, d))          //benignPrevalence
	ret = append(ret, r.formatRate(r.bentotal, malknown)) //Propbenign
	if age {
		ret = append(ret, DASH)
		ret = append(ret, avgAge(r.age, r.agetotal))      //AverageAge(months)
		ret = append(ret, avgAge(r.cancerage, r.catotal)) //AvgAgeNeoplasia(months)
	}
	if sex {
		ret = append(ret, DASH)
		ret = append(ret, strconv.Itoa(r.male))         //Male
		ret = append(ret, strconv.Itoa(r.malecancer))   //MaleNeoplasia
		ret = append(ret, strconv.Itoa(r.malemal))      //MaleMalignant
		ret = append(ret, strconv.Itoa(r.female))       //Female
		ret = append(ret, strconv.Itoa(r.femalecancer)) //FemaleNeoplasia
		ret = append(ret, strconv.Itoa(r.femalemal))    //FemaleMalignant
	}
	ret = append(ret, DASH)
	ret = append(ret, strconv.Itoa(r.grandtotal)) //RecordsFromAllSources
	ret = append(ret, strconv.Itoa(r.allcancer))  //NeoplasiaFromAllSources
	ret = append(ret, strconv.Itoa(r.necropsy))   //Necropsies
	ret = append(ret, r.setsources())             //Sources
	if den {
		ret = append(ret, nt) //NoTissueInfo
	}
	return ret
}

func (r *Record) cancerMeasures(allrecords bool, age, sex, mal, service string) {
	// Adds cancer measures
	r.allcancer++
	if mal == "1" {
		r.maltotal++
	} else if mal == "0" {
		r.bentotal++
	}
	if allrecords {
		r.cancer++
		f, err := strconv.ParseFloat(age, 64)
		if err == nil && f >= 0.0 {
			r.cancerage += f
			r.catotal++
		}
		if sex == "male" {
			r.malecancer++
		} else if sex == "female" {
			r.femalecancer++
		}
		if mal == "1" {
			r.malignant++
			if sex == "male" {
				r.malemal++
			} else if sex == "female" {
				r.femalemal++
			}
		} else if mal == "0" {
			r.benign++
		}
	}
}

func (r *Record) nonCancerMeasures(allrecords bool, age, sex, nec, service, aid string) {
	// Adds non-cancer meaures
	r.grandtotal++
	if allrecords {
		// Add to total and grandtotal
		r.total++
		f, err := strconv.ParseFloat(age, 64)
		if err == nil && f >= 0.0 {
			r.age += f
			r.agetotal++
		}
		if sex == "male" {
			r.male++
		} else if sex == "female" {
			r.female++
		}
		if nec == "1" {
			r.necropsy++
		}
	}
	r.sources.Add(aid)
}

func (r *Record) addTotal(n int) {
	// Adds n to total and grandtotal
	r.grandtotal += n
	r.total += n
}

func (r *Record) Add(v *Record) {
	// Adds v values to record
	r.age += v.age
	r.allcancer += v.allcancer
	r.benign += v.benign
	r.bentotal += v.bentotal
	r.cancer += v.cancer
	r.cancerage += v.cancerage
	r.female += v.female
	r.femalecancer += v.femalecancer
	r.grandtotal += v.grandtotal
	r.male += v.male
	r.malecancer += v.malecancer
	r.malignant += v.malignant
	r.maltotal += v.maltotal
	r.necropsy += v.necropsy
	r.total += v.total
	for _, i := range v.sources.ToStringSlice() {
		r.sources.Add(i)
	}
}

func (r *Record) Copy() *Record {
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
	c.sources = r.sources.Copy()
	c.total = r.total
	return c
}

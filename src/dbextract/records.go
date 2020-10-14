// Defines record struct and methods

package dbextract

import (
	"fmt"
	"github.com/icwells/simpleset"
	"strconv"
	"strings"
)

type Record struct {
	Age          float64
	allcancer    int
	Benign       int
	bentotal     int
	Cancer       int
	Cancerage    float64
	Female       int
	Femalecancer int
	grandtotal   int
	Lifehistory  []string
	Male         int
	Malecancer   int
	Malignant    int
	maltotal     int
	Necropsy     int
	Sources      *simpleset.Set
	Taxonomy     []string
	Total        int
}

func NewRecord() *Record {
	// Initializes new record struct
	r := new(Record)
	r.Sources = simpleset.NewStringSet()
	return r
}

func (r *Record) setTaxonomy(taxonomy []string) {
	// Stores taxonomy
	r.Taxonomy = taxonomy
}

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

func (r *Record) String() string {
	// Returns formatted string of record attributes
	ret := fmt.Sprintf("\nSpecies: %s\n", r.Taxonomy[6])
	ret += fmt.Sprintf("Total: %d\n", r.Total)
	ret += fmt.Sprintf("Cancer Records: %d", r.Cancer)
	return ret
}

func (r *Record) CalculateAvgAges() {
	// Sets average age and average cancer age
	r.Age = avgAge(r.Age, r.Total)
	r.Cancerage = avgAge(r.Cancerage, r.Cancer)
}

func (r *Record) formatRate(n, d int) string {
	// Divides n by d and returns formatted string
	var v float64
	if d != 0 {
		v = float64(n) / float64(d)
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func (r *Record) setSources() string {
	// Returns number of unique sources
	if r.Sources.Length() > 0 {
		return strconv.Itoa(r.Sources.Length())
	}
	return "0"
}

func (r *Record) CalculateRates(id, name string, lh bool) []string {
	// Returns string slice of rates
	var ret []string
	r.CalculateAvgAges()
	if len(id) > 0 {
		ret = append(ret, id)
	}
	if len(r.Taxonomy) > 0 {
		ret = append(ret, r.Taxonomy...)
	}
	if len(name) > 0 {
		ret = append(ret, name)
	}
	//"AdultRecords,CancerRecords,CancerRate,AverageAge(months),AvgAgeCancer(months),Male,Female\n"
	ret = append(ret, strconv.Itoa(r.Total))
	ret = append(ret, strconv.Itoa(r.Cancer))
	ret = append(ret, r.formatRate(r.Cancer, r.grandtotal))
	ret = append(ret, strconv.Itoa(r.Malignant))
	ret = append(ret, r.formatRate(r.Malignant, r.grandtotal))
	ret = append(ret, r.formatRate(r.maltotal, r.allcancer))
	ret = append(ret, strconv.Itoa(r.Benign))
	ret = append(ret, r.formatRate(r.Benign, r.grandtotal))
	ret = append(ret, r.formatRate(r.bentotal, r.allcancer))
	ret = append(ret, strconv.FormatFloat(r.Age, 'f', 2, 64))
	ret = append(ret, strconv.FormatFloat(r.Cancerage, 'f', 2, 64))
	ret = append(ret, strconv.Itoa(r.Male))
	ret = append(ret, strconv.Itoa(r.Female))
	ret = append(ret, strconv.Itoa(r.Malecancer))
	ret = append(ret, strconv.Itoa(r.Femalecancer))
	ret = append(ret, strconv.Itoa(r.Necropsy))
	ret = append(ret, r.setSources())
	for idx, i := range ret {
		// Replace -1 with NA
		if strings.Split(i, ".")[0] == "-1" {
			ret[idx] = "NA"
		}
	}
	if lh {
		ret = append(ret, r.Lifehistory...)
	}
	return ret
}

func (r *Record) Add(v *Record) {
	// Add values from v to r
	r.Total += v.Total
	r.Age += v.Age
	r.Male += v.Male
	r.Female += v.Female
	r.Cancer += v.Cancer
	r.Cancerage += v.Cancerage
	r.Malecancer += v.Malecancer
	r.Femalecancer += v.Femalecancer
	r.Malignant += v.Malignant
	r.Benign += v.Benign
	r.Necropsy += v.Necropsy
	r.grandtotal += v.grandtotal
	r.allcancer += v.allcancer
	r.maltotal += v.maltotal
	r.bentotal += v.bentotal
	for _, i := range v.Sources.ToStringSlice() {
		r.Sources.Add(i)
	}
}

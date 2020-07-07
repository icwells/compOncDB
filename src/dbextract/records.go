// Defines record struct and methods

package dbextract

import (
	"fmt"
	"strconv"
	"strings"
)

type Record struct {
	Taxonomy     []string
	Total        int
	Age          float64
	Male         int
	Female       int
	Cancer       int
	Cancerage    float64
	Malecancer   int
	Femalecancer int
	Malignant    int
	Lifehistory  []string
}

func NewRecord(taxonomy []string) *Record {
	// Initializes new record struct
	r := new(Record)
	r.Taxonomy = taxonomy
	return r
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
		v = float64(n)/float64(d)
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func (r *Record) CalculateRates(id string, lh bool) []string {
	// Returns string slice of rates
	var ret []string
	r.CalculateAvgAges()
	if id != "" {
		ret = append(ret, id)
	}
	ret = append(ret, r.Taxonomy...)
	//"AdultRecords,CancerRecords,CancerRate,AverageAge(months),AvgAgeCancer(months),Male,Female\n"
	ret = append(ret, strconv.Itoa(r.Total))
	ret = append(ret, strconv.Itoa(r.Cancer))
	ret = append(ret, r.formatRate(r.Cancer, r.Total))
	ret = append(ret, strconv.Itoa(r.Malignant))
	ret = append(ret, r.formatRate(r.Malignant, r.Total))
	ret = append(ret, r.formatRate(r.Malignant, r.Cancer))
	ret = append(ret, strconv.FormatFloat(r.Age, 'f', 2, 64))
	ret = append(ret, strconv.FormatFloat(r.Cancerage, 'f', 2, 64))
	ret = append(ret, strconv.Itoa(r.Male))
	ret = append(ret, strconv.Itoa(r.Female))
	ret = append(ret, strconv.Itoa(r.Malecancer))
	ret = append(ret, strconv.Itoa(r.Femalecancer))
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

func getRecKeys(records map[string]*Record) string {
	// Returns string of taxa_ids
	first := true
	var buffer strings.Builder
	for k := range records {
		if first == false {
			// Write name with preceding comma
			buffer.WriteByte(',')
			buffer.WriteString(k)
		} else {
			buffer.WriteString(k)
			first = false
		}
	}
	return buffer.String()
}

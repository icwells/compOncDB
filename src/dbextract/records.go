// Defines record struct and methods

package dbextract

import (
	"fmt"
	"strconv"
	"strings"
)

type Record struct {
	Taxonomy     []string
	Infant       float64
	Total        int
	Age          float64
	Male         int
	Female       int
	Cancer       int
	Cancerage    float64
	Adult        int
	Malecancer   int
	Femalecancer int
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
	r.Age = avgAge(r.Age, r.Adult)
	r.Cancerage = avgAge(r.Cancerage, r.Cancer)
}

/*func (r *Record) ToSlice(id string) []string {
	// Returns string slice of values for upload to table
	var ret []string
	r.CalculateAvgAges()
	ret = append(ret, id)
	ret = append(ret, strconv.Itoa(r.Total))
	ret = append(ret, strconv.FormatFloat(r.Age, 'f', -1, 64))
	ret = append(ret, strconv.Itoa(r.Adult))
	ret = append(ret, strconv.Itoa(r.Male))
	ret = append(ret, strconv.Itoa(r.Female))
	ret = append(ret, strconv.Itoa(r.Cancer))
	ret = append(ret, strconv.FormatFloat(r.Cancerage, 'f', -1, 64))
	ret = append(ret, strconv.Itoa(r.Malecancer))
	ret = append(ret, strconv.Itoa(r.Femalecancer))
	return ret
}*/

func (r *Record) CalculateRates(id string, lh bool) []string {
	// Returns string slice of rates
	var ret []string
	r.CalculateAvgAges()
	if id != "" {
		ret = append(ret, id)
	}
	ret = append(ret, r.Taxonomy...)
	//"AdultRecords,CancerRecords,CancerRate,AverageAge(months),AvgAgeCancer(months),Male,Female\n"
	ret = append(ret, strconv.Itoa(r.Adult))
	ret = append(ret, strconv.Itoa(r.Cancer))
	ret = append(ret, strconv.FormatFloat(float64(r.Cancer)/float64(r.Adult), 'f', 2, 64))
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

/*func (r *Record) SetRecord(row []string) {
	// Reads values from Totals table entry
	r.Total, _ = strconv.Atoi(row[1])
	r.Age, _ = strconv.ParseFloat(row[2], 64)
	r.Adult, _ = strconv.Atoi(row[3])
	r.Male, _ = strconv.Atoi(row[4])
	r.Female, _ = strconv.Atoi(row[5])
	r.Cancer, _ = strconv.Atoi((row[6]))
	r.Cancerage, _ = strconv.ParseFloat(row[7], 64)
	r.Malecancer, _ = strconv.Atoi((row[8]))
	r.Femalecancer, _ = strconv.Atoi((row[9]))
}*/

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

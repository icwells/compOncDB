// Defines record struct and methods

package dbupload

import (
	"bytes"
	"fmt"
	"strconv"
)

type Record struct {
	Species      string
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
}

func avgAge(n float64, d int) string {
	// Returns string of n/d
	var ret string
	if n > 0.0 && d > 0 {
		age := n / float64(d)
		ret = strconv.FormatFloat(age, 'f', -1, 64)
	} else {
		ret = "-1"
	}
	return ret
}

func (r *Record) String() string {
	// Returns formatted string of record attributes
	ret := fmt.Sprintf("\nSpecies: %s\n", r.Species)
	ret += fmt.Sprintf("Total: %d\n", r.Total)
	ret += fmt.Sprintf("Cancer Records: %d", r.Cancer)
	return ret
}

func (r *Record) getAvgAge() string {
	// Returns string of avg age
	return avgAge(r.Age, r.Adult)
}

func (r *Record) getCancerAge() string {
	// Returns string of average cancer record age
	return avgAge(r.Cancerage, r.Cancer)
}

func (r *Record) ToSlice(id string) []string {
	// Returns string slice of values for upload to table
	var ret []string
	ret = append(ret, id)
	ret = append(ret, strconv.Itoa(r.Total))
	ret = append(ret, r.getAvgAge())
	ret = append(ret, strconv.Itoa(r.Adult))
	ret = append(ret, strconv.Itoa(r.Male))
	ret = append(ret, strconv.Itoa(r.Female))
	ret = append(ret, strconv.Itoa(r.Cancer))
	ret = append(ret, r.getCancerAge())
	ret = append(ret, strconv.Itoa(r.Malecancer))
	ret = append(ret, strconv.Itoa(r.Femalecancer))
	return ret
}

func (r *Record) CalculateRates() []string {
	// Returns string slice of rates
	//"ScientificName,AdultRecords,CancerRecords,CancerRate,AverageAge(months),AvgAgeCancer(months),Male,Female\n"
	ret := []string{r.Species}
	ret = append(ret, strconv.Itoa(r.Adult))
	ret = append(ret, strconv.Itoa(r.Cancer))
	// Calculate rates
	rate := float64(r.Cancer) / float64(r.Adult)
	// Append rates to slice and return
	ret = append(ret, strconv.FormatFloat(rate, 'f', 2, 64))
	ret = append(ret, strconv.FormatFloat(r.Age, 'f', 2, 64))
	ret = append(ret, strconv.FormatFloat(r.Cancerage, 'f', 2, 64))
	ret = append(ret, strconv.Itoa(r.Male))
	ret = append(ret, strconv.Itoa(r.Female))
	ret = append(ret, strconv.Itoa(r.Malecancer))
	ret = append(ret, strconv.Itoa(r.Femalecancer))
	return ret
}

func (r *Record) SetRecord(row []string) {
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
}

func InMapRec(m map[string]*Record, s string) bool {
	// Return true if s is a key in m
	_, ret := m[s]
	return ret
}

func GetRecKeys(records map[string]*Record) string {
	// Returns string of taxa_ids
	first := true
	buffer := bytes.NewBufferString("")
	for k, _ := range records {
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

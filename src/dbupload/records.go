// Defines record struct and methods

package dbupload

import (
	"bytes"
	"fmt"
	"github.com/icwells/dbIO"
	"strconv"
)

type Record struct {
	species      string
	infant       float64
	total        int
	age          float64
	male         int
	female       int
	cancer       int
	cancerage    float64
	adult        int
	malecancer   int
	femalecancer int
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
	ret := fmt.Sprintf("\nSpecies: %s\n", r.species)
	ret += fmt.Sprintf("Total: %d\n", r.total)
	ret += fmt.Sprintf("Cancer Records: %d", r.cancer)
	return ret
}

func (r *Record) getAvgAge() string {
	// Returns string of avg age
	return avgAge(r.age, r.adult)
}

func (r *Record) getCancerAge() string {
	// Returns string of average cancer record age
	return avgAge(r.cancerage, r.cancer)
}

func (r *Record) toSlice(id string) []string {
	// Returns string slice of values for upload to table
	var ret []string
	ret = append(ret, id)
	ret = append(ret, strconv.Itoa(r.total))
	ret = append(ret, r.getAvgAge())
	ret = append(ret, strconv.Itoa(r.adult))
	ret = append(ret, strconv.Itoa(r.male))
	ret = append(ret, strconv.Itoa(r.female))
	ret = append(ret, strconv.Itoa(r.cancer))
	ret = append(ret, r.getCancerAge())
	ret = append(ret, strconv.Itoa(r.malecancer))
	ret = append(ret, strconv.Itoa(r.femalecancer))
	return ret
}

func (r *dbupload.Record) calculateRates() []string {
	// Returns string slice of rates
	//"ScientificName,AdultRecords,CancerRecords,CancerRate,AverageAge(months),AvgAgeCancer(months),Male:Female\n"
	ret := []string{r.species}
	ret = append(ret, strconv.Itoa(r.adult))
	ret = append(ret, strconv.Itoa(r.cancer))
	// Calculate rates
	rate := float64(r.cancer) / float64(r.adult)
	// Append rates to slice and return
	ret = append(ret, strconv.FormatFloat(rate, 'f', 2, 64))
	ret = append(ret, strconv.FormatFloat(r.age, 'f', 2, 64))
	ret = append(ret, strconv.FormatFloat(r.cancerage, 'f', 2, 64))
	ret = append(ret, strconv.Itoa(r.male))
	ret = append(ret, strconv.Itoa(r.female))
	ret = append(ret, strconv.Itoa(r.malecancer))
	ret = append(ret, strconv.Itoa(r.femalecancer))
	return ret
}

func (r *dbupload.Record) setRecord(row []string) {
	// Reads values from Totals table entry
	r.total, _ = strconv.Atoi(row[1])
	r.age, _ = strconv.ParseFloat(row[2], 64)
	r.adult, _ = strconv.Atoi(row[3])
	r.male, _ = strconv.Atoi(row[4])
	r.female, _ = strconv.Atoi(row[5])
	r.cancer, _ = strconv.Atoi((row[6]))
	r.cancerage, _ = strconv.ParseFloat(row[7], 64)
}

// Summarizes the number of species with a minimum number of records for each sex

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"strconv"
	"time"
)

var (
	ages     = kingpin.Flag("ages", "Require non-zero age values.").Default("false").Bool()
	approved = kingpin.Flag("approved", "Only counts approved necropsies.").Default("false").Bool()
	outfile  = kingpin.Flag("outfile", "Name of output file.").Short('o').Required().String()
	user     = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type sexTotals struct {
	female int
	male   int
	total  int
}

type sexSummary struct {
	ages     bool
	approved bool
	columns  [][]string
	header   string
	min      map[int]int
	species  map[string]*sexTotals
	steps    []int
	total    map[int]int
}

func newSexSummary() *sexSummary {
	// Returns initialized struct
	s := new(sexSummary)
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	s.ages = *ages
	s.approved = *approved
	s.columns = db.GetColumns("Records", []string{"Species", "Sex", "age_months", "Approved", "Necropsy"})
	s.header = "Min,TotalSpecies,SpeciesWithMaleAndFemale"
	s.min = make(map[int]int)
	s.species = make(map[string]*sexTotals)
	s.steps = []int{10, 15, 20, 25, 30, 40, 45, 50}
	s.total = make(map[int]int)
	for _, i := range s.steps {
		s.min[i] = 0
		s.total[i] = 0
	}
	return s
}

func (s *sexSummary) write() {
	// Writes results to file
	fmt.Println("\tWriting results to file...")
	var res [][]string
	for k, v := range s.min {
		res = append(res, []string{strconv.Itoa(k), strconv.Itoa(s.total[k]), strconv.Itoa(v)})
	}
	iotools.WriteToCSV(*outfile, s.header, res)
}

func (s *sexSummary) getTotals() {
	// Counts number of species at each step
	fmt.Println("\tCalculating species minimums...")
	for _, v := range s.species {
		for _, i := range s.steps {
			if v.total >= i {
				s.total[i]++
				if v.male >= i && v.female >= 1 {
					s.min[i]++
				}
			}
		}
	}
}

func (s *sexSummary) pass(age, app, nec string) bool {
	// Returns true if approved is false or app and nec are true
	if !s.ages || age != "-1.00" {
		if !s.approved {
			return true
		} else if app == "-1" && nec == "-1" {
			return true
		}
	}
	return false
}

func (s *sexSummary) setSpecies() {
	// Creates entries for species
	fmt.Println("\n\tSetting species slice...")
	for _, i := range s.columns {
		sp := i[0]
		if _, ex := s.species[sp]; !ex {
			s.species[sp] = new(sexTotals)
		}
		if s.pass(i[2], i[3], i[4]) {
			sex := i[1]
			s.species[sp].total++
			if sex == "male" {
				s.species[sp].male++
			} else if sex == "female" {
				s.species[sp].female++
			}
		}
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	s := newSexSummary()
	s.setSpecies()
	s.getTotals()
	s.write()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

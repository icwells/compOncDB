// Summarizes basic statistics about the database

package dbextract

import (
	"fmt"
	"github.com/icwells/dbIO"
	"strconv"
)

func getRow(name string, num, den int) []string {
	// Returns string slice of name, numerator, and percent
	ret := []string{name}
	ret = append(ret, string(num))
	if den > 0 {
		percent := (float64(num)/float64(den)) * 100
		ret = append(ret, strconv.FormatFloat(percent, 'f', 2, 64) + "%")
	}
	return ret
}

type summary struct {
	total	int
	male	int
	female	int
	age		int
	mass	int
	hyper	int
	mal		int
	benign	int
	nec		int
	taxa	int
	com		int
	hist	int
}

func (s *summary) toSlice() [][]string {
	// Calculates percents and returns slice of string slices
	var ret [][]string
	ret = append(ret, getRow("total", s.total, 0))
	ret = append(ret, getRow("male", s.male, s.total))
	ret = append(ret, getRow("female", s.female, s.total))
	ret = append(ret, getRow("entries with ages", s.age, s.total))
	ret = append(ret, getRow("cancer", s.mass, s.total))
	ret = append(ret, getRow("hyperplasia", s.hyper, s.total))
	ret = append(ret, getRow("malignant", s.mal, s.total))
	ret = append(ret, getRow("benign", s.benign, s.total))
	ret = append(ret, getRow("necropsies", s.nec, s.total))
	ret = append(ret, getRow("taxonomies", s.taxa, 0))
	ret = append(ret, getRow("taxonomies with common names", s.com, s.taxa))
	ret = append(ret, getRow("taxonomies with life history data", s.hist, s.taxa))
	return ret
}

func (s *summary) setTotals(db *dbIO.DBIO) {
	// Queries database for total number of occurances
	s.total = db.Count("Patients", "ID", "", "", "", true)
	s.male = db.Count("Patients", "Sex", "*", "=", "male", false)
	s.female = db.Count("Patients", "Sex", "*", "=", "female", false)
	s.age = db.Count("Patients", "Age", "*", ">=", "0", false)
	s.mass = db.Count("Diagnosis", "Masspresent", "*", "=", "1", false)
	s.hyper = db.Count("Diagnosis", "Hyperplasia", "*", "=", "1", false)
	s.nec = db.Count("Diagnosis", "Necropsy", "*", "=", "1", false)
	s.mal = db.Count("Tumor_relation", "Malignant", "*", "=", "1", false)
	s.benign = db.Count("Tumor_relation", "Malignant", "*", "=", "0", false)
	s.taxa = db.Count("Taxonomy", "taxa_id", "", "", "", true)
	s.com = db.Count("Common", "taxa_id", "", "", "", true)
	s.hist = db.Count("Life_history", "taxa_id", "", "", "", true)
}

func GetSummary(db *dbIO.DBIO) [][]string {
	// Returns summary statistics from database
	fmt.Println("\n\tGenerating database summary statistics...")
	s := new(summary)
	s.setTotals(db)
	return s.toSlice()
}

// Summarizes basic statistics about the database

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/simpleset"
	"strconv"
	"strings"
)

func getRow(name string, num, den int) []string {
	// Returns string slice of name, numerator, and percent
	ret := []string{name}
	ret = append(ret, strconv.Itoa(num))
	if den > 0 {
		percent := (float64(num) / float64(den)) * 100
		ret = append(ret, fmt.Sprintf("%.2f%%", percent))
	}
	return ret
}

type summary struct {
	total  int
	infant int
	adult  int
	male   int
	female int
	age    int
	mass   int
	hyper  int
	mal    int
	benign int
	nec    int
	taxa   int
	path   int
	tmass  int
	com    int
	hist   int
}

func (s *summary) toSlice() [][]string {
	// Calculates percents and returns slice of string slices
	var ret [][]string
	ret = append(ret, getRow("total", s.total, 0))
	ret = append(ret, getRow("infant records", s.infant, s.total))
	ret = append(ret, getRow("adult records", s.adult, s.total))
	ret = append(ret, getRow("male", s.male, s.total))
	ret = append(ret, getRow("female", s.female, s.total))
	ret = append(ret, getRow("entries with ages", s.age, s.total))
	ret = append(ret, getRow("neoplasia", s.mass, s.total))
	ret = append(ret, getRow("hyperplasia", s.hyper, s.total))
	ret = append(ret, getRow("malignant", s.mal, s.total))
	ret = append(ret, getRow("benign", s.benign, s.total))
	ret = append(ret, getRow("necropsies", s.nec, s.total))
	ret = append(ret, getRow("taxonomies", s.taxa, 0))
	ret = append(ret, getRow("taxonomies with pathology records", s.path, s.taxa))
	ret = append(ret, getRow("taxonomies with cancer records", s.tmass, s.taxa))
	ret = append(ret, getRow("taxonomies with common names", s.com, s.taxa))
	ret = append(ret, getRow("taxonomies with life history data", s.hist, s.taxa))
	return ret
}

func (s *summary) setCancerTaxa(db *dbIO.DBIO) {
	// Identifies number of unique species with cancer records
	ids := simpleset.NewStringSet()
	rows := db.GetRows("Diagnosis", "Masspresent", "1", "ID")
	for _, i := range rows {
		ids.Add(i[0])
	}
	tids := db.GetRows("Patient", "ID", strings.Join(ids.ToStringSlice(), ","), "taxa_id")
	ids = simpleset.NewStringSet()
	for _, i := range tids {
		ids.Add(i[0])
	}
	s.tmass = ids.Length()
}

func (s *summary) getNumAdult(db *dbIO.DBIO) {
	// Gets total adult and infant records
	s.adult = db.Count("Patient", "Infant", "Infant", "=", "0", false)
	s.infant = db.Count("Patient", "Infant", "Infant", "=", "1", false)
}

func (s *summary) setTotals(db *dbIO.DBIO) {
	// Queries database for total number of occurances
	s.total = db.Count("Patient", "", "ID", "", "", true)
	s.getNumAdult(db)
	s.male = db.Count("Patient", "Sex", "*", "=", "male", false)
	s.female = db.Count("Patient", "Sex", "*", "=", "female", false)
	s.age = db.Count("Patient", "Age", "*", ">=", "0", false)
	s.mass = db.Count("Diagnosis", "Masspresent", "*", "=", "1", false)
	s.hyper = db.Count("Diagnosis", "Hyperplasia", "*", "=", "1", false)
	s.nec = db.Count("Diagnosis", "Necropsy", "*", "=", "1", false)
	s.mal = db.Count("Tumor", "Malignant", "*", "=", "1", false)
	s.benign = db.Count("Tumor", "Malignant", "*", "=", "0", false)
	s.taxa = db.Count("Taxonomy", "", "taxa_id", "", "", true)
	s.path = db.Count("Patient", "", "taxa_id", "", "", true)
	s.com = db.Count("Common", "", "taxa_id", "", "", true)
	s.hist = db.Count("Life_history", "", "taxa_id", "", "", true)
	s.setCancerTaxa(db)
}

func GetSummary(db *dbIO.DBIO) [][]string {
	// Returns summary statistics from database
	codbutils.GetLogger().Println("Generating database summary statistics...")
	s := new(summary)
	s.setTotals(db)
	return s.toSlice()
}

// Summarizes completeness life history table

package dbextract

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/simpleset"
	"log"
	"strconv"
	"strings"
)

type lifeHist struct {
	all       bool
	db        *dbIO.DBIO
	diagnosis map[string][]int
	logger    *log.Logger
	res       *dataframe.Dataframe
	taxa      map[string][]string
	taxaids   string
}

func newLifeHist(db *dbIO.DBIO, all bool) *lifeHist {
	// Returns initialized struct
	l := new(lifeHist)
	l.all = all
	l.db = db
	l.diagnosis = make(map[string][]int)
	l.logger = codbutils.GetLogger()
	l.res, _ = dataframe.NewDataFrame(-1)
	l.res.SetHeader(codbutils.LifeHistorySummaryHeader())
	l.setTaxa()
	l.setDiagnses()
	return l
}

func (l *lifeHist) setDiagnses() {
	// Stores number of neoplasia and malignant records
	l.logger.Println("Getting number of records per species...")
	var patients []string
	ids := codbutils.ToMap(l.db.GetRows("Patient", "taxa_id", l.taxaids, "ID,taxa_id"))
	for k, v := range ids {
		patients = append(patients, k)
		i := v[0]
		if _, ex := l.diagnosis[i]; !ex {
			l.diagnosis[i] = []int{0, 0, 0}
		}
		l.diagnosis[i][2]++
	}
	tumor := codbutils.ToMap(l.db.GetRows("Tumor", "ID", strings.Join(patients, ","), "ID,Malignant"))
	for _, i := range l.db.GetRows("Diagnosis", "ID", strings.Join(patients, ","), "ID,Masspresent") {
		tid := ids[i[0]][0]
		if i[1] == "1" {
			l.diagnosis[tid][0]++
		}
		if tumor[tid][0] == "1" {
			l.diagnosis[tid][1]++
		}
	}
}

func (l *lifeHist) setTaxa() {
	// Sets taxonomy map and
	l.logger.Println("Getting taxa ids...")
	col := strings.Split(l.db.Columns["Taxonomy"], ",")
	// Remove source column
	col = col[:len(col)-1]
	if l.all {
		l.taxaids = strings.Join(l.db.GetColumnText("Taxonomy", "taxa_id"), ",")
		l.taxa = codbutils.ToMap(l.db.GetColumns("Taxonomy", col))
	} else {
		// Subset ids to exclude entries without records
		s := simpleset.NewStringSet()
		for _, i := range l.db.GetColumnText("Patient", "taxa_id") {
			s.Add(i)
		}
		l.taxaids = strings.Join(s.ToStringSlice(), ",")
		l.taxa = codbutils.ToMap(l.db.GetRows("Taxonomy", "taxa_id", l.taxaids, strings.Join(col, ",")))
	}
}

func (l *lifeHist) summarize() {
	// Stores y/n if value is set
	l.logger.Println("Summarizing table...")
	for k, v := range codbutils.ToMap(l.db.GetRows("Life_history", "taxa_id", l.taxaids, "*")) {
		row := append([]string{k}, l.taxa[k]...)
		var complete int
		for _, i := range v {
			if i == "-1" {
				row = append(row, "n")
			} else {
				row = append(row, "y")
				complete++
			}
		}
		p := float64(complete) / float64(len(v)) * 100
		row = append(row, strconv.FormatFloat(p, 'f', 2, 64))
		for _, i := range l.diagnosis[k] {
			row = append(row, strconv.Itoa(i))
		}
		l.res.AddRow(row)
	}
}

func LifeHistorySummary(db *dbIO.DBIO, all bool) *dataframe.Dataframe {
	// Returns life history database summarized for completeness
	l := newLifeHist(db, all)
	l.logger.Println("Summarizing life history table...")
	l.summarize()
	l.logger.Printf("Summarized %d rows from life history table.\n", l.res.Length())
	return l.res
}

// Summarizes completeness life history table

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/simpleset"
	"strconv"
	"strings"
)

type lifeHist struct {
	all       bool
	db        *dbIO.DBIO
	diagnosis map[string][]int
	res       *dataframe.Dataframe
	taxa      map[string][]string
	taxaids   string
}

func newLifeHist(db *dbIO.DBIO, all bool) *lifeHist {
	// Returns initialized struct
	l := new(lifeHist)
	l.all = all
	l.diagnosis = make(map[string][]int)
	l.db = db
	l.res, _ = dataframe.NewDataFrame(-1)
	l.res.SetHeader(codbutils.LifeHistorySummaryHeader())
	l.setTaxa()
	l.setDiagnses()
	return l
}

func (l *lifeHist) setDiagnses() {
	// Stores number of neoplasia and malignant records
	var patients []string
	ids := codbutils.ToMap(l.db.GetRows("Taxonomy", "taxa_id", l.taxaids, "ID,taxa_id"))
	for k := range ids {
		patients = append(patients, k)
	}
	for _, i := range l.db.GetRows("Tumor", "ID", strings.Join(patients, ","), "ID,Malignant") {
		tid := ids[i[0]][0]
		if _, ex := l.diagnosis[tid]; !ex {
			l.diagnosis[tid] = []int{0, 0}
		}
		l.diagnosis[tid][0]++
		if i[1] == "1" {
			l.diagnosis[tid][1]++
		}
	}
}

func (l *lifeHist) setTaxa() {
	// Sets taxonomy map and
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
		row = append(row, strconv.Itoa(l.diagnosis[k][0]))
		row = append(row, strconv.Itoa(l.diagnosis[k][1]))
		l.res.AddRow(row)
	}
}

func LifeHistorySummary(db *dbIO.DBIO, all bool) *dataframe.Dataframe {
	// Returns life history database summarized for completeness
	fmt.Println("\n\tSummarizing life history table...")
	l := newLifeHist(db, all)
	l.summarize()
	fmt.Printf("\tSummarized %d rows from life history table.\n", l.res.Length())
	return l.res
}

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
	all     bool
	db      *dbIO.DBIO
	res     *dataframe.Dataframe
	taxa    map[string][]string
	taxaids []string
}

func newLifeHist(db *dbIO.DBIO, all bool) *lifeHist {
	// Returns initialized struct
	l := new(lifeHist)
	l.all = all
	l.db = db
	l.res, _ = dataframe.NewDataFrame(-1)
	l.setTaxa()
	return l
}

func (l *lifeHist) setHeader(head []string) {
	// Stores dataframe header
	col := strings.Split(l.db.Columns["Life_history"], ",")
	// Remove redundant taxa_id column
	col = col[1:]
	head = append(head, col...)
	l.res.SetHeader(append(head, "%Complete"))
}

func (l *lifeHist) setTaxa() {
	// Sets taxonomy map and
	col := strings.Split(l.db.Columns["Taxonomy"], ",")
	// Remove source column
	col = col[:len(col)-1]
	if l.all {
		l.taxaids = l.db.GetColumnText("Taxonomy", "taxa_id")
		l.taxa = codbutils.ToMap(l.db.GetColumns("Taxonomy", col))
	} else {
		// Subset ids to exclude entries without records
		s := simpleset.NewStringSet()
		for _, i := range l.db.GetColumnText("Patient", "taxa_id") {
			s.Add(i)
		}
		l.taxaids = s.ToStringSlice()
		l.taxa = codbutils.ToMap(l.db.GetRows("Taxonomy", "taxa_id", strings.Join(l.taxaids, ","), strings.Join(col, ",")))
	}
	l.setHeader(col)
}

func (l *lifeHist) summarize() {
	// Stores y/n if value is set
	for k, v := range codbutils.ToMap(l.db.GetRows("Life_history", "taxa_id", strings.Join(l.taxaids, ","), "*")) {
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
		p := float64(complete)/float64(len(v))*100
		row = append(row, strconv.FormatFloat(p, 'f', 2, 64))
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

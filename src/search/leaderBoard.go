// Returns top cancer locations and the top species and types associated with them.

package search

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"log"
	"sort"
	"strconv"
	"strings"
)

type species struct {
	cancer int
	common string
	name   string
	typ    int
}

type leaderboard struct {
	df			*dataframe.Dataframe
	list        []*species
	logger  	*log.Logger
	table		*dataframe.Dataframe
	taxa        map[string]*species
	top			[]string
	typ         string
}

func newLeaderBoard(db *dbIO.DBIO, typ string) *leaderboard {
	// Initializes new struct
	l := new(leaderboard)
	l.df, _ = dataframe.NewDataFrame(-1)
	l.df.SetHeader([]string{"Species", "Commmon Name", "NeoplasiaRecords", "TypeRecords"})
	l.logger = codbutils.GetLogger()
	l.table, _ = SearchRecords(db, l.logger, "Infant!=1,Masspresent=1,Type!=NA", false, false)
	l.taxa = make(map[string]*species)
	l.top = make([]string, 5)
	l.typ = typ
	return l
}

func (l *leaderboard) countRecords() {
	// Counts tissue types
	l.logger.Print("Counting neoplasia records...")
	for idx := range l.table.Rows {
		sp, _ := l.table.GetCell(idx, "Species")
		typ, _ := l.table.GetCell(idx, "Type")
		if _, ex := l.taxa[sp]; !ex {
			l.taxa[sp] = new(species)
			l.taxa[sp].common, _ = l.table.GetCell(idx, "common_name")
			l.taxa[sp].name = sp
		}
		l.taxa[sp].cancer++
		for _, i := range strings.Split(typ, ";") {
			if i == l.typ {
				l.taxa[sp].typ++
			}
		}
	}
}

func (l *leaderboard) Len() int {
	return len(l.taxa)
}

func (l *leaderboard) Less(i, j int) bool {
	return l.list[i].typ > l.list[j].typ
}

func (l *leaderboard) Swap(i, j int) {
	l.list[i], l.list[j] = l.list[j], l.list[i]
}

func (l *leaderboard) sortRecords() {
	// Stores adult tumor records
	l.logger.Print("Sorting neoplasia records...")
	for _, v := range l.taxa {
		l.list = append(l.list, v)
	}
	sort.Sort(l)
	for _, v := range l.list[:10] {
		l.df.AddRow([]string{v.name, v.common, strconv.Itoa(v.cancer), strconv.Itoa(v.typ)})
	}
}

func LeaderBoard(db *dbIO.DBIO, typ string) *dataframe.Dataframe {
	// Returns top cancer locations and the top species and types associated with them.
	l := newLeaderBoard(db, typ)
	l.countRecords()
	l.sortRecords()
	return l.df
}

// Returns top cancer locations and the top species and types associated with them.

package search

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/simpleset"
	"log"
	"strconv"
	"strings"
)

type location struct {
	name 	string
	species	map[string]int
	total	int
	types	map[string]int
}

func newLocation(loc string) *location {
	// Returns new location counter
	l := new(location)
	l.name = loc
	l.species = make(map[string]int)
	l.types = make(map[string]int)
	return l
}

func (l *location) add(species, typ string) {
	// Adds species and type to location counter
	l.total++
	if species != "NA" {
		if _, ex := l.species[species]; !ex {
			l.species[species] = 0
		}
		l.species[species]++
	}
	if typ != "NA" {
		if _, ex := l.types[typ]; !ex {
			l.types[typ] = 0
		}
		l.types[typ]++
	}
}

func (l *location) toSlice() []string {
	// Returns struct as string slice
	var typ, sp string
	var tcount, scount int
	// Get Most common type
	for k, v := range l.types {
		if v > tcount {
			tcount = v
			typ = k
		}
	}
	// Get Most common species
	for k, v := range l.species {
		if v > scount {
			scount = v
			sp = k
		}
	}
	return []string{l.name, strconv.Itoa(l.total), typ, strconv.Itoa(tcount), sp, strconv.Itoa(scount)}
}

//--------------------------------------------------------------------------------------

type leaderboard struct {
	df			*dataframe.Dataframe
	locations	map[string]*location
	logger  	*log.Logger
	min         int
	table		*dataframe.Dataframe
	taxa        *simpleset.Set
	top			[]string
}

func newLeaderBoard(db *dbIO.DBIO, min int) *leaderboard {
	// Initializes new struct
	l := new(leaderboard)
	l.df, _ = dataframe.NewDataFrame(-1)
	l.df.SetHeader([]string{"Location", "LocationTotal", "TopType", "TypeTotal", "TopSpecies", "SpeciesTotal"})
	l.locations = make(map[string]*location)
	l.logger = codbutils.GetLogger()
	l.min = min
	l.table, _ = SearchRecords(db, l.logger, "Approved=1", false, false)
	l.taxa = simpleset.NewStringSet()
	l.top = make([]string, 5)
	return l
}

func (l *leaderboard) getUnique(loc, typ string) [][]string {
	// Returns unique type and location pairs
	var ret [][]string
	d := ";"
	set := simpleset.NewStringSet()
	types := strings.Split(typ, d)
	// Get unique pairs
	for idx, i := range strings.Split(loc, d) {
		if i != "NA" {
			set.Add(strings.Join([]string{i, types[idx]}, d))
		}
	}
	for _, i := range set.ToStringSlice() {
		ret = append(ret, strings.Split(i, d))
	}
	return ret
}

func (l *leaderboard) minTaxa() {
	// Identifies species with more than min records
	count := make(map[string]int)
	for idx := range l.table.Rows {
		sp, _ := l.table.GetCell(idx, "Species")
		if _, ex := count[sp]; !ex {
			count[sp] = 0
		}
		count[sp]++
	}
	for k, v := range count {
		if v >= l.min {
			l.taxa.Add(k)
		}
	}
}

func (l *leaderboard) countRecords() {
	// Counts tissue types
	l.logger.Print("Counting neoplasia records...")
	for idx := range l.table.Rows {
		if mp, _ := l.table.GetCell(idx, "Masspresent"); mp == "1" {
			sp, _ := l.table.GetCell(idx, "Species")
			if ex, _ := l.taxa.InSet(sp); ex {
				loc, _ := l.table.GetCell(idx, "Location")
				typ, _ := l.table.GetCell(idx, "Type")
				for _, i := range l.getUnique(loc, typ) {
					if _, ex := l.locations[i[0]]; !ex {
						l.locations[i[0]] = newLocation(i[0])
					}
					l.locations[i[0]].add(sp, i[1])
				}
			}
		}
	}
}

func (l *leaderboard) sort(s string) {
	// Adds name to list if in top 5
	for idx, i := range l.top {
		if i == "" {
			l.top[idx] = s
		} else {
			if l.locations[s].total > l.locations[i].total {
				for j := len(l.top) - 1; j > idx; j-- {
					// Shift lower entries to the right
					l.top[j] = l.top[j-1]
				}
				// Insert entry in order
				l.top[idx] = s
				break
			}
		}
	}
}

func (l *leaderboard) sortRecords() {
	// Stores adult tumor records
	l.logger.Print("Sorting neoplasia records...")
	for k := range l.locations {
		l.sort(k)
	}
	for _, v := range l.top {
		l.df.AddRow(l.locations[v].toSlice())
	}
}

func LeaderBoard(db *dbIO.DBIO, min int) *dataframe.Dataframe {
	// Returns top cancer locations and the top species and types associated with them.
	l := newLeaderBoard(db, min)
	l.minTaxa()
	l.countRecords()
	l.sortRecords()
	return l.df
}

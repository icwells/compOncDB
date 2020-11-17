// Tests searcher functions

package search

import (
	"flag"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/simpleset"
	"testing"
)

var (
	password     = flag.String("password", "", "MySQL password.")
	user         = flag.String("user", "", "MySQL username.")
)

func getTestSearcher() *searcher {
	// Returns database connection
	flag.Parse()
	c := codbutils.SetConfiguration(*user, false)
	db, _ := dbIO.Connect(c.Host, c.Database, c.User, *password)
	return newSearcher(db, codbutils.GetLogger(), false)
}

func setIDs(ids [][]string) *simpleset.Set {
	ret := simpleset.NewStringSet()
	ids = append(ids, []string{"16", "17", "18"})
	for _, row := range ids {
		for _, i := range row {
			ret.Add(i)
		}
	}
	return ret
}

func filterIDs(ids *simpleset.Set, match []string) *simpleset.Set {
	ret := simpleset.NewStringSet()
	for _, i := range match {
		if ex, _ := ids.InSet(i); ex {
			ret.Add(i)
		}
	}
	return ret
}

func TestFilterIDs(t *testing.T) {
	// tests filter ids algorithm
	s := getTestSearcher()
	input := [][]string {
		{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"},
		{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"},
		{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
		{"1", "2", "3", "4", "5", "6", "7", "8"},
		{"1", "2", "3", "4"},
		{"1"},
	}
	s.ids = setIDs(input)
	for idx, row := range input {
		act := filterIDs(s.ids, row)
		if len(row) != act.Length() {
			t.Errorf("%d: Set length %d does not equal expected: %d", idx, act.Length(), len(row))
		} else {
			for _, i := range row {
				if ex, _ := act.InSet(i); !ex {
					t.Errorf("%d: %s not found in ids set.", idx, i)
					break
				}
			}
		}
		s.ids = act
	}
}

// Tests searcher functions

package search

import (
	"flag"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	//"github.com/icwells/simpleset"
	"strings"
	"testing"
)

var (
	password = flag.String("password", "", "MySQL password.")
	user     = flag.String("user", "", "MySQL username.")
)

func getTestConnection() *dbIO.DBIO {
	// Returns database connection
	flag.Parse()
	c := codbutils.SetConfiguration(*user, false)
	db, _ := dbIO.Connect(c.Host, c.Database, c.User, *password)
	db.GetTableColumns()
	return db
}

func TestLocations(t *testing.T) {
	// Tests location results when multiple variables are submitted
	db := getTestConnection()
	input := []string{"Location=uterus, Sex=male", "Location=ovary, Sex=male", "Location=testis, Sex=female", "Location=mammary, Class=Reptilia"}
	for _, i := range input {
		var count int
		eval := codbutils.SetOperations(db.Columns, i)
		act, _ := SearchRecords(db, codbutils.GetLogger(), i, false, false)
		for key := range act.Rows {
			e := eval[0]
			if a, _ := act.GetCell(key, e[0].Column); strings.Contains(a, e[0].Value) {
				count++
			} else if a, _ := act.GetCell(key, e[1].Column); a != e[1].Value {
				count++
			}
		}
		if count > 0 {
			e := eval[0]
			t.Errorf("Found %d records where %s does not equal %s or %s does not equal %s.", count, e[0].Column, e[0].Value, e[1].Column, e[1].Value)
			break
		}
	}
}

/*func setIDs(ids [][]string) *simpleset.Set {
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
	s := newSearcher(getTestConnection(), codbutils.GetLogger())
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
}*/

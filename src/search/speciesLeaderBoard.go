// Returns caancer type frequency by species

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

type tumorTypes struct {
	counts map[string]int
	types  []string
}

func newTumorTypes(db *dbIO.DBIO) *tumorTypes {
	// Returns initialized struct
	t := new(tumorTypes)
	t.count = make(map[string]int)
	for _, i := range db.GetColumn("Records", "Type") {
		if _, ex := t.count[i]; !ex {
			t.count[i] = 0
		}
		t.count[i]++
	}
	return t
}

func SpeciesLeaderBoard(db *dbIO.DBIO, min int) {
	// Returns caancer type frequency by species

}

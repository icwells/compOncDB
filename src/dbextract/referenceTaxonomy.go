// Returns merged common and taxonomy tables

package dbextract

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/simpleset"
	"strings"
)

type taxaMerger struct {
	header []string
	taxa   map[string][]string
	common [][]string
	com    *simpleset.Set
}

func newTaxaMerger(db *dbIO.DBIO) *taxaMerger {
	// Initializes and populates struct fields
	t := new(taxaMerger)
	t.header = strings.Split(db.Columns["Taxonomy"], ",")
	t.header[0] = "Common"
	t.taxa = codbutils.ToMap(db.GetTable("Taxonomy"))
	t.common = db.GetTable("Common")
	t.com = simpleset.NewStringSet()
	return t
}

func (t *taxaMerger) merge() *dataframe.Dataframe {
	// Merges common and taxonomy tables
	ret, _ := dataframe.NewDataFrame(-1)
	ret.SetHeader(t.header)
	for _, i := range t.common {
		tid := i[0]
		name := i[1]
		if taxa, ex := t.taxa[tid]; ex {
			ret.AddRow(append([]string{name}, taxa...))
			t.com.Add(tid)
		}
	}
	// Add any taxonomies without common names
	for k, v := range t.taxa {
		if ex, _ := t.com.InSet(k); !ex {
			ret.AddRow(append([]string{""}, v...))
		}
	}
	return ret
}

func GetReferenceTaxonomy(db *dbIO.DBIO) *dataframe.Dataframe {
	// Returns merged common and taxonomy tables
	t := newTaxaMerger(db)
	return t.merge()
}

// Returns merged common and taxonomy tables

package dbextract

import (
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/simpleset"
	"strings"
)

type taxaMerger struct {
	header []string
	taxa   map[string][]string
	common map[string][]string
	com    *simpleset.Set
}

func newTaxaMerger(db *dbIO.DBIO) *taxaMerger {
	// Initializes and populates struct fields
	t := new(taxaMerger)
	col := db.Columns["Taxonomy"]
	t.header = strings.Split(col, ",")
	t.header[0] = "Common"
	t.taxa = dbupload.ToMap(db.GetTable("Taxonomy"))
	t.common = dbupload.ToMap(db.GetTable("Common"))
	t.com = simpleset.NewStringSet()
	return t
}

func (t *taxaMerger) merge() *dataframe.Dataframe {
	// Merges common and taxonomy tables
	ret, _ := dataframe.NewDataFrame(-1)
	ret.SetHeader(t.header)
	for k, v := range t.common {
		t.com.Add(k)
		if taxa, ex := t.taxa[k]; ex {
			ret.AddRow(append([]string{v[0]}, taxa...))
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

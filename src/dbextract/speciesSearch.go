// Finds taxonomy matches for external species names

package dbextract

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/strarray"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"log"
	"sort"
)

type speciesSearcher struct {
	found   int
	list    []string
	logger  *log.Logger
	species map[string]string
	taxa    map[string][]string
}

func newSpeciesSearcher(db *dbIO.DBIO) speciesSearcher {
	// Initializes new searcher
	var s speciesSearcher
	// Get map of scientific and common species names
	s.logger = codbutils.GetLogger()
	s.species = dbupload.GetTaxaIDs(db, true)
	s.taxa = db.GetTableMap("Taxonomy")
	for k := range s.species {
		s.list = append(s.list, k)
	}
	return s
}

func (s *speciesSearcher) getTaxonomy(ch chan []string, n string) {
	// Attempts to find match for input name
	var ret []string
	k := strarray.TitleCase(n)
	id, ex := s.species[k]
	if ex == false {
		// Attempt fuzzy search if there is no literal match
		matches := fuzzy.RankFindFold(k, s.list)
		if matches.Len() > 0 {
			sort.Sort(matches)
			if matches[0].Distance <= 2 {
				k = matches[0].Target
				id = s.species[k]
				ex = true
			}
		}
	}
	if ex == true {
		ret = []string{n, k}
		ret = append(ret, s.taxa[id]...)
		s.found++
	}
	ch <- ret
}

func SearchSpeciesNames(db *dbIO.DBIO, names []string) *dataframe.Dataframe {
	// Finds taxonomies for input terms
	ch := make(chan []string)
	s := newSpeciesSearcher(db)
	s.logger.Println("Searching for taxonomy matches...")
	ret, _ := dataframe.NewDataFrame(-1)
	ret.SetHeader(append([]string{"Term", "MatchedName"}, db.Columns["Taxonomy"][1:]))
	for _, i := range names {
		go s.getTaxonomy(ch, i)
		row := <-ch
		if len(row) > 0 {
			ret.AddRow(row)
		}
	}
	s.logger.Printf("Found taxonomy matches for %d of %d queries.\n", s.found, len(names))
	return ret
}

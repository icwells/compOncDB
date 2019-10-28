// Finds taxonomy matches for external species names

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/strarray"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
)

type speciesSearcher struct {
	species map[string]string
	list    []string
	taxa    map[string][]string
	found   int
}

func newSpeciesSearcher(db *dbIO.DBIO) speciesSearcher {
	// Initializes new searcher
	var s speciesSearcher
	// Get map of scientific and common species names
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
	ret := dataframe.NewDataFrame(-1)
	ch := make(chan []string)
	ret.SetHeader([]string{"Term", "MatchedName", "Kingdom", "Phylum", "Class", "Order", "Family", "Genus", "Species", "Source"})
	fmt.Print("\n\tSearching for taxonomy matches...\n")
	s := newSpeciesSearcher(db)
	for _, i := range names {
		go s.getTaxonomy(ch, i)
		row := <-ch
		if len(row) > 0 {
			ret.AddRow(row)
		}
	}
	fmt.Printf("\tFound taxonomy matches for %d of %d queries.\n", s.found, len(names))
	return ret
}

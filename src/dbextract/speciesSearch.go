// Finds taxonomy matches for external species names

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
	"strings"
)

func titleCase(t string) string {
	// Manually converts term to title case (strings.Title is buggy)
	var query []string
	s := strings.Split(t, " ")
	for _, i := range s {
		if len(i) > 1 {
			// Skip stray characters
			query = append(query, strings.ToUpper(string(i[0]))+strings.ToLower(i[1:]))
		}
	}
	return strings.Join(query, " ")
}

type speciesSearcher struct {
	species map[string]string
	list    []string
	taxa    map[string][]string
	found	int
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
	id, ex := s.species[n]
	if ex == true {
		ret = s.taxa[id]
	} else {
		// Attempt fuzzy search if there is no literal match
		matches := fuzzy.RankFindFold(n, s.list)
		if len(matches) > 0 {
			sort.Sort(matches)
			if matches[0].Distance <= 3 {
				ret = s.taxa[s.species[matches[0].Target]]
				found++
			}
		}
	}
	ch <- ret
}

func SearchSpeciesNames(db *dbIO.DBIO, names []string) ([][]string, string) {
	// Finds taxonomies for input terms
	var ret [][]string
	ch := make(chan []string)
	header := "Term,MatchedName,Kingdom,Phylum,Class,Order,Family,Genus,Species,Source"
	fmt.Print("\n\tSearching for taxonomy matches...\n")
	s := newSpeciesSearcher(db)
	for _, i := range names {
		go s.getTaxonomy(ch, i)
		row := <-ch
		ret = append(ret, row)
	}
	fmt.Printf("\tFound taxonomy matches for %d of %d queries.\n", s.found, len(names))
	return ret, header
}

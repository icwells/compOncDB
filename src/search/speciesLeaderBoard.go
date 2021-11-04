// Returns caancer type frequency by species

package search

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"log"
	"sort"
	"strconv"
)

type taxaTypes struct {
	cancer   int
	id       string
	taxonomy []string
	total    int
	types    map[string]int
}

func newTaxaTypes(tid string, total int, taxonomy []string, types map[string]int) *taxaTypes {
	// Returns initialized struct
	t := new(taxaTypes)
	t.id = tid
	t.taxonomy = taxonomy
	t.total = total
	t.types = make(map[string]int)
	for k := range types {
		t.types[k] = 0
	}
	return t
}

func (t *taxaTypes) incrementType(typ string) {
	// Increments matching type counter
	t.cancer++
	if _, ex := t.types[typ]; ex {
		t.types[typ]++
	}
}

func (t *taxaTypes) toSlice(types []string) []string {
	// Returns string formatted for printing
	ret := []string{t.id}
	ret = append(ret, t.taxonomy...)
	ret = append(ret, strconv.Itoa(t.total))
	ret = append(ret, strconv.Itoa(t.cancer))
	ret = append(ret, strconv.FormatFloat(float64(t.cancer) / float64(t.total), 'f', -1, 64))
	for _, i := range types {
		ret = append(ret, strconv.FormatFloat(float64(t.types[i]) / float64(t.cancer), 'f', -1, 64))
	}
	return ret
}

//----------------------------------------------------------------------------

type speciesBoard struct {
	header []string
	list   []*taxaTypes
	logger *log.Logger
	min    int
	sorted []string
	table  *dataframe.Dataframe
	taxa   map[string]*taxaTypes
	types  map[string]int
}

func newSpeciesBoard(db *dbIO.DBIO, min int) *speciesBoard {
	// Returns new struct
	var msg string
	s := new(speciesBoard)
	s.logger = codbutils.GetLogger()
	s.min = min
	s.logger.Printf("Calculating tumor type frequency for species with at least %d entries...\n", s.min)
	// Infant and life history = false
	s.table, msg = SearchRecords(db, s.logger, "Approved=1,Necropsy=1", false, false)
	s.logger.Println(msg)
	s.getTumorTypes()
	s.setTaxa()
	s.setHeader()
	return s
}

func (s *speciesBoard) toDF() *dataframe.Dataframe {
	// Formats types into dataframe
	ret, _ := dataframe.NewDataFrame(-1)
	ret.SetHeader(s.header)
	for _, i := range s.list {
		if err := ret.AddRow(i.toSlice(s.sorted)); err != nil {
			panic(err)
		}
	}
	return ret
}

func (s *speciesBoard) Len() int {
	return len(s.taxa)
}

func (s *speciesBoard) Less(i, j int) bool {
	return s.list[i].cancer > s.list[j].cancer
}

func (s *speciesBoard) Swap(i, j int) {
	s.list[i], s.list[j] = s.list[j], s.list[i]
}

func (s *speciesBoard) sort() {
	// Sorts records slice by number of records
	s.logger.Println("Sorting results...")
	for _, v := range s.taxa {
		s.list = append(s.list, v)
	}
	sort.Sort(s)
}


func (s *speciesBoard) countSpeciesTypes() {
	// Counts tumor types per species
	s.logger.Println("Determining neoplasia type frequencies...")
	for idx := range s.table.Rows {
		tid, _ := s.table.GetCell(idx, "taxa_id")
		if v, ex := s.taxa[tid]; ex {
			i, _ := s.table.GetCell(idx, "Type")
			v.incrementType(i)
		}
	}
}

func (s *speciesBoard) setHeader() {
	// Sets header for output file
	h := codbutils.NewHeaders()
	s.header = h.Taxonomy[:len(h.Taxonomy) - 1]
	s.header = append(s.header, []string{"TotalRecords", "NeoplasiaRecords", "NeoplasiaPrevalence"}...)
	for k := range s.types {
		s.sorted = append(s.sorted, k)
	}
	// Sort types by most common first
	sort := true
	for sort {
		sort = false
		for idx, i := range s.sorted[:len(s.sorted) - 1] {
			v := s.sorted[idx + 1]
			if s.types[i] < s.types[v] {
				s.sorted[idx], s.sorted[idx + 1] = v, i
			}
		}
	}
	s.header = append(s.header, s.sorted...)
}

func (s *speciesBoard) getTumorTypes() {
	// Returns initialized struct
	s.types = make(map[string]int)
	for idx := range s.table.Rows {
		i, _ := s.table.GetCell(idx, "Type")
		if i != "NA" {
			if _, ex := s.types[i]; !ex {
				s.types[i] = 0
			}
			s.types[i]++
		}
	}
}

func (s *speciesBoard) setTaxa() {
	// Identifies species with more than min records
	s.taxa = make(map[string]*taxaTypes)
	count := make(map[string]int)
	taxa := make(map[string][]string)
	for idx := range s.table.Rows {
		tid, _ := s.table.GetCell(idx, "taxa_id")
		if _, ex := count[tid]; !ex {
			count[tid] = 0
			for _, i := range []string{"Kingdom", "Phylum", "Class", "Orders", "Family", "Genus", "Species", "common_name"} {
				v, _ := s.table.GetCell(idx, i)
				taxa[tid] = append(taxa[tid], v)
			}
		}
		count[tid]++
	}
	for k, v := range count {
		if v >= s.min {
			s.taxa[k] = newTaxaTypes(k, v, taxa[k], s.types)
		}
	}
	s.logger.Printf("Found %d species with at least %d entries.\n", len(s.taxa), s.min)
}

func SpeciesLeaderBoard(db *dbIO.DBIO, min int) *dataframe.Dataframe {
	// Returns caancer type frequency by species
	s := newSpeciesBoard(db, min)
	s.countSpeciesTypes()
	s.sort()
	return s.toDF()
}

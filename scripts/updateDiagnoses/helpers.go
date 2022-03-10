// Contains helper structs for diagnosis checker

package main

import (
	"github.com/icwells/simpleset"
	"strconv"
)

var LOCATIONS = []string{"abdomen","adrenal cortex","adrenal medulla","bile duct","bladder","blood","bone","bone marrow","brain","carotid body","cartilage","colon","dendritic cell","duodenum","esophagus","fat","fibrous","gall bladder","gland","glial cell","hair follicle","heart","iris","kidney","larynx","liver","lung","lymph nodes","mammary","mast cell","meninges","mesothelium","myxomatous tissue","NA","nerve cell","neuroendocrine","neuroepithelial","nose","notochord","oral","ovary","oviduct","pancreas","parathyroid gland","peripheral nerve sheath","pigment cell","pituitary gland","pnet","prostate","pupil","skin","small intestine","smooth muscle","spinal cord","spleen","stomach","striated muscle","synovium","testis","thyroid","trachea","transitional epithelium","uterus","vulva","widespread"}

type columns struct {
	comments    string
	common      string
	hyperplasia string
	id          string
	location    string
	malignant   string
	masspresent string
	service     string
	sex         string
	species     string
	tid         string
	tissue      string
	typ         string
}

func newColumns() *columns {
	// Returns initialized struct
	c := new(columns)
	c.comments = "Comments"
	c.common = "common_name"
	c.hyperplasia = "Hyperplasia"
	c.id = "ID"
	c.location = "Location"
	c.malignant = "Malignant"
	c.masspresent = "Masspresent"
	c.service = "service_name"
	c.sex = "Sex"
	c.species = "Species"
	c.tid = "taxa_id"
	c.tissue = "Tissue"
	c.typ = "Type"
	return c
}

type species struct {
	common  string
	id      string
	name    string
	novel   int
	updated int
}

func newSpecies(id, name, common string) *species {
	// Returns initialized species struct
	s := new(species)
	s.common = common
	s.id = id
	s.name = name
	return s
}

func (s *species) addNovel() {
	// Increments novel counter
	s.novel++
}

func (s *species) addUpdated() {
	// Increments updated counter
	s.updated++
}

func (s *species) toSlice() []string {
	// Returns struct formatted as string slice
	ret := []string{s.id, s.name, s.common}
	ret = append(ret, strconv.Itoa(s.updated))
	ret = append(ret, strconv.Itoa(s.novel))
	return ret
}

func (l *lzDiagnosis) speciesSlice() [][]string {
	// Returns species formatted as slice
	var ret [][]string
	for _, v := range l.taxa {
		if v.updated > 0 || v.novel > 0 {
			ret = append(ret, v.toSlice())
		}
	}
	return ret
}

func (l *lzDiagnosis) setLocations() {
	// Stores locations in set
	l.locations = simpleset.NewStringSet()
	for _, i := range LOCATIONS {
		l.locations.Add(i)
	}
}

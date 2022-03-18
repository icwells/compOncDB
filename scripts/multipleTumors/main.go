// Deifnes struct and methods for examining multiple tumor hits

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/diagnoses"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/simpleset"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"strings"
	"time"
)

var (
	D = ";"
	LOCATIONS = []string{"abdomen","adrenal cortex","adrenal medulla","bile duct","bladder","blood","bone","bone marrow","brain","carotid body","cartilage","colon","dendritic cell","duodenum","esophagus","fat","fibrous","gall bladder","gland","glial cell","hair follicle","heart","iris","kidney","larynx","liver","lung","lymph nodes","mammary","mast cell","meninges","mesothelium","myxomatous tissue","NA","nerve cell","neuroendocrine","neuroepithelial","nose","notochord","oral","ovary","oviduct","pancreas","parathyroid gland","peripheral nerve sheath","pigment cell","pituitary gland","pnet","prostate","pupil","skin","small intestine","smooth muscle","spinal cord","spleen","stomach","striated muscle","synovium","testis","thyroid","trachea","transitional epithelium","uterus","vulva","widespread"}
	outfile  = kingpin.Flag("outfile", "Optional path to output file. Prints proposed changes to file instead of updating database.").Short('o').Default("").String()
	user     = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type tumor struct {
	location	string
	malignant	string
	tissue		string
	typ			string
}

func newTumor(location, malignant, tissue, typ string) *tumor {
	// Returns initialized struct for inidividual tumor
	t := new(tumor)
	t.location = location
	t.malignant = malignant
	t.tissue = tissue
	t.typ = typ
	return t
}

func (t *tumor) equals(v *tumor) bool {
	// Returns true if all attributes are equal
	if t.location == v.location && t.tissue == v.tissue && t.typ == v.typ {
		return true
	}
	return false
}

type diagnosis struct {
	hyperplasia	string
	id			string
	masspresent	string
	primary		string
	tumors		[]*tumor
}

func newDiagnosis(id, hyperplasia, masspresent, primary string) *diagnosis {
	// Returns initialized struct for record's diagnosis
	d := new(diagnosis)
	d.id = id
	d.hyperplasia = hyperplasia
	d.masspresent = masspresent
	d.primary = primary
	return d
}

func (d *diagnosis) addTumors(location, malignant, tissue, typ string) {
	// Splits diagnoses and adds unique tumor info
	l := strings.Split(location, D)
	t := strings.Split(tissue, D)
	for idx, i := range strings.Split(typ, D) {
		d.tumors = append(d.tumors, newTumor(l[idx], malignant, t[idx], i))
	}
}

func (d *diagnosis) equals(v *diagnosis) bool {
	// Returns true if both records contain some tumors
	if len(d.tumors) == len(v.tumors) {
		for idx, i := range d.tumors {
			if !i.equals(v.tumors[idx]) {
				return false
			}
		}
		return true
	}
	return false
}

//----------------------------------------------------------------------------

type multipleTumors struct {
	db          *dbIO.DBIO
	hyperplasia int
	locations   *simpleset.Set
	logger      *log.Logger
	match       diagnoses.Matcher
	metastasis	string
	neoplasia   int
	records     map[string]*diagnosis
	summary     [][]string
	update      bool
}

func newMultipleTumors() *multipleTumors {
	// Return new struct
	m := new(multipleTumors)
	m.logger = codbutils.GetLogger()
	m.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	m.logger.Println("Initializing struct...")
	m.match = diagnoses.NewMatcher(m.logger)
	m.records = make(map[string]*diagnosis)
	if *outfile == "" {
		m.update = true
	}
	m.setLocations()
	return m
}

func (m *multipleTumors) setLocations() {
	// Stores locations in set
	m.locations = simpleset.NewStringSet()
	for _, i := range LOCATIONS {
		m.locations.Add(i)
	}
}

func (m *multipleTumors) setRecords() {
	// Stores records in diagnosis struct
	records, msg := search.SearchRecords(m.db, m.logger, "Comments!=NA", true, false)
	m.logger.Println(msg)
	for i := range records.Iterate() {
		typ, _ := i.GetCell("Type")
		if strings.Contains(typ, D) {
			id, _ := i.GetCell("ID")
			hyp, _ := i.GetCell("Hyperplasia")
			mal, _ := i.GetCell("Malignant")
			mass, _ := i.GetCell("Masspresent")
			prim, _ := i.GetCell("Primary")
			loc, _ := i.GetCell("Location")
			tissue, _ := i.GetCell("Tissue")
			// Make new entry and add tumor
			m.records[id] = newDiagnosis(id, hyp, mass, prim)
			m.records[id].addTumors(loc, mal, tissue, typ)
		}
	}
}

/*func (m *multipleTumors) write() {
	// Writes summary to file if outfile is given
	if !m.update {
		header := "Species,Comments,Masspresent,ProposedMP,Hyperplasia,ProposedHyp,Type,ProposedType,Tissue,ProposedTissue,Location,ProposedLoc"
		iotools.WriteToCSV(*outfile, header, m.summary)
	}
}*/

func main() {
	start := time.Now()
	kingpin.Parse()
	m := newMultipleTumors()
	m.setRecords()
	m.checkRecords()
	m.write()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

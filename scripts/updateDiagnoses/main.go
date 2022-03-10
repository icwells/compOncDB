// Double checks London Zoo diagnoses

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
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	outfile = kingpin.Flag("outfile", "Optional path to output file. Prints proposed changes to file instead of updating database.").Short('o').Default("").String()
	user    = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type lzDiagnosis struct {
	col         *columns
	db          *dbIO.DBIO
	diag        *dataframe.Dataframe
	hyperplasia int
	locations   *simpleset.Set
	logger      *log.Logger
	match       diagnoses.Matcher
	neoplasia   int
	summary     [][]string
	tables      map[string]string
	taxa        map[string]*species
	update      bool
}

func newLZDiagnosis() *lzDiagnosis {
	// Return new struct
	var msg string
	diag := "Diagnosis"
	tumor := "Tumor"
	l := new(lzDiagnosis)
	l.logger = codbutils.GetLogger()
	l.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	l.logger.Println("Initializing struct...")
	l.col = newColumns()
	l.diag, msg = search.SearchRecords(l.db, l.logger, "Comments!=NA", true, false)
	l.logger.Println(msg)
	l.match = diagnoses.NewMatcher(l.logger)
	l.tables = map[string]string{
		l.col.hyperplasia: diag,
		l.col.location:    tumor,
		l.col.malignant:   tumor,
		l.col.masspresent: diag,
		l.col.tissue:      tumor,
		l.col.typ:         tumor,
	}
	l.taxa = make(map[string]*species)
	if *outfile == "" {
		l.update = true
	}
	l.setLocations()
	return l
}

func (l *lzDiagnosis) updateCell(column, id, val string) {
	// Updates individual attribute of record
	if l.update {
		l.db.UpdateRow(l.tables[column], column, val, l.col.id, "=", id)
	}
}

func (l *lzDiagnosis) getComments(i *dataframe.Series) string {
	// Returns relevant comments for parsing
	ret, _ := i.GetCell(l.col.comments)
	if service, _ := i.GetCell(l.col.service); service == "NWZP" {
		for idx, i := range ret {
			if unicode.IsLetter(i) && !unicode.IsUpper(i) {
				if idx >= 4 {
					return ret[:idx-1]
				} else {
					break
				}
			}
		}
	}
	return ret
}

func (l *lzDiagnosis) getTaxaID(i *dataframe.Series) string {
	// Returns taxa_id and initializes new species struct if needed
	ret, _ := i.GetCell(l.col.tid)
	if _, ex := l.taxa[ret]; !ex {
		com, _ := i.GetCell(l.col.common)
		sp, _ := i.GetCell(l.col.species)
		l.taxa[ret] = newSpecies(ret, sp, com)
	}
	return ret
}

func (l *lzDiagnosis) checkRecord(i *dataframe.Series) bool {
	// Checks individual record and updates if needed
	var ret bool
	var hyperplasia int
	neoplasia := 1
	if typ, _ := i.GetCell(l.col.typ); !strings.Contains(typ, ";") {
		comments := l.getComments(i)
		sex, _ := i.GetCell(l.col.sex)
		if tumor, tissue, location, malignant := l.match.GetTumor(comments, sex, true); !strings.Contains(tumor, ";") {
			tid := l.getTaxaID(i)
			hyp, _ := i.GetCellInt(l.col.hyperplasia)
			loc, _ := i.GetCell(l.col.location)
			mal, _ := i.GetCell(l.col.malignant)
			mp, _ := i.GetCellInt(l.col.masspresent)
			tis, _ := i.GetCell(l.col.tissue)
			if tumor != "NA" {
				if tumor == "hyperplasia" {
					hyperplasia = 1
					neoplasia = 0
					l.hyperplasia++
				}
				if tumor != typ {
					l.updateCell(l.col.typ, i.Name, tumor)
					ret = true
				}
				if hyp != hyperplasia {
					l.updateCell(l.col.hyperplasia, i.Name, strconv.Itoa(hyperplasia))
					ret = true
				}
				if mp != neoplasia {
					l.updateCell(l.col.masspresent, i.Name, strconv.Itoa(neoplasia))
					ret = true
					if mp == 0 {
						l.neoplasia++
					}
				}
			}
			if location != "NA" && location != loc {
				if ex, _ := l.locations.InSet(location); ex {
					l.updateCell(l.col.location, i.Name, location)
					ret = true
				} else {
					tissue = "NA"
				}
			}
			if tissue != "NA" && tissue != tis {
				l.updateCell(l.col.tissue, i.Name, tissue)
				ret = true
			}
			if malignant != "-1" && malignant != mal {
				l.updateCell(l.col.malignant, i.Name, malignant)
				ret = true
			}
			if ret {
				row := []string{l.taxa[tid].name, comments, strconv.Itoa(mp), strconv.Itoa(neoplasia), strconv.Itoa(hyp), strconv.Itoa(hyperplasia), typ, tumor, tis, tissue, loc, location}
				l.summary = append(l.summary, row)
				if mp == 0 && neoplasia != 0 {
					l.taxa[tid].addNovel()
				} else {
					l.taxa[tid].addUpdated()
				}
			}
		}
	}
	return ret
}

func (l *lzDiagnosis) checkRecords() {
	// Updates life history table with converted values
	var count, updated int
	l.logger.Println("Comparing diagnoses...")
	for i := range l.diag.Iterate() {
		count++
		if l.checkRecord(i) {
			updated++
		}
		fmt.Printf("\tChecked %d of %d terms.\r", count, l.diag.Length())
	}
	fmt.Println()
	l.logger.Printf("Updated %d records.", updated)
	l.logger.Printf("Found %d new neoplasia records and %d new hyperplasia records.", l.neoplasia, l.hyperplasia)
	l.logger.Printf("Updated %d existing neoplasia records.", updated-l.neoplasia-l.hyperplasia)
}

func (l *lzDiagnosis) write() {
	// Writes summary to file if outfile is given
	if !l.update {
		header := "Species,Comments,Masspresent,ProposedMP,Hyperplasia,ProposedHyp,Type,ProposedType,Tissue,ProposedTissue,Location,ProposedLoc"
		iotools.WriteToCSV(*outfile, header, l.summary)
		spfile := strings.Replace(*outfile, ".csv", ".Species.csv", 1)
		iotools.WriteToCSV(spfile, "taxa_id,Species,Common,Updated,New", l.speciesSlice())
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	l := newLZDiagnosis()
	l.checkRecords()
	l.write()
	l.logger.Printf("Finished. Runtime: %s", time.Since(start))
}

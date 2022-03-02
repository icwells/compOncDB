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

type columns struct {
	comments    string
	hyperplasia string
	id          string
	location    string
	malignant   string
	masspresent string
	service		string
	sex         string
	tissue      string
	typ         string
}

func newColumns() *columns {
	// Returns initialized struct
	c := new(columns)
	c.comments = "Comments"
	c.hyperplasia = "Hyperplasia"
	c.id = "ID"
	c.location = "Location"
	c.malignant = "Malignant"
	c.masspresent = "Masspresent"
	c.service = "service_name"
	c.sex = "Sex"
	c.tissue = "Tissue"
	c.typ = "Type"
	return c
}

type lzDiagnosis struct {
	col         *columns
	db          *dbIO.DBIO
	diag        *dataframe.Dataframe
	hyperplasia int
	logger      *log.Logger
	match       diagnoses.Matcher
	neoplasia   int
	summary     [][]string
	tables      map[string]string
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
		l.col.masspresent: diag,
		l.col.tissue:      tumor,
		l.col.typ:         tumor,
	}
	if *outfile == "" {
		l.update = true
	}
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
					return ret[:idx - 1]
				} else {
					break
				}
			}
		}
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
			if tissue != "NA" && tissue != tis {
				l.updateCell(l.col.tissue, i.Name, tissue)
				ret = true
			}
			if location != "NA" && location != loc {
				l.updateCell(l.col.location, i.Name, location)
				ret = true
			}
			if malignant != "-1" && malignant != mal {
				l.updateCell(l.col.malignant, i.Name, malignant)
				ret = true
			}
			if ret {
				row := []string{comments, strconv.Itoa(mp), strconv.Itoa(neoplasia), strconv.Itoa(hyp), strconv.Itoa(hyperplasia), typ, tumor, tis, tissue, loc, location}
				l.summary = append(l.summary, row)
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
		header := "Comments,Masspresent,ProposedMP,Hyperplasia,ProposedHyp,Type,ProposedType,Tissue,ProposedTissue,Location,ProposedLoc"
		iotools.WriteToCSV(*outfile, header, l.summary)
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

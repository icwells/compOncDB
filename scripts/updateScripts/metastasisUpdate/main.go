// Reruns metastasis regex on patient comments

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/diagnoses"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/dbIO"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"strings"
	"time"
)

var (
	D    = ";"
	user = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type metUpdater struct {
	db     *dbIO.DBIO
	logger *log.Logger
	match  diagnoses.Matcher
	met    string
	//records		map[string]*record
	table   string
	updates []string
}

func newMetUpdater() *metUpdater {
	// Initializes struct
	m := new(metUpdater)
	m.logger = codbutils.GetLogger()
	m.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	m.logger.Println("Initializing struct...")
	m.match = diagnoses.NewMatcher(m.logger)
	m.met = "Metastasis"
	m.table = "Diagnosis"
	return m
}

func (m *metUpdater) update() {
	// Updates diagnosis table
	m.logger.Println("Updating diagnosis table...")
	for idx, i := range m.updates {
		if !m.db.UpdateRow(m.table, m.met, "1", "ID", "=", i) {
			break
		}
		fmt.Printf("\tUpdated %d of %d records.\r", idx+1, len(m.updates))
	}
	fmt.Println()
}

func (m *metUpdater) checkLocations(loc string) bool {
	// Returns true if locations contains unique entries
	if strings.Contains(loc, D) {
		locs := strings.Split(loc, D)
		for _, i := range locs[:len(locs)-2] {
			for _, j := range locs[:len(locs)-1] {
				if i != j {
					return true
				}
			}
		}
		return false
	}
	return false
}

func (m *metUpdater) setRecords() {
	// Stores metastasis records
	diag, msg := search.SearchRecords(m.db, m.logger, "Masspresent=1", true, false)
	m.logger.Println(msg)
	for i := range diag.Iterate() {
		if met, _ := i.GetCell(m.met); met != "1" {
			comments, _ := i.GetCell("Comments")
			loc, _ := i.GetCell("Location")
			if match := m.match.BinaryMatch(m.match.Metastasis, comments); match == "1" {
				m.updates = append(m.updates, i.Name)
			} else if m.checkLocations(loc) {
				m.updates = append(m.updates, i.Name)
			}
		}
	}
	m.logger.Printf("Found %d records to update.", len(m.updates))
}

func main() {
	start := time.Now()
	kingpin.Parse()
	m := newMetUpdater()
	m.setRecords()
	m.update()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

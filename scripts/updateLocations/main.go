// Updates tumor locations

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/diagnoses"
	"github.com/icwells/dbIO"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"time"
)

var user = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()

type locations struct {
	db      *dbIO.DBIO
	ids     []string
	matcher diagnoses.Matcher
	records map[string]string
	tumor   map[string][][]string
	terms   []string
}

func newLocations() *locations {
	// Return new struct
	l := new(locations)
	l.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	fmt.Println("\n\tInitializing struct...")
	l.matcher = diagnoses.NewMatcher(codbutils.GetLogger())
	l.records = make(map[string]string)
	l.tumor = make(map[string][][]string)
	l.setTumor()
	return l
}

func (l *locations) setTumor() {
	// Returns map of all tumor entries per ID in 2d slice
	for _, row := range l.db.GetTable("Tumor") {
		id := row[0]
		if row[4] == "NA" {
			l.ids = append(l.ids, id)
		}
		if _, ex := l.tumor[id]; !ex {
			// Add new entry
			var rows [][]string
			l.tumor[id] = append(rows, row[3:])
		} else {
			l.tumor[id] = append(l.tumor[id], row[3:])
		}
	}
}

func (l *locations) update(id, orig, typ, loc string) {
	// Updates entry with new location
	var s string
	if orig == typ {
		s = fmt.Sprintf("UPDATE Tumor SET Location = '%s' WHERE ID = '%s' and Type = '%s' limit 1;", loc, id, orig)
	} else {
		// Update type as well
		s = fmt.Sprintf("UPDATE Tumor SET Location = '%s', Type = '%s' WHERE ID = '%s' and Type = '%s' limit 1;", loc, typ, id, orig)
	}
	cmd, err := l.db.DB.Prepare(s)
	if err != nil {
		panic(err)
	} else {
		_, err = cmd.Exec()
		cmd.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (l *locations) removeRow(id string, idx int) {
	// Removes row from tumor map entry
	rows := l.tumor[id]
	if len(rows) == 1 {
		delete(l.tumor, id)
	} else {
		if idx == 0 {
			rows = rows[1:]
		} else if idx == len(rows)-1 {
			rows = rows[:idx]
		} else {
			rows = append(rows[:idx], rows[idx+1:]...)
		}
	}
}

func (l *locations) checkTypes(id, typ, loc string) string {
	// Returns true if a location should be updated
	if rows, ex := l.tumor[id]; ex {
		for idx, i := range rows {
			if i[1] == "NA" {
				if i[0] == typ || i[0] == "NA" {
					l.removeRow(id, idx)
					return i[0]
				}
			}
		}
	}
	return ""
}

func (l *locations) checkLocations() {
	// Determines if records are infant records
	var count int
	fmt.Println("\tIdentifying tumor locations...")
	rows := l.db.GetRows("Patient", "ID", strings.Join(l.ids, ","), "ID,Sex,Comments")
	for rdx, row := range rows {
		fmt.Printf("\r\tProcessed %d of %d lines.", rdx+1, len(rows))
		id := row[0]
		comment := strings.ToLower(row[2])
		if typ, loc, _ := l.matcher.GetTumor(comment, row[1], true); loc != "NA" {
			types := strings.Split(typ, ";")
			for idx, i := range strings.Split(loc, ";") {
				if i != "NA" {
					t := types[idx]
					if orig := l.checkTypes(id, t, i); orig != "" {
						l.update(id, orig, t, i)
						count++
					}
				}
			}
		}
	}
	fmt.Printf("\n\tUpdated %d locations.\n", count)
}

func main() {
	start := time.Now()
	kingpin.Parse()
	l := newLocations()
	l.checkLocations()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

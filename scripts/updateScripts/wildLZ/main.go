// Sets London Zoo records to wild if they have an XT code in the deathbooks

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/simpleset"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	//"strconv"
	"strings"
	"time"
)

var (
	common     = kingpin.Flag("common", "Path to London Zoo common names file.").Short('c').Required().String()
	deathbooks = kingpin.Flag("deathbooks", "Path to London Zoo deathbooks file.").Short('d').Required().String()
	scientific = kingpin.Flag("scientific", "Path to London Zoo scientific names file.").Short('s').Required().String()
	user       = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type record struct {
	common	string
	date	string
	set		bool
	sex		string
	species	string
}

func newRecord(line []string) *record {
	// Returnsnew struct
	r := new(record)
	r.common = strings.TrimSpace(line[4])
	r.date = strings.TrimSpace(line[1])
	r.setSex(line[6])
	r.setSpecies(line[5])
	r.isSet()
	return r
}

func (r *record) isSet() {
	// Stores true if at least two fields are not empty
	var count int
	for _, i := range []string{r.common, r.date, r.sex, r.species} {
		if i != "" {
			count++
		}
	}
	if count > 2 {
		r.set = true
	}
}

func (r *record) setSpecies(species string) {
	// Formats species name
	r.species = strings.TrimSpace(species)
	sp := strings.Split(r.species, " ")
	if len(sp) > 1 {
		if sp[1] == "sp" || sp[1] == "sp." || sp[1] == "spp" {
			r.species = sp[0]
		}
	}
}

func (r *record) setSex(sex string) {
	// Stores record's sex
	sex = strings.ToLower(strings.TrimSpace(sex))
	if sex == "m" {
		r.sex = "male"
	} else if sex == "f" {
		r.sex = "female"
	} else {
		r.sex = "NA"
	}
}

func (r *record) equals(common, date, sex, species string) bool {
	// Returns true if records are equal
	if r.sex == sex && r.date == date {
		if r.common == common || r.species == species {
			return true
		}
	}
	return false
}

type wild struct {
	db      *dbIO.DBIO
	logger  *log.Logger
	pids    *simpleset.Set
	sids    *simpleset.Set
	records []*record
	target  string
}

func newWild() *wild {
	// Returns initialized converter struct
	w := new(wild)
	w.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	w.logger = codbutils.GetLogger()
	w.pids = simpleset.NewStringSet()
	w.sids = simpleset.NewStringSet()
	w.target = "XT"
	w.setRecords()
	return w
}

func (w *wild) setRecords() {
	// Stores source_id to ID map
	w.logger.Println("Reading deathbooks file...")
	reader, _ := iotools.YieldFile(*deathbooks, false)
	for i := range reader {
		if strings.TrimSpace(i[2]) == w.target {
			if r := newRecord(i); r.set {
				w.records = append(w.records, r)
			}
		}
	}
	w.logger.Printf("Found %d source records.", len(w.records))
}

func (w *wild) isWild(common, date, sex, species string) bool {
	// Compares record against wild records
	for _, i := range w.records {
		if i.equals(common, date, sex, species) {
			return true
		}
	}
	return false
}

func (w *wild) setSourceIDs() {
	// Stores source ids from london zoo upload files
	w.logger.Println("Getting source IDs...")
	for _, file := range []string{*scientific, *common} {
		reader, h := iotools.YieldFile(file, true)
		for i := range reader {
			if w.isWild(i[h["Name"]], i[h["Date"]], i[h["Sex"]], i[h["Species"]]) {
				w.sids.Add(i[h["ID"]])
			}
		}
	}
	w.logger.Printf("Found %d source IDs.", w.sids.Length())
}

func (w *wild) update() {
	// Updates wild values for target source_ids
	var count int
	w.logger.Println("Updating target records....")
	for _, i := range w.pids.ToStringSlice() {
		count++
		if !w.db.UpdateRow("Patient", "Wild", "1", "ID", "=", i) {
			panic(i)
		}
		fmt.Printf("\tUpdated %d of %d records.\r", count, w.pids.Length())
	}
	fmt.Println()
}

func (w *wild) setPatientIDs() {
	// Stores id for target LZ records
	w.logger.Println("Getting target patient IDs....")
	for _, i := range w.db.GetRows("Records", "service_name", "LZ", "ID,source_id,common_name") {
		if ex, _ := w.sids.InSet(i[1]); ex {
			w.pids.Add(i[0])
		}
	}
	w.logger.Printf("Found %d wild patient IDs.", w.pids.Length())
}

func main() {
	kingpin.Parse()
	w := newWild()
	w.setSourceIDs()
	w.setPatientIDs()
	w.update()
	w.logger.Printf("Finished. Runtime: %s\n\n", time.Since(w.db.Starttime))
}

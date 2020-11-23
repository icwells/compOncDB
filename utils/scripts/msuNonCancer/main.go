// Removes msu non cancer records

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/dbIO"
	"github.com/icwells/simpleset"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
)

var (
	user = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type trimmer struct {
	db  *dbIO.DBIO
	ids *simpleset.Set
}

func newTrimmer() *trimmer {
	// Returns initialized converter struct
	t := new(trimmer)
	t.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false))
	t.ids = simpleset.NewStringSet()
	t.setIDs()
	return t
}

func (t *trimmer) setIDs() {
	// Stores source_id to ID map
	fmt.Println("\n\tGetting non-cancer MSU IDs...")
	for _, i := range t.db.GetRows("Source", "service_name", "MSU", "ID") {
		t.ids.Add(i[0])
	}
	for _, i := range t.db.GetRows("Diagnosis", "Masspresent", "1", "ID") {
		// Remove cancer record IDs
		if ex, _ := t.ids.InSet(i[0]); ex {
			t.ids.Pop(i[0])
		}
	}
	fmt.Printf("\tFound %d non-cancer MSU IDs.\n", t.ids.Length())
}

func (t *trimmer) removeNonCancer() {
	// Removes non-cancer records
	fmt.Println("\tRemoving non-cancer records...")
	t.db.DeleteRows("Patient", "ID", t.ids.ToStringSlice())
	dbextract.AutoCleanDatabase(t.db)
	codbutils.UpdateTimeStamp(t.db)
}

func main() {
	start := time.Now()
	kingpin.Parse()
	t := newTrimmer()
	t.removeNonCancer()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

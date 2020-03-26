// Writes all database contents to gzipped tarball

package dbextract

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"os"
	"os/exec"
	"time"
)

type dbCompress struct {
	dir    string
	db     *dbIO.DBIO
	stamp  string
	tables []string
	temp   string
}

func newDbCompress(db *dbIO.DBIO, outdir string) *dbCompress {
	// Initializes new struct
	d := new(dbCompress)
	d.db = db
	d.dir, _ = iotools.FormatPath(outdir, false)
	d.temp, _ = iotools.FormatPath(outdir+"tmp/", true)
	d.setDateStamp()
	d.tables = []string{"Accounts", "Common", "Denominators", "Life_history", "Taxonomy", "Unmatched"}
	return d
}

func (d *dbCompress) setDateStamp() {
	// Stores date stamp
	d.stamp = time.Now().Format("2006-01-02")
}

func (d *dbCompress) compressDir() {
	// Compresses temp directory
	fmt.Println("\tCompressing output directory...")
	comp := exec.Command("tar", "-czf", fmt.Sprintf("%scomparativeOncology_%s.tar.gz", d.dir, d.stamp), d.temp)
	err := comp.Run()
	if err == nil {
		os.Remove(d.temp)
	} else {
		fmt.Printf("\n\t[Error] Failed to compress tables: %v", err)
	}
}

func (d *dbCompress) getOutfile(name string) string {
	// Returns formatted output file name
	return fmt.Sprintf("%s%s_%s.csv", d.temp, name, d.stamp)
}

func (d *dbCompress) writeTables() {
	// Writes tables to outdir
	fmt.Println("\n\tExtracting database tables...")
	df, _ := SearchColumns(d.db, "", codbutils.SetOperations(d.db.Columns, "ID > 0"), false, false)
	df.ToCSV(d.getOutfile("Records"))
	for _, i := range d.tables {
		table := d.db.GetTable(i)
		codbutils.WriteResults(d.getOutfile(i), d.db.Columns[i], table)
	}
}

func DumpDatabase(db *dbIO.DBIO, outdir string) {
	// Writes all database contents to gzipped tarball
	d := newDbCompress(db, outdir)
	d.writeTables()
	d.compressDir()
}

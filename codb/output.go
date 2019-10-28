// Contains functions for interacting with slq database

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"net/http"
	"strings"
	"time"
)

type Output struct {
	User    string
	Update  string
	Flash   string
	File    string
	Outfile string
	Count   string
	db      *dbIO.DBIO
}

func newOutput(user, pw, ut string) (*Output, error) {
	// Returns empty output struct
	var err error
	o := new(Output)
	o.User = user
	o.Update = strings.Replace(ut, "UTC", "Eastern Time", 1)
	if pw != "" {
		o.db, err = dbIO.Connect(C.config.Host, C.config.Database, o.User, pw)
		if err == nil {
			o.db.GetTableColumns()
		}
	}
	return o, err
}

func newFlash(msg string) *Output {
	// Returns output with flash error message
	o, _ := newOutput("", "", "")
	o.Flash = msg
	return o
}

func (o *Output) getTempFile(name string) {
	// Stores path to named file in tmp directory
	t := time.Now()
	fmt.Println(t)
	stamp := t.Format(time.RFC3339)
	// Trim timestamp to minutes
	stamp = stamp[:strings.LastIndex(stamp, "-")]
	stamp = stamp[:strings.LastIndex(stamp, ":")]
	o.File = fmt.Sprintf("%s.%s.csv", name, stamp)
	o.Outfile = fmt.Sprintf("/tmp/%s", o.File)
}

func (o *Output) summary(w http.ResponseWriter, password string) {
	// Returns general database summary
	o.getTempFile("databaseSummary")
	header := "Field,Total,%"
	codbutils.WriteResults(o.Outfile, header, dbextract.GetSummary(o.db))
	C.renderTemplate(w, C.temp.result, o)
}

func (o *Output) searchDB(f *SearchForm) {
	// Searches database for results
	res := dataframe.NewDataFrame(-1)
	// Search for column/value match
	for _, v := range f.eval {
		r := dbextract.SearchColumns(o.db, o.User, f.Table, v, f.Count, f.Infant)
		if res.Length() == 0 {
			res = r
		} else {
			// Append successive results to results slice
			for _, i := range r.Rows {
				res.Rows = append(res.Rows, i)
			}
		}
	}
	if f.Count == false && res.Length() >= 1 {
		o.getTempFile(o.User)
		res.ToCSV(o.Outfile)
	} else {
		o.Count = fmt.Sprintf("\tFound %d records matching search criteria.\n", res.Length())
	}
}

func (o *Output) extractFromDB(w http.ResponseWriter, r *http.Request, password string) {
	// Extracts data to outfile/stdout
	var f *SearchForm
	f, o.Flash = setSearchForm(r, o.db.Columns)
	if o.Flash == "" {
		if len(f.Table) > 0 {
			// Extract entire table
			table := o.db.GetTable(f.Table)
			o.getTempFile(o.User)
			codbutils.WriteResults(o.Outfile, o.db.Columns[f.Table], table)
		} else if f.Cancerrate == true {
			// Extract cancer rates
			var e []codbutils.Evaluation
			o.getTempFile(fmt.Sprintf("cancerRates.min%d", f.Min))
			for _, v := range f.eval {
				e = v
			}
			rates := dbextract.GetCancerRates(o.db, f.Min, f.Necropsy, e)
			rates.ToCSV(o.Outfile)
		} else {
			o.searchDB(f)
		}
		C.renderTemplate(w, C.temp.result, o)
	} else {
		// Return to search page with flash message
		C.renderTemplate(w, C.temp.search, o)
	}
}

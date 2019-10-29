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
	Search  bool
	db      *dbIO.DBIO
	pw      string
	w       http.ResponseWriter
	r       *http.Request
}

func newOutput(w http.ResponseWriter, r *http.Request, user, pw, ut string) (*Output, error) {
	// Returns empty output struct
	var err error
	o := new(Output)
	o.w = w
	o.r = r
	o.User = user
	o.Update = strings.Replace(ut, "UTC", "Eastern Time", 1)
	if pw != "" {
		o.pw = pw
		o.db, err = dbIO.Connect(C.config.Host, C.config.Database, o.User, o.pw)
		if err == nil {
			o.db.GetTableColumns()
		}
	}
	return o, err
}

func newFlash(w http.ResponseWriter, msg string) *Output {
	// Returns output with flash error message
	var o Output
	o.w = w
	o.Flash = msg
	return &o
}

func (o *Output) getTempFile(name string) {
	// Stores path to named file in tmp directory
	t := time.Now()
	stamp := t.Format(time.RFC3339)
	// Trim timestamp to minutes
	stamp = stamp[:strings.LastIndex(stamp, "-")]
	stamp = stamp[:strings.LastIndex(stamp, ":")]
	o.File = fmt.Sprintf("%s.%s.csv", name, stamp)
	o.Outfile = fmt.Sprintf("/tmp/%s", o.File)
}

func (o *Output) summary() {
	// Returns general database summary
	o.getTempFile("databaseSummary")
	header := "Field,Total,%"
	codbutils.WriteResults(o.Outfile, header, dbextract.GetSummary(o.db))
	C.renderTemplate(C.temp.result, o)
}

func (o *Output) cancerRates() {
	// Calculates cancer rates for matching criteria
	opt := setOptions(o.r)
	eval, msg, pass := setEvaluation(o.r, o.db.Columns, "0", "0")
	if msg == "" || !strings.Contains(msg, "Accounts") {
		var e []codbutils.Evaluation
		if pass {
			// Skip empty evaluations
			e = append(e, eval)
		}
		o.getTempFile(fmt.Sprintf("cancerRates.min%d", opt.Min))
		rates := dbextract.GetCancerRates(o.db, opt.Min, opt.Necropsy, e)
		rates.ToCSV(o.Outfile)
		C.renderTemplate(C.temp.result, o)
	} else {
		// Return to menu page with flash message
		o.Flash = msg
		C.renderTemplate(C.temp.menu, o)
	}
}

func (o *Output) extractTable() {
	// Extracts given table from the database
	name := strings.TrimSpace(o.r.PostForm.Get("Table"))
	if name != "" && name != "Empty" {
		table := o.db.GetTable(name)
		o.getTempFile(name)
		codbutils.WriteResults(o.Outfile, o.db.Columns[name], table)
		C.renderTemplate(C.temp.result, o)
	} else {
		o.Flash = "Please select a table to extract."
		C.renderTemplate(C.temp.menu, o)
	}
}

func (o *Output) searchDB() {
	// Searches database for results
	var eval map[string][]codbutils.Evaluation
	opt := setOptions(o.r)
	eval, o.Flash = checkEvaluations(o.r, o.db.Columns)
	if o.Flash == "" {
		res := dataframe.NewDataFrame(-1)
		// Search for column/value match
		for _, v := range eval {
			r := dbextract.SearchColumns(o.db, "", v, opt.Count, opt.Infant)
			if res.Length() == 0 {
				res = r
			} else {
				// Append successive results to results slice
				for _, i := range r.Rows {
					res.Rows = append(res.Rows, i)
				}
			}
		}
		if opt.Count == false && res.Length() >= 1 {
			o.getTempFile(o.User)
			res.ToCSV(o.Outfile)
		} else {
			o.Count = fmt.Sprintf("\tFound %d records matching search criteria.\n", res.Length())
		}
		o.Search = true
		C.renderTemplate(C.temp.result, o)
	} else {
		// Return to search page with flash message
		C.renderTemplate(C.temp.search, o)
	}
}

func (o *Output) routePost(source string) {
	// Sends post data to appropriate function
	o.r.ParseForm()
	switch source {
	case C.u.summary:
		o.summary()
	case C.u.rates:
		o.cancerRates()
	case C.u.table:
		o.extractTable()
	case C.u.search:
		o.searchDB()
	}
}

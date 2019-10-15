// Contains functions for interacting with slq database

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/dbIO"
	"net/http"
	"strings"
	"time"
)

func ping(user, password string) bool {
	// Returns true if credentials are valid
	return dbIO.Ping(C.config.Host, C.config.Database, user, password)
}

func changePassword(r *http.Request, user, password string) string {
	// Changes suer password or returns flash message
	var ret string
	db, err := dbIO.Connect(C.config.Host, C.config.Database, user, password)
	if err == nil {
		r.ParseForm()
		newpw := r.PostForm.Get("password")
		confpw := r.PostForm.Get("newpassword")
		if newpw != confpw {
			ret = "Passwords do not match."
		} else {
			cmd := fmt.Sprintf("SET PASSWORD = PASSWORD('%s')", newpw)
			_, er := db.DB.Exec(cmd)
			if er != nil {
				ret = er.Error()
			}
		}
	} else {
		// Convert error to string
		ret = err.Error()
	}
	return ret
}

type Output struct {
	User    string
	Flash   string
	File    string
	Outfile string
	Count   string
}

func newOutput(user string) *Output {
	// Returns empty output struct
	o := new(Output)
	o.User = user
	return o
}

func newFlash(msg string) *Output {
	// Returns output with flash error message
	o := newOutput("")
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

func (o *Output) searchDB(db *dbIO.DBIO, f *SearchForm, user string) {
	// Searches database for results
	var res [][]string
	var header string
	// Search for column/value match
	for _, v := range f.eval {
		r, h := dbextract.SearchColumns(db, o.User, f.Table, v, f.Count, f.Infant)
		// Append successive results to results slice
		res = append(res, r...)
		if header == "" {
			// Only record header once
			header = h
		}
	}
	if f.Count == false && len(res) >= 1 {
		o.getTempFile(user)
		codbutils.WriteResults(o.Outfile, header, res)
	} else {
		o.Count = fmt.Sprintf("\tFound %d records matching search criteria.\n", len(res))
	}
}

func extractFromDB(r *http.Request, user, password string) (*Output, error) {
	// Extracts data to outfile/stdout
	ret := newOutput(user)
	db, err := dbIO.Connect(C.config.Host, C.config.Database, ret.User, password)
	if err == nil {
		var f *SearchForm
		db.GetTableColumns()
		f, ret.Flash = setSearchForm(r, db.Columns)
		if ret.Flash == "" {
			if f.Dump == true {
				// Extract entire table
				table := db.GetTable(f.Table)
				ret.getTempFile(user)
				codbutils.WriteResults(ret.Outfile, db.Columns[f.Table], table)
			} else if f.Summary == true {
				ret.getTempFile("databaseSummary")
				header := "Field,Total,%"
				codbutils.WriteResults(ret.Outfile, header, dbextract.GetSummary(db))
			} else if f.Cancerrate == true {
				// Extract cancer rates
				ret.getTempFile(fmt.Sprintf("cancerRates.min%d", f.Min))
				header := "Kingdom,Phylum,Class,Orders,Family,Genus,ScientificName,TotalRecords,CancerRecords,CancerRate,"
				header += "AverageAge(months),AvgAgeCancer(months),Male,Female,MaleCancer,FemaleCancer"
				codbutils.WriteResults(ret.Outfile, header, dbextract.GetCancerRates(db, f.Min, f.Necropsy))
			} else {
				ret.searchDB(db, f, user)
			}
		}
	}
	return ret, err
}

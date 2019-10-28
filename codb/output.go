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

func ping(user, password string) (bool, string) {
	// Returns true if credentials are valid
	var update string
	ret := dbIO.Ping(C.config.Host, C.config.Database, user, password)
	if ret {
		db, _ := dbIO.Connect(C.config.Host, C.config.Database, user, password)
		db.GetTableColumns()
		update = db.LastUpdate().Format(time.RFC822)

	}
	return ret, update
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
	Update  string
	Flash   string
	File    string
	Outfile string
	Count   string
}

func newOutput(user, ut string) *Output {
	// Returns empty output struct
	o := new(Output)
	o.User = user
	o.Update = strings.Replace(ut, "UTC", "Eastern Time", 1)
	return o
}

func newFlash(msg string) *Output {
	// Returns output with flash error message
	o := newOutput("", "")
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

func (o *Output) searchDB(db *dbIO.DBIO, f *SearchForm) {
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
		o.getTempFile(o.User)
		codbutils.WriteResults(o.Outfile, header, res)
	} else {
		o.Count = fmt.Sprintf("\tFound %d records matching search criteria.\n", len(res))
	}
}

func (o *Output) extractFromDB(r *http.Request, password string) error {
	// Extracts data to outfile/stdout
	db, err := dbIO.Connect(C.config.Host, C.config.Database, o.User, password)
	if err == nil {
		var f *SearchForm
		db.GetTableColumns()
		f, o.Flash = setSearchForm(r, db.Columns)
		if o.Flash == "" {
			if len(f.Table) > 0 {
				// Extract entire table
				table := db.GetTable(f.Table)
				o.getTempFile(o.User)
				codbutils.WriteResults(o.Outfile, db.Columns[f.Table], table)
			} else if f.Summary == true {
				o.getTempFile("databaseSummary")
				header := "Field,Total,%"
				codbutils.WriteResults(o.Outfile, header, dbextract.GetSummary(db))
			} else if f.Cancerrate == true {
				// Extract cancer rates
				var e []codbutils.Evaluation
				o.getTempFile(fmt.Sprintf("cancerRates.min%d", f.Min))
				header := "Kingdom,Phylum,Class,Orders,Family,Genus,ScientificName,TotalRecords,CancerRecords,CancerRate,"
				header += "AverageAge(months),AvgAgeCancer(months),Male,Female,MaleCancer,FemaleCancer"
				for _, v := range f.eval {
					e = v
				}
				codbutils.WriteResults(o.Outfile, header, dbextract.GetCancerRates(db, f.Min, f.Necropsy, e))
			} else {
				o.searchDB(db, f)
			}
		}
	}
	return err
}

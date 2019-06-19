// Contains functions for interacting with slq database

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/dbIO"
	"strings"
	"time"
)

func ping(user, password string) bool {
	// Returns true if credentials are valid
	return dbIO.Ping(C.config.Host, C.config.Database, user, password)
}

type Output struct {
	User    string
	Flash   string
	File    string
	Outfile string
	Results []string
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
	stamp := strings.Replace(t.Format(time.Stamp), " ", "_", -1)
	o.File = fmt.Sprintf("%s_%s.csv", name, stamp)
	o.Outfile = fmt.Sprintf("/tmp/%s", o.File)
}

func (o *Output) searchDB(db *dbIO.DBIO, f *SearchForm, user string) {
	// Searches database for results
	var res [][]string
	var header string
	// Search for column/value match
	if f.Table == "" {
		// Set count to false to allow searching of results
		res, header = dbextract.SearchColumns(db, o.User, eval, f.Count, f.Infant)
	} else {
		res, header = dbextract.SearchSingleTable(db, f.Table, o.User, f.Column, f.Operator, f.Value, f.Infant)
	}
	if f.Count == false && len(res) >= 1 {
		o.getTempFile(user)
		codbutils.WriteResults(o.Outfile, header, res)
	} else {
		o.Results = []string{fmt.Sprintf("\tFound %d records where %s is %s.\n", len(res), f.Column, f.Value)}
	}
}

func extractFromDB(f *SearchForm, user, password string) (*Output, error) {
	// Extracts data to outfile/stdout
	ret := newOutput(user)
	db, err := dbIO.Connect(C.config.Host, C.config.Database, ret.User, password)
	if err == nil {
		db.GetTableColumns()
		if f.Dump == true {
			// Extract entire table
			table := db.GetTable(f.Table)
			ret.getTempFile(user)
			codbutils.WriteResults(ret.Outfile, db.Columns[f.Table], table)
		} else if f.Summary == true {
			ret.Results = []string{"Field\tTotal\t%\n"}
			summary := dbextract.GetSummary(db)
			for _, i := range summary {
				ret.Results = append(ret.Results, strings.Join(i, "\t"))
			}
		} else if f.Cancerrate == true {
			// Extract cancer rates
			ret.getTempFile(fmt.Sprintf("cancerRates.min%d.csv", f.Min))
			header := "Kingdom,Phylum,Class,Orders,Family,Genus,ScientificName,TotalRecords,CancerRecords,CancerRate,"
			header += "AverageAge(months),AvgAgeCancer(months),Male,Female,MaleCancer,FemaleCancer"
			codbutils.WriteResults(ret.Outfile, header, dbextract.GetCancerRates(db, f.Min, f.Necropsy))
		} else {
			ret.searchDB(db, f, user)
		}
	}
	return ret, err
}

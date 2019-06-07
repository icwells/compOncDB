// Contains functions for interacting with slq database

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/dbIO"
	"strings"
)

func ping(user, password string) bool {
	// Returns true if credentials are valid
	return dbIO.Ping(C.config.Host, C.config.Database, user, password)
}

type Output struct {
	User     string
	File     bool
	Outfile  string
	Filename string
	Results  []string
}

func newOutput(user string) *Output {
	// Returns empty output struct
	o := new(Output)
	o.User = user
	return o
}

func (o *Output) getTempFile(name string) {
	// Returns path to named file in tmp directory
	o.Outfile = fmt.Sprintf("/tmp/%s.csv", name)
	o.Filename = fmt.Sprintf("%s.csv", name)
	o.File = true
}

func (o *Output) searchDB(db *dbIO.DBIO, f SearchForm) {
	// Searches database for results
	var res [][]string
	var header string
	if f.Taxon == true {
		// Extract all data for a given species
		// Get single term (replace underscores in case terms are copied and pasted)
		names := []string{strings.Replace(f.Value, "_", " ", -1)}
		res, header = dbextract.SearchTaxonomicLevels(db, names, o.User, f.Column, f.Count, f.Common, f.Infant)
	} else if len(f.Operator) >= 1 {
		// Search for column/value match
		if f.Table == "" {
			tables := codbutils.GetTable(db.Columns, f.Column)
			// Set count to false to allow searching of results
			res, header = dbextract.SearchColumns(db, tables, o.User, f.Column, f.Operator, f.Value, false, f.Common, f.Infant)
		} else {
			res, header = dbextract.SearchSingleTable(db, f.Table, o.User, f.Column, f.Operator, f.Value, f.Common, f.Infant)
		}
		/*if len(e) > 1 {
			res = codbutils.FilterSearchResults(header, e[1:], res)
		}*/
	}
	if f.Count == false && len(res) >= 1 {
		o.getTempFile(f.Value)
		codbutils.WriteResults(o.Outfile, header, res)
	} else {
		o.Results = []string{fmt.Sprintf("\tFound %d records where %s is %s.\n", len(res), f.Column, f.Value)}
	}
}

func extractFromDB(f SearchForm, user, password string) *Output {
	// Extracts data to outfile/stdout
	ret := newOutput(user)
	db := dbIO.Connect(C.config.Host, C.config.Database, ret.User, password)
	db.GetTableColumns()
	if f.Dump == true {
		// Extract entire table
		table := db.GetTable(f.Table)
		ret.getTempFile(f.Table)
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
		ret.searchDB(db, f)
	}
	return ret
}

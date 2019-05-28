// Contains functions for interacting with slq database

package main

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
)

func ping(name, pw string) bool {
	// Returns true if credentials are valid
	return dbIO.Ping(CONFIG.Host, CONFIG.Database, user, pw)
}

func extractFromDB() {
	// Extracts data to outfile/stdout (all input variables are global)
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, false))
	if *dump != "nil" {
		// Extract entire table
		table := db.GetTable(*dump)
		codbutils.WriteResults(*outfile, db.Columns[*dump], table)
	} else if *sum == true {
		summary := dbextract.GetSummary(db)
		codbutils.WriteResults(*outfile, "Field,Total,%\n", summary)
	} else if *cr == true {
		// Extract cancer rates
		header := "Kingdom,Phylum,Class,Orders,Family,Genus,ScientificName,TotalRecords,CancerRecords,CancerRate,"
		header += "AverageAge(months),AvgAgeCancer(months),Male,Female,MaleCancer,FemaleCancer"
		rates := dbextract.GetCancerRates(db, *min, *nec)
		codbutils.WriteResults(*outfile, header, rates)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
}

func searchDB() {
	// Performs search functions on database
	var res [][]string
	var header string
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*config, *user, false))
	if *taxon != "nil" {
		// Extract all data for a given species
		var names []string
		if iotools.Exists(*taxon) == true {
			names = codbutils.ReadList(*taxon, *col)
		} else {
			// Get single term
			if strings.Contains(*taxon, "_") == true {
				names = []string{strings.Replace(*taxon, "_", " ", -1)}
			} else {
				names = []string{*taxon}
			}
		}
		res, header = dbextract.SearchTaxonomicLevels(db, names, *user, *level, *count, *com, *infant)
		fmt.Printf("\tFound %d records where %s is %s.\n", len(res), *level, *taxon)
	} else if *eval != "nil" {
		// Search for column/value match
		e := codbutils.SetOperations(*eval)
		if *table == "nil" {
			count := *count
			tables := codbutils.GetTable(db.Columns, e[0].Column)
			if len(e) > 1 {
				// Set count to false to allow searching of results
				count = false
			}
			res, header = dbextract.SearchColumns(db, tables, *user, e[0].Column, e[0].Operator, e[0].Value, count, *com, *infant)
		} else {
			res, header = dbextract.SearchSingleTable(db, *table, *user, e[0].Column, e[0].Operator, e[0].Value, *com, *infant)
		}
		fmt.Printf("\tFound %d records where %s is %s.\n", len(res), e[0].Column, e[0].Value)
		if len(e) > 1 {
			res = codbutils.FilterSearchResults(header, e[1:], res)
		}
	} else if *taxonomies == true {
		names := codbutils.ReadList(*infile, *col)
		res, header = dbextract.SearchSpeciesNames(db, names)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	if *count == false && len(res) >= 1 {
		codbutils.WriteResults(*outfile, header, res)
	}
}

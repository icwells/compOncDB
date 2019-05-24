// Contains functions for interacting with slq database

package main

import (
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/compOncDB/src/dbupload"
	"github.com/icwells/dbIO"
)

func newDatabase() time.Time {
	// Creates new database and tables
	c := setConfiguration(false)
	db := dbIO.CreateDatabase("", c.database, *user)
	db.NewTables(c.tables)
	return db.Starttime
}

func ping(name, pw string) bool {
	// Returns true if credentials are valid
}

func extractFromDB() time.Time {
	// Extracts data to outfile/stdout (all input variables are global)
	db := connectToDatabase(setConfiguration(false))
	if *dump != "nil" {
		// Extract entire table
		table := db.GetTable(*dump)
		if *outfile != "nil" {
			iotools.WriteToCSV(*outfile, db.Columns[*dump], table)
		} else {
			printArray(db.Columns[*dump], table)
		}
	} else if *sum == true {
		summary := dbextract.GetSummary(db)
		writeResults(*outfile, "Field,Total,%\n", summary)
	} else if *cr == true {
		// Extract cancer rates
		header := "Kingdom,Phylum,Class,Orders,Family,Genus,ScientificName,TotalRecords,CancerRecords,CancerRate,"
		header += "AverageAge(months),AvgAgeCancer(months),Male,Female,MaleCancer,FemaleCancer"
		rates := dbextract.GetCancerRates(db, *min, *nec)
		writeResults(*outfile, header, rates)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	return db.Starttime
}

func searchDB() time.Time {
	// Performs search functions on database
	var res [][]string
	var header string
	db := connectToDatabase(setConfiguration(false))
	if *taxon != "nil" {
		// Extract all data for a given species
		var names []string
		if iotools.Exists(*taxon) == true {
			names = readList(*taxon, *col)
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
		e := setOperations(*eval)
		if *table == "nil" {
			count := *count
			tables := getTable(db.Columns, e[0].column)
			if len(e) > 1 {
				// Set count to false to allow searching of results
				count = false
			}
			res, header = dbextract.SearchColumns(db, tables, *user, e[0].column, e[0].operator, e[0].value, count, *com, *infant)
		} else {
			res, header = dbextract.SearchSingleTable(db, *table, *user, e[0].column, e[0].operator, e[0].value, *com, *infant)
		}
		fmt.Printf("\tFound %d records where %s is %s.\n", len(res), e[0].column, e[0].value)
		if len(e) > 1 {
			res = filterSearchResults(header, e[1:], res)
		}
	} else if *taxonomies == true {
		names := readList(*infile, *col)
		res, header = dbextract.SearchSpeciesNames(db, names)
	} else {
		fmt.Print("\n\tPlease enter a valid command.\n\n")
	}
	if *count == false && len(res) >= 1 {
		writeResults(*outfile, header, res)
	}
	return db.Starttime
}


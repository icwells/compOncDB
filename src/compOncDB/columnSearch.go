// This script contains functions for searching tables for a given column/value combination

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/icwells/go-tools/strarray"
	"os"
	"strings"
)

func searchPairedTables(db *sql.DB, tables []string, column, value string) [][]string {
	// Returns cancatentaed results from paired tables
	t1 := dbIO.GetRows(db, tables[0], column, value, "*")

func searchColumns(db *sql.DB, col map[string]string, tables []string, column, value string) [][]string {
	// Determines search procedure
	var ret []string
	switch table[0] {
		// Start with potential mutliple entries
		case "Patient":
			switch column {
				case "ID":
				default:
			}
		case "Source":
			switch column {
				case "ID":
				default:
			}
		case "Tumor_relation":
			if len(tables) == 1 {

			} else {

			}
		case "Taxonomy":
			

		case "Common":

		case "Life_history":

		case "Totals":

		case "Diagnosis":
			
		case "Tumor":

		case "Accounts":
			if *user == "root" {
				// Return both tables
			} else {
				// Return single table
			}

		default:
			// Get matches from single table
			ret = dbIO.GetRows(db, tables[0], column, value, "*")
	}
	
}

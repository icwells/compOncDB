// This script will search values from a given mysql table

package dbIO

import (
	"database/sql"
	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

func searchText(ch chan []string, db *sql.DB, matcher *search.Matcher, col []string, table, column, term string) {
	// Searches slice for text match to single term
	var key string
	for _, i := range col {
		if matcher.EqualString(term, i) == true {
			key = i
			break
		}
	}
	rows := GetRows(db, table, column, key)
	ch <- rows[0]
}

func SearchColumnText(db *sql.DB, table, column string, terms []string) [][]string {
	// Searches given table for a match to the term in the given column
	var rows [][]string
	ch := make(chan []string)
	col := GetColumnText(db, table, column)
	matcher := search.New(language.English)
	for _, i := range terms {
		go searchText(ch, db, matcher, col, table, column, i)
		ret := <-ch
		rows = append(rows, ret)
	}
	return rows
}

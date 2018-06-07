// This script will search values from a given mysql table

package dbIO

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/text/search"
)

func searchText(ch chan string, matcher *Matcher, col *[]string, term string) {
	// Searches slice for text match to single term
	var key string
	for _, i := range col {
		if matcher.EqualString(term, i) == true {
			key = i
			break
		}
	}
	row := getRow(table, column, key)
	ch <- row
}

func SearchColumnText(db *DB, table, column string, terms []string) []string {
	// Searches given table for a match to the term in the given column
	var rows []string
	ch := make(chan string)
	col := getColumnText(db, table, column)
	matcher := search.New("AmericanEnglish")
	for _, i := range terms {
		go searchText(ch, matcher, *col, i)
		ret := <-ch
		rows = append(rows, ret)
	}
	return rows
}

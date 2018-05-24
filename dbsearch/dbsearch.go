// This script will search values from a given mysql table

package dbsearch

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/text/search"
)

func getRow(table, column, key string) string {
	// Returns row with key in column
	var r string
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s;", table, column, key)
	row := db.QueryRow(sql)
	err := row.Scan(&r)
	if err != nil {
		fmt.Fprintf("\n\t[Error] Reading %s where %s == %s: %v", table, column, key, err)
	}
	return r
}

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

func getColumnText(db *DB, table, column string) []string {
	// Returns slice of all entries in column
	var col []string
	sql := fmt.Sprintf("SELECT %s FROM %s;", column, table)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Fprintf("\n\t[Error] Extracting %s column from %s: %v", column, table, err)
	}
	defer rows.Close()
	for rows.Next() {
		var val string
		// Assign data to val while checking err
		if err := rows.Scan(&val); err != nil {
			fmt.Fprintf("\n\t[Error] Reading %s from %s: %v", column, table, err)
		}
		col = append(col, val)
	}
	return col
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

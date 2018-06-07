// This script contains general functions for extracting data from a database

package dbIO

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/icwells/go-tools/iotools"
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

func getColumnInt(db *DB, table, column string) []int {
	// Returns slice of all entries in column of integers
	var col []int
	sql := fmt.Sprintf("SELECT %s FROM %s;", column, table)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Fprintf("\n\t[Error] Extracting %s column from %s: %v", column, table, err)
	}
	defer rows.Close()
	for rows.Next() {
		var val int
		// Assign data to val while checking err
		if err := rows.Scan(&val); err != nil {
			fmt.Fprintf("\n\t[Error] Reading %s from %s: %v", column, table, err)
		}
		col = append(col, val)
	}
	return col
}

func getColumnText(db *DB, table, column string) []string {
	// Returns slice of all entries in column of text
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

func getTable(db *DB, table string) []string {
	// Returns contents of table
	var col []string
	sql := fmt.Sprintf("SELECT * FROM %s ;", table)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Fprintf("\n\t[Error] Extracting %s: %v", table, err)
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

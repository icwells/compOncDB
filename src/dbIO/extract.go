// This script contains general functions for extracting data from a database

package dbIO

import (
	"database/sql"
	"fmt"
)

func GetCount(db *sql.DB, table string) int {
	// Returns number of rows from table
	var n int
	cmd := fmt.Sprintf("SELECT COUNT(*) FROM %s;", table)
	val := db.QueryRow(cmd)
	err := val.Scan(&n)
	if err != nil {
		fmt.Printf("\n\t[Error] Determining number of rows from %s: %v\n\n", table, err)
	}
	return n
}

func GetMax(db *sql.DB, table, column string) int {
	// Returns maximum number from given column
	var m int
	n := GetCount(db, table)
	if n > 0 {
		cmd := fmt.Sprintf("SELECT MAX(%s) FROM %s;", column, table)
		val := db.QueryRow(cmd)
		err := val.Scan(&m)
		if err != nil {
			fmt.Printf("\n\t[Error] Determining maximum value from %s in %s: %v\n\n", column, table, err)
		}
	} else {
		m = n
	}
	return m
}

func GetRow(db *sql.DB, table, column, key string) string {
	// Returns row with key in column
	var r string
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s;", table, column, key)
	row := db.QueryRow(sql)
	err := row.Scan(&r)
	if err != nil {
		fmt.Printf("\n\t[Error] Reading %s where %s == %s: %v", table, column, key, err)
	}
	return r
}

func GetColumnInt(db *sql.DB, table, column string) []int {
	// Returns slice of all entries in column of integers
	var col []int
	sql := fmt.Sprintf("SELECT %s FROM %s;", column, table)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("\n\t[Error] Extracting %s column from %s: %v", column, table, err)
	}
	defer rows.Close()
	for rows.Next() {
		var val int
		// Assign data to val while checking err
		if err := rows.Scan(&val); err != nil {
			fmt.Printf("\n\t[Error] Reading %s from %s: %v", column, table, err)
		}
		col = append(col, val)
	}
	return col
}

func GetColumnText(db *sql.DB, table, column string) []string {
	// Returns slice of all entries in column of text
	var col []string
	sql := fmt.Sprintf("SELECT %s FROM %s;", column, table)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("\n\t[Error] Extracting %s column from %s: %v", column, table, err)
	}
	defer rows.Close()
	for rows.Next() {
		var val string
		// Assign data to val while checking err
		if err := rows.Scan(&val); err != nil {
			fmt.Printf("\n\t[Error] Reading %s from %s: %v", column, table, err)
		}
		col = append(col, val)
	}
	return col
}

func GetColumns(db *sql.DB, table, columns []string) []string {
	// Returns slice of slices of all entries in given columns of text
	var col [][]string
	sql := fmt.Sprintf("SELECT %s FROM %s;", strings.Join(columns, ","), table)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("\n\t[Error] Extracting columns from %s: %v", column, table, err)
	}
	defer rows.Close()
	for rows.Next() {
		var val []string
		// Assign data to val while checking err
		if err := rows.Scan(&val); err != nil {
			fmt.Printf("\n\t[Error] Reading %s from %s: %v", column, table, err)
		}
		col = append(col, val)
	}
	return col
}

func GetTable(db *sql.DB, table string) [][]string {
	// Returns contents of table
	var tbl [][]string
	sql := fmt.Sprintf("SELECT * FROM %s ;", table)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("\n\t[Error] Extracting %s: %v", table, err)
	}
	defer rows.Close()
	for rows.Next() {
		var val []string
		// Assign data to val while checking err
		if err := rows.Scan(&val); err != nil {
			fmt.Printf("\n\t[Error] Extracting %s: %v", table, err)
		}
		tbl = append(tbl, val)
	}
	return tbl
}

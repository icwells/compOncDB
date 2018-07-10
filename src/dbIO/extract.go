// This script contains general functions for extracting data from a database

package dbIO

import (
	"database/sql"
	"fmt"
	"strings"
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

func GetRows(db *sql.DB, table, column, key string) [][]string {
	// Returns rows with key in column
	var ret [][]string
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s;", table, column, key)
	rows := db.QueryRows(sql)
	if err != nil {
		fmt.Printf("\n\t[Error] Extracting rows from %s: %v", table, err)
	}
	defer rows.Close()
	for rows.Next() {
		var val []string
		// Assign data to val while checking err
		if err := rows.Scan(&val); err != nil {
			fmt.Printf("\n\t[Error] Reading row from %s: %v", table, err)
		}
		ret = append(ret, val)
	}
	return ret
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

func GetColumns(db *sql.DB, table string, columns []string) [][]string {
	// Returns slice of slices of all entries in given columns of text
	var col [][]string
	sql := fmt.Sprintf("SELECT %s FROM %s;", strings.Join(columns, ","), table)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("\n\t[Error] Extracting columns from %s: %v", table, err)
	}
	defer rows.Close()
	for rows.Next() {
		var val []string
		// Assign data to val while checking err
		if err := rows.Scan(&val); err != nil {
			fmt.Printf("\n\t[Error] Reading row from %s: %v", table, err)
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

func GetTableMap(db *sql.DB, table string) map[string][]string {
	// Returns table as a map with id as the key
	tbl := make(map[int][]string)
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
		tbl[val[0]] = val[1:]
	}
	return tbl
}

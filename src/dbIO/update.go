// This script contains general functions for extracting data from a database

package dbIO

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/icwells/go-tools/strarray"
	"strings"
)

func TruncateTable(db *sql.DB, table string) {
	// Clears all table contents
	cmd, err := db.Prepare(fmt.Sprintf("TRUNCATE TABLE %s;", table))
	if err != nil {
		fmt.Printf("\t[Error] Formatting command to truncate table %s: %v\n", table, err)
	} else {
		_, err = cmd.Exec()
		if err != nil {
			fmt.Printf("\t[Error] Truncating table %s: %v\n", table, err)
		}
	}
}

func columnEqualTo(columns string, values [][]string) []string {
	// Matches columns to inner slice by index, returns empty slice if indeces are not equal
	var ret []string
	col := strings.Split(columns, ",")
	for _, val := range values {
		if len(val) == len(col) {
			first := true
			buffer := bytes.NewBufferString("")
			for idx, i := range val {
				if first == false {
					// Write seperating comma
					buffer.WriteByte(',')
				}
				buffer.WriteString(col[idx])
				buffer.WriteByte('=')
				buffer.WriteByte('\'')
				buffer.WriteString(i)
				buffer.WriteByte('\'')
			}
			ret = append(ret, buffer.String())
		}
	}
	return ret
}

func UpdateRow(db *sql.DB, table, columns, target, key string, values [][]string) int {
	// Updates rows where target = key with values (matched to columns)
	ret := 0
	val := columnEqualTo(columns, values)
	for _, i := range val {
		cmd, err := db.Prepare(fmt.Sprintf("UPDATE %s SET %s WHERE %s = %s;", table, i, target, key))
		if err != nil {
			fmt.Printf("\t[Error] Preparing update for %s: %v\n", table, err)
		} else {
			_, err = cmd.Exec()
			cmd.Close()
			if err != nil {
				fmt.Printf("\t[Error] Updating row(s) from %s: %v\n", table, err)
			} else {
				ret++
			}
		}
	}
	return ret
}

func DeleteRow(db *sql.DB, table, column, value string) {
	// Deletes row(s) from database where column name = given value
	cmd, err := db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE %s = %s;", table, column, value))
	if err != nil {
		fmt.Printf("\t[Error] Preparing deletion from %s: %v\n", table, err)
	} else {
		_, err = cmd.Exec()
		cmd.Close()
		if err != nil {
			fmt.Printf("\t[Error] Deleting row(s) from %s: %v\n", table, err)
		}
	}
}

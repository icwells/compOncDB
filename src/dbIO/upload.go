// These functions will upload data to a database

package dbIO

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"golang.org/x/text/search"
	"os"
	"strings"
)

func updateDB(db *DB, table, columns, values string) {
	// Adds new rows to table
	//(values must be formatted for single/multiple rows before calling function)
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", table, column, values)
	_, err := db.Exec(sql)
	if err != nil {
		fmt.Fprintf("\t[Error] uploading to %s: %v", table, err)
		return 0
	}
	return 1
}

func formatValues(v []string, n int) string {
	// Organizes input data into n rows for upload
	if len(v)%n != 0 {
		fmt.Fprintf("\t[Error] Slice is not the correct length to fit into table: %v", v)
		return ""
	}
	buffer := bytes.NewBufferString("(")
	l := len(v) / n
	c := 0
	// Iterate through in blocks of row length
	for i := 0; i < len(v); i++ {
		if c == 0 {
			if i != 0 {
				// Write preceding comma and new open field
				buffer.WriteString(",(")
			}
		}
		// Write entry
		buffer.WriteString(v[i+c])
		if c == l {
			// Print close
			buffer.WriteString(")")
			// Reset counter
			c = 0
		} else {
			// Print seprating comma
			buffer.WriteString(",")
		}
		c++
	}
	buffer.WriteString(";")
	values := buffer.String()
	return values
}

func readColumns(types bool) map[string]string {
	// Build map of column statements
	var columns map[string]string
	var table string
	f := iotools.OpenFile("tableColumns.txt")
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if line[0] == '#' {
			// Get table names
			table = strings.TrimSpace(line[1:])
		} else {
			// Get columns for given table
			var col string
			if types == true {
				col = strings.TrimSpace(line)
			} else {
				c := strings.Split(line, " ")
				col = c[0]
			}
			if strarray.InMapStr(columns, table) == true {
				columns[table] = columns[table] + "," + col
			} else {
				columns[table] = col
			}
		}
	}
	return columns
}

func newTables(db *DB) {
	// Initializes new tables
	fmt.Println("\n\tInitializing new tables...")
	columns := readColumns(true)
	for k, v := range columns {
		sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", table, column, values)
		cmd := fmt.Sprintf("CREATE TABLE %s(%s);", k, v)
		_, err := db.Exec(cmd)
		if err == nil {
			fmt.Fprintf("\t[Error] Creating table {}. Exiting.\n", k)
			os.Exit(1)
		}
	}
}

// These functions will upload data to a database

package dbIO

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"os"
	"strings"
)

func UpdateDB(db *sql.DB, table, columns, values string) int {
	// Adds new rows to table
	//(values must be formatted for single/multiple rows before calling function)
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", table, columns, values)
	_, err := db.Exec(sql)
	if err != nil {
		fmt.Printf("\t[Error] uploading to %s: %v", table, err)
		return 0
	}
	return 1
}

func FormatMap(data map[string][]string) string {
	// Formats a map of string slices for upload
	buffer := bytes.NewBufferString("")
	count := 0
	length := len(data) - 1
	for _, val := range data {
		buffer.WriteString("(")
		c := 0
		l := len(val) - 1
		for _, v := range val {
			// Add row entries
			buffer.WriteString(v)
			if c != l {
				buffer.WriteString(",")
			}
			c++
		}
		buffer.WriteString(")")
		if count != length {
			buffer.WriteString(",")
		}
		count++
	}
	buffer.WriteString(";")
	values := buffer.String()
	return values
}

func FormatSlice(data [][]string) string {
	// Organizes input data into n rows for upload
	buffer := bytes.NewBufferString("")
	length := len(data) - 1
	for idx, row := range data {
		buffer.WriteString("(")
		l := len(row) - 1
		for i, v := range row {
			// Add row entries
			buffer.WriteString(v)
			if i != l {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString(")")
		if idx != length {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(";")
	values := buffer.String()
	return values
}

func ReadColumns(infile string, types bool) map[string]string {
	// Build map of column statements
	var columns map[string]string
	var table string
	f := iotools.OpenFile(infile)
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

func NewTables(db *sql.DB, infile string) {
	// Initializes new tables
	fmt.Println("\n\tInitializing new tables...")
	columns := ReadColumns(infile, true)
	for k, v := range columns {
		cmd := fmt.Sprintf("CREATE TABLE %s(%s);", k, v)
		_, err := db.Exec(cmd)
		if err == nil {
			fmt.Printf("\t[Error] Creating table {}. Exiting.\n", k)
			os.Exit(1)
		}
	}
}

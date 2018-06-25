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
	cmd := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", table, columns, values)
	_, err := db.Exec(cmd)
	if err != nil {
		fmt.Println(cmd)
		os.Exit(3)
		//fmt.Printf("\t[Error] Uploading to %s: %v", table, err)
		return 0
	}
	return 1
}

func FormatMap(data map[string][]string) string {
	// Formats a map of string slices for upload
	buffer := bytes.NewBufferString("")
	first := true
	for _, val := range data {
		f := true
		if first == false {
			// Add sepearating comma
			buffer.WriteByte(',')
		}
		buffer.WriteByte('(')
		for _, v := range val {
			if f == false {
				buffer.WriteByte(',')
			}
			// Wrap in back ticks to preserve spaces and reserved characters
			buffer.WriteByte('`')
			buffer.WriteString(v)
			buffer.WriteByte('`')
			f = false
		}
		buffer.WriteByte(')')
		first = false
	}
	return buffer.String()
}

func FormatSlice(data [][]string) string {
	// Organizes input data into n rows for upload
	buffer := bytes.NewBufferString("")
	for idx, row := range data {
		if idx != 0 {
			buffer.WriteByte(',')
		}
		buffer.WriteByte('(')
		for i, v := range row {
			if i != 0 {
				buffer.WriteByte(',')
			}
			// Wrap in back ticks to preserve spaces and reserved characters
			buffer.WriteByte('`')
			buffer.WriteString(v)
			buffer.WriteByte('`')
		}
		buffer.WriteByte(')')
	}
	return buffer.String()
}

func ReadColumns(infile string, types bool) map[string]string {
	// Build map of column statements
	columns := make(map[string]string)
	var table string
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if len(line) >= 3 {
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
					columns[table] = columns[table] + ", " + col
				} else {
					columns[table] = col
				}
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
		cmd := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(%s);", k, v)
		_, err := db.Exec(cmd)
		if err != nil {
			fmt.Printf("\n%s\n\n", cmd)
			fmt.Printf("\t[Error] Creating table %s. %v\n\n", k, err)
			os.Exit(1)
		}
	}
}

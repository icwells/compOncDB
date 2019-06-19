// Contains functions for convertng slice of string slices to map

package codbutils

import (
	"bufio"
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"os"
	"path"
	"strings"
)

func getutils() string {
	// Returns path to utils directory
	return path.Join(iotools.GetGOPATH(), "src/github.com/icwells/compOncDB/utils")
}

func getAbsPath(f string) string {
	// Prepends GOPATH to file name if needed
	if !strings.Contains(f, string(os.PathSeparator)) {
		f = path.Join(getutils(), f)
	}
	if iotools.Exists(f) == false {
		fmt.Printf("\n\t[Error] Cannot find %s file. Exiting.\n", f)
		os.Exit(1)
	}
	return f
}

type Configuration struct {
	Host     string
	Database string
	User     string
	Testdb   string
	Tables   string
	Test     bool
}

func SetConfiguration(config, user string, test bool) Configuration {
	// Gets setting from config.txt
	var c Configuration
	c.Test = test
	c.User = user
	f := iotools.OpenFile(getAbsPath(config))
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		s := strings.Split(string(scanner.Text()), "=")
		for idx, i := range s {
			s[idx] = strings.TrimSpace(i)
		}
		switch s[0] {
		case "host":
			c.Host = s[1]
		case "database":
			c.Database = s[1]
		case "test_database":
			c.Testdb = s[1]
		case "table_columns":
			c.Tables = getAbsPath(s[1])
		}
	}
	return c
}

func ConnectToDatabase(c Configuration) *dbIO.DBIO {
	// Manages call to Connect and GetTableColumns
	d := c.Database
	if c.Test == true {
		d = c.Testdb
	}
	db, err := dbIO.Connect(c.Host, d, c.User, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1000)
	}
	db.GetTableColumns()
	return db
}

func GetTable(tables map[string]string, col string) []string {
	// Determines which table column is in
	var ret []string
	col = strings.ToLower(col)
	if col == "id" {
		// Return tables for uid
		ret = []string{"Patient", "Source", "Diagnosis", "Tumor_relation"}
	} else if strings.Contains(col, "_id") == false {
		if strings.Contains(col, "_") == false {
			col = strings.Title(col)
		}
		// Iterate through available column names
		for k, val := range tables {
			for _, i := range strings.Split(val, ",") {
				i = strings.TrimSpace(i)
				if col == i {
					ret = append(ret, k)
					break
				}
			}
		}
	} else {
		// Return multiple tables for ids
		if col == "taxa_id" {
			ret = []string{"Patient", "Taxonomy", "Common", "Totals", "Life_history"}
		} else if col == "account_id" {
			ret = []string{"Source", "Accounts"}
		} else if col == "source_id" {
			ret = append(ret, "Patient")
		}
	}
	if len(ret) == 0 {
		fmt.Printf("\n\t[Error] Cannot find table with column %s. Exiting.\n\n", col)
		os.Exit(999)
	}
	return ret
}

func ReadList(infile string, idx int) []string {
	// Reads list of queries from file
	set := strarray.NewSet()
	var d string
	first := true
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := iotools.GetScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			s := strings.Split(line, d)
			if len(s) > idx {
				// Replace underscores if present
				name := strings.Replace(s[idx], "_", " ", -1)
				name = strings.TrimSpace(name)
				if len(name) > 1 {
					set.Add(name)
				}
			}
		} else {
			d = iotools.GetDelim(line)
			first = false
		}
	}
	return set.ToSlice()
}

func printArray(header string, table [][]string) {
	// Prints slice of string slcies to screen
	head := strings.Split(header, ",")
	// Wrap in newlines
	fmt.Println()
	fmt.Println(strings.Join(head, "\t"))
	for _, row := range table {
		fmt.Println(strings.Join(row, "\t"))
	}
	fmt.Println()
}

func WriteResults(outfile, header string, table [][]string) {
	// Wraps calls to writeCSV/printArray
	if len(table) > 0 {
		if outfile != "nil" {
			iotools.WriteToCSV(outfile, header, table)
		} else {
			printArray(header, table)
		}
	}
}

func DeleteEntries(d *dbIO.DBIO, tables []string, column, value string) {
	// Deletes matches from appropriate tables
	t := strings.Join(tables, ", ")
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\tAre you sure you want to delete all records from %s where %s equals %s (Enter Y for yes)? ", t, column, value)
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToUpper(input)) == "Y" {
		fmt.Println("\tProceeding with deletion...")
		for _, i := range tables {
			d.DeleteRow(i, column, value)
		}
	} else {
		fmt.Println("\tSkipping deletion.")
	}
}

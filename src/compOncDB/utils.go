// Contains functions for convertng slice of string slices to map

package main

import (
	"bufio"
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"os"
	"strings"
)

type configuration struct {
	host     string
	database string
	testdb   string
	tables   string
	test     bool
}

func setConfiguration(test bool) configuration {
	// Gets setting from config.txt
	var c configuration
	c.test = test
	if iotools.Exists(*config) == false {
		fmt.Print("\n\t[Error] Cannot find config file. Exiting.\n")
		os.Exit(1)
	}
	f := iotools.OpenFile(*config)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		s := strings.Split(string(scanner.Text()), "=")
		for idx, i := range s {
			s[idx] = strings.TrimSpace(i)
		}
		switch s[0] {
		case "host":
			c.host = s[1]
		case "database":
			c.database = s[1]
		case "test_database":
			c.testdb = s[1]
		case "table_columns":
			c.tables = s[1]
		}
	}
	return c
}

func connectToDatabase(c configuration) *dbIO.DBIO {
	// Manages call to Connect and GetTableColumns
	d := c.database
	if c.test == true {
		d = c.testdb
	}
	db := dbIO.Connect(c.host, d, *user)
	db.GetTableColumns()
	return db
}

func getOperation(eval string) (string, string, string) {
	// Splits eval into column, operator, value
	found := false
	var column, op, value string
	operators := []string{"==", ">=", "<=", "=", ">", "<"}
	for _, i := range operators {
		if strings.Contains(eval, i) == true {
			op = i
			if op == "==" {
				// Convert to single equals sign for sql
				op = "="
			}
			s := strings.Split(eval, i)
			if len(s) == 2 {
				// Only store properly formed queries
				column = strings.TrimSpace(s[0])
				value = strings.TrimSpace(s[1])
				found = true
			}
			break
		}
	}
	if found == false {
		fmt.Print("\n\t[Error] Please supply a valid evaluation argument. Exiting.\n\n")
		os.Exit(1001)
	}
	return column, op, value
}

func getTable(tables map[string]string, col string) []string {
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

func readList(infile string, idx int) []string {
	// Reads list of queries from file
	var ret []string
	var d string
	first := true
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := iotools.GetScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			s := strings.Split(line, d)
			if len(s) >= idx {
				// Replace underscores if present
				name := strings.Replace(s[idx], "_", " ", -1)
				ret = append(ret, strings.TrimSpace(name))
			}
		} else {
			d = iotools.GetDelim(line)
			first = false
		}
	}
	return ret
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

func writeResults(outfile, header string, table [][]string) {
	// Wraps calls to writeCSV/printArray
	if len(table) > 0 {
		if outfile != "nil" {
			iotools.WriteToCSV(outfile, header, table)
		} else {
			printArray(header, table)
		}
	}
}

func deleteEntries(d *dbIO.DBIO, tables []string, column, value string) {
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

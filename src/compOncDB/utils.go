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
		} else if col == "tumor_id" {
			ret = []string{"Tumor_relation", "Tumor"}
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

func readList(infile string) []string {
	// Reads list of queries from file
	var ret []string
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := iotools.GetScanner(f)
	for input.Scan() {
		line := string(input.Text())
		// Replace underscores if present
		line = strings.Replace(line, "_", " ", -1)
		ret = append(ret, strings.TrimSpace(line))
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

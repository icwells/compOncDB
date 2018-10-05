// Contains functions for convertng slice of string slices to map

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"os"
	"strings"
)

func typeof(v interface{}) string {
	// Returns string of object type
	return fmt.Sprintf("%T", v)
}

func toMap(t [][]string) map[string][]string {
	// Converts slice of string slices to map with first element as key
	// If slice is two columns wide, it will append the second item to map entry
	m := make(map[string][]string)
	for _, i := range t {
		_, ex := m[i[0]]
		if ex == false {
			if typeof(i[1]) == "string" {
				// Create new slice
				m[i[0]] = []string{i[1]}
			} else {
				m[i[0]] = i[1:]
			}
		} else if typeof(i[1]) == "string" && strarray.InSliceStr(m[i[0]], i[1]) == false {
			// Append new stirng element
			m[i[0]] = append(m[i[0]], i[1])
		}
	}
	return m
}

func mapOfMaps(t [][]string) map[string]map[string]string {
	// Converts table to map of maps for easier searching
	ret := make(map[string]map[string]string)
	for _, row := range t {
		if m, ex := ret[row[1]]; ex == true {
			if _, e := m[row[2]]; e == false {
				// Add to existing map
				ret[row[1]][row[2]] = row[0]
			}
		} else {
			// Make new sub-map
			ret[row[1]] = make(map[string]string)
			ret[row[1]][row[2]] = row[0]
		}
	}
	return ret
}

func entryMap(t [][]string) map[string]string {
	// Converts pair of columns to map for easier searching
	m := make(map[string]string)
	for _, i := range t {
		if _, ex := m[i[1]=; ex == false {
			m[i[1]] = i[0]
		}
	}
	return m
}

func getTable(tables map[string]string, col string) []string {
	// Determines which table column is in
	var ret []string
	col = strings.ToLower(col)
	if strings.Contains(col, "_id") == false {
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
	} else if col == "id" {
		// Return tables for uid
		ret = []string{"Patient", "Source", "Diagnosis", "Tumor_relation"}
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
		os.Exit(100)
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
	fmt.Print(strings.Join(head, "\t"))
	for _, row := range table {
		fmt.Println(strings.Join(row, "\t"))
	}
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

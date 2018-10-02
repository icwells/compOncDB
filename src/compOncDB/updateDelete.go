// This script contains functions for updating/deleting values from the database

package main

import (
	"bufio"
	"database/sql"
	"dbIO"
	"fmt"
	"os"
	"strings"
)

func deleteEntries(db *sql.DB, col map[string]string, tables []string, column, value string) {
	// Deletes matches from appropriate tables
	t := strings.Join(tables, ", ")
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\tAre you sure you want to delete all records from %s where %s equals %s (Enter Y for yes)? ", t, column, value)
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToUpper(input)) == "Y" {
		fmt.Println("\tProceeding with deletion...")
		for _, i := range tables {
			dbIO.DeleteRow(db, i, column, value)
		}
	} else {
		fmt.Println("\tSkipping deletion.")
	}
}

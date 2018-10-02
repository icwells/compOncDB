// This script contains functions for searching tables for a given column/value combination

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/icwells/go-tools/strarray"
	"os"
	"strings"
)

func searchColumns(db *sql.DB, col map[string]string, tables []string, column, value string) [][]string {
	// Searches given tables for matches
	var ret []string
	
}

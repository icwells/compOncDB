// This script will perform black box tests on compOncDB

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	//"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	"testing"
)

type searchTest struct {
	term   	string
	table	string
	total  	int
	single 	int
}

type searchCases struct {
	db    *sql.DB
	col   map[string]string
	cases []searchTest
}

func newSearchCases() searchCases {
	// Initializes new search struct
	var c searchCases
	c.db, _, _ = dbIO.Connect(DB, "guest")
	c.col = dbIO.ReadColumns("../../bin/" + COL, false)
	c.cases = []searchTest{
		{"Name=coyote", "Common", 209, 1},
		{"Sex==male", "Patient", 20000, 20000},
	}
	return c
}

func TestSearchColumns(t *testing.T) {
	// Perform tests on the column search functions
	c := newSearchCases()
	for _, i := range c.cases {
		column, op, value := getOperation(i.term)
		tables := getTable(c.col, column)
		res, _ := searchColumns(c.db, c.col, tables, column, op, value)
		fmt.Println(res)
		if len(res) != i.total {
			msg := fmt.Sprintf("Term %s failed column search. Expected: %d. Actual: %d.", i.term, i.total, len(res))
			t.Error(msg)
		}
	}
}

func TestSearchSingleTable(t *testing.T) {
	// Tests single table search function
	c := newSearchCases()
	for _, i := range c.cases {
		column, op, value := getOperation(i.term)
		res, _ := searchColumns(c.db, c.col, []string{i.table}, column, op, value)
		if len(res) != i.single {
			msg := fmt.Sprintf("Term %s failed single table search. Expected: %d. Actual: %d.", i.term, i.single, len(res))
			t.Error(msg)
		}
	}
}

/*func TestSearchTaxonomicLevels(t *testing.T) {
	// Tests taxa search functions
	c := newSearchCases()
	for _, i := range c.cases {
		column, op, value := getOperation(i.term)
		tables := getTable(col, column)

		res, _ := searchColumns(c.db, c.col, []string{value})
		if len(res) != i.total {
			msg := fmt.Sprintf("Term %s failed column search. Expected: %d. Actual: %d.", i.term, i.total, len(res))
			t.Error(msg)
		}
	}
}
}*/

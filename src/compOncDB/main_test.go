// This script will perform black box tests on compOncDB

package main

import (
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type searchTest struct {
	term	string
	total	int
	single	int
}

type searchCases struct {
	db		*sql.DB
	col		map[string]string
	cases	[]searchTest
}

func newSearchCases() searchCases {
	// Initializes new search struct
	var c searchCases
	c.db, _, _ = dbIO.Connect(DB, "root")
	c.col = dbIO.ReadColumns(COL, false)
	c.cases = []searchTest{
		{"Name==coyote", 209, 1}
		}
	return c
}

func TestSearchColumns(t *testing.T) {
	// TPerform tests on the column search functions
	c := newSearchCases()
	for _, i := range c.cases {
		column, op, value := getOperation(i.term)
		tables := getTable(col, column)
		res, _ := searchColumns(c.db, c.col, tables, column, op, value)
		if len(res) != i.total {
			msg := fmt.Sprintf("Term %s failed. Expected: %d. Actual: %d." i.term, i.total, len(res))
			t.Error(msg)
		}
	}
}

func TestSearchSingleTable(t *testing.T) {
	// Tests single table search function
	
}

func TestSearchTaxonomicLevels(t *testing.T) {
	// Tests taxa search functions
}

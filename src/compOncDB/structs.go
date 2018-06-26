// This script contains structs used for the comparative oncology database

package main

import (
	"bufio"
	"database/sql"
	"dbIO"
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"strconv"
	"strings"
)

type Patient struct (
	ID			int
	Sex			string
	Age			int
	Castrated	int
	taxa_id		int
	source_id	int
	Species		string
	Date		string
	Comments	string

)

type Diagnosis struct (
	ID				int
	masspresent		int
	metastasis_id	int
}

type TumorRelation struct (
	ID			int
	tumor_id	int
)

type Source struct (
	ID				int
	service_name	string
	account_id		int
	submitter_name	string
)

type Columns struct (
	sex			int
	age 		int
	castrated	int
	species		int
	date		int
	comments	int
	location	int
	tumor		int
	metastasis	int
	service		int
	submitter	int
	account		int
}

func (c *Columns) getIndeces(line string) {
	// Assigns columns indeces to struct
	s := strings.Split(line, ",")
}

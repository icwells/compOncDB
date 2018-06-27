// This script contains structs used for the comparative oncology database

package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Patient -------------------------------------------------------------------

type Patient struct {
	ID        int
	Sex       string
	Age       int
	Castrated int
	taxa_id   int
	source_id int
	Species   string
	Date      string
	Comments  string
}

func (p *Patient) setPatient(line []string, c Columns, id, tid int) {
	// Reads data from line to struct
	var err error
	p.ID = id
	p.taxa_id = tid
	p.Sex = line[c.sex]
	p.Age, err = strconv.Atoi(line[c.age])
	if err != nil {
		fmt.Printf("\t[Error] Could not convert %s to int: %v\n", line[c.age], err)
	}
	/*p.Castrated, err = strconv.Atoi(line[c./////])
	if err != nil {
		/////////////////////////////////////////////////////////
	}*/
	p.source_id, err = strconv.Atoi(line[c.id])
	if err != nil {
		fmt.Printf("\t[Error] Could not convert %s to int: %v\n", line[c.id], err)
	}
	p.Species = line[c.species]
	p.Date = line[c.date]
	p.Comments = line[c.Comments]
}

func (p *Patient) toSlice() []string {
	// Returns string slice
	return []string{strconv.Itoa(p.ID), p.Sex, strconv.Itoa(p.Age), strconv.Itoa(p.Castrated), strconv.Itoa(p.taxa_id),
		strconv.Itoa(p.source_id), p.Species, p.Date, p.Comments}
}

// Diagnosis -----------------------------------------------------------------

type Diagnosis struct {
	ID            int
	masspresent   int
	metastasis_id int
}

func (d *Diagnosis) setDiag(id, mass, meta int) {
	// Sets diagnosis variables
	d.ID = id
	d.masspresent = mass
	d.metastasis_id = meta
}

func (d *Diagnosis) toSlice() []string {
	// Returns string slice
	return []string{strconv.Itoa(d.ID), strconv.Itoa(d.masspresent), strconv.Itoa(d.metastasis_id)}
}

// TumorRelation -------------------------------------------------------------

type TumorRelation struct {
	ID       int
	tumor_id int
}

func (t *TumorRelation) setRelation(id, tid int) {
	// Sets tumor relation variables
	t.ID = id
	tumor_id = tid
}

func (t *TumorRelation) toSlice() []string {
	// Returns string slice
	return []string{strconv.Itoa(t.ID), strconv.Itoa(t.tumor_id)}
}

// Source --------------------------------------------------------------------

type Source struct {
	ID           int
	service_name string
	account_id   int
}

func (s *Source) setSource(id, aid int, service string) {
	// Sets source variables
	s.ID = id
	s.service_name = service
	s.account_id = aid
}

func (s *Source) toSlice() []string {
	// Returns string slice
	return []string{strconv.Itoa(ts.ID), s.service_name, strconv.Itoa(s.account_id), s.submitter_name}
}

// Columns -------------------------------------------------------------------

type Columns struct {
	id         int
	sex        int
	age        int
	castrated  int
	species    int
	date       int
	comments   int
	location   int
	tumor      int
	metastasis int
	submitter  int
	account    int
}

func (c *Columns) setIndeces(line []string) {
	// Assigns columns indeces to struct
	line = strings.ToLower(line)
	s := strings.Split(line, ",")
	for idx, i := range s {
		if i == "access#" || i == "id" {
			c.id = idx
		} else if i == "sex" {
			c.sex = idx
		} else if i == "age" || i == "age(months)" {
			c.age = idx
		} else if i == "castrated" {
			c.castrated = idx
		} else if i == "scientificname" {
			c.species = idx
		} else if i == "date" {
			c.date = idx
		} else if i == "diagnosis" {
			c.comments = idx
		} else if i == "location" {
			c.location = idx
		} else if i == "cancertype" {
			c.tumor = idx
		} else if i == "metastasis" {
			c.metastasis = idx
		} else if i == "client" {
			c.submitter = idx
		} else if i == "account" {
			c.account = idx
		}
	}
}

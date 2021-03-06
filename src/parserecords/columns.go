// Defines a struct for extracting and storing column data from a file header

package parserecords

import (
	"strings"
)

type columns struct {
	account    int
	age        int
	castrated  int
	code       int
	comments   int
	common     int
	date       int
	days       int
	id         int
	location   int
	malignant  int
	max        int
	metastasis int
	necropsy   int
	patient    int
	primary    int
	sex        int
	species    int
	submitter  int
	typ        int
	year       int
}

func newColumns() columns {
	// Returns new struct with columns set to -1
	var c columns
	c.id = -1
	c.species = -1
	c.common = -1
	c.age = -1
	c.days = -1
	c.sex = -1
	c.castrated = -1
	c.location = -1
	c.typ = -1
	c.primary = -1
	c.metastasis = -1
	c.malignant = -1
	c.necropsy = -1
	c.date = -1
	c.comments = -1
	c.account = -1
	c.submitter = -1
	c.code = -1
	c.patient = -1
	return c
}

func (c *columns) maxIndex(idx int) {
	// Replaces max if idx is greater
	if idx > c.max {
		c.max = idx
	}
}

func (c *columns) setColumns(header []string) {
	// Stores column indeces
	for idx, i := range header {
		i = strings.TrimSpace(i)
		i = strings.Replace(i, " ", "", -1)
		if i == "ID" || i == "OriginID" || i == "Access#" || i == "UID" {
			if c.id < 0 {
				// Only store first column field (later columns tend to be source ids)
				c.id = idx
				c.maxIndex(idx)
			}
		} else if i == "CommonName" || i == "Breed" || i == "PT_Name" {
			c.common = idx
			c.maxIndex(idx)
		} else if i == "ScientificName" || i == "BinomialScientific" {
			c.species = idx
			c.maxIndex(idx)
		} else if i == "Age(months)" || i == "Age" {
			c.age = idx
			c.maxIndex(idx)
		} else if i == "Days" {
			c.days = idx
			c.maxIndex(idx)
		} else if i == "Sex" {
			c.sex = idx
			c.maxIndex(idx)
		} else if i == "Castrated" {
			c.castrated = idx
			c.maxIndex(idx)
		} else if i == "Location" || i == "Tissue" {
			c.location = idx
			c.maxIndex(idx)
		} else if i == "CancerType" || i == "Type" {
			c.typ = idx
			c.maxIndex(idx)
		} else if i == "PrimaryTumor" || i == "Primary" {
			c.primary = idx
			c.maxIndex(idx)
		} else if i == "Metastasis" || i == "Metastatic" {
			c.metastasis = idx
			c.maxIndex(idx)
		} else if i == "Malignant" {
			c.malignant = idx
			c.maxIndex(idx)
		} else if i == "Necropsy" || i == "DeathviaCancerY/N" {
			c.necropsy = idx
			c.maxIndex(idx)
		} else if strings.Contains(i, "Date") == true {
			c.date = idx
			c.maxIndex(idx)
		} else if i == "Diagnosis" || i == "Comments" || i == "Description" {
			c.comments = idx
			c.maxIndex(idx)
		} else if i == "Account" {
			c.account = idx
			c.maxIndex(idx)
		} else if i == "Client" || i == "Owner" || i == "InstitutionID" {
			c.submitter = idx
			c.maxIndex(idx)
		} else if i == "Code" || i == "CancerY/N" {
			c.code = idx
			c.maxIndex(idx)
		} else if i == "Patient" || i == "Name" {
			c.patient = idx
			c.maxIndex(idx)
		} else if i == "Year" {
			c.year = idx
			c.maxIndex(idx)
		}
	}
}

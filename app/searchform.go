// Defines search struct for compOncDB web app

package main

import (
	"net/http"
	"strconv"
)

func setInt(v string) int {
	// Parses range value to integer
	ret := 50
	if i, err := strconv.Atoi(v); err == nil {
		ret = i
	}
	return ret
}

func setBool(v string) bool {
	// Parses checkbox value to true/false
	ret := false
	if len(v) >= 1 {
		ret = true
	}
	return ret
}

type SearchForm struct {
	Column     string
	Operator   string
	Value      string
	Taxon      bool
	Table      string
	Dump       bool
	Summary    bool
	Cancerrate bool
	Min        int
	Necropsy   bool
	Common     bool
	Count      bool
	Infant     bool
}

func setSearchForm(request *http.Request) SearchForm {
	// Populates struct from request data
	var s SearchForm
	s.Column = request.FormValue("Column")
	s.Operator = request.FormValue("Operator")
	s.Value = request.FormValue("Value")
	s.Taxon = setBool(request.FormValue("Taxon"))
	s.Table = request.FormValue("Table")
	s.Dump = setBool(request.FormValue("Dump"))
	s.Summary = setBool(request.FormValue("Summary"))
	s.Cancerrate = setBool(request.FormValue("Cancerrate"))
	s.Min = setInt(request.FormValue("Min"))
	s.Necropsy = setBool(request.FormValue("Necropsy"))
	s.Common = setBool(request.FormValue("Common"))
	s.Count = setBool(request.FormValue("Count"))
	s.Infant = setBool(request.FormValue("Infant"))
	return s
}

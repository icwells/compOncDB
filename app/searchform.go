// Defines search struct for compOncDB web app

package main

import (
	"github.com/gorilla/schema"
	"net/http"
	//"strconv"
)

/*func setInt(v string) int {
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
}*/

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

func setSearchForm(r *http.Request) *SearchForm {
	// Populates struct from request data
	s := new(SearchForm)
	decoder := schema.NewDecoder()
	r.ParseForm()
	decoder.Decode(s, r.PostForm)
	/*s.Column = r.PostForm.Get("Column")
	s.Operator = r.PostForm.Get("Operator")
	s.Value = r.PostForm.Get("Value")
	s.Taxon = setBool(r.PostForm.Get("Taxon"))
	s.Table = r.PostForm.Get("Table")
	s.Dump = setBool(r.PostForm.Get("Dump"))
	s.Summary = setBool(r.PostForm.Get("Summary"))
	s.Cancerrate = setBool(r.PostForm.Get("Cancerrate"))
	s.Min = setInt(r.PostForm.Get("Min"))
	s.Necropsy = setBool(r.PostForm.Get("Necropsy"))
	s.Common = setBool(r.PostForm.Get("Common"))
	s.Count = setBool(r.PostForm.Get("Count"))
	s.Infant = setBool(r.PostForm.Get("Infant"))*/
	return s
}

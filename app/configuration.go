// Defines configuration struct and methods

package main

import (
	"fmt"
	"github.com/gorilla/schema"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/iotools"
	"html/template"
	"net/http"
	"path"
)

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
	fmt.Println(s.Column, s.Operator, s.Value)
	return s
}

//----------------------------------------------------------------------------

type configuration struct {
	name       string
	appdir     string
	login      string
	search     string
	searchtemp string
	output     string
	resulttemp string
	logintemp  string
	logout     string
	static     string
	tmpl       string
	templates  *template.Template
	config     codbutils.Configuration
}

func setConfiguration() *configuration {
	// Returns pointer to initialized configuration struct
	var c configuration
	c.name = "session"
	c.appdir = path.Join(iotools.GetGOPATH(), "src/github.com/icwells/compOncDB/app")
	c.login = "/codb"
	c.search = "/codb/search"
	c.output = "/codb/results"
	c.logout = "/codb/logout"
	c.static = "/static/"
	c.tmpl = "templates/*.html"
	c.logintemp = "login"
	c.searchtemp = "search"
	c.resulttemp = "result"
	c.templates = template.Must(template.ParseGlob(c.tmpl))
	c.config = codbutils.SetConfiguration("config.txt", "", false)
	return &c
}

func (c *configuration) renderTemplate(w http.ResponseWriter, tmpl string, out *Output) {
	// Renders template and handles errrors
	err := c.templates.ExecuteTemplate(w, tmpl, out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

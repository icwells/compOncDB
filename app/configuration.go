// Defines configuration struct and methods

package main

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"html/template"
	"net/http"
)

type configuration struct {
	name       string
	source     string
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
	c.source = "/codb"
	c.search = "/codb/search"
	c.output = "/codb/results"
	c.logout = "/codb/logout"
	c.static = "/static/"
	c.tmpl = "templates/*.html"
	c.logintemp = "login.html"
	c.searchtemp = "search.html"
	c.resulttemp = "result.html"
	c.templates = template.Must(template.ParseGlob(c.tmpl))
	c.config = codbutils.SetConfiguration("config.txt", "", false)
	return &c
}

func (c *configuration) newCookie() *http.Cookie {
	// Populates cookie struct from configuration
	return &http.Cookie{
		Name:  c.name,
		Value: "",
		Path:  c.source,
	}
}

func (c *configuration) renderTemplate(w http.ResponseWriter, tmpl string, out *Output) {
	// Renders template and handles errrors
	err := c.templates.ExecuteTemplate(w, tmpl, out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

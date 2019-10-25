// Defines configuration struct and methods

package main

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/iotools"
	"html/template"
	"net/http"
	"path"
)

type configuration struct {
	name       string
	appdir     string
	source     string
	login      string
	search     string
	searchtemp string
	output     string
	get        string
	resulttemp string
	logintemp  string
	logout     string
	newpw      string
	changepw   string
	changetemp string
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
	c.source = "/codb/"
	c.login = "/codb/login"
	c.search = "/codb/search/"
	c.output = "/codb/results/"
	c.get = "/codb/get/"
	c.logout = "/codb/logout"
	c.newpw = "/codb/newpassword"
	c.changepw = "/codb/changepassword"
	c.static = "/static/"
	c.tmpl = "templates/*.html"
	c.logintemp = "login"
	c.searchtemp = "search"
	c.resulttemp = "result"
	c.changetemp = "changepassword"
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

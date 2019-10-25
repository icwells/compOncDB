// Defines configuration struct and methods

package main

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/iotools"
	"html/template"
	"net/http"
	"path"
)

type urls struct {
	source     string
	login      string
	search     string
	output     string
	get        string
	logout     string
	newpw      string
	changepw   string
	static     string
}

func setURLs() *urls {
	// Stores url stems
	u := new(urls)
	u.source = "/codb/"
	u.login = "/codb/login"
	u.search = "/codb/search/"
	u.output = "/codb/results/"
	u.get = "/codb/get/"
	u.logout = "/codb/logout"
	u.newpw = "/codb/newpassword"
	u.changepw = "/codb/changepassword"
	u.static = "/static/"
	return u
}

type temps struct {
	
}

type configuration struct {
	name       string
	appdir     string
	u		   *urls

	searchtemp string
	resulttemp string
	logintemp  string
	changetemp string
	tmpl       string

	templates  *template.Template
	config     codbutils.Configuration
}

func setConfiguration() *configuration {
	// Returns pointer to initialized configuration struct
	var c configuration
	c.name = "session"
	c.appdir = path.Join(iotools.GetGOPATH(), "src/github.com/icwells/compOncDB/app")
	c.u = setURLs()
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

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
	changepw string
	get      string
	lifehist string
	login    string
	logout   string
	menu     string
	newpw    string
	output   string
	reftaxa  string
	source   string
	static   string
	summary  string
	table    string
}

func setURLs() *urls {
	// Stores url stems
	u := new(urls)
	u.changepw = "/codb/changepassword"
	u.get = "/codb/get/"
	u.lifehist = "/codb/lifehistory/"
	u.login = "/codb/login"
	u.logout = "/codb/logout"
	u.menu = "/codb/menu/"
	u.newpw = "/codb/newpassword"
	u.output = "/codb/results/"
	u.reftaxa = "/codb/referencetaxonomy/"
	u.source = "/codb/"
	u.static = "/static/"
	u.summary = "/codb/summary/"
	u.table = "/codb/extractTable/"
	return u
}

type temps struct {
	source string
	login  string
	change string
	menu   string
	search string
	result string
}

func setTemps() *temps {
	// Stores tmeplate names
	t := new(temps)
	t.source = "templates/*.html"
	t.login = "login"
	t.change = "changepassword"
	t.menu = "menu"
	t.search = "search"
	t.result = "result"
	return t
}

type configuration struct {
	name      string
	appdir    string
	u         *urls
	temp      *temps
	templates *template.Template
	config    codbutils.Configuration
}

func setConfiguration() *configuration {
	// Returns pointer to initialized configuration struct
	var c configuration
	c.name = "session"
	c.appdir = path.Join(iotools.GetGOPATH(), "src/github.com/icwells/compOncDB/app")
	c.u = setURLs()
	c.temp = setTemps()
	c.templates = template.Must(template.ParseGlob(c.temp.source))
	c.config = codbutils.SetConfiguration("", false)
	return &c
}

func (c *configuration) renderTemplate(tmpl string, out *Output) {
	// Renders template and handles errrors
	err := c.templates.ExecuteTemplate(out.w, tmpl, out)
	if err != nil {
		http.Error(out.w, err.Error(), http.StatusInternalServerError)
	}
}

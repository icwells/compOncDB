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
	source   string
	login    string
	menu     string
	search   string
	summary  string
	output   string
	get      string
	logout   string
	newpw    string
	changepw string
	static   string
}

func setURLs() *urls {
	// Stores url stems
	u := new(urls)
	u.source = "/codb/"
	u.login = "/codb/login"
	u.menu = "/codb/menu/"
	u.search = "/codb/search/"
	u.summary = "/codb/summary/"
	u.output = "/codb/results/"
	u.get = "/codb/get/"
	u.logout = "/codb/logout"
	u.newpw = "/codb/newpassword"
	u.changepw = "/codb/changepassword"
	u.static = "/static/"
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

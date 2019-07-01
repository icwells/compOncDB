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
	"strconv"
	"strings"
)

type SearchForm struct {
	eval       []codbutils.Evaluation
	Table      string
	Dump       bool
	Summary    bool
	Cancerrate bool
	Min        int
	Necropsy   bool
	Count      bool
	Infant     bool
}

func (s *SearchForm) setEvaluation(n string, r *http.Request, columns map[string]string) (bool, string) {
	// Populates evalutaiton struct if matching term is found
	var e codbutils.Evaluation
	var msg string
	ret := false
	c := "Column"
	o := "Operator"
	v := "Value"
	e.Column = strings.TrimSpace(r.PostForm.Get(c + n))
	e.Operator = strings.TrimSpace(r.PostForm.Get(o + n))
	e.Value = strings.TrimSpace(r.PostForm.Get(v + n))
	if e.Column != "" || e.Value != "" {
		if e.Column != "" && e.Operator != "" && e.Value != "" {
			// Assign table and id type
			msg = e.SetTable(columns, false)
			if msg == "" {
				if e.Table != "Accounts" {
					s.eval = append(s.eval, e)
					ret = true
				} else {
					msg = "Cannot access Accounts table."
				}
			}
		} else {
			msg = "Please specify column and value fields."
		}
	}
	return ret, msg
}

func (s *SearchForm) checkEvaluations(r *http.Request, columns map[string]string) string {
	// Reads in variable number of search parameters
	var msg string
	found := true
	count := 0
	for found == true {
		// Keep checking until count exceeds number of fields
		found, msg = s.setEvaluation(strconv.Itoa(count), r, columns)
		count++
	}
	return msg
}

func setSearchForm(r *http.Request, columns map[string]string) (*SearchForm, string) {
	// Populates struct from request data
	s := new(SearchForm)
	decoder := schema.NewDecoder()
	r.ParseForm()
	decoder.Decode(s, r.PostForm)
	msg := s.checkEvaluations(r, columns)
	return s, msg
}

func (s *SearchForm) String() string {
	// Returns formatted string of options
	var ret strings.Builder
	ret.WriteString("Search Criteria:\n")
	for _, i := range s.eval {
		ret.WriteString(fmt.Sprintf("\t\t%s.%s %s %s return %s\n", i.Table, i.Column, i.Operator, i.Value, i.ID))
	}
	ret.WriteString(fmt.Sprintf("Table:\t%s\n", s.Table))
	ret.WriteString(fmt.Sprintf("Dump:\t%v\n", s.Dump))
	ret.WriteString(fmt.Sprintf("Summary:\t%v\n", s.Summary))
	ret.WriteString(fmt.Sprintf("Cancerrate:\t%v\n", s.Cancerrate))
	ret.WriteString(fmt.Sprintf("Min:\t%d\n", s.Min))
	ret.WriteString(fmt.Sprintf("Necropsy:\t%v\n", s.Necropsy))
	ret.WriteString(fmt.Sprintf("Count:\t%v\n", s.Count))
	ret.WriteString(fmt.Sprintf("Infant:\t%v\n", s.Infant))
	return ret.String()
}

//----------------------------------------------------------------------------

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

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
	eval       map[string][]codbutils.Evaluation
	Table      string
	Dump       bool
	Summary    bool
	Cancerrate bool
	Min        int
	Necropsy   bool
	Count      bool
	Infant     bool
}

func (s *SearchForm) setEvaluation(r *http.Request, columns map[string]string, search, n string) (bool, string) {
	// Populates evalutaiton struct if matching term is found
	var e codbutils.Evaluation
	var msg string
	ret := false
	e.Table = strings.TrimSpace(r.PostForm.Get(search + "Table" + n))
	e.Column = strings.TrimSpace(r.PostForm.Get(search + "Column" + n))
	e.Operator = strings.TrimSpace(r.PostForm.Get(search + "Operator" + n))
	e.Value = strings.TrimSpace(r.PostForm.Get(search + "Value" + n))
	slct := strings.TrimSpace(r.PostForm.Get(search + "Select" + n))
	if e.Value != "" || slct != "" {
		if e.Value == "" {
			// Assign selected value to input value
			e.Value = slct
		}
		if e.Table != "" && e.Column != "" && e.Operator != "" && e.Value != "" {
			e.SetIDType(columns)
			if e.Table != "Accounts" {
				s.eval[search] = append(s.eval[search], e)
				ret = true
			} else {
				msg = "Cannot access Accounts table."
			}
		} else {
			msg = "Please supply valid table, column, and value."
		}
	}
	return ret, msg
}

func (s *SearchForm) checkEvaluations(r *http.Request, columns map[string]string) string {
	// Reads in variable number of search parameters
	var msg string
	for i := 0; i < 10; i++ {
		// Loop through all possible searches since deletions might disrupt numerical order
		found := true
		count := 0
		for found == true {
			// Keep checking until count exceeds number of fields
			found, msg = s.setEvaluation(r, columns, strconv.Itoa(i), strconv.Itoa(count))
			count++
		}
	}
	return msg
}

func setSearchForm(r *http.Request, columns map[string]string) (*SearchForm, string) {
	// Populates struct from request data
	s := new(SearchForm)
	s.eval = make(map[string][]codbutils.Evaluation)
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
	for k, v := range s.eval {
		ret.WriteString(fmt.Sprintf("\tSearch %s:\n", k))
		for _, i := range v {
			ret.WriteString(fmt.Sprintf("\t\t%s.%s %s %s return %s\n", i.Table, i.Column, i.Operator, i.Value, i.ID))
		}
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

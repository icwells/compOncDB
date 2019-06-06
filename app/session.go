// Defines session struct and methods

package main

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/iotools"
	"html/template"
	"net/http"
	"path"
)

type Session struct {
	name       string
	User       string
	password   string
	source     string
	search     string
	searchtemp *template.Template
	output     string
	appdir     string
	resulttemp *template.Template
	login      string
	logintemp  *template.Template
	logout     string
	template   string
	static     string
	config     codbutils.Configuration
}

func (s *Session) setTemplates() {
	// Stores templates for rendering
	base := path.Join(s.appdir, "/templates/base.html")
	login := path.Join(s.appdir, "/templates/login.html")
	search := path.Join(s.appdir, "/templates/search.html")
	result := path.Join(s.appdir, "/templates/result.html")
	s.logintemp = template.Must(template.ParseFiles(base, login))
	s.searchtemp = template.Must(template.ParseFiles(base, search))
	s.resulttemp = template.Must(template.ParseFiles(base, result))
}

func setSession() *Session {
	// Returns pointer to initialized session struct
	var s Session
	s.name = "session"
	s.source = "/codb"
	s.search = "/codb/search"
	s.output = "/codb/results"
	s.login = "/codb/login"
	s.logout = "/codb/logout"
	s.appdir = path.Join(iotools.GetGOPATH(), "src/github.com/icwells/compOncDB/app")
	s.static = "/static/"
	s.config = codbutils.SetConfiguration("config.txt", "", false)
	s.setTemplates()
	return &s
}

func (s *Session) newCookie() *http.Cookie {
	// Populates cookie struct from session
	return &http.Cookie{
		Name:  s.name,
		Value: "",
		Path:  s.source,
	}
}

func (s *Session) storeSession(w http.ResponseWriter) {
	// Stores session info in cookie
	value := map[string]string{
		"name":     s.User,
		"password": s.password,
	}
	if encoded, err := cookieHandler.Encode(s.name, value); err == nil {
		cookie := s.newCookie()
		cookie.Value = encoded
		http.SetCookie(w, cookie)
	}
}

func (s *Session) getCredentials(r *http.Request) {
	// Stores username and password from cookie
	if cookie, err := r.Cookie(s.name); err == nil {
		value := make(map[string]string)
		// Extract cookie.Value to value map
		if err = cookieHandler.Decode(s.name, cookie.Value, &value); err == nil {
			s.User = value["name"]
			s.password = value["password"]
		}
	}
}

func (s *Session) renderTemplate(w http.ResponseWriter, tmpl *template.Template, out *Output) {
	// Renders template and handles errrors
	err := tmpl.Execute(w, out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

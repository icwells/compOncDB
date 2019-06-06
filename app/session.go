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
	searchtemp string
	output     string
	appdir     string
	resulttemp string
	login      string
	logintemp  string
	logout     string
	static     string
	tmpl       string
	templates  *template.Template
	config     codbutils.Configuration
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
	s.tmpl = path.Join(s.appdir, "templates/*.html")
	s.logintemp = "login.html"
	s.searchtemp = "search.html"
	s.resulttemp = "result.html"
	s.templates = template.Must(template.ParseGlob(s.tmpl))
	s.config = codbutils.SetConfiguration("config.txt", "", false)
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

func (s *Session) renderTemplate(w http.ResponseWriter, tmpl string, out *Output) {
	// Renders template and handles errrors
	err := s.templates.ExecuteTemplate(w, tmpl, out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

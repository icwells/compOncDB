// Main script for compOncDB web user interface

package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
	"path"
)

var (
	cookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	)
	S = setSession()
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serves login page
	S.renderTemplate(w, S.logintemp, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Handles login
	redirect := S.source
	S.User = r.FormValue("name")
	S.password = r.FormValue("password")
	if S.User != "" && S.password != "" {
		// Check credentials
		if ping() {
			S.storeSession(w)
			redirect = S.search
		}
	}
	http.Redirect(w, r, redirect, 302)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clears session and returns to login page
	http.SetCookie(w, S.newCookie())
	http.Redirect(w, r, S.source, 302)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	// Renders search form
	S.renderTemplate(w, S.searchtemp, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Reads search form
	form := setSearchForm(r)
	out := extractFromDB(form)
	S.renderTemplate(w, S.resulttemp, out)
}

func main() {
	// Register handler functions
	r := mux.NewRouter()

	fileserver := http.FileServer(http.Dir(path.Join(S.appdir, S.static)))
	r.PathPrefix(S.static).Handler(http.StripPrefix(S.static, fileserver))
	r.HandleFunc(S.source, indexHandler).Methods("GET")
	r.HandleFunc(S.login, loginHandler).Methods("POST")
	r.HandleFunc(S.logout, logoutHandler).Methods("POST")
	r.HandleFunc(S.search, formHandler).Methods("GET")
	r.HandleFunc(S.search, searchHandler).Methods("POST")
	// Serve and log errors to terminal
	http.Handle(S.source, r)
	//log.Fatal(http.ListenAndServe(S.config.Host + ":8080", nil))
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

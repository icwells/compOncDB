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
	ROUTER = mux.NewRouter()
	S      = setSession()
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serves login page
	S.renderTemplate(w, S.logintemp, newOutput())
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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Reads search form
	if r.Method == http.MethodPost {
		// Parse search form
		form := setSearchForm(r)
		out := extractFromDB(form)
		S.renderTemplate(w, S.resulttemp, out)
	} else {
		// Render search form
		S.renderTemplate(w, S.searchtemp, newOutput())
	}
}

func main() {
	// Register handler functions
	fileserver := http.FileServer(http.Dir(path.Join(S.appdir, S.static)))
	ROUTER.HandleFunc(S.source, indexHandler).Methods("GET")
	ROUTER.HandleFunc(S.login, loginHandler).Methods("POST")
	ROUTER.HandleFunc(S.logout, logoutHandler).Methods("POST")
	ROUTER.HandleFunc(S.search, searchHandler)
	// Serve and log errors to terminal
	http.Handle(S.static, http.StripPrefix(S.static, fileserver))
	http.Handle(S.source, ROUTER)
	//log.Fatal(http.ListenAndServe(S.config.Host + ":8080", nil))
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

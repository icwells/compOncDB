// Main script for compOncDB web user interface

package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
	"strings"
)

var (
	cookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	)
	ROUTER = mux.NewRouter()
	S      = setSession()
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Handles login r
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
	ROUTER.HandleFunc(S.source, loginHandler).Methods("POST")
	ROUTER.HandleFunc(S.logout, logoutHandler).Methods("POST")
	ROUTER.HandleFunc(S.search, searchHandler)
	// Serve and log errors to terminal
	http.Handle(S.static, http.StripPrefix(S.static, http.FileServer(http.Dir(strings.Replace(S.static, "/", "", 2)))))
	http.Handle(S.source, ROUTER)
	log.Fatal(http.ListenAndServe(S.config.Host, nil))
}

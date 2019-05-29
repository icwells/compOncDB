// Main script for compOncDB web user interface

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/icwells/compOncDB/src/codbutils"
	"html/template"
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
		if form.search == true {

		} else {

		}
	} else {
		// Render search form
		t, err := template.ParseFile(fmt.Sprintf("%s%s.html", S.template, S.search))
		if err == nil {
			t.Execute(w, S)
		}
		//http.Redirect(w, r, S.search, 302)
	}
}

func main() {
	// Register handler functions
	router.HandleFunc(S.source, loginHandler).Methods("POST")
	router.HandleFunc(S.logout, logoutHandler).Methods("POST")
	router.HandleFunc(S.search, searchHandler)
	// Serve and log errors to terminal
	http.Handle(S.static, http.StripPrefix(S.static, http.FileServer(http.Dir(strings.Replace(S.static, "/", "", 2)))))
	http.Handle(S.source, ROUTER)
	log.Fatal(http.ListenAndServe(S.host, nil))
}

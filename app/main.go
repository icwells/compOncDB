// Main script for compOncDB web user interface

package main

import (
	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var (
	STORE = sessions.NewCookieStore(securecookie.GenerateRandomKey(64))
	C     = setConfiguration()
)

func getCredentials(r *http.Request) (string, string) {
	// Stores username and password from cookie
	session, _ := STORE.Get(r, C.name)
	name, ex := session.Values["username"]
	if !ex {
		return "", ""
	}
	pw, e := session.Values["password"]
	if !e {
		return "", ""
	}
	return name.(string), pw.(string)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serves login page
	//w.Header().Set("Content-Type", "text/css; charset=utf-8")
	C.renderTemplate(w, C.logintemp, newOutput(""))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Handles login
	redirect := C.source
	r.ParseForm()
	user := r.PostForm.Get("name")
	pw := r.PostForm.Get("password")
	if user != "" && pw != "" {
		// Check credentials
		if ping(user, pw) {
			// Store cookie
			session, _ := STORE.Get(r, C.name)
			session.Values["username"] = user
			// encrypt later on
			session.Values["password"] = pw
			session.Save(r, w)
			redirect = C.search
		}
	}
	http.Redirect(w, r, redirect, 302)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clears session and returns to login page
	//http.SetCookie(w, C.newCookie())
	http.Redirect(w, r, C.source, 302)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	// Renders search form (newOutput supplies username)
	user, _ := getCredentials(r)
	C.renderTemplate(w, C.searchtemp, newOutput(user))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Reads search form
	user, pw := getCredentials(r)
	form := setSearchForm(r)
	out := extractFromDB(form, user, pw)
	C.renderTemplate(w, C.resulttemp, out)
}

func main() {
	// Initilaize multiplexer and fileserver
	r := mux.NewRouter()
	fs := http.FileServer(rice.MustFindBox("static").HTTPBox())
	r.PathPrefix(C.static).Handler(http.StripPrefix(C.static, fs))
	// Register handler functions
	r.HandleFunc(C.source, indexHandler).Methods("GET")
	r.HandleFunc(C.source, loginHandler).Methods("POST")
	r.HandleFunc(C.logout, logoutHandler).Methods("POST")
	r.HandleFunc(C.search, formHandler).Methods("GET")
	r.HandleFunc(C.search, searchHandler).Methods("POST")
	// Serve and log errors to terminal
	http.Handle(C.source, r)
	//log.Fatal(http.ListenAndServe(C.config.Host + ":8080", nil))
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

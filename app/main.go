// Main script for compOncDB web user interface

package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var (
	STORE = sessions.NewCookieStore(securecookie.GenerateRandomKey(64))
	C = setConfiguration()
)

func getCredentials(r *http.Request) {
	// Stores username and password from cookie
	session, err := C.store.Get(r, C.name)
	if err == nil {
		value := make(map[string]string)
		// Extract cookie.Value to value map
		if err = cookieHandler.Decode(c.name, cookie.Value, &value); err == nil {
			c.User = value["name"]
			c.password = value["password"]
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serves login page
	C.renderTemplate(w, C.logintemp, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Handles login
	redirect := C.source
	r.ParseForm()
	user = r.PostForm.Get("name")
	ps = r.PostForm.Get("password")
	if C.User != "" && C.password != "" {
		// Check credentials
		if ping() {
			// Store cookie
			session, _ := store.Get(r, C.name)
			session.Values["username"] = user
			// bcrypt later on
			session.Values["password"] = pw
			session.Save(r, w)
			redirect = C.search
		}
	}
	http.Redirect(w, r, redirect, 302)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clears session and returns to login page
	http.SetCookie(w, C.newCookie())
	http.Redirect(w, r, C.source, 302)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	// Renders search form
	C.renderTemplate(w, C.searchtemp, newOutput())
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Reads search form
	form := setSearchForm(r)
	out := extractFromDB(form)
	C.renderTemplate(w, C.resulttemp, out)
}

func main() {
	// Register handler functions
	r := mux.NewRouter()
	fileserver := http.FileServer(http.Dir("." + C.static))
	r.PathPrefix(C.static).Handler(http.StripPrefix(C.static, fileserver))
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

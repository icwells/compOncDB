// Main script for compOncDB web user interface

package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"html/template"
	"log"
	"net/http"
)

var (
	cookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	)
	ROUTER  = mux.NewRouter()
	SOURCE  = "/"
	SEARCH  = "/search"
	OUTPUT  = "/output"
	LOGIN   = "/login"
	LOGOUT  = "/logout"
	SESSION = "session"
)

func newCookie() *http.Cookie {
	// Returns empty cookie struct
	return 	&http.Cookie{
		Name:  SESSION,
		Value: "",
		Path:  SOURCE,
	}
}

func setSession(name, pw string, response http.Response) {
	// Stores session info
	value := map[string]string{
		"name":     name,
		"password": pw,
	}
	if encoded, err := cookieHandler.Encode(SESSION, value); err == nil {
		cookie := newCookie()
		cookie.Value = encoded
		http.SetCookie(response, cookie)
	}
}

func getCredentials(request *http.Request) (string, string) {
	// Returns username and password from cookie
	var name, pw string
	if cookie, err := request.Cookie(SESSION); err == nil {
		value := make(map[string]string)
		// Extract cookie.Value to value map
		if err = cookieHandler.Decode(SESSION, cookie.Value, &value); err == nil {
			name = value["name"]
			pw = value["password"]
		}
	}
	return name, pw
}

func indexHandler(response http.ResponseWriter, request *http.Request) {
	fileServer := http.FileServer(http.Dir("static/"))
	fileServer.ServeHTTP(response, request)
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	// Handles login request
	redirect := SOURCE
	name := request.FormValue("name")
	pw := request.FormValue("password")
	if name != "" && pw != "" {
		// Check credentials
		if ping(name, pw) {
			setSession(name, pw, response)
			redirect = SEARCH
		}
	}
	http.Redirect(response, request, redirect, 302)
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	// Clears session and returns to login page
	http.SetCookie(response, newCookie())
	http.Redirect(response, request, SOURCE, 302)
}

func searchHandler(response http.ResponseWriter, request *http.Request) {
	// Reads search form
}

func main() {
	// Register handler functions
	router.HandleFunc(SOURCE, indexHandler)
	router.HandleFunc(LOGIN, loginHandler).Methods("POST")
	router.HandleFunc(LOGOUT, logoutHandler).Methods("POST")
	router.HandleFunc(SEARCH, searchHandler)
	// Serve and log errors to terminal
	http.Handle("/", ROUTER)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

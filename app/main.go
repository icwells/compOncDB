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
	STORE = sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	C     = setConfiguration()
)

func getCredentials(r *http.Request) (string, string) {
	// Reads username and password from cookie
	var user, password string
	session, _ := STORE.Get(r, C.name)
	name, ex := session.Values["username"]
	if ex {
		user = name.(string)
	}
	pw, e := session.Values["password"]
	if e {
		password = pw.(string)
	}
	return user, password
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serves login page
	user, pw := getCredentials(r)
	if user != "" && pw != "" && ping(user, pw) {
		// Forward to search form if logged in
		http.Redirect(w, r, C.search, http.StatusFound)
	} else {
		// Render login template
		C.renderTemplate(w, C.logintemp, newOutput(""))
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Handles login
	pass := false
	r.ParseForm()
	user := r.PostForm.Get("name")
	pw := r.PostForm.Get("password")
	if user != "" && pw != "" {
		// Check credentials
		if ping(user, pw) {
			// Store cookie
			session, _ := STORE.Get(r, C.name)
			session.Values["username"] = user
			session.Values["password"] = pw
			session.Save(r, w)
			pass = true
		}
	}
	if pass {
		http.Redirect(w, r, C.search, http.StatusFound)
	} else {
		C.renderTemplate(w, C.logintemp, newFlash("Username or password is incorrect."))
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clears session and returns to login page
	session, _ := STORE.Get(r, C.name)
	session.Values["username"] = ""
	session.Values["password"] = ""
	session.Save(r, w)
	http.Redirect(w, r, C.login, http.StatusFound)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	// Renders search form (newOutput supplies username)
	user, _ := getCredentials(r)
	if user != "" {
		C.renderTemplate(w, C.searchtemp, newOutput(user))
	} else {
		C.renderTemplate(w, C.logintemp, newFlash("Please login to access database."))
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Reads search form
	user, pw := getCredentials(r)
	if user != "" && pw != "" {
		form := setSearchForm(r)
		out, err := extractFromDB(form, user, pw)
		if err == nil {
			C.renderTemplate(w, C.resulttemp, out)
		} else {
			// Return to login page if an error is encoutered
			C.renderTemplate(w, C.logintemp, newFlash(err.Error()))
		}
	} else {
		C.renderTemplate(w, C.logintemp, newFlash("Please login to access database."))
	}
}

func main() {
	// Initilaize multiplexer and fileserver
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("." + C.static))
	http.Handle(C.static, http.StripPrefix(C.static, fs))
	// Register handler functions
	r.HandleFunc(C.login, indexHandler).Methods(http.MethodGet)
	r.HandleFunc(C.login, loginHandler).Methods(http.MethodPost)
	r.HandleFunc(C.logout, logoutHandler).Methods(http.MethodGet)
	r.HandleFunc(C.search, formHandler).Methods(http.MethodGet)
	r.HandleFunc(C.search, searchHandler).Methods(http.MethodPost)
	// Serve and log errors to terminal
	http.Handle("/", r)
	//log.Fatal(http.ListenAndServe(C.config.Host + ":8080", nil))
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

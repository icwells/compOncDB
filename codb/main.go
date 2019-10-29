// Main script for compOncDB web user interface

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/http"
)

var (
	STORE = sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	C     = setConfiguration()
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serves login page
	login := true
	user, pw, _ := getCredentials(w, r)
	if user != "" && pw != "" {
		if pass, _ := ping(user, pw); pass {
			// Forward to search form if logged in
			http.Redirect(w, r, C.u.menu, http.StatusFound)
			login = false
		}
	}
	if login {
		// Render login template
		o, _ := newOutput(w, r, "", "", "")
		C.renderTemplate(C.temp.login, o)
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
		p, ut := ping(user, pw)
		if p {
			// Store cookie
			session, _ := STORE.Get(r, C.name)
			session.Values["timestamp"] = getTimestamp()
			session.Values["username"] = user
			session.Values["password"] = pw
			session.Values["updatetime"] = ut
			session.Save(r, w)
			pass = true
		}
	}
	if pass {
		http.Redirect(w, r, C.u.menu, http.StatusFound)
	} else {
		C.renderTemplate(C.temp.login, newFlash(w, "Username or password is incorrect."))
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clears session and returns to login page
	clearSession(w, r)
	http.Redirect(w, r, C.u.source, http.StatusFound)
}

func changeHandler(w http.ResponseWriter, r *http.Request) {
	// Renders change password form
	handleRender(w, r, C.temp.change, C.temp.login, "Please login to access database.")
}

func passwordHandler(w http.ResponseWriter, r *http.Request) {
	// Renders change password form
	msg := "Please login to access database."
	template := C.temp.login
	user, pw, _ := getCredentials(w, r)
	if user != "" && pw != "" {
		// Redirect to same page if an error occurs
		template = C.temp.change
		msg = changePassword(r, user, pw)
		if msg == "" {
			// Logout and return to login page
			msg = "Successfully changed password."
			template = C.temp.login
			clearSession(w, r)
		}
	}
	C.renderTemplate(template, newFlash(w, msg))
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	// Renders menu page
	handleRender(w, r, C.temp.menu, C.temp.login, "Please login to access database.")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	// Renders search form
	handleRender(w, r, C.temp.search, C.temp.login, "Please login to access database.")
}

func summaryHandler(w http.ResponseWriter, r *http.Request) {
	// Performs and renders database summary
	handlePost(w, r, C.u.summary)
}

func cancerRateHandler(w http.ResponseWriter, r *http.Request) {
	// Handles cancer rate calculations
	handlePost(w, r, C.u.rates)
}

func tableDumpHandler(w http.ResponseWriter, r *http.Request) {
	// Handles full table extraction
	handlePost(w, r, C.u.table)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Reads search form
	handlePost(w, r, C.u.search)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Serves output files for download
	user, pw, _ := getCredentials(w, r)
	if user != "" && pw != "" {
		vars := mux.Vars(r)
		http.ServeFile(w, r, fmt.Sprintf("/tmp/%s", vars["filename"]))
	} else {
		C.renderTemplate(C.temp.login, newFlash(w, "Please login to access database."))
	}
}

func main() {
	// Initilaize multiplexer and fileserver
	var (
		host = kingpin.Flag("host", "Host IP (default is localHost).").Short('h').Default("127.0.0.1").String()
		port = kingpin.Flag("port", "Host port (default is 8080).").Short('p').Default("8080").String()
	)
	// Set max cookie age to 4 hours
	STORE.MaxAge(14400)
	kingpin.Parse()
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("." + C.u.static))
	http.Handle(C.u.static, http.StripPrefix(C.u.static, fs))
	// Register handler functions
	r.HandleFunc(C.u.source, indexHandler).Methods(http.MethodGet)
	r.HandleFunc(C.u.login, loginHandler).Methods(http.MethodPost)
	r.HandleFunc(C.u.logout, logoutHandler).Methods(http.MethodGet)
	r.HandleFunc(C.u.changepw, changeHandler).Methods(http.MethodGet)
	r.HandleFunc(C.u.newpw, passwordHandler).Methods(http.MethodPost)
	r.HandleFunc(C.u.menu, menuHandler).Methods(http.MethodGet)
	r.HandleFunc(C.u.rates, cancerRateHandler).Methods(http.MethodPost)
	r.HandleFunc(C.u.table, tableDumpHandler).Methods(http.MethodPost)
	r.HandleFunc(C.u.search, formHandler).Methods(http.MethodGet)
	r.HandleFunc(C.u.summary, summaryHandler).Methods(http.MethodGet)
	r.HandleFunc(C.u.output, searchHandler).Methods(http.MethodPost)
	r.HandleFunc(C.u.get+"{filename}", downloadHandler).Methods(http.MethodGet)
	// Serve and log errors to terminal
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", *host, *port), nil))
}

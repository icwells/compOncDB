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
	"time"
)

var (
	STORE = sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	C     = setConfiguration()
)

func clearSession(w http.ResponseWriter, r *http.Request) {
	// Deletes username and password cookies
	session, _ := STORE.Get(r, C.name)
	session.Values["timestamp"] = ""
	session.Values["username"] = ""
	session.Values["password"] = ""
	session.Save(r, w)
}

func getTimestamp() string {
	// Returns date and time as string
	return time.Now().Format(time.RFC1123Z)
}

func updateTimestamp(w http.ResponseWriter, r *http.Request) {
	// Updates timestamp cookie
	session, _ := STORE.Get(r, C.name)
	session.Values["timestamp"] = getTimestamp()
	session.Save(r, w)
}

func checkTimestamp(stamp string) bool {
	// Requires login after one hour of inactivity
	timestamp, err := time.Parse(time.RFC1123Z, stamp)
	if err == nil {
		return time.Since(timestamp) < time.Hour
	} else {
		return false
	}
}

func getCredentials(w http.ResponseWriter, r *http.Request) (string, string) {
	// Reads username and password from cookie
	var user, password string
	session, _ := STORE.Get(r, C.name)
	stamp, exists := session.Values["timestamp"]
	if (exists && checkTimestamp(stamp.(string))) || exists == false {
		// Proceed if stamp has been updated in the last hour
		name, ex := session.Values["username"]
		if ex {
			user = name.(string)
			pw, e := session.Values["password"]
			if e {
				password = pw.(string)
				updateTimestamp(w, r)
			}
		}
	}
	return user, password
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serves login page
	user, pw := getCredentials(w, r)
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
			session.Values["timestamp"] = getTimestamp()
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
	clearSession(w, r)
	http.Redirect(w, r, C.source, http.StatusFound)
}

func changeHandler(w http.ResponseWriter, r *http.Request) {
	// Renders change password form
	user, _ := getCredentials(w, r)
	if user != "" {
		C.renderTemplate(w, C.changetemp, newOutput(user))
	} else {
		C.renderTemplate(w, C.logintemp, newFlash("Please login to access database."))
	}
}

func passwordHandler(w http.ResponseWriter, r *http.Request) {
	// Renders change password form
	msg := "Please login to access database."
	template := C.logintemp
	user, pw := getCredentials(w, r)
	if user != "" && pw != "" {
		// Redirect to same page if an error occurs
		template = C.changetemp
		msg = changePassword(r, user, pw)
		if msg == "" {
			// Logout and return to login page
			msg = "Successfully changed password."
			template = C.logintemp
			clearSession(w, r)
		}
	}
	C.renderTemplate(w, template, newFlash(msg))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	// Renders search form (newOutput supplies username)
	user, _ := getCredentials(w, r)
	if user != "" {
		C.renderTemplate(w, C.searchtemp, newOutput(user))
	} else {
		C.renderTemplate(w, C.logintemp, newFlash("Please login to access database."))
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Reads search form
	user, pw := getCredentials(w, r)
	if user != "" && pw != "" {
		out, err := extractFromDB(r, user, pw)
		if err == nil {
			if out.Flash != "" {
				// Return to search page with flash message
				C.renderTemplate(w, C.searchtemp, out)
			} else {
				C.renderTemplate(w, C.resulttemp, out)
			}
		} else {
			// Return to login page if an error is encoutered (error occurs at connection)
			C.renderTemplate(w, C.logintemp, newFlash(err.Error()))
		}
	} else {
		C.renderTemplate(w, C.logintemp, newFlash("Please login to access database."))
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Serves output files for download
	user, pw := getCredentials(w, r)
	if user != "" && pw != "" {
		vars := mux.Vars(r)
		http.ServeFile(w, r, fmt.Sprintf("/tmp/%s", vars["filename"]))
	} else {
		C.renderTemplate(w, C.logintemp, newFlash("Please login to access database."))
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
	fs := http.FileServer(http.Dir("." + C.static))
	http.Handle(C.static, http.StripPrefix(C.static, fs))
	// Register handler functions
	r.HandleFunc(C.source, indexHandler).Methods(http.MethodGet)
	r.HandleFunc(C.login, loginHandler).Methods(http.MethodPost)
	r.HandleFunc(C.logout, logoutHandler).Methods(http.MethodGet)
	r.HandleFunc(C.changepw, changeHandler).Methods(http.MethodGet)
	r.HandleFunc(C.newpw, passwordHandler).Methods(http.MethodPost)
	r.HandleFunc(C.search, formHandler).Methods(http.MethodGet)
	r.HandleFunc(C.output, searchHandler).Methods(http.MethodPost)
	r.HandleFunc(C.get+"{filename}", downloadHandler).Methods(http.MethodGet)
	// Serve and log errors to terminal
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", *host, *port), nil))
}

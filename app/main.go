// Main script for compOncDB web user interface

package main

import (
	"fmt"
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
	C.renderTemplate(w, C.logintemp, newOutput(""))
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
		C.renderTemplate(w, C.logintemp, newFlash("", "Username or password is incorrect."))
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clears session and returns to login page
	session, _ := STORE.Get(r, C.name)
	session.Values["username"] = ""
	session.Values["password"] = ""
	session.Save(r, w)
	http.Redirect(w, r, C.source, http.StatusFound)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	// Renders search form (newOutput supplies username)
	user, _ := getCredentials(r)
	if user != "" {
		fmt.Println(user)
		C.renderTemplate(w, C.searchtemp, newOutput(user))
	} else {
		C.renderTemplate(w, C.logintemp, newFlash("", "Please login to access database."))
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Reads search form
	user, pw := getCredentials(r)
	if user != "" && pw != "" {
		form := setSearchForm(r)
		out := extractFromDB(form, user, pw)
		C.renderTemplate(w, C.resulttemp, out)
	} else {
		C.renderTemplate(w, C.logintemp, newFlash("", "Please login to access database."))
	}
}

/*func staticCache(h http.Handler) http.Handler {
	// Adds cache time and content type to static files
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ct string
		path := r.URL.Path[1:]
		fmt.Println(path)
		if strings.HasSuffix(path, ".css") {
			ct = "text/css"
		} else if strings.HasSuffix(path, "png") {
			ct = "image/png"
		}
		w.Header().Add("Content-Type", ct)
		// 2 Days
        w.Header().Set("Cache-Control", "max-age=172800")
        h.ServeHTTP(w, r)
    })
}*/

func main() {
	// Initilaize multiplexer and fileserver
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("." + C.static))
	r.PathPrefix(C.static).Handler(http.StripPrefix(C.static, fs))
	// Register handler functions
	r.HandleFunc(C.source, indexHandler).Methods(http.MethodGet)
	r.HandleFunc(C.source, loginHandler).Methods(http.MethodPost)
	r.HandleFunc(C.logout, logoutHandler).Methods(http.MethodPost)
	r.HandleFunc(C.search, formHandler).Methods(http.MethodGet)
	r.HandleFunc(C.search, searchHandler).Methods(http.MethodPost)
	// Serve and log errors to terminal
	http.Handle(C.source, r)
	//log.Fatal(http.ListenAndServe(C.config.Host + ":8080", nil))
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

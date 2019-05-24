// Main script for compOncDB web user interface

package main

import (
	"log"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"net/http"
)

var(
	cookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	)
	router = mux.NewRouter()
)

func main() {
	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	}
	fileServer := http.FileServer(http.Dir("static/"))
	// Register functions
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/search", searchHandler)
	router.HandleFunc("/output", outputHandler)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	// Serve and log errors to terminal
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

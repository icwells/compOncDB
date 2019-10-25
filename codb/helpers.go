// Defines helper functions for handlers

package main

import (
	"net/http"
	"time"
)

func clearSession(w http.ResponseWriter, r *http.Request) {
	// Deletes username and password cookies
	session, _ := STORE.Get(r, C.name)
	session.Values["timestamp"] = ""
	session.Values["username"] = ""
	session.Values["password"] = ""
	session.Values["updatetime"] = ""
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

func getCredentials(w http.ResponseWriter, r *http.Request) (string, string, string) {
	// Reads username, password, and last update from cookie
	var user, password, update string
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
			ut, exists := session.Values["updatetime"]
			if exists {
				update = ut.(string)
			}
		}
	}
	return user, password, update
}

func handleRender(w http.ResponseWriter, r *http.Request, target, def, msg string) {
	// Handles basic credential check and redirect
	user, _, update := getCredentials(w, r)
	if user != "" {
		C.renderTemplate(w, target, newOutput(user, update))
	} else {
		C.renderTemplate(w, def, newFlash(msg))
	}
}

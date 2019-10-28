// Defines helper functions for handlers

package main

import (
	"fmt"
	"github.com/icwells/dbIO"
	"net/http"
	"time"
)

func ping(user, password string) (bool, string) {
	// Returns true if credentials are valid
	var update string
	ret := dbIO.Ping(C.config.Host, C.config.Database, user, password)
	if ret {
		db, _ := dbIO.Connect(C.config.Host, C.config.Database, user, password)
		db.GetTableColumns()
		update = db.LastUpdate().Format(time.RFC822)

	}
	return ret, update
}

func changePassword(r *http.Request, user, password string) string {
	// Changes suer password or returns flash message
	var ret string
	db, err := dbIO.Connect(C.config.Host, C.config.Database, user, password)
	if err == nil {
		r.ParseForm()
		newpw := r.PostForm.Get("password")
		confpw := r.PostForm.Get("newpassword")
		if newpw != confpw {
			ret = "Passwords do not match."
		} else {
			cmd := fmt.Sprintf("SET PASSWORD = PASSWORD('%s')", newpw)
			_, er := db.DB.Exec(cmd)
			if er != nil {
				ret = er.Error()
			}
		}
	} else {
		// Convert error to string
		ret = err.Error()
	}
	return ret
}

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
		o, err := newOutput(user, "", update)
		if err == nil {
			C.renderTemplate(w, target, o)
		} else {
			o.Flash = err.Error()
			C.renderTemplate(w, def, o)
		}
	} else {
		C.renderTemplate(w, def, newFlash(msg))
	}
}

// Defines SearchForm

package main

import (
	"github.com/gorilla/schema"
	"github.com/icwells/compOncDB/src/codbutils"
	"net/http"
	"strconv"
	"strings"
)

type Options struct {
	Min      int
	Necropsy bool
	Count    bool
	Infant   bool
	Print    bool
}

func setOptions(r *http.Request) *Options {
	// Returns commonly used options
	opt := new(Options)
	decoder := schema.NewDecoder()
	decoder.Decode(opt, r.PostForm)
	return opt
}

func setEvaluation(r *http.Request, columns map[string]string, search, n string) (codbutils.Evaluation, string) {
	// Populates evalutaiton struct if matching term is found
	var e codbutils.Evaluation
	var msg string
	e.Table = strings.TrimSpace(r.PostForm.Get(search + "Table" + n))
	e.Column = strings.TrimSpace(r.PostForm.Get(search + "Column" + n))
	e.Operator = strings.TrimSpace(r.PostForm.Get(search + "Operator" + n))
	e.Value = strings.TrimSpace(r.PostForm.Get(search + "Value" + n))
	slct := strings.TrimSpace(r.PostForm.Get(search + "Select" + n))
	if e.Value != "" || slct != "" {
		if e.Value == "" {
			// Assign selected value to input value
			e.Value = slct
		}
		if e.Table != "" && e.Table != "Empty" && e.Column != "" && e.Column != "Empty" && e.Operator != "" && e.Value != "" {
			e.SetIDType(columns)
			if e.Table == "Accounts" {
				msg = "Cannot access Accounts table."
			}
		} else {
			msg = "Please supply valid table, column, and value."
		}
	}
	return e, msg
}

func checkEvaluations(r *http.Request, columns map[string]string) (map[string][]codbutils.Evaluation, string) {
	// Reads in variable number of search parameters
	var msg string
	var e codbutils.Evaluation
	eval := make(map[string][]codbutils.Evaluation)
	for i := 0; i < 10; i++ {
		// Loop through all possible searches since deletions might disrupt numerical order
		found := true
		count := 0
		for found == true {
			// Keep checking until count exceeds number of fields
			search := strconv.Itoa(i)
			e, msg = setEvaluation(r, columns, search, strconv.Itoa(count))
			if msg != "" {
				found = false
			} else {
				eval[search] = append(eval[search], e)
				count++
			}
		}
	}
	return eval, msg
}

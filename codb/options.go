// Defines SearchForm

package main

import (
	"fmt"
	"github.com/gorilla/schema"
	"github.com/icwells/compOncDB/src/codbutils"
	"net/http"
	"strconv"
	"strings"
)

type Options struct {
	AllTaxa     bool
	Cancerrate  bool
	Infant      bool
	Key         string
	Lifehistory bool
	Min         int
	Necropsy    bool
	Print       bool
}

func setOptions(r *http.Request) *Options {
	// Returns commonly used options
	opt := new(Options)
	decoder := schema.NewDecoder()
	decoder.Decode(opt, r.PostForm)
	fmt.Println(opt.Key)
	return opt
}

func setEvaluation(r *http.Request, columns map[string]string, search, n string) (codbutils.Evaluation, string, bool) {
	// Populates evalutaiton struct if matching term is found
	pass := false
	var e codbutils.Evaluation
	var msg string
	e.Table = strings.TrimSpace(r.PostForm.Get(search + "Table" + n))
	e.Column = strings.TrimSpace(r.PostForm.Get(search + "Column" + n))
	e.Operator = strings.TrimSpace(r.PostForm.Get(search + "Operator" + n))
	e.Value = strings.TrimSpace(r.PostForm.Get(search + "Value" + n))
	slct := strings.TrimSpace(r.PostForm.Get(search + "Select" + n))
	if e.Table == "Accounts" {
		msg = "Cannot access Accounts table."
	} else {
		for _, i := range []string{e.Table, e.Column, e.Operator} {
			if len(i) >= 2 && i != "Empty" {
				pass = true
			}
		}
		if pass {
			if e.Value == "" && slct != "" {
				// Assign selected value to input value
				e.Value = slct
			}
			e.SetIDType(columns)
		} else {
			msg = "Please supply valid table, column, and value."
		}
	}
	return e, msg, pass
}

func checkEvaluations(r *http.Request, columns map[string]string) ([][]codbutils.Evaluation, string) {
	// Reads in variable number of search parameters
	var ret [][]codbutils.Evaluation
	var msg string
	var e codbutils.Evaluation
	eval := make(map[string][]codbutils.Evaluation)
	for i := 0; i < 10; i++ {
		// Loop through all possible searches since deletions might disrupt numerical order
		found := true
		count := 0
		for found == true {
			// Keep checking until count exceeds number of fields
			var m string
			search := strconv.Itoa(i)
			e, m, found = setEvaluation(r, columns, search, strconv.Itoa(count))
			if found {
				eval[search] = append(eval[search], e)
				count++
			} else if len(eval) == 0 {
				msg = m
			}
		}
	}
	for _, v := range eval {
		ret = append(ret, v)
	}
	return ret, msg
}

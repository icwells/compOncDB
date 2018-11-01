// This script will perform white box tests on parseRecords diagnosis functions

package main

import (
	"fmt"
	"testing"
)

func TestCountNA(t *testing.T) {
	// Tests coutn NA method
	nas := []struct {
		row      []string
		found    bool
		complete bool
	}{
		{[]string{"1", "12", "male", "Y", "Liver", "neoplasm", "N", "N", "N", "Y"}, true, true},
		{[]string{"2", "12", "female", "NA", "NA", "NA", "NA", "NA", "NA", "NA"}, true, false},
		{[]string{"3", "12", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA"}, false, false},
		{[]string{"NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA"}, false, false},
	}
	for _, i := range nas {
		found, complete := countNA(i.row)
		if found != i.found || complete != i.complete {
			msg := fmt.Sprintf("countNA returned %v, %v instead of %v, %v.", found, complete, i.found, i.complete)
			t.Error(msg)
		}
	}
}

func TestCheckAge(t *testing.T) {
	// 
}

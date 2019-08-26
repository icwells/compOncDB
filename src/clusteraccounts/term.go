// Struct for storing term info

package clusteraccounts

import (
	"strconv"
	"strings"
)

var (
	ZOOS       = []string{"zoo", "aquarium", "museum"}
	INSTITUTES = []string{"center", "institute", "service", "research", "rescue"}
)

type term struct {
	query, name, id   string
	length, zoo, inst int
}

func newTerm(q, n, id string) *term {
	// Initializes term struct
	if id == "NA" {
		id = ""
	}
	t := new(term)
	t.query = q
	t.name = n
	t.id = ""
	t.length = strings.Count(n, " ") + 1
	t.zoo = 0
	t.inst = 0
	return t
}

func (t *term) toSlice() []string {
	// Returns slice for map entry
	return []string{t.name, strconv.Itoa(t.zoo), strconv.Itoa(t.inst)}
}

func (t *term) setType() {
	// Sets 1/0 for zoo/institute columns
	found := false
	n := strings.ToLower(t.name)
	for _, i := range INSTITUTES {
		if strings.Contains(n, i) {
			t.inst = 1
			found = true
			break
		}
	}
	if found == false {
		for _, i := range ZOOS {
			if strings.Contains(n, i) {
				t.zoo = 1
				break
			}
		}
	}
}

func (t *term) getType() string {
	// Returns type as string
	if t.zoo == 1 {
		return "zoo"
	} else if t.inst == 1 {
		return "inst"
	} else {
		return "other"
	}
}

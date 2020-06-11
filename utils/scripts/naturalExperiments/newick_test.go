// Test functions for newick tree

package main

import (
	"testing"
)

// Returns newick tree for testing
func testTree() string {
	return "(A:0.1,B:0.2,(C:0.3,D:0.4)E:0.5)F:0.1;"
}

func testCases() map[string]string {
	return map[string]string{
		"A": "A:0.1",
		"B": "B:0.2",
		"C": "C:0.3",
		"D": "D:0.4",
		"E": "E:0.5",
		"F": "F:0.1",
	}
}

func TestNewTree(t *testing.T) {
	tree := FromString(testTree())
	cases := testCases()
	root := tree.Root.String()
	if root != cases["F"] {
		t.Errorf("Root node %s does not equal %s", root, cases["F"])
	} else if tree.Root.Ancestor != nil {
		t.Error("Root node has ancestor node.")
	}
	for k, exp := range cases {
		v, ex := tree.nodes[k]
		if !ex {
			t.Errorf("%s not found in nodes map", k)
		} else if v.String() != exp {
			t.Errorf("Actual node value %s does not equal expected: %s", v.String(), exp)
		} else if k != "F" && tree.nodes[k].Ancestor == nil {
			t.Errorf("Node %s ancestor was not stored.", k)
			//} else {
			//t.Error(k, tree.nodes[k].Descendants)
		}
	}
}

func TestDivergence(t *testing.T) {
	tree := FromString(testTree())
	cases := []struct {
		a    string
		b    string
		dist float64
	}{
		{"A", "B", 0.2},
		{"E", "B", 0.5},
		{"C", "D", 0.4},
		{"C", "E", 0.3},
		{"F", "D", 0.9},
	}
	for _, i := range cases {
		act := tree.Divergence(i.a, i.b)
		if act != i.dist {
			t.Errorf("Actual distance between %s and %s %f does not equal expected: %f", i.a, i.b, act, i.dist)
		}
	}
}

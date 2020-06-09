// Defines structs for reading and storing Newick trees

package main

import (
	"github.com/icwells/go-tools/iotools"
	"strconv"
	"strings"
)

// Node stores data for each node of the tree.
type Node struct {
	Ancestor    *Node
	Descendants []*Node
	Length      float64
	Name        string
}

// NewNode returns new node struct.
func NewNode(name string, length float64, descendants []*Node) *Node {
	n := new(Node)
	n.Length = length
	n.Name = name
	for _, i := range descendants {
		n.AddDescendant(i)
	}
	return n
}

// AddDescendant appends a new descendant to the node.
func (n *Node) AddDescendant(d *Node) {
	d.Ancestor = n
	n.Descendants = append(n.Descendants, d)
}

// IsLeaf returns true if node has no descendants
func (n *Node) IsLeaf() bool {
	if len(n.Descendants) == 0 {
		return true
	}
	return false
}

// NewickTree stores nodes for parsing.
type NewickTree struct {
	nodes map[string]*Node
	root  *Node
}

// NewTree returns a Newick tree struct from the given string
func NewTree(tree string) *NewickTree {
	var err error
	t := new(NewickTree)
	t.nodes = make(map[string]*Node)
	t.root = t.parseNodes(tree)
	return t
}

// parseName returns the node name and length.
func (t *NewickTree) parseName(s string) (string, float64) {
	var length float64
	var name string
	if strings.Contains(s, ":") {
		n := strings.Split(s, ":")
		name = strings.TrimSpace(n[0])
		if len(n) > 1 {
			length, _ = strconv.ParseFloat(n[1], 64)
		}
	}
	return name, length
}

func (t *NewickTree) parseSiblings(s string) <-chan *Node {
	var level int
	var builder strings.Builder
	ch := make(chan *Node)
	// Remove special-case of trailing chars
	for _, c := range s + "," {
		if c == ',' && level == 0 {
			// Recursively submits entries on the same level
			go func() {
				ch <- t.parseNodes(builder.String())
				builder.Reset()
				close(ch)
			}()
		} else {
			if c == '(' {
				level++
			} else if c == ')' {
				level--
			}
			builder.WriteRune(c)
		}
	}
	return ch
}

// parseNodes parses string into node structs.
func (t *NewickTree) parseNodes(s string) *Node {
	var err error
	var descendants []*Node
	parts := strings.Split(s, ")")
	label := s
	if len(parts) > 1 {
		for d := range t.parseSiblings(strings.Join(parts[:len(parts)-1][1:], ")")) {
			descendants = append(descendants, d)
		}
		label = parts[len(parts)-1]
	}
	name, length := t.parseName(label)
	t.nodes[name] = NewNode(name, length, descendants)
	return t.nodes[name]
}

// Divergeance returns the sum of lengths between two nodes.
/*func (t *NewickTree) Divergeance(a, b string) float64 {
	var ret float64

	return ret
}*/

// FromString returns a Newick tree from the given string
func FromString(tree string) *NewickTree {
	tree = strings.Replace(strings.TrimSpace(tree), ";", "", 1)
	return NewTree(tree)
}

// FromFile reads a single Newick tree from the given file.
func FromFile(infile string) *NewickTree {
	var line string
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := iotools.GetScanner(f)
	for input.Scan() {
		line = strings.TrimSpace(string(input.Text()))
		break
	}
	return FromString(line)
}

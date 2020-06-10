// Defines structs for reading and storing Newick trees

package main

import (
	//"fmt"
	"github.com/icwells/go-tools/iotools"
	"strconv"
	"strings"
)

// NewickTree stores nodes for parsing.
type NewickTree struct {
	nodes map[string]*Node
	root  *Node
}

// NewTree returns a Newick tree struct from the given string
func NewTree(tree string) *NewickTree {
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

func (t *NewickTree) parseSiblings(s string) []*Node {
	var level int
	var ret []*Node
	var builder strings.Builder
	ch := make(chan *Node)
	// Remove special-case of trailing chars
	for _, c := range s + "," {
		if c == ',' && level == 0 {
			// Recursively submits entries on the same level
			go func() {
				ch <- t.parseNodes(builder.String())
			}()
			d := <-ch
			if d != nil {
				ret = append(ret, d)
			}
			builder.Reset()
		} else {
			if c == '(' {
				level++
			} else if c == ')' {
				level--
			}
			builder.WriteRune(c)
		}
	}
	close(ch)
	return ret
}

// parseNodes parses string into node structs.
func (t *NewickTree) parseNodes(s string) *Node {
	var descendants []*Node
	parts := strings.Split(s, ")")
	label := s
	if len(parts) > 1 {
		// Recusively append descendants
		for _, d := range t.parseSiblings(strings.Join(parts[:len(parts)-1], ")")[1:]) {
			descendants = append(descendants, d)
		}
		label = parts[len(parts)-1]
	}
	name, length := t.parseName(label)
	t.nodes[name] = NewNode(name, length, descendants)
	return t.nodes[name]
}

// walkBack traverses the tree in reverse, starting from given node.
func (t *NewickTree) walkBack(name string) []*Node {
	var ret []*Node
	n := t.nodes[name]
	for n.Name != t.root.Name {
		ret = append([]*Node{n}, ret...)
		n = n.Ancestor
	}
	return append([]*Node{t.root}, ret...)
}

// Divergeance returns the sum of lengths between two nodes.
func (t *NewickTree) Divergeance(a, b string) float64 {
	var ret float64
	apath := t.walkBack(a)
	bpath := t.walkBack(b)
	for idx, i := range apath {
		if i.Name != bpath[idx].Name {
			// Record where paths diverge
			apath = apath[idx:]
			bpath = bpath[idx:]
			break
		}
	}
	for _, path := range [][]*Node{apath, bpath} {
		for _, i := range path {
			ret += i.Length
		}
	}
	return ret
}

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

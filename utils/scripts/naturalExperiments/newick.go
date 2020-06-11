// Defines structs for reading and storing Newick trees

package main

import (
	//"fmt"
	"math"
	"strconv"
	"strings"
)

// NewickTree stores nodes for parsing.
type NewickTree struct {
	idx   int
	nodes map[string]*Node
	Root  *Node
}

// NewTree returns a Newick tree struct from the given string
func NewTree(tree string) *NewickTree {
	t := new(NewickTree)
	t.nodes = make(map[string]*Node)
	t.Root = t.parseNodes(tree)
	return t
}

// parseName returns the node name and length.
func (t *NewickTree) parseName(s string) (string, float64) {
	var length float64
	var name string
	s = strings.Replace(s, "'", "", -1)
	if strings.Contains(s, ":") {
		n := strings.Split(s, ":")
		name = strings.TrimSpace(n[0])
		if len(n) > 1 {
			length, _ = strconv.ParseFloat(n[1], 64)
		}
	} else if l, err := strconv.ParseFloat(name, 64); err == nil {
		length = l
		t.idx++
		name = strconv.Itoa(t.idx)
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
	if n, ex := t.nodes[name]; ex {
		for n.Ancestor != nil {
			ret = append([]*Node{n}, ret...)
			n = n.Ancestor
		}
		ret = append([]*Node{t.Root}, ret...)
	}
	return ret
}

// totalLength returns the length of a given branch.
func (t *NewickTree) totalLength(s []*Node) float64 {
	var ret float64
	for _, i := range s {
		if i.Name != t.Root.Name {
			ret += i.Length
		}
	}
	return ret
}

// Divergence returns the sum of lengths between two nodes.
func (t *NewickTree) Divergence(a, b string) float64 {
	var ret float64
	apath := t.walkBack(a)
	bpath := t.walkBack(b)
	l := len(apath)
	if len(bpath) < l {
		l = len(bpath)
	}
	if l > 0 {
		for idx := 0; idx < l; idx++ {
			if apath[idx].Name != bpath[idx].Name {
				// Record where paths diverge
				apath = apath[idx:]
				bpath = bpath[idx:]
				ret = math.Max(t.totalLength(apath), t.totalLength(bpath))
				break
			} else if idx == l-1 {
				// Account for differnce in length: last entry in one will be an ancestor in the other
				if apath[idx].Name == a {
					ret = t.totalLength(bpath[idx+1:])
				} else if bpath[idx].Name == b {
					ret = t.totalLength(apath[idx+1:])
				}
			}
		}
	}
	return ret
}

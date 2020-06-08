// Defines structs for reading and storing Newick trees

package main

import (
	"github.com/icwells/go-tools/iotools"
	"strings"
)

// Node stores data for each node of the tree.
type Node struct {
	Ancestor string
	Descendants []*Node
	Length float64
	Name string
}

// NewNode returns new node struct.
func NewNode(length float64) *Node {
	n := new(Node)
	n.Length = length
	n.Name = name
	return n
}

// NewickTree stores nodes for parsing.
type NewickTree struct {
	keys map[string][]int
	nodes [][]*Node
}

// NewTree returns a Newick tree struct from the given string
func NewTree(tree string) (*NewickTree, error) {
	t := new(NewickTree)
	t.keys = make(map[string][]int)
	err := t.parseNodes(tree)
	return t, err
}

def _parse_name_and_length(s):
    length = None
    if ':' in s:
        s, length = s.split(':', 1)
    return s or None, length or None


def _parse_siblings(s, **kw):
    bracket_level = 0
    current = []

    # trick to remove special-case of trailing chars
    for c in (s + ","):
        if c == "," and bracket_level == 0:
            yield parse_node("".join(current), **kw)
            current = []
        else:
            if c == "(":
                bracket_level += 1
            elif c == ")":
                bracket_level -= 1
            current.append(c)

// parseNodes parses string into node structs.
func (t *NewickTree) parseNodes(s) {
	var err error
    parts = strings.Split(s, ")")
    if len(parts) == 1 {
        descendants, label = [], s
    } else {
        if parts[0] != '(' {
            err = fmt.Error("unmatched braces %s", parts[0][:100])
		}
        descendants = list(_parse_siblings(')'.join(parts[:-1])[1:], **kw))
        label = parts[-1]
    name, length = _parse_name_and_length(label)
    return Node.create(name=name, length=length, descendants=descendants, **kw)
}

// FromString returns a Newick tree from the given string
func FromString(tree string) (*NewickTree, error) {
	tree = strings.Replace(strings.TrimSpace(tree), ";", "", 1)
	return NewTree(tree)
}

// FromFile reads a single Newick tree from the given file.
func FromFile(infile string) (*NewickTree, error) {
	var line string
	iotools.CheckFile(infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := iotools.GetScanner()
	defer input.Close()
	for input.Scan() {
		line = strings.TrimSpace(string(input.Text()))
		break
	}
	return FromString(line)
}

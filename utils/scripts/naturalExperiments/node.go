// Defines node struct for newick tree

package main

import (
	"fmt"
	"strconv"
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

// String returns node anme and length as a string
func (n *Node) String() string {
	return fmt.Sprintf("%s:%s", n.Name, strconv.FormatFloat(n.Length, 'f', -1, 64))
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

// Walk traverses tree starting from this node.
/*func (n *Node) Walk() <-chan *Node {
	ch := make(chan *Node)
	ch <- n
	go func() {
		for _, i := range n.Descendants {
			for j := range i.Walk() {
				ch <- j
				if j.IsLeaf() {
					break
				}
			}
		}
		close(ch)
	}()
	return ch
}*/

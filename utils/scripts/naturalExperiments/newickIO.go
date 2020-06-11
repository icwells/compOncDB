// Conctains input/output functions for newick tree

package main

import (
	"github.com/icwells/go-tools/iotools"
	"strings"
)

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

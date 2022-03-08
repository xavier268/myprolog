package inter

import (
	"fmt"
	"strings"
)

// Node is the node/tree structure for prolog objects.
type Node struct {
	name     string  // printable name
	constant bool    // true if there are no Variable in the whole tree. Hence, cloning can be avoided.
	args     []*Node // Children nodes
}

// lastArg last child of a node, nil if no child.
func (n *Node) lastArg() *Node {
	if len(n.args) == 0 {
		return nil
	}
	return n.args[len(n.args)-1]
}

// printPrefix used to pint node trees.
const printPrefix = "  "

// StringIndent provides a human readeable indented form for the object.
func (n *Node) StringIndent() string {
	return n.stringPrefix("")
}

// stringPrefix prints node with the given prefix.
func (n *Node) stringPrefix(pfx string) string {
	if len(n.args) == 0 {
		return n.name + " "
	} else {
		var b strings.Builder
		fmt.Fprintf(&b, "%s (\n%s ", n.name, pfx)
		for _, nn := range n.args {
			fmt.Fprintf(&b, "%s ", nn.stringPrefix(pfx+printPrefix))
		}
		fmt.Fprintf(&b, "\n%s ) ", pfx)
		return b.String()
	}
}

func (n *Node) String() string {
	if len(n.args) == 0 {
		return n.name + " "
	} else {
		var b strings.Builder
		fmt.Fprintf(&b, "%s ( ", n.name)
		for _, nn := range n.args {
			fmt.Fprintf(&b, "%s", nn)
		}
		fmt.Fprint(&b, ") ")
		return b.String()
	}
}

// StringType describes the type of node.
func (n *Node) StringType() string {
	var b string
	if isVariable(n.name) {
		b += "V"
	} else {
		b += "."
	}
	if isNumber(n.name) {
		b += "N"
	} else {
		b += "."
	}
	if n.constant {
		b += "C"
	} else {
		b += "."
	}
	return b
}

// dump display detailled structure of tree.
func (n *Node) dump() {

	fmt.Printf("%p  %s %10s : ", n, n.StringType(), n.name)

	for _, s := range n.args {
		fmt.Printf("%s ", s.name)
	}
	fmt.Println()
	for _, a := range n.args {
		a.dump()
	}
}

// ------------------------------------------------------------------

// WalkFunction used for Walk.
type WalkFunction func(*Node) error

// Walk the tree, applying the provided function to each node.
// If function returns an error, terminate immediately.
func (n *Node) Walk(f WalkFunction) error {
	if err := f(n); err != nil {
		return err
	}
	for _, nn := range n.args {
		if err := nn.Walk(f); err != nil {
			return err
		}
	}
	return nil
}

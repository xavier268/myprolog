package inter

import (
	"fmt"
	"strings"
)

// Node is the node/tree structure for prolog objects.
type Node struct {
	name string  // printable name
	args []*Node // Children nodes
}

// lastArg last child of a node, nil if no child.
func (n *Node) lastArg() *Node {
	if len(n.args) == 0 {
		return nil
	}
	return n.args[len(n.args)-1]
}

// StringPretty provides an idented human readeable string
func (n *Node) StringPretty() string {
	return n.stringPrefix("")
}

// stringPrefix prints node with the given prefix.
func (n *Node) stringPrefix(pfx string) string {
	const inc = "   "
	if len(n.args) == 0 {
		return n.name + " "
	} else {
		var b strings.Builder
		fmt.Fprintf(&b, "%s (\n%s", n.name, pfx+inc)
		for _, nn := range n.args {
			fmt.Fprintf(&b, "%s", nn.stringPrefix(pfx+inc))
		}
		fmt.Fprintf(&b, "\n%s) ", pfx)
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
	return b
}

// dumpTree display detailled structure of tree.
func (n *Node) DumpTree(withHeader bool) {
	if withHeader {
		dumpHeader()
	}
	var dedup map[*Node]bool = nil
	if len(n.args) != 0 {
		dedup = make(map[*Node]bool)
	}
	n.dumpTree(dedup)
}

func dumpHeader() {
	fmt.Println("pointer     flags       name / arity :  args ....")
}

// using dedup to deduplicate nodes.
func (n *Node) dumpTree(dedup map[*Node]bool) {

	if dedup == nil || !dedup[n] { // no dedup or not alredy seen ?
		if dedup != nil {
			dedup[n] = true // mark for next time
		}
		// dump node
		fmt.Printf("%p  %s %10s / %5d : ", n, n.StringType(), n.name, len(n.args))
		for _, s := range n.args {
			fmt.Printf("%s{%p} ", s.name, s)
		}
		fmt.Println()
	}
	// process childs
	for _, a := range n.args {
		a.dumpTree(dedup)
	}
}

/*
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
*/

// Equal recursively compares nodes for equality, based on their names.
// Actual memory address may differ.
func (n *Node) Equal(p *Node) bool {
	if p == nil {
		return n == nil
	}
	if n.name != p.name || len(n.args) != len(p.args) {
		return false
	}

	for i := 0; i < len(n.args); i++ {
		if !n.args[i].Equal(p.args[i]) {
			return false
		}
	}
	return true
}

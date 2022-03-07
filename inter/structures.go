package inter

import (
	"fmt"
	"strings"
)

// PTypes are the types of the leaves in the prolog tree.
type PType byte

const (
	Compound PType = iota
	Atom
	Number
	Variable
)

func (p PType) String() string {
	switch p {
	case Compound:
		return "Compound"
	case Atom:
		return "Atom"
	case Number:
		return "Number"
	case Variable:
		return "Variable"
	default:
		return "UNKNOWN"
	}
}

// Node is the node/tree structure for prolog objects.
type Node struct {
	name  string // printable name
	ptype PType  // type of object
	// immutable bool    // immutable means that there are no Variable in object. Hence, cloning is avoided.
	args []*Node // Children nodes
}

// String provides a human readeable form for the object.
func (n *Node) String() string {

	if n.ptype != Compound {
		return n.name
	} else {
		var b strings.Builder
		fmt.Fprintf(&b, "%s(", n.name)
		for i, nn := range n.args {
			fmt.Fprintf(&b, "%s", nn.String())
			if i < len(n.args)-1 {
				fmt.Fprint(&b, ",")
			}
		}
		fmt.Fprint(&b, ")")
		return b.String()
	}

}

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

// Package node implements the generic data structures and their related utilities.
package node

import (
	"fmt"
	"strconv"
	"strings"
)

type Node struct {
	load     any // load type defines the node type, should never be nil !
	children []*Node
}

func NewNode(name string) *Node {
	if len(name) == 0 {
		return nil
	}
	n := new(Node)
	if name == "_" {
		n.load = Underscore{}
		return n
	}
	if name[0] >= 'A' && name[0] <= 'Z' { // it's a variable node !
		n.load = Variable{name, 0}
		return n
	}
	f, err := strconv.ParseFloat(name, 64) // test if its a number ?
	if err == nil {
		n.load = Number{name, f}
		return n
	}
	// else, its a string !
	n.load = name
	return n
}

func (n *Node) GetLoad() any {
	if n == nil {
		return nil
	}
	return n.load
}

// Dump prints a detailled view of the node and its children, showing the type of each node.
func (n *Node) Dump() string {
	var sb strings.Builder
	n.dump(&sb, "")
	return sb.String()
}

func (n *Node) dump(sb *strings.Builder, prefix string) {

	if n == nil || n.load == nil {
		fmt.Fprintf(sb, "%s<Nil>", prefix)
		return
	}
	switch v := n.load.(type) {
	case Variable:
		if v.nsp == 0 {
			fmt.Fprintf(sb, "%s<Var>%s", prefix, v.name)
		} else {
			fmt.Fprintf(sb, "%s<Var>%s#%d", prefix, v.name, v.nsp)
		}
	case Underscore:
		fmt.Fprintf(sb, "%s<Uds>_", prefix)
	case Number:
		fmt.Fprintf(sb, "%s<Num>%v", prefix, v.value)
	case string:
		fmt.Fprintf(sb, "%s<Str>%s", prefix, v)
	default:
		panic("unimplemented load type for node")
	}

	prefix = prefix + "\t"
	for _, c := range n.children {
		fmt.Fprintln(sb)
		c.dump(sb, prefix)
	}

}

func (n *Node) String() string {
	if n == nil {
		return ""
	}
	var sb strings.Builder
	n.string(&sb)
	return sb.String()
}

func (n *Node) string(sb *strings.Builder) {
	if n == nil {
		return
	}
	switch v := n.load.(type) {
	case Variable:
		if v.nsp == 0 {
			fmt.Fprint(sb, v.name)
		} else {
			fmt.Fprintf(sb, "%s#%d", v.name, v.nsp)
		}
	case Underscore:
		fmt.Fprint(sb, "_")
	case Number:
		fmt.Fprint(sb, v.value)
	case string:
		fmt.Fprint(sb, v)
	default:
		panic("unimplemented load type for node")
	}
}

func (n *Node) Equal(m *Node) bool {
	if n == nil && m == nil {
		return true
	}
	if m == nil || n == nil {
		return false
	}
	if n.load != m.load {
		return false
	}
	if len(n.children) != len(m.children) {
		return false
	}
	for i, c := range n.children {
		if eq := c.Equal(m.children[i]); !eq {
			return false
		}
	}
	return true
}

func (n *Node) IsLeaf() bool {
	if n == nil {
		return true
	}
	return len(n.children) == 0
}

// Variable start with a capital letter A-Z.
// Variables have namespace versions, to differentiate instatiation of local variables.
type Variable struct {
	name string
	nsp  int // namespace version of variable
}

type Underscore struct{}

type Number struct {
	name  string
	value float64
}

func (n *Node) Clone() *Node {
	if n == nil {
		return nil
	}
	m := new(Node)
	m.load = n.load
	for _, c := range n.children {
		m.Add(c.Clone())
	}
	return m
}

func (n *Node) Add(children ...*Node) {
	if n == nil || len(children) == 0 {
		return
	}
	if _, ok := n.load.(Variable); ok {
		panic("trying to add children to a Variable node")
	}
	n.children = append(n.children, children...)
}

// Walk the Node, applying f recursively to all node loads, stopping immediately if error.
// (ie : generate an error if you want to break the walk ...)
func (n *Node) Walk(f func(load any) error) error {

	if n == nil || n.load == nil {
		return nil
	}
	if err := f(n.load); err != nil {
		return err
	}
	for _, c := range n.children {
		if err := c.Walk(f); err != nil {
			return err
		}
	}
	return nil
}

// LastArg last child of a node, nil if no child.
func (n *Node) LastArg() *Node {
	if len(n.children) == 0 {
		return nil
	}
	return n.children[len(n.children)-1]
}

// Returns the number of children.
func (n *Node) NbChildren() int {
	if n == nil {
		return 0
	}
	return len(n.children)
}

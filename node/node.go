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

func NewStringNode(s string) *Node {
	n := new(Node)
	n.load = String(s)
	return n
}

func NewNode(name string) *Node {
	if len(name) == 0 || name == "nil" {
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

	kw, err := NewKeyword(name)
	if err == nil {
		n.load = kw
		return n
	}

	f, err := strconv.ParseFloat(name, 64) // test if its a number ?
	if err == nil {
		n.load = Number{name, f}
		return n
	}
	// else, its just a string !
	n.load = String(name)
	return n
}

func (n *Node) GetLoad() any {
	if n == nil {
		return nil
	}
	return n.load
}

func (n *Node) GetChildren() []*Node {
	if n == nil || len(n.children) == 0 {
		return nil
	}
	return n.children
}

func (n *Node) GetChild(i int) *Node {
	if n == nil || i >= len(n.children) {
		return nil
	}
	return n.children[i]
}

func (n *Node) String() string {
	if n == nil {
		return " nil"
	}
	var sb strings.Builder
	n.string(&sb)
	return sb.String()
}

func (n *Node) ChildrenAllowed() bool {
	if n == nil {
		return false
	}
	switch n.GetLoad().(type) {
	case Variable, Underscore, Number:
		return false
	default:
		return true
	}
}

func (n *Node) string(sb *strings.Builder) {
	if n == nil {
		fmt.Fprint(sb, " nil")
		return
	}
	switch v := n.load.(type) {
	case Variable:
		if v.nsp == 0 {
			fmt.Fprintf(sb, " %s", v.name)
		} else {
			fmt.Fprintf(sb, " %s#%d", v.name, v.nsp)
		}
	case Underscore:
		fmt.Fprint(sb, " _")
	case Number:
		fmt.Fprintf(sb, " %v", v.value)
	case Keyword:
		fmt.Fprintf(sb, " %s", v.String())
	case String:
		fmt.Fprintf(sb, " %s", v)
	default:
		fmt.Printf("Debug : %#v %t\n", v, n.load)
		panic("unimplemented load type for node")
	}

	if len(n.children) != 0 {
		fmt.Fprint(sb, " (")
		for _, c := range n.children {
			c.string(sb)
		}
		fmt.Fprint(sb, " )")
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

type String string

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
	if len(children) == 0 {
		return
	}
	if n == nil || !n.ChildrenAllowed() {
		panic("trying to add children to a Variable node")
	}
	n.children = append(n.children, children...)
}

// Walk the Node, applying f recursively to all node loads, stopping immediately if error.
// (ie : generate an error if you want to break the walk ...)
// Caution : f can be called on a nil node ...
func (n *Node) Walk(f func(*Node) error) error {

	if err := f(n); err != nil {
		return err
	}
	if n == nil {
		return nil
	}
	for _, c := range n.children {
		if err := c.Walk(f); err != nil {
			return err
		}
	}
	return nil
}

// LastChild last child of a node, nil if no child.
func (n *Node) LastChild() *Node {
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

// ReplaceChild replaces the child in place. Replacing with nil is allowed.
func (n *Node) ReplaceChild(i int, c *Node) {
	if n == nil || i >= len(n.children) {
		return
	}
	n.children[i] = c
}

func (n *Node) RemoveChild(i int) {
	if n == nil || i >= len(n.children) {
		return
	}
	if i+1 < len(n.children) {
		n.children = append(n.children[:i], n.children[i+1:]...)
	} else {
		n.children = n.children[:i]
	}
}

func (n *Node) SwapChildren(i, j int) {
	if n == nil || i >= len(n.children) || j >= len(n.children) || i == j {
		return
	}
	n.children[i], n.children[j] = n.children[j], n.children[i]
}

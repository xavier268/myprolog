package parser

import (
	"fmt"
	"strings"
)

var (
	INTFORMAT   = "%d"
	FLOATFORMAT = "%.3e"
)

type Term interface { // Term is the most general form of a term
	String() string // String() is the string representation of the entire term
	// Strings are neither quoted nor escaped internally, they are stored without the start/end " or `
	Pretty() string // Pretty() is the string representation of the term, pretty printing lists and rules and queries.
	Clone() Term    // Clone() returns a deep copy of the term
}

var _ Term = Atom{}
var _ Term = Integer{}
var _ Term = Float{}
var _ Term = String{}
var _ Term = Variable{}
var _ Term = Underscore{}
var _ Term = CompoundTerm{}
var _ Term = Char{}

type Variable struct { // a named variable
	Name string
	Nsp  int // name space
}

// Clone implements Term.
func (t Variable) Clone() Term {
	return Variable{
		Name: t.Name,
		Nsp:  t.Nsp,
	}
}

// Pretty implements Term.
func (t Variable) Pretty() string {
	return t.String()
}

func (v Variable) String() string {
	if v.Nsp > 0 {
		return fmt.Sprintf("%v$%d", v.Name, v.Nsp)
	}
	return v.Name
}

type Underscore struct{}

// Clone implements Term.
func (Underscore) Clone() Term {
	return Underscore{}
}

// Pretty implements Term.
func (u Underscore) Pretty() string {
	return u.String()
}

func (u Underscore) String() string {
	return "_"
}

type String struct {
	Value string
}

// Clone implements Term.
func (t String) Clone() Term {
	return String{
		Value: t.Value,
	}
}

// Pretty implements AtomicTerm.
func (t String) Pretty() string {
	return t.String()
}

func (s String) String() string {
	return fmt.Sprintf("%q", s.Value) // quote the string
}

type Atom struct {
	Value string
}

// Clone implements Term.
func (t Atom) Clone() Term {
	return Atom{
		Value: t.Value,
	}
}

// Pretty implements AtomicTerm.
func (t Atom) Pretty() string {
	return t.String()
}

func (s Atom) String() string {
	return s.Value // do NOT quote the name of an Atom
}

type Integer struct {
	Value int
}

// Clone implements Term.
func (t Integer) Clone() Term {
	return Integer{
		Value: t.Value,
	}
}

// Pretty implements AtomicTerm.
func (t Integer) Pretty() string {
	return t.String()
}

func (i Integer) String() string {
	return fmt.Sprintf(INTFORMAT, i.Value)
}

type Float struct {
	Value float64
}

// Clone implements Term.
func (t Float) Clone() Term {
	return Float{
		Value: t.Value,
	}
}

// Pretty implements AtomicTerm.
func (t Float) Pretty() string {
	return t.String()
}

func (f Float) String() string {
	return fmt.Sprintf(FLOATFORMAT, f.Value)
}

// a compound term is a Term with children.
// A compound term withoutout children remains a compound term, different from an Atom.
type CompoundTerm struct {
	Functor  string
	Children []Term
}

func (t CompoundTerm) Clone() Term {
	var children = make([]Term, 0, len(t.Children))
	for _, child := range t.Children {
		children = append(children, child.Clone())
	}
	return CompoundTerm{
		Functor:  t.Functor,
		Children: children,
	}
}

// Pretty implements Term.
func (c CompoundTerm) Pretty() string {
	switch c.Functor {
	case "dot":
		return prettyList(c)
	case "query":
		var sb strings.Builder
		fmt.Fprintf(&sb, "?- ")
		for i, child := range c.Children {
			if child == nil {
				fmt.Println(START_RED, "WARNING : unexpected nil children", END_RED)
			} else {
				fmt.Fprintf(&sb, "%s", child.Pretty()) // caution ! Need to pretty inside the tree also !
			}
			if i < len(c.Children)-1 {
				fmt.Fprintf(&sb, ", ")
			}
		}
		fmt.Fprint(&sb, " ")
		return sb.String()
	case "rule":
		var sb strings.Builder
		if len(c.Children) == 0 {
			fmt.Println(START_RED, "WARNING : unexpected rule without head", END_RED)
		}
		for i, child := range c.Children {
			if child == nil {
				fmt.Println(START_RED, "WARNING : unexpected nil children", END_RED)
			} else {
				fmt.Fprintf(&sb, "%s", child.Pretty()) // caution ! Need to pretty inside the tree also !
			}
			if i == 0 {
				fmt.Fprintf(&sb, " :- ")
			} else {
				if i < len(c.Children)-1 {
					fmt.Fprintf(&sb, ", ")
				}
			}
		}
		fmt.Fprint(&sb, " ")
		return sb.String()

	default:
		var sb strings.Builder
		fmt.Fprintf(&sb, "%s(", c.Functor)
		for i, child := range c.Children {
			if child == nil {
				fmt.Println(START_RED, "WARNING : unexpected nil children", END_RED)
			} else {
				fmt.Fprintf(&sb, "%s", child.Pretty()) // caution ! Need to pretty inside the tree also !
			}
			if i < len(c.Children)-1 {
				fmt.Fprintf(&sb, ", ")
			}
		}
		fmt.Fprint(&sb, ")")
		return sb.String()
	}

}

func (c CompoundTerm) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s(", c.Functor)
	for i, child := range c.Children {
		fmt.Fprintf(&sb, "%s", child)
		if i < len(c.Children)-1 {
			fmt.Fprintf(&sb, ", ")
		}
	}
	fmt.Fprint(&sb, ")")
	return sb.String()
}

type Char struct {
	Char rune
}

// Clone implements Term.
func (t Char) Clone() Term {
	return Char{
		Char: t.Char,
	}
}

// Pretty implements Term.
func (c Char) Pretty() string {
	return c.String()
}

// Char implements Term.
func (c Char) String() string {
	return fmt.Sprintf("%q", c.Char)
}

// ---------------------------

// create a new list with the provided terms
func newList(terms ...Term) CompoundTerm {
	if len(terms) == 0 {
		return CompoundTerm{Functor: "dot"}
	}
	return CompoundTerm{
		Functor:  "dot",
		Children: []Term{terms[0], newList(terms[1:]...)},
	}
}

// pretty print a list.
// is is assumed to be a dot functor.
func prettyList(c CompoundTerm) string {

	if c.Functor != "dot" {
		panic("not a dot functor")
	}
	if len(c.Children) == 0 {
		return "[]"
	}
	if len(c.Children) == 1 {
		return "[" + c.Children[0].Pretty() + "|]"
	}

	// Assume it is a regular list, and try to build its string representation in sb.
	sb := new(strings.Builder)
	islist := true // if not a List, use this flag to mark it when breaking out of the loop.
	c2 := c.Children[1]

	fmt.Fprintf(sb, "[%s", c.Children[0].Pretty())

	for islist {

		dc2, ok := c2.(CompoundTerm)

		if !ok || dc2.Functor != "dot" { // not a list !
			islist = false
			break
		}

		if len(dc2.Children) == 0 { // end of list !
			fmt.Fprintf(sb, "]")
			return sb.String()
		}

		if len(dc2.Children) == 1 { // not a list !
			islist = false
			break
		}

		// still a list, continue
		fmt.Fprintf(sb, ", %s", dc2.Children[0].Pretty())
		c2 = dc2.Children[1]
	}

	if islist {
		panic("internal error - should never occur")
	}

	// ok, it's not a list, but its a dot with two children, so we can just print them.
	return "[" + c.Children[0].Pretty() + "|" + c.Children[1].Pretty() + "]"

}

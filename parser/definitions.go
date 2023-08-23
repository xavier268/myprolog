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
	// Strings are neither quoted nor escaped, they are stored without the start/end " or `
}

type Variable struct { // a named variable
	Name string
	Nsp  int // name space
}

var _ Term = &Variable{}

func (v *Variable) String() string {
	if v.Nsp > 0 {
		return fmt.Sprintf("%v$%d", v.Name, v.Nsp)
	}
	return v.Name
}

type Underscore struct{}

var _ Term = &Underscore{}

func (u *Underscore) String() string {
	return "_"
}

type AtomicTerm interface {
	Term
}

type String struct {
	Value string
}

var _ AtomicTerm = new(String)

func (s *String) String() string {
	return s.Value
}

type Atom struct {
	Value string
}

var _ AtomicTerm = new(Atom)

func (s *Atom) String() string {
	return fmt.Sprintf("%q", s.Value) // quote the string
}

type Integer struct {
	Value int
}

var _ AtomicTerm = new(Integer)

func (i *Integer) String() string {
	return fmt.Sprintf(INTFORMAT, i.Value)
}

type Float struct {
	Value float64
}

var _ AtomicTerm = new(Float)

func (f *Float) String() string {
	return fmt.Sprintf(FLOATFORMAT, f.Value)
}

// a compound term is a Term with children.
// A compound term withoutout children remains a compound term, different from an Atom.
type CompoundTerm struct {
	Functor  string
	Children []Term
}

var _ Term = &CompoundTerm{}

func (c *CompoundTerm) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s(", c.Functor)
	for i, child := range c.Children {
		fmt.Fprintf(&sb, "%s", child)
		if i == len(c.Children)-1 {
			fmt.Fprintf(&sb, ")")
		} else {
			fmt.Fprintf(&sb, ", ")
		}
	}
	return sb.String()
}

type Char struct {
	Char rune
}

// Char implements Term.
func (c *Char) String() string {
	return fmt.Sprintf("%q", c.Char)
}

var _ Term = new(Char)

// create a new list with the provided terms
func newList(terms ...Term) *CompoundTerm {
	if len(terms) == 0 {
		return &CompoundTerm{Functor: "dot"}
	}
	if len(terms) == 1 {
		return &CompoundTerm{Functor: "dot", Children: terms}
	}
	return &CompoundTerm{
		Functor:  "dot",
		Children: []Term{terms[0], newList(terms[1:]...)},
	}
}

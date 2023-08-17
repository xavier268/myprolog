package config

import (
	"fmt"
	"strings"
)

type Term interface { // Term is the most gerenic level of data
	String() string // String() is the string representation of the term
}

type Variable struct { // a named variable
	name string
	nsp  int // name space
}

var _ Term = &Variable{}

func (v *Variable) String() string {
	if v.nsp > 0 {
		return fmt.Sprintf("%v$%d", v.name, v.nsp)
	}
	return v.name
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
	return s.Value
}

type Integer struct {
	Value int
}

var _ AtomicTerm = new(Integer)

func (i *Integer) String() string {
	return fmt.Sprintf(INTFORMAT, i.Value)
}

/*
type Float struct {
	Value float64
}

var _ AtomicTerm = new(Float)

func (f *Float) String() string {
	return fmt.Sprintf(FLOATFORMAT, f.Value)
}
*/

type CompoundTerm struct { // a compound term is a Term with children
	Functor  *Atom
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

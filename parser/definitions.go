package parser

import (
	"fmt"
	"math"
	"strings"
)

var (
	FLOATFORMAT = "%.2f"
)

type Term interface { // Term is the most general form of a term
	String() string // String() is the string representation of the entire term
	// Strings are neither quoted nor escaped internally, they are stored without the start/end " or `
	Pretty() string        // Pretty() is the string representation of the term, pretty printing lists and rules and queries.
	Clone() Term           // Clone() returns a deep copy of the term
	CloneNsp(nsp int) Term // CloneNsp() returns a deep copy of the term with a new name space
}

var _ Term = Atom{}
var _ Term = Number{}
var _ Term = String{}
var _ Term = Variable{}
var _ Term = Underscore{}
var _ Term = CompoundTerm{}

var (
	// Min Number value
	MinNumber = Number{
		Num:        int(math.MinInt64),
		Den:        1,
		normalized: true,
	}
	// Max Number value
	MaxNumber = Number{
		Num:        int(math.MaxInt64),
		Den:        1,
		normalized: true,
	}
	ZeroNumber = Number{
		Num:        0,
		Den:        1,
		normalized: true,
	}
)

// Number are immutable
type Number struct { // numbers are representented as rational  Num/Den
	Num        int
	Den        int
	normalized bool
}

// Check if 2 numbers are equals.
func (n Number) Eq(t Term) bool {
	if nt, ok := t.(Number); ok {
		return n.Num*nt.Den == n.Den*nt.Num
	}
	return false
}

// Check if n is strictly less than r
func (n Number) Less(r Number) bool {
	panic("not implemented")
}

// Normalize the internal representation of a number.
func (n Number) Normalize() Number {
	if n.normalized {
		return n
	}
	if n.Num == 0 {
		return Number{
			Num:        0,
			Den:        1,
			normalized: true}
	}
	if n.Den == 0 {
		return Number{
			Num:        1,
			Den:        0,
			normalized: true,
		}
	}
	p := Gcd(n.Num, n.Den)
	if n.Den > 0 {
		return Number{
			Num:        n.Num / p,
			Den:        n.Den / p,
			normalized: true,
		}
	}
	if n.Den < 0 {
		return Number{
			Num:        -n.Num / p,
			Den:        -n.Den / p,
			normalized: true,
		}
	}
	panic("code should be unreacheable")
}

func (n Number) IsInteger() bool {
	return n.Normalize().Den == 1
}

func (n Number) IsZero() bool {
	return n.Num == 0 // note, 0/0 will normalized to 0/1, ie 0
}

func (n Number) Minus() Number {
	n = n.Normalize()
	if n.Den == 0 {
		return n
	}
	return Number{
		Num:        -n.Num,
		Den:        n.Den,
		normalized: true,
	}
}

// Clone implements Term.
func (n Number) Clone() Term {
	return n
}

// CloneNsp implements Term.
func (n Number) CloneNsp(nsp int) Term {
	return n
}

// Pretty implements Term. Tries to be clever ...
func (n Number) Pretty() string {
	i := n.Normalize()
	if i.Den == 1 {
		return fmt.Sprintf("%d", i.Num)
	}
	if i.Num == 0 { // ensure 0/0, default zero value,  is 0
		return "0"
	}
	if i.Den == 0 {
		return "NaN"
	}
	return fmt.Sprintf("%d/%d", i.Num, i.Den)
}

// String implements Term. Dump internal data.
func (n Number) String() string {
	return fmt.Sprintf("%d/%d", n.Num, n.Den)
}

// Variable are immutable
type Variable struct { // a named variable
	Name string
	Nsp  int // name space
}

// CloneNsp implements Term.
func (t Variable) CloneNsp(nsp int) Term {
	return Variable{
		Name: t.Name,
		Nsp:  nsp,
	}
}

// Clone implements Term.
func (t Variable) Clone() Term {
	return t
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

// CloneNsp implements Term.
func (Underscore) CloneNsp(nsp int) Term {
	return Underscore{}
}

// Clone implements Term.
func (Underscore) Clone() Term {
	return Underscore{}
}

// Pretty implements Term.
func (u Underscore) Pretty() string {
	return "_"
}

func (u Underscore) String() string {
	return "_"
}

// String is immutable
type String struct {
	Value string
}

// CloneNsp implements Term.
func (t String) CloneNsp(nsp int) Term {
	return String{
		Value: t.Value,
	}
}

// Clone implements Term.
func (t String) Clone() Term {
	return t
}

// Pretty implements AtomicTerm.
func (t String) Pretty() string {
	return t.String()
}

func (s String) String() string {
	return fmt.Sprintf("%q", s.Value) // quote the string
}

// Atom is immutable
type Atom struct {
	Value string
}

// CloneNsp implements Term.
func (t Atom) CloneNsp(nsp int) Term {
	return t
}

// Clone implements Term.
func (t Atom) Clone() Term {
	return t
}

// Pretty implements AtomicTerm.
func (s Atom) Pretty() string {
	return s.Value
}

func (s Atom) String() string {
	return s.Value // do NOT quote the name of an Atom
}

// a compound term is a Term with children.
// A compound term withoutout children remains a compound term, different from an Atom.
type CompoundTerm struct {
	Functor  string
	Children []Term
}

// CloneNsp implements Term.
func (t CompoundTerm) CloneNsp(nsp int) Term {
	var children = make([]Term, 0, len(t.Children))
	for _, child := range t.Children {
		children = append(children, child.CloneNsp(nsp))
	}
	return CompoundTerm{
		Functor:  t.Functor,
		Children: children,
	}
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

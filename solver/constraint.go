package solver

import (
	"fmt"

	"github.com/xavier268/myprolog/parser"
)

// a Constraint is immutable
type Constraint interface {
	// deep copy
	Clone() Constraint
	// Check will check validity of constraint, clean it or simplify it.
	// It will return the constraint itself,  possibly modified, or nil if constraint should be ignored.
	// An error means the constraint is impossible to satisfy (e.g. positive occur check, empty number interval, ...)
	Check() (Constraint, error)
	// simplify c, taking into account the calling constraint (unchanged).
	// If changed, replace c by all of cc (possibly empty, to just suppress c).
	Simplify(c Constraint) (cc []Constraint, changed bool, err error)
	String() string
}

var _ Constraint = VarIsCompoundTerm{}
var _ Constraint = VarIsString{}
var _ Constraint = VarIsNumber{}
var _ Constraint = VarIsVar{}
var _ Constraint = VarIsAtom{}

var ErrInvalidConstraintNaN = fmt.Errorf("invalid constraint (NaN)")
var ErrInvalidConstraintEmptyRange = fmt.Errorf("invalid constraint, specified range is empty")
var ErrInvalidConstraintSimplify = fmt.Errorf("incompatible constraints detectted when simplifying")
var ErrPositiveOccur = fmt.Errorf("positive occur check")

type VarIsNumber struct {
	V           Variable
	Min         Number // minimum acceptable Number
	Max         Number // maximum acceptable Number
	IntegerOnly bool   // accept only integers
}

// String implements Constraint.
// Constraints are assumed already checked and normalized.
func (c VarIsNumber) String() string {
	if c.Min.Eq(c.Max) {
		return c.V.Pretty() + " = " + c.Min.Pretty()
	}
	if c.IntegerOnly { // only integer
		if c.Min.Eq(parser.MinNumber) && c.Max.Eq(parser.MaxNumber) {
			return c.V.Pretty() + " is an integer"
		}
		if c.Min.Eq(parser.MinNumber) {
			return c.V.Pretty() + " is an integer and " + c.V.Pretty() + " <= " + c.Max.Pretty()
		}
		if c.Max.Eq(parser.MaxNumber) {
			return c.V.Pretty() + " is an integer and " + c.Min.Pretty() + " <= " + c.V.Pretty()
		}
		return c.V.Pretty() + " is an integer and " + c.Min.Pretty() + " <= " + c.V.Pretty() + " <= " + c.Max.Pretty()
	} else { //  not necessarily an integer ...
		if c.Min.Eq(parser.MinNumber) && c.Max.Eq(parser.MaxNumber) {
			return c.V.Pretty() // that should have been cleaned before !
		}
		if c.Min.Eq(parser.MinNumber) {
			return c.V.Pretty() + " <= " + c.Max.Pretty()
		}
		if c.Max.Eq(parser.MaxNumber) {
			return c.Min.Pretty() + " <= " + c.V.Pretty()
		}
		return c.Min.Pretty() + " <= " + c.V.Pretty() + " <= " + c.Max.Pretty()
	}
}

// Check implements Constraint.
func (v VarIsNumber) Check() (Constraint, error) {
	if v.V.Name == "" {
		return nil, nil // ignore silently
	}
	if v.Min.Den == 0 || v.Max.Den == 0 { // ensure numbers are not NaN
		return nil, ErrInvalidConstraintNaN
	}
	if v.Max.Less(v.Min) { // ensure range is not empty because of limits relative positions
		return nil, ErrInvalidConstraintEmptyRange
	}
	if v.IntegerOnly { // Special case for integers only

		v.Max = v.Max.Floor() // convert limits to nearest integer
		v.Min = v.Min.Ceil()  // convert limits to nearest integer

		if v.Max.Eq(v.Min) { // range contains a single,  integer value
			return v, nil
		}

		if v.Max.Less(v.Min) { // range is empty
			return nil, ErrInvalidConstraintEmptyRange
		}
	}
	return v, nil
}

// Clone implements Constraint.
func (c VarIsNumber) Clone() Constraint {
	return c
}

// Simplify implements Constraint.
func (VarIsNumber) Simplify(c Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

type VarIsAtom struct {
	V Variable
	A Atom
}

// String implements Constraint.
func (c VarIsAtom) String() string {
	return c.V.Pretty() + " =" + c.A.Pretty()
}

// Check implements Constraint.
func (c VarIsAtom) Check() (Constraint, error) {
	if c.V.Name == "" {
		return nil, nil // ignore silently
	}
	if c.A.Value == "" {
		panic("invalid atom constraint, atom value should not be the empty string")
	}
	return c, nil
}

// Simplify by replacing c2 with the set, possibly empty, of cc.
// c1 remains unchanged, and is never part of cc.
// Assume check was performed on both c1 and c2.
// If changed is false, ignore cc, keep c2 as is.
// If changed is true, remove c2 and replace it by all of cc (possibly empty, to just suppress c2).
func (c1 VarIsAtom) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {
	switch c2 := c2.(type) {
	case VarIsAtom:
		if c1.V.Eq(c2.V) { // same variable
			if c1.A.Value == c2.A.Value { // same atom
				return nil, true, nil // remove, duplicated.
			} else {
				return nil, false, ErrInvalidConstraintSimplify
			}
		} else { // different variables
			return nil, false, nil // no change, keep all
		}
	case VarIsString, VarIsNumber, VarIsVar:
		return nil, false, nil // no change, keep all
	case VarIsCompoundTerm:
		if !FindVar(c1.V, c2.T) {
			return nil, false, nil // no change, keep all
		}
		// here, we try to substitute c2/Atom in c1/Term ?
		t3, found := ReplaceVar(c1.V, c2.T, c1.A)
		if !found {
			return nil, false, nil // no change, keep all
		}
		c3 := VarIsCompoundTerm{
			V: c2.V,
			T: t3}

		return []Constraint{c3}, true, nil // no change, keep all

	default:
		panic("unreacheable code")
	}
}

// Clone implements Constraint.
func (c VarIsAtom) Clone() Constraint {
	return c
}

// Constraint for X = term
type VarIsCompoundTerm struct {
	V Variable
	T Term
}

// String implements Constraint.
func (VarIsCompoundTerm) String() string {
	panic("unimplemented")
}

// Check implements Constraint.
func (c VarIsCompoundTerm) Check() (Constraint, error) {
	if c.V.Name == "" {
		return nil, nil // ignore silently
	}
	if FindVar(c.V, c.T) {
		return nil, ErrPositiveOccur
	}
	return c, nil
}

// Simplify implements Constraint.
func (VarIsCompoundTerm) Simplify(c Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

// Clone implements Constraint.
func (c VarIsCompoundTerm) Clone() Constraint {
	return VarIsCompoundTerm{
		V: Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		T: c.T,
	}
}

type VarIsString struct {
	V Variable
	S string
}

// String implements Constraint.
func (c VarIsString) String() string {
	return fmt.Sprintf("%s = %q", c.V.Pretty(), c.S)
}

// Check implements Constraint.
func (c VarIsString) Check() (Constraint, error) {
	if c.V.Name == "" {
		return nil, nil // silently ignore
	}
	return c, nil
}

// Simplify implements Constraint.
func (VarIsString) Simplify(c Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

// Clone implements Constraint.
func (c VarIsString) Clone() Constraint {
	return c
}

type VarIsVar struct {
	V Variable
	W Variable
}

// String implements Constraint.
func (c VarIsVar) String() string {
	return c.V.Pretty() + " = " + c.W.Pretty()
}

// Check implements Constraint.
// There is a cannonical order of variables, with Y = Y means Y is latest (highest Nsp)
func (c VarIsVar) Check() (Constraint, error) {
	if c.V.Eq(c.W) { // Ignore X=X silently
		return nil, nil
	} else {
		// Put in canonical order, to facilitate substitution and dedup. Ensure in V = W,   V appeared later than W (nsp >)
		if c.V.Nsp < c.W.Nsp || (c.V.Nsp == c.W.Nsp && c.V.Name < c.W.Name) {
			return VarIsVar{c.W, c.V}, nil
		} else {
			return c, nil
		}
	}
}

// Simplify implements Constraint.
func (VarIsVar) Simplify(c Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

// Clone implements Constraint.
func (c VarIsVar) Clone() Constraint {
	return c
}

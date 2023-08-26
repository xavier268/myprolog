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
var _ Constraint = VarIsInteger{}

var ErrInvalidConstraintNaN = fmt.Errorf("invalid constraint (NaN)")
var ErrInvalidConstraintEmptyRange = fmt.Errorf("invalid constraint, specified range is empty")

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
	}
	// Now, not necessarily an integer ...
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
}

// Check implements Constraint.
func (v VarIsNumber) Check() (Constraint, error) {
	if v.Min.Den == 0 || v.Max.Den == 0 { // ensure numbers are not NaN
		return nil, ErrInvalidConstraintNaN
	}
	if v.Max.Less(v.Min) { // ensure range is not empty because of limits relative positions
		return nil, ErrInvalidConstraintEmptyRange
	}
	if v.IntegerOnly { // Special case for integers only

		if v.Max.Eq(v.Min) && !v.Max.IsInteger() { // range contains a single, non integer value
			return nil, ErrInvalidConstraintEmptyRange
		}

		if v.Max.IsInteger() || v.Min.IsInteger() { // range contains at least one integer value as a limit
			return v, nil
		}
		// TODO - avoid using float64 ?
		if int(float64(v.Max.Num)/float64(v.Max.Den))-int(float64(v.Min.Num)/float64(v.Min.Den)) <= 0 { // range contains only non integer value
			return nil, ErrInvalidConstraintEmptyRange
		}
	}
	return v, nil
}

// Clone implements Constraint.
func (VarIsNumber) Clone() Constraint {
	panic("unimplemented")
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
func (VarIsAtom) String() string {
	panic("unimplemented")
}

// Check implements Constraint.
func (VarIsAtom) Check() (Constraint, error) {
	panic("unimplemented")
}

// Simplify implements Constraint.
func (VarIsAtom) Simplify(c Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

// Clone implements Constraint.
func (c VarIsAtom) Clone() Constraint {
	return VarIsAtom{
		V: Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		A: c.A,
	}
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
func (VarIsCompoundTerm) Check() (Constraint, error) {
	panic("unimplemented")
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

// Constraint for var that should resolve to an Integer in the given range
type VarIsInteger struct {
	V   Variable
	Min int // minimum acceptable value, included.
	Max int // max acceptable value, included.
}

// String implements Constraint.
func (VarIsInteger) String() string {
	panic("unimplemented")
}

// Check implements Constraint.
func (VarIsInteger) Check() (Constraint, error) {
	panic("unimplemented")
}

// Simplify implements Constraint.
func (VarIsInteger) Simplify(c Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

func (c VarIsInteger) Clone() Constraint {
	return VarIsInteger{
		V: Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		Min: c.Min,
		Max: c.Max,
	}
}

type VarIsString struct {
	V Variable
	S string
}

// String implements Constraint.
func (VarIsString) String() string {
	panic("unimplemented")
}

// Check implements Constraint.
func (VarIsString) Check() (Constraint, error) {
	panic("unimplemented")
}

// Simplify implements Constraint.
func (VarIsString) Simplify(c Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

// Clone implements Constraint.
func (c VarIsString) Clone() Constraint {
	return VarIsString{
		V: Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		S: c.S,
	}
}

type VarIsVar struct {
	V Variable
	W Variable
}

// String implements Constraint.
func (VarIsVar) String() string {
	panic("unimplemented")
}

// Check implements Constraint.
func (VarIsVar) Check() (Constraint, error) {
	panic("unimplemented")
}

// Simplify implements Constraint.
func (VarIsVar) Simplify(c Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

// Clone implements Constraint.
func (c VarIsVar) Clone() Constraint {
	return VarIsVar{
		V: Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		W: Variable{
			Name: c.W.Name,
			Nsp:  c.W.Nsp,
		},
	}
}

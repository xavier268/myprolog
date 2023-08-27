package solver

import (
	"github.com/xavier268/myprolog/parser"
)

type VarIsNumber struct {
	V           Variable
	Min         Number // minimum acceptable Number
	Max         Number // maximum acceptable Number
	IntegerOnly bool   // accept only integers
}

var _ Constraint = VarIsNumber{}

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

// Clone implements Constraint.
func (c VarIsNumber) Clone() Constraint {
	return c
}

// Check will check validity of constraint, clean it or simplify it.
// It will return the constraint itself,  possibly modified, or nil if constraint should be ignored (
// CAUTION : return can be nil, despite a nil error !
// An error means the constraint is impossible to satisfy (e.g. positive occur check, empty number interval, ...)
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
			return v, nil // return the modified constraint, X = 2
		}

		if v.Max.Less(v.Min) { // range is empty
			return nil, ErrInvalidConstraintEmptyRange
		}
	}
	return v, nil // v was assumed to be checked, so no further check.
}

// Simplify implements Constraint.
func (c1 VarIsNumber) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {
	switch c2 := c2.(type) {
	case VarIsVar:
		if c2.V.Eq(c1.V) {
			c3 := VarIsNumber{
				V:           c2.W,
				Min:         c1.Min,
				Max:         c1.Max,
				IntegerOnly: c1.IntegerOnly,
			}
			// c2, c3 is valid, since input had been checked.
			return []Constraint{c2, c3}, true, nil
		}
		if c2.W.Eq(c1.V) {
			c3 := VarIsNumber{
				V:           c2.V,
				Min:         c1.Min,
				Max:         c1.Max,
				IntegerOnly: c1.IntegerOnly,
			}
			return []Constraint{c2, c3}, true, nil
		}
		return nil, false, nil // no change
	case VarIsNumber:
		if !c1.V.Eq(c2.V) {
			return nil, false, nil // no change
		}
		if c1.Min.Eq(c2.Min) && c1.Max.Eq(c2.Max) && c1.IntegerOnly == c2.IntegerOnly {
			return nil, true, nil // identical, remove
		}
		changed := false
		c3 := c2
		if c2.Min.Less(c1.Min) {
			c3.Min = c1.Min
			changed = true
		}
		if c2.Max.Greater(c1.Max) {
			c3.Max = c1.Max
			changed = true
		}
		if c1.IntegerOnly && !c2.IntegerOnly {
			c3.IntegerOnly = true
			changed = true
		}
		if !changed {
			return nil, false, nil // no change
		}
		c4, err := c3.Check() // required to respect the garantee that all output are checked, assumming inputs are.
		if err != nil {
			return nil, false, err
		}
		if c4 == nil {
			return nil, true, nil // remove
		}
		return []Constraint{c4}, true, nil
	case VarIsAtom:
		if c2.V.Eq(c1.V) {
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // no change
	case VarIsCompoundTerm:
		if c2.V.Eq(c1.V) {
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // no change
	case VarIsString:
		if c2.V.Eq(c1.V) {
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // no change
	default:
		panic("code not reachable")
	}
}

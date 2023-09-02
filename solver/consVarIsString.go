package solver

import "fmt"

type VarIsString struct {
	V Variable
	S String
}

var _ Constraint = VarIsString{}

func (c VarIsString) GetV() Variable {
	return c.V
}

// String implements Constraint.
func (c VarIsString) String() string {
	return fmt.Sprintf("%s = %q", c.V.Pretty(), c.S)
}

// Clone implements Constraint.
func (c VarIsString) Clone() Constraint {
	return c
}

// Check will check validity of constraint, clean it or simplify it.
// It will return the constraint itself,  possibly modified, or nil if constraint should be ignored (
// CAUTION : return can be nil, despite a nil error !
// An error means the constraint is impossible to satisfy (e.g. positive occur check, empty number interval, ...)
func (c VarIsString) Check() (Constraint, error) {
	if c.V.Name == "" {
		return nil, nil // silently ignore
	}
	return c, nil
}

// simplify c2, taking into account the calling constraint c1 (unchanged).
// if error, other results are not significant and should not be used.
// If changed, replace c by all of cc (possibly empty, to just suppress c).
// If changed is false, cc has no meaning.
// Input constraints MUST BE checked
// Output constraints are checkedutput.
func (c1 VarIsString) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {
	switch c2 := c2.(type) {
	case VarIsString:
		if c1.V.Eq(c2.V) { // same variable
			if c1.S == c2.S { // same string
				return nil, true, nil // remove, duplicated.
			} else {
				return nil, false, ErrInvalidConstraintSimplify
			}
		} else { // different variables
			return nil, false, nil // no change, keep all
		}

	case VarIsAtom:
		if c1.V.Eq(c2.V) { // same variable {
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // no change, keep all

	case VarIsCompoundTerm:
		if c1.V.Eq(c2.V) { // same variable {
			return nil, false, ErrInvalidConstraintSimplify
		}
		if !FindVar(c1.V, c2.T) {
			return nil, false, nil // no change, keep all
		}
		// here, we try to substitute c1.V by c1.S in c2.T
		t3, found := ReplaceVar(c1.V, c2.T, c1.S)
		if !found {
			return nil, false, nil // no change, keep all
		}
		c3 := VarIsCompoundTerm{
			V: c2.V,
			T: t3}
		c4, err := c3.Check() // verify, possible positive occur check ?
		if err == nil {       // c4 could also be nil
			if c4 == nil {
				return nil, true, nil // remove.
			}
			return []Constraint{c4}, true, nil
		} else {
			return nil, false, err
		}
	case VarIsVar:
		if c1.V.Eq(c2.V) { // same variable
			c3 := VarIsString{
				V: c2.W,
				S: c1.S}
			return []Constraint{c3}, true, nil // c1.V substituted by c1.A
		}
		if c1.V.Eq(c2.W) { // same variable
			c3 := VarIsString{
				V: c2.V,
				S: c1.S}
			return []Constraint{c3}, true, nil // c1.V substituted by c1.A
		}
		return nil, false, nil // no change, keep all
	case VarIsNum:
		if c1.V.Eq(c2.V) { // same variable
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // no change, keep all
	case VarLTNum:
		if c1.V.Eq(c2.V) { // same variable
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // no change, keep all
	case VarGTNum:
		if c1.V.Eq(c2.V) { // same variable
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // no change, keep all
	case VarGTENum:
		if c1.V.Eq(c2.V) { // same variable
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // no change, keep all
	case VarLTENum:
		if c1.V.Eq(c2.V) { // same variable
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // no change, keep all
	case VarINT:
		if c1.V.Eq(c2.V) { // same variable
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // no change, keep all

	default:
		_ = c2 // keep the compiler happy
		panic("internal error - unimplemented case")
	}

}

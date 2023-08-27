package solver

import "fmt"

// Constraint for X = term
type VarIsCompoundTerm struct {
	V Variable
	T Term
}

var _ Constraint = VarIsCompoundTerm{}

// String implements Constraint.
func (v VarIsCompoundTerm) String() string {
	return v.V.Pretty() + " = " + v.T.Pretty()
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

// Check will check validity of constraint, clean it or simplify it.
// It will return the constraint itself,  possibly modified, or nil if constraint should be ignored (
// CAUTION : return can be nil, despite a nil error !
// An error means the constraint is impossible to satisfy (e.g. positive occur check, empty number interval, ...)
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
	fmt.Println("VarIsCompoundTerm.Simplify error :", ErrNotImplemented)
	return nil, false, ErrNotImplemented
}

package solver

import "fmt"

// Constraint for X = term
type VarIsCompoundTerm struct {
	V Variable
	T Term
}

// String implements Constraint.
func (v VarIsCompoundTerm) String() string {
	return v.V.Pretty() + " = " + v.T.Pretty()
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
	fmt.Println("VarIsCompoundTerm.Simplify error :", ErrNotImplemented)
	return nil, false, ErrNotImplemented
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

package solver

import "fmt"

type VarIsString struct {
	V Variable
	S string
}

var _ Constraint = VarIsString{}

// String implements Constraint.
func (c VarIsString) String() string {
	return fmt.Sprintf("%s = %q", c.V.Pretty(), c.S)
}

// Clone implements Constraint.
func (c VarIsString) Clone() Constraint {
	return c
}

// Check implements Constraint.
func (c VarIsString) Check() (Constraint, error) {
	if c.V.Name == "" {
		return nil, nil // silently ignore
	}
	return c, nil
}

// Simplify implements Constraint.
// Returned constraints are checked, if input was.
func (VarIsString) Simplify(c Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

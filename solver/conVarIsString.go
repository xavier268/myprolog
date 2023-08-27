package solver

import "fmt"

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

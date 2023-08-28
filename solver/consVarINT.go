package solver

// Force V to be any integer number
// NaN is not considered a valid number.
type VarINT struct{ V Variable }

var _ Constraint = VarINT{}

// String implements Constraint.
func (c VarINT) String() string {
	return c.V.Pretty() + " is an integer"
}

// Check implements Constraint.
func (c VarINT) Check() (Constraint, error) {
	return c, nil
}

// Clone implements Constraint.
func (c VarINT) Clone() Constraint {
	return c
}

// Simplify implements Constraint.
func (c1 VarINT) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

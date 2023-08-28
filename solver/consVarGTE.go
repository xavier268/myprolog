package solver

// Ensure V >= Value
type VarGTE struct { // greater than or equal to
	V     Variable
	Value Number
}

var _ Constraint = VarGTE{}

// String implements Constraint.
func (c VarGTE) String() string {
	return c.V.Pretty() + " >= " + c.Value.Pretty()
}

// Check implements Constraint.
func (c VarGTE) Check() (Constraint, error) {
	c.Value = c.Value.Normalize()
	if c.Value.IsNaN() {
		return nil, ErrNaN
	}
	return c, nil
}

// Clone implements Constraint.
func (c VarGTE) Clone() Constraint {
	return c
}

// Simplify implements Constraint.
func (c1 VarGTE) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

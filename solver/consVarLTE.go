package solver

// Ensure V <= Value
type VarLTE struct { // less than or equal to
	V     Variable
	Value Number
}

var _ Constraint = VarLTE{}

// String implements Constraint.
func (c VarLTE) String() string {
	return c.V.Pretty() + " <= " + c.Value.Pretty()
}

// Check implements Constraint.
func (c VarLTE) Check() (Constraint, error) {
	c.Value = c.Value.Normalize()
	if c.Value.IsNaN() {
		return nil, ErrNaN
	}
	return c, nil
}

// Clone implements Constraint.
func (c VarLTE) Clone() Constraint {
	return c
}

// Simplify implements Constraint.
func (c1 VarLTE) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

package solver

// Ensure V < Value
type VarLT struct { // less than
	V     Variable
	Value Number
}

var _ Constraint = VarLT{}

// String implements Constraint.
func (c VarLT) String() string {
	return c.V.Pretty() + " < " + c.Value.Pretty()
}

// Check implements Constraint.
func (c VarLT) Check() (Constraint, error) {
	c.Value = c.Value.Normalize()
	if c.Value.IsNaN() {
		return nil, ErrNaN
	}
	return c, nil
}

// Clone implements Constraint.
func (c VarLT) Clone() Constraint {
	return c
}

// Simplify implements Constraint.
func (c1 VarLT) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

package solver

// Ensure V > Value
type VarGT struct { // greater than
	V     Variable
	Value Number
}

var _ Constraint = VarGT{}

// String implements Constraint.
func (c VarGT) String() string {
	return c.V.Pretty() + " > " + c.Value.Pretty()
}

// Check implements Constraint.
func (c VarGT) Check() (Constraint, error) {
	c.Value = c.Value.Normalize()
	if c.Value.IsNaN() {
		return nil, ErrNaN
	}
	return c, nil
}

// Clone implements Constraint.
func (c VarGT) Clone() Constraint {
	return c
}

// Simplify implements Constraint.
func (c1 VarGT) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {
	panic("unimplemented")
}

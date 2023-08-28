package solver

// Ensure V = Number Value
type VarEQ struct { // =
	V     Variable
	Value Number
}

var _ Constraint = VarEQ{}

// String implements Constraint.
func (c VarEQ) String() string {
	return c.V.Pretty() + " = " + c.Value.Pretty()
}

// Check implements Constraint.
func (c VarEQ) Check() (Constraint, error) {
	c.Value = c.Value.Normalize()
	if c.Value.IsNaN() {
		return nil, ErrNaN
	}
	return c, nil
}

// Clone implements Constraint.
func (c VarEQ) Clone() Constraint {
	return c
}

// Simplify implements Constraint.
func (c1 VarEQ) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {
	switch c2 := c2.(type) {
	case VarEQ:
		panic("unimplement")
	case VarLT:
		panic("unimplement")
	case VarGT:
		panic("unimplement")
	case VarGTE:
		panic("unimplement")
	case VarLTE:
		panic("unimplement")
	case VarINT:
		panic("unimplement")
	case VarIsVar:
		panic("unimplement")
	case VarIsAtom:
		panic("unimplement")
	case VarIsString:
		panic("unimplement")
	case VarIsCompoundTerm:
		panic("unimplement")
	default:
		_ = c2 // keep the compiler happy
		panic("internal error - unimplemented case")
	}
}

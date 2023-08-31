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
		if c1.V == c2.V {
			if c1.Value.Eq(c2.Value) {
				return nil, true, nil // ignore duplicate
			} else {
				return nil, false, ErrInvalidConstraintEmptyRange // contradiction
			}
		}
		return nil, false, nil // keep, no change
	case VarLT:
		if c1.V == c2.V {
			if c1.Value.Less(c2.Value) {
				return nil, true, nil // ignore duplicate
			} else {
				return nil, false, ErrInvalidConstraintEmptyRange // contradiction
			}
		}
		return nil, false, nil // keep, no change
	case VarGT:
		if c1.V == c2.V {
			if c1.Value.Greater(c2.Value) {
				return nil, true, nil // ignore duplicate
			} else {
				return nil, false, ErrInvalidConstraintEmptyRange // contradiction
			}
		}
		return nil, false, nil // keep, no change
	case VarGTE:
		if c1.V == c2.V {
			if c1.Value.Greater(c2.Value) || c1.Value.Eq(c2.Value) {
				return nil, true, nil // ignore duplicate
			} else {
				return nil, false, ErrInvalidConstraintEmptyRange // contradiction
			}
		}
		return nil, false, nil // keep, no change
	case VarLTE:
		if c1.V == c2.V {
			if c1.Value.Less(c2.Value) || c1.Value.Eq(c2.Value) {
				return nil, true, nil // ignore duplicate
			} else {
				return nil, false, ErrInvalidConstraintEmptyRange // contradiction
			}
		}
		return nil, false, nil // keep, no change
	case VarINT:
		if c1.V == c2.V {
			if c1.Value.IsInteger() {
				return nil, true, nil // ignore duplicate
			} else {
				return nil, false, ErrInvalidConstraintEmptyRange // contradiction
			}
		}
		return nil, false, nil // keep, no change
	case VarIsVar:
		if c1.V.Eq(c2.V) {
			c3 := VarEQ{
				V:     c2.W,
				Value: c1.Value}
			return []Constraint{c3}, true, nil
		}
		if c1.V.Eq(c2.W) {
			c3 := VarEQ{
				V:     c2.V,
				Value: c1.Value}
			return []Constraint{c3}, true, nil
		}
		return nil, false, nil // keep, no change
	case VarIsAtom:
		if c1.V == c2.V {
			return nil, false, ErrInvalidConstraintSimplify // contradiction
		}
		return nil, false, nil // keep, no change
	case VarIsString:
		if c1.V == c2.V {
			return nil, false, ErrInvalidConstraintSimplify // contradiction
		}
		return nil, false, nil // keep, no change
	case VarIsCompoundTerm:
		if c1.V == c2.V {
			return nil, false, ErrInvalidConstraintSimplify // contradiction
		}
		// lets see if we can substitute ?
		t3, found := ReplaceVar(c1.V, c2.T, c1.Value)
		if found {
			c3, err := VarIsCompoundTerm{
				V: c2.V,
				T: t3,
			}.Check() // verify
			if err != nil {
				return nil, false, err // report err
			}
			if c3 == nil {
				return nil, true, nil // remove
			} else {
				return []Constraint{c3}, true, nil // change
			}
		}

		return nil, false, nil // keep

	default:
		panic("internal error - unimplemented case")
	}
}

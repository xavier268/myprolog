package solver

// Ensure V >= Value
type VarGTENum struct { // greater than or equal to
	V     Variable
	Value Number
}

var _ Constraint = VarGTENum{}

func (c VarGTENum) GetV() Variable {
	return c.V
}

// String implements Constraint.
func (c VarGTENum) String() string {
	return c.Value.Pretty() + " <= " + c.V.Pretty()
}

// Check implements Constraint.
func (c VarGTENum) Check() (Constraint, error) {
	c.Value = c.Value.Normalize()
	if c.Value.IsNaN() {
		return nil, ErrNaN
	}
	return c, nil
}

// Clone implements Constraint.
func (c VarGTENum) Clone() Constraint {
	return c
}

// Simplify implements Constraint.
func (c1 VarGTENum) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {
	switch c2 := c2.(type) {
	case VarEQNum:
		if c1.V == c2.V {
			if c1.Value.Greater(c2.Value) {
				return nil, false, ErrInvalidConstraintEmptyRange // contradiction
			}
		}
		return nil, false, nil // keep, no change
	case VarLTNum:
		if c1.V == c2.V {
			if c1.Value.Greater(c2.Value) || c1.Value.Eq(c2.Value) {
				return nil, false, ErrInvalidConstraintEmptyRange // contradiction
			}
		}
		return nil, false, nil // keep, no change
	case VarGTNum:
		if c1.V == c2.V {
			if c1.Value.Greater(c2.Value) {
				return nil, true, nil // ignore, duplicate
			}
		}
		return nil, false, nil // keep, no change
	case VarGTENum:
		if c1.V == c2.V {
			if c1.Value.Greater(c2.Value) || c1.Value.Eq(c2.Value) {
				return nil, true, nil // ignore, duplicate
			} else {
				return nil, false, nil // keep, no change
			}
		}
		return nil, false, nil // keep, no change
	case VarLTENum:
		if c1.V == c2.V {
			if c1.Value.Greater(c2.Value) {
				return nil, false, ErrInvalidConstraintEmptyRange // contradiction	, nil // ignore duplicate
			}
			if c1.Value == c2.Value {
				c3 := VarEQNum(c1)
				return []Constraint{c3}, true, nil
			}
		}
		return nil, false, nil // keep, no change
	case VarINT:
		return nil, false, nil // keep, no change
	case VarIsVar:
		if c1.V == c2.V {
			c3 := VarGTENum{c2.W, c1.Value}
			return []Constraint{c2, c3}, true, nil
		}
		if c1.V == c2.W {
			c3 := VarGTENum{c2.V, c1.Value}
			return []Constraint{c2, c3}, true, nil
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
		return nil, false, nil // keep, no change
	default:
		panic("case unimplemented")
	}
}

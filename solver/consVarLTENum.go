package solver

// Ensure V <= Value
type VarLTENum struct { // less than or equal to
	V     Variable
	Value Number
}

var _ Constraint = VarLTENum{}

func (c VarLTENum) GetV() Variable {
	return c.V
}

// String implements Constraint.
func (c VarLTENum) String() string {
	return c.V.Pretty() + " <= " + c.Value.Pretty()
}

// Check implements Constraint.
func (c VarLTENum) Check() (Constraint, error) {
	c.Value = c.Value.Normalize()
	if c.Value.IsNaN() {
		return nil, ErrNaN
	}
	return c, nil
}

// Clone implements Constraint.
func (c VarLTENum) Clone() Constraint {
	return c
}

// Simplify implements Constraint.
func (c1 VarLTENum) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {
	switch c2 := c2.(type) {
	case VarIsNum:
		if c1.V == c2.V {
			if c1.Value.Less(c2.Value) {
				return nil, false, ErrInvalidConstraintEmptyRange // contradiction
			}
		}
		return nil, false, nil // keep, no change
	case VarLTNum:
		if c1.V == c2.V {
			if c1.Value.Less(c2.Value) {
				return nil, true, nil // remove
			}
		}
		return nil, false, nil // keep, no change
	case VarGTNum: // interval
		if c1.V == c2.V {
			if c1.Value.Floor().Less(c2.Value) || c1.Value == c2.Value {
				return nil, false, ErrInvalidConstraintEmptyRange // remove
			}
		}
		return nil, false, nil // keep, no change
	case VarGTENum: // interval
		if c1.V == c2.V {
			if c1.Value.Less(c2.Value) {
				return nil, false, ErrInvalidConstraintEmptyRange
			}
			if c1.Value == c2.Value {
				c3 := VarIsNum(c1)
				return []Constraint{c3}, true, nil
			}
		}
		return nil, false, nil // keep, no change
	case VarLTENum:
		if c1.V == c2.V {
			if c1.Value.Less(c2.Value) || c1.Value == c2.Value {
				return nil, true, nil // remove
			}
		}
		return nil, false, nil // keep, no change
	case VarINT:
		return nil, false, nil // keep, no change
	case VarIsVar:
		if c1.V == c2.V {
			c3 := VarLTENum{c2.W, c1.Value}
			return []Constraint{c3}, true, nil
		}
		if c1.V == c2.W {
			c3 := VarLTENum{c2.V, c1.Value}
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
		return nil, false, nil // keep, no change
	default:
		panic("case unimplemented")
	}
}

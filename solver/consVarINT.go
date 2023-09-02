package solver

// Force V to be any integer number
// NaN is not considered a valid number.
type VarINT struct{ V Variable }

var _ Constraint = VarINT{}

func (c VarINT) GetV() Variable {
	return c.V
}

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
	switch c2 := c2.(type) {
	case VarIsNum:
		if c1.V == c2.V {
			if c2.Value.IsInteger() {
				return nil, false, nil // keep, no change
			} else {
				return nil, false, ErrInvalidConstraintSimplify
			}
		}
		return nil, false, nil // keep, no change
	case VarLTNum:
		if c1.V == c2.V { // X is INT & X < a
			if c2.Value.IsInteger() { // a is int

				c3 := VarLTENum{V: c2.V, Value: c2.Value.Minus(OneNumber)} // c3: X <= a-1
				return []Constraint{c3}, true, nil
			} else { // a is NOT an int
				c3 := VarLTENum{V: c2.V, Value: c2.Value.Floor()} // X <= floor(a)
				return []Constraint{c3}, true, nil
			}
		}
		return nil, false, nil // keep, no change
	case VarGTNum:
		if c1.V == c2.V { // X is INT & X > a
			if c2.Value.IsInteger() { // a is int

				c3 := VarLTENum{V: c2.V, Value: c2.Value.Plus(OneNumber)} // c3: X >= a+1
				return []Constraint{c3}, true, nil
			} else { // a is NOT an int
				c3 := VarLTENum{V: c2.V, Value: c2.Value.Ceil()} // X >= ceiling(a)
				return []Constraint{c3}, true, nil
			}
		}
		return nil, false, nil // keep, no change
	case VarGTENum:
		if c1.V == c2.V { // X is INT & X > a
			if c2.Value.IsInteger() { // a is int
				return nil, false, nil // keep, no change
			} else { // a is NOT an int
				c3 := VarGTENum{V: c2.V, Value: c2.Value.Ceil()} // X >= ceiling(a)
				return []Constraint{c3}, true, nil
			}
		}
		return nil, false, nil // keep, no change
	case VarLTENum:
		if c1.V == c2.V { // X is INT & X <= a
			if c2.Value.IsInteger() { // a is int
				return nil, false, nil // keep, no change
			} else { // a is NOT an int
				c3 := VarLTENum{V: c2.V, Value: c2.Value.Floor()} // X <= floor(a)
				return []Constraint{c3}, true, nil
			}
		}
		return nil, false, nil // keep, no change
	case VarINT:
		return nil, true, nil // remove, duplicate
	case VarIsVar:
		if c1.V == c2.V {
			c3 := VarINT{c2.W}
			return []Constraint{c2, c3}, true, nil
		}
		if c1.V == c2.W {
			c3 := VarINT{c2.V}
			return []Constraint{c2, c3}, true, nil
		}
		return nil, false, nil // keep, no change
	case VarIsAtom:
		if c1.V == c2.V {
			return nil, false, ErrInvalidConstraintSimplify // conflict
		}
		return nil, false, nil // keep, no change
	case VarIsString:
		if c1.V == c2.V {
			return nil, false, ErrInvalidConstraintSimplify // conflict
		}
		return nil, false, nil // keep, no change
	case VarIsCompoundTerm:
		if c1.V == c2.V {
			return nil, false, ErrInvalidConstraintSimplify // conflict
		}
		return nil, false, nil // keep, no change
	default:
		panic("cases not implemented")
	}
}

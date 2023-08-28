package solver

// Constraint for X = term
type VarIsCompoundTerm struct {
	V Variable
	T Term
}

var _ Constraint = VarIsCompoundTerm{}

// String implements Constraint.
func (v VarIsCompoundTerm) String() string {
	if v.T == nil {
		return v.V.Pretty() + " = <nil>" // should never happen outside tests
	}
	return v.V.Pretty() + " = " + v.T.Pretty()
}

// Clone implements Constraint.
func (c VarIsCompoundTerm) Clone() Constraint {
	return VarIsCompoundTerm{
		V: c.V,
		T: c.T.Clone(),
	}
}

// Check will check validity of constraint, clean it or simplify it.
// It will return the constraint itself,  possibly modified, or nil if constraint should be ignored (
// CAUTION : return can be nil, despite a nil error !
// An error means the constraint is impossible to satisfy (e.g. positive occur check, empty number interval, ...)
func (c VarIsCompoundTerm) Check() (Constraint, error) {
	if c.V.Name == "" || c.T == nil {
		return nil, nil // ignore silently
	}
	if FindVar(c.V, c.T) {
		return nil, ErrPositiveOccur
	}
	return c, nil
}

// Simplify implements Constraint.
func (c1 VarIsCompoundTerm) Simplify(c Constraint) (cc []Constraint, changed bool, err error) {
	switch c2 := c.(type) {
	case VarIsCompoundTerm:
		panic("not implemented - need to unify)")
	case VarIsAtom:
		if c1.V == c2.V {
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // keep, no change
	case VarIsString:
		if c1.V == c2.V {
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // keep, no change

	case VarIsVar:
		if c1.V == c2.V {
			// substitute
			c3 := VarIsCompoundTerm{
				V: c2.W,
				T: c1.T,
			}
			c4, err := c3.Check()
			if err != nil {
				return nil, false, err
			}
			if c4 == nil {
				return nil, true, nil // remove
			}
			return []Constraint{c4}, true, nil // update
		}
		if c1.V == c2.W {
			// substitute
			c3 := VarIsCompoundTerm{
				V: c2.V,
				T: c1.T,
			}
			c4, err := c3.Check()
			if err != nil {
				return nil, false, err
			}
			if c4 == nil {
				return nil, true, nil // remove
			}
			return []Constraint{c4}, true, nil // update
		}
		return nil, false, nil // keep, no change
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

	default:
		_ = c2 // keep the compiler happy
		panic("internal error - unimplemented case")
	}

}

package solver

type VarIsVar struct {
	V Variable
	W Variable
}

var _ Constraint = VarIsVar{}

// String implements Constraint.
func (c VarIsVar) String() string {
	return c.V.Pretty() + " = " + c.W.Pretty()
}

// Clone implements Constraint.
func (c VarIsVar) Clone() Constraint {
	return c
}

// Check will check validity of constraint, clean it or simplify it.
// It will return the constraint itself,  possibly modified, or nil if constraint should be ignored (
// CAUTION : return can be nil, despite a nil error !
// An error means the constraint is impossible to satisfy (e.g. positive occur check, empty number interval, ...)
// There is a cannonical order of variables, with Y = Y means Y is latest (highest Nsp)
func (c VarIsVar) Check() (Constraint, error) {
	if c.V.Eq(c.W) { // Ignore X=X silently
		return nil, nil
	} else {
		// Put in canonical order, to facilitate substitution and dedup. Ensure in V = W,   V appeared later than W (nsp >)
		if c.V.Less(c.W) {
			return VarIsVar{c.W, c.V}, nil
		} else {
			return c, nil
		}
	}
}

// simplify c2, taking into account the calling constraint c1 (unchanged).
// if error, other results are not significant and should not be used.
// If changed, replace c by all of cc (possibly empty, to just suppress c).
// If changed is false, cc has no meaning.
// Input constraints MUST BE checked
// Output constraints are checked
func (c1 VarIsVar) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {

	switch c2 := c2.(type) {
	case VarIsAtom, VarIsString:
		// no action
		return nil, false, nil // no change
	case VarIsVar:
		if c1.V.Eq(c2.V) && c1.W.Eq(c2.W) {
			return nil, true, nil // duplicate
		}
		if c1.V.Eq(c2.W) && c1.W.Eq(c2.V) {
			return nil, true, nil // duplicate, should never happen if correctly ordered (checked) before.
		}
		if c1.V.Eq(c2.W) { // substitute !
			c3 := VarIsVar{
				c2.V,
				c1.W,
			}
			c4, err := c3.Check()
			if err != nil {
				return nil, false, err
			}
			if c4 != nil {
				return []Constraint{c4}, true, nil
			} else {
				return nil, true, nil // remove
			}
		}
		return nil, false, nil // no change
	case VarIsCompoundTerm:
		tt, found := ReplaceVar(c1.V, c2.T, c1.W)
		if !found {
			return nil, false, nil // no change
		}
		c3 := VarIsCompoundTerm{c2.V, tt}
		c4, err := c3.Check()
		if err != nil {
			return nil, false, err
		}
		return []Constraint{c4}, true, nil

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

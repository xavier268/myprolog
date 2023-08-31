package solver

type VarIsVar struct {
	V Variable
	W Variable
}

var _ Constraint = VarIsVar{}

func (c VarIsVar) GetV() Variable {
	return c.V
}

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
	case VarEQ:
		if c1.V == c2.V { // same value - can unify both contents !
			c3 := VarEQ{
				V:     c1.W,
				Value: c2.Value,
			}
			return []Constraint{c3}, true, nil
		}
		return nil, false, nil // no change
	case VarLT:
		return nil, false, nil // keep, no change - will be handled in the other direction
	case VarGT:
		return nil, false, nil // keep, no change - will be handled in the other direction
	case VarGTE:
		return nil, false, nil // keep, no change - will be handled in the other direction
	case VarLTE:
		return nil, false, nil // keep, no change - will be handled in the other direction
	case VarINT:
		return nil, false, nil // keep, no change - will be handled in the other direction
	case VarIsAtom:
		if c1.V == c2.V { // same value - can unify both contents !
			c3 := VarIsAtom{
				V: c1.W,
				A: c2.A,
			}
			return []Constraint{c3}, true, nil
		}
		return nil, false, nil // no change
	case VarIsString:
		if c1.V == c2.V { // same value - can unify both contents !
			c3 := VarIsString{
				V: c1.W,
				S: c2.S,
			}
			return []Constraint{c3}, true, nil
		}
		return nil, false, nil // no change
	case VarIsVar:
		if c1.V.Eq(c2.V) && c1.W.Eq(c2.W) {
			return nil, true, nil // duplicate
		}
		// no further attempts. Special situation like X=Y, X=a => Y=a will be better handled as X=a, X=Y => Y=a, removing X=Y !
		return nil, false, nil // no change
	case VarIsCompoundTerm: // X=Y, f(a,X,Y) => eliminate X , ie retun f(a,Y, Y)
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

	default:
		panic("internal error - unimplemented case")
	}
}

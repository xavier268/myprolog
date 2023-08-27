package solver

type VarIsAtom struct {
	V Variable
	A Atom
}

var _ Constraint = VarIsAtom{}

// String implements Constraint.
func (c VarIsAtom) String() string {
	return c.V.Pretty() + " = " + c.A.Pretty()
}

// Clone implements Constraint.
func (c VarIsAtom) Clone() Constraint {
	return c
}

// Check will check validity of constraint, clean it or simplify it.
// It will return the constraint itself,  possibly modified, or nil if constraint should be ignored (
// CAUTION : return can be nil, despite a nil error !
// An error means the constraint is impossible to satisfy (e.g. positive occur check, empty number interval, ...)
func (c VarIsAtom) Check() (Constraint, error) {
	if c.V.Name == "" {
		return nil, nil // ignore silently
	}
	if c.A.Value == "" {
		panic("invalid atom constraint, atom value should not be the empty string")
	}
	return c, nil
}

// simplify c2, taking into account the calling constraint c1 (unchanged).
// if error, other results are not significant and should not be used.
// If changed, replace c by all of cc (possibly empty, to just suppress c).
// If changed is false, cc has no meaning.
// Input constraints MUST BE checked
// Output constraints are checked
func (c1 VarIsAtom) Simplify(c2 Constraint) (cc []Constraint, changed bool, err error) {

	switch c2 := c2.(type) {
	case VarIsAtom:
		if c1.V.Eq(c2.V) { // same variable
			if c1.A.Value == c2.A.Value { // same atom
				return nil, true, nil // remove, duplicated.
			} else {
				return nil, false, ErrInvalidConstraintSimplify
			}
		} else { // different variables
			return nil, false, nil // no change, keep all
		}
	case VarIsString:
		if c1.V.Eq(c2.V) { // same variable {
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // no change, keep all
	case VarIsNumber:
		if c1.V.Eq(c2.V) { // same variable {
			return nil, false, ErrInvalidConstraintSimplify
		}
		return nil, false, nil // no change, keep all
	case VarIsCompoundTerm:
		if c1.V.Eq(c2.V) { // same variable {
			return nil, false, ErrInvalidConstraintSimplify
		}
		if !FindVar(c1.V, c2.T) {
			return nil, false, nil // no change, keep all
		}
		// here, we try to substitute c1.V by c1.A in c2.T
		t3, found := ReplaceVar(c1.V, c2.T, c1.A)
		if !found {
			return nil, false, nil // no change, keep all
		}
		c3 := VarIsCompoundTerm{
			V: c2.V,
			T: t3}
		c4, err := c3.Check() // verify, possible positive occur check ?
		if err == nil {       // c4 could also be nil
			if c4 == nil {
				return nil, true, nil // remove
			}
			return []Constraint{c4}, true, nil
		} else {
			return nil, false, err
		}

	case VarIsVar:
		if c1.V.Eq(c2.V) { // same variable
			c3 := VarIsAtom{
				V: c2.W,
				A: c1.A}
			return []Constraint{c3}, true, nil // c1.V substituted by c1.A
		}
		if c1.V.Eq(c2.W) { // same variable
			c3 := VarIsAtom{
				V: c2.V,
				A: c1.A}
			return []Constraint{c3}, true, nil // c1.V substituted by c1.A
		}
		return nil, false, nil // no change, keep all
	default:
		panic("case not implemented")
	}
}

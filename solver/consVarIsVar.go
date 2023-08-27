package solver

import "fmt"

type VarIsVar struct {
	V Variable
	W Variable
}

// String implements Constraint.
func (c VarIsVar) String() string {
	return c.V.Pretty() + " = " + c.W.Pretty()
}

// Check implements Constraint.
// There is a cannonical order of variables, with Y = Y means Y is latest (highest Nsp)
func (c VarIsVar) Check() (Constraint, error) {
	if c.V.Eq(c.W) { // Ignore X=X silently
		return nil, nil
	} else {
		// Put in canonical order, to facilitate substitution and dedup. Ensure in V = W,   V appeared later than W (nsp >)
		if c.V.Nsp < c.W.Nsp || (c.V.Nsp == c.W.Nsp && c.V.Name < c.W.Name) {
			return VarIsVar{c.W, c.V}, nil
		} else {
			return c, nil
		}
	}
}

// Simplify implements Constraint.
func (VarIsVar) Simplify(c Constraint) (cc []Constraint, changed bool, err error) {
	fmt.Println("VarIsVar.Simplify error :", ErrNotImplemented)
	return nil, false, ErrNotImplemented
}

// Clone implements Constraint.
func (c VarIsVar) Clone() Constraint {
	return c
}

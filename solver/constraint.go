package solver

//a Constraint is immutable
type Constraint interface {
	// deep copy
	Clone() Constraint
	// Check will check validity of constraint, clean it or simplify it.
	// It will return the constraint itself,  nil if constraint should not be added.
	// An error means the constraint is impossible to satisfy (e.g. positive occur check, empty number interval, ...)
	Check() (Constraint, error)
	// simplify c, taking into account the calling constraint (unchanged).
	// If a nil constraint with the changed flag set is returned, it means c can be safely removed.
	Simplify(c Constraint) (cc Constraint, changed bool, err error)
}

var _ Constraint = VarIsCompoundTerm{}
var _ Constraint = VarIsString{}
var _ Constraint = VarIsChar{}
var _ Constraint = VarIsInteger{}
var _ Constraint = VarIsFloat{}
var _ Constraint = VarIsVar{}
var _ Constraint = VarIsAtom{}

type VarIsAtom struct {
	V Variable
	A Atom
}

// Check implements Constraint.
func (VarIsAtom) Check() (Constraint, error) {
	panic("unimplemented")
}

// Simplify implements Constraint.
func (VarIsAtom) Simplify(c Constraint) (cc Constraint, changed bool, err error) {
	panic("unimplemented")
}

// Clone implements Constraint.
func (c VarIsAtom) Clone() Constraint {
	return VarIsAtom{
		V: Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		A: c.A,
	}
}

// Constraint for X = term
type VarIsCompoundTerm struct {
	V Variable
	T Term
}

// Check implements Constraint.
func (VarIsCompoundTerm) Check() (Constraint, error) {
	panic("unimplemented")
}

// Simplify implements Constraint.
func (VarIsCompoundTerm) Simplify(c Constraint) (cc Constraint, changed bool, err error) {
	panic("unimplemented")
}

// Clone implements Constraint.
func (c VarIsCompoundTerm) Clone() Constraint {
	return VarIsCompoundTerm{
		V: Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		T: c.T,
	}
}

// Constraint for var that should resolve to an Integer in the given range
type VarIsInteger struct {
	V   Variable
	Min int // minimum acceptable value, included.
	Max int // max acceptable value, included.
}

// Check implements Constraint.
func (VarIsInteger) Check() (Constraint, error) {
	panic("unimplemented")
}

// Simplify implements Constraint.
func (VarIsInteger) Simplify(c Constraint) (cc Constraint, changed bool, err error) {
	panic("unimplemented")
}

func (c VarIsInteger) Clone() Constraint {
	return VarIsInteger{
		V: Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		Min: c.Min,
		Max: c.Max,
	}
}

type VarIsFloat struct {
	V   Variable
	Min float64
	Max float64
}

// Check implements Constraint.
func (VarIsFloat) Check() (Constraint, error) {
	panic("unimplemented")
}

// Simplify implements Constraint.
func (VarIsFloat) Simplify(c Constraint) (cc Constraint, changed bool, err error) {
	panic("unimplemented")
}

// Clone implements Constraint.
func (c VarIsFloat) Clone() Constraint {
	return VarIsFloat{
		V: Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		Min: c.Min,
		Max: c.Max,
	}
}

type VarIsString struct {
	V Variable
	S string
}

// Check implements Constraint.
func (VarIsString) Check() (Constraint, error) {
	panic("unimplemented")
}

// Simplify implements Constraint.
func (VarIsString) Simplify(c Constraint) (cc Constraint, changed bool, err error) {
	panic("unimplemented")
}

// Clone implements Constraint.
func (c VarIsString) Clone() Constraint {
	return VarIsString{
		V: Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		S: c.S,
	}
}

type VarIsChar struct {
	V Variable
	C rune
}

// Check implements Constraint.
func (VarIsChar) Check() (Constraint, error) {
	panic("unimplemented")
}

// Simplify implements Constraint.
func (VarIsChar) Simplify(c Constraint) (cc Constraint, changed bool, err error) {
	panic("unimplemented")
}

// Clone implements Constraint.
func (c VarIsChar) Clone() Constraint {
	return VarIsChar{
		V: Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		C: c.C,
	}
}

type VarIsVar struct {
	V Variable
	W Variable
}

// Check implements Constraint.
func (VarIsVar) Check() (Constraint, error) {
	panic("unimplemented")
}

// Simplify implements Constraint.
func (VarIsVar) Simplify(c Constraint) (cc Constraint, changed bool, err error) {
	panic("unimplemented")
}

// Clone implements Constraint.
func (c VarIsVar) Clone() Constraint {
	return VarIsVar{
		V: Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		W: Variable{
			Name: c.W.Name,
			Nsp:  c.W.Nsp,
		},
	}
}

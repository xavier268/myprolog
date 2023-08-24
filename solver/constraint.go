package solver

//a Constraint is immutable
type Constraint interface {
	Clone() Constraint
}

// Constraint for X = term
type VarEqCons struct {
	V *Variable
	T Term
}

var _ Constraint = VarEqCons{}

// Clone implements Constraint.
func (c VarEqCons) Clone() Constraint {
	return VarEqCons{
		V: &Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		T: c.T,
	}
}

// Constraint for var that should resolve to an Integer in the given range
type VarInteger struct {
	V   *Variable
	Min int // minimum acceptable value, included.
	Max int // max acceptable value, included.
}

var _ Constraint = VarInteger{}

func (c VarInteger) Clone() Constraint {
	return VarInteger{
		V: &Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		Min: c.Min,
		Max: c.Max,
	}
}

// Attempt to simplify constratint list.
// Return error if an incompatibility was detected.
func SimplifyConstraints(constraints []Constraint) ([]Constraint, error) {
	panic("unimplemented")
}

package solver

import "github.com/xavier268/myprolog/solver"

//a Constraint is immutable
type Constraint interface {
	Clone() Constraint
}

var _ Constraint = VarIsCompoundTerm{}
var _ Constraint = VarIsString{}
var _ Constraint = VarIsChar{}
var _ Constraint = VarIsInteger{}
var _ Constraint = VarIsFloat{}
var _ Constraint = VarIsVar{}
var _ Constraint = VarIsAtom{}

type VarIsAtom struct {
	V *Variable
	A Atom
}

// Clone implements Constraint.
func (c VarIsAtom) Clone() solver.Constraint {
	return VarIsAtom{
		V: &Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		A: c.A,
	}
}

// Constraint for X = term
type VarIsCompoundTerm struct {
	V *Variable
	T Term
}

// Clone implements Constraint.
func (c VarIsCompoundTerm) Clone() Constraint {
	return VarIsCompoundTerm{
		V: &Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		T: c.T,
	}
}

// Constraint for var that should resolve to an Integer in the given range
type VarIsInteger struct {
	V   *Variable
	Min int // minimum acceptable value, included.
	Max int // max acceptable value, included.
}

func (c VarIsInteger) Clone() Constraint {
	return VarIsInteger{
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

type VarIsFloat struct {
	V   *Variable
	Min float64
	Max float64
}

// Clone implements Constraint.
func (c VarIsFloat) Clone() solver.Constraint {
	return VarIsFloat{
		V: &Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		Min: c.Min,
		Max: c.Max,
	}
}

type VarIsString struct {
	V *Variable
	S string
}

// Clone implements Constraint.
func (c VarIsString) Clone() solver.Constraint {
	return VarIsString{
		V: &Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		S: c.S,
	}
}

type VarIsChar struct {
	V *Variable
	C rune
}

// Clone implements Constraint.
func (c VarIsChar) Clone() solver.Constraint {
	return VarIsChar{
		V: &Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		C: c.C,
	}
}

type VarIsVar struct {
	V *Variable
	W *Variable
}

// Clone implements Constraint.
func (c VarIsVar) Clone() solver.Constraint {
	return VarIsVar{
		V: &Variable{
			Name: c.V.Name,
			Nsp:  c.V.Nsp,
		},
		W: &Variable{
			Name: c.W.Name,
			Nsp:  c.W.Nsp,
		},
	}
}

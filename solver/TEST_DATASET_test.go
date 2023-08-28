package solver

// data set and predefined variables for testing.

var (
	X  = Variable{Name: "X", Nsp: 0}
	X1 = Variable{Name: "X", Nsp: 1}

	Y  = Variable{Name: "Y", Nsp: 0}
	Y2 = Variable{Name: "Y", Nsp: 2}

	Z  = Variable{Name: "Z", Nsp: 0}
	Z3 = Variable{Name: "Z", Nsp: 3}
)

var TEST_VAR_IS_ATOM = []Constraint{
	VarIsAtom{
		V: X,
		A: Atom{Value: "toto"},
	},
	VarIsAtom{
		V: Y2,
		A: Atom{Value: "tata"},
	},
}

var TEST_VAR_IS_VAR = []Constraint{
	VarIsVar{X, X},
	VarIsVar{Y, X},
	VarIsVar{X, Y2},
	VarIsVar{Y, Y2},
	VarIsVar{Z3, X},
	VarIsVar{Z3, Y2},
}

var TEST_EMPTY_CONSTRAINTS = []Constraint{
	nil,
	VarEQ{},
	VarGT{},
	VarGTE{},
	VarLT{},
	VarLTE{},
	VarINT{},
	VarIsAtom{},
	VarIsString{},
	VarIsCompoundTerm{},
	VarIsCompoundTerm{V: X},
	VarIsCompoundTerm{V: X, T: nil},
}

func cleanAllConstraints(cc1 []Constraint) (cc2 []Constraint) {
	for _, c1 := range cc1 {
		c2, err := c1.Check()
		if err != nil || c2 == nil {
			continue
		}
		cc2 = append(cc2, c2)
	}
	return cc2
}

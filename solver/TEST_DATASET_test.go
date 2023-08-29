package solver

import (
	"fmt"
	"strings"
	"testing"
)

// data set and predefined variables for testing.

var (
	X  = Variable{Name: "X", Nsp: 0}
	X1 = Variable{Name: "X", Nsp: 1}

	Y  = Variable{Name: "Y", Nsp: 0}
	Y2 = Variable{Name: "Y", Nsp: 2}

	Z  = Variable{Name: "Z", Nsp: 0}
	Z3 = Variable{Name: "Z", Nsp: 3}
)

var TEST_NUMBERS = []Number{
	ZeroNumber,

	OneNumber,
	OneNumber.ChSign(),

	{Num: 1, Den: 2, Normalized: true},
	{Num: -1, Den: 2, Normalized: true},

	{Num: 3, Den: 2, Normalized: true},
	{Num: -3, Den: 2, Normalized: true},

	{Num: 3, Den: 1, Normalized: true},
	{Num: -3, Den: 1, Normalized: true},
}

var TEST_VAR_IS_ATOM = []Constraint{
	VarIsAtom{
		V: X,
		A: Atom{Value: "toto"},
	},
	VarIsAtom{
		V: Y2,
		A: Atom{Value: "tata"},
	},
	VarIsAtom{
		V: Z3,
		A: Atom{Value: "titi"},
	},
}

var TEST_VAR_IS_STRING = []Constraint{
	VarIsString{},
	VarIsString{X, String{}},
	VarIsString{Y, String{Value: "a string"}},
	VarIsString{Z3, String{Value: "another string"}},
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
	VarIsCompoundTerm{V: X, T: X},
	VarIsCompoundTerm{V: X, T: CompoundTerm{Functor: "foo", Children: []Term{}}},
}

func cleanAllConstraints(cc1 []Constraint) (cc2 []Constraint) {
	for _, c1 := range cc1 {
		if c1 == nil {
			continue
		}
		c2, err := c1.Check()
		if err != nil || c2 == nil {
			continue
		}
		cc2 = append(cc2, c2)
	}
	return cc2
}

// generate all possible number constraints using provided numbers
func generateConstraintsFromNumbers(nn []Number) (cc []Constraint) {

	cc = make([]Constraint, 0, 5*len(nn)+6)

	cc = append(cc, VarIsVar{Y, X})
	cc = append(cc, VarINT{Y})
	cc = append(cc, VarIsVar{Z, Y})

	cc = append(cc, VarIsAtom{Y, Atom{Value: "titi"}})
	cc = append(cc, VarIsString{Y, String{Value: "toto"}})
	cc = append(cc, VarIsCompoundTerm{Y, CompoundTerm{Functor: "tata", Children: []Term{Y, X}}})

	for _, n := range nn {

		cc = append(cc, VarEQ{X, n})
		cc = append(cc, VarLT{X, n})
		cc = append(cc, VarLTE{X, n})
		cc = append(cc, VarGT{X, n})
		cc = append(cc, VarGTE{X, n})

	}

	return cc

}

// Simplify each constraint against the others, two by two.
// no cleaning is done on constraints, except skipping nils.
func runConstraintSimplify2x2Test(t *testing.T, tCons []Constraint) *strings.Builder {

	sb := new(strings.Builder)

	for i, c1 := range tCons {

		fmt.Fprintln(sb)
		if c1 == nil {
			continue
		}
		fmt.Fprintf(sb, "\n\n%3d\tusing : %v\n", i, c1)
		for j, c2 := range tCons {
			if c2 == nil {
				continue
			}
			fmt.Fprintf(sb, "\n%3d-%3d\tto simplify :\t%v    --> ", i, j, c2)
			cc, changed, err := c1.Simplify(c2)
			if err != nil {
				fmt.Fprintf(sb, "\n%3d-%3d\t\t\t\t\tERROR   : %v", i, j, err)
			} else {
				if changed {
					if len(cc) == 0 {
						fmt.Fprintf(sb, "\n%3d-%3d\t\t\t\t\t\t\tREMOVE", i, j)
					} else {
						fmt.Fprintf(sb, "\n%3d-%3d\t\t\t\t\t\t\tREPLACE BY", i, j)
						for _, c := range cc {
							fmt.Fprintf(sb, "\n%3d-%3d\t\t\t\t\t\t\t\t\t%v", i, j, c)
						}
					}
				} else {
					fmt.Fprintf(sb, "\n%3d-%3d\t\t\t\t\t\t\tNO CHANGE, KEEP AS IS", i, j)
				}
			}
		}
	}
	return sb
}

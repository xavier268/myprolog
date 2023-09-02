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
	X2 = Variable{Name: "X", Nsp: 2}

	Y  = Variable{Name: "Y", Nsp: 0}
	Y1 = Variable{Name: "Y", Nsp: 1}
	Y2 = Variable{Name: "Y", Nsp: 2}

	Z  = Variable{Name: "Z", Nsp: 0}
	Z2 = Variable{Name: "Z", Nsp: 2}
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
		V: Y,
		A: Atom{Value: "tutu"},
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
	VarIsNum{},
	VarGTNum{},
	VarGTENum{},
	VarLTNum{},
	VarLTENum{},
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

		cc = append(cc, VarIsNum{X, n})
		cc = append(cc, VarLTNum{X, n})
		cc = append(cc, VarLTENum{X, n})
		cc = append(cc, VarGTNum{X, n})
		cc = append(cc, VarGTENum{X, n})

	}

	return cc

}

// Simplify each constraint against the others, two by two.
// no cleaning is done on constraints, except skipping nils.
func runConstraintSimplify2x2Test(t *testing.T, cc []Constraint) *strings.Builder {

	do := func(sb *strings.Builder, i int, j int, cc []Constraint) {
		c1 := cc[i]
		c2 := cc[j]
		if c1 == nil || c2 == nil {
			return
		}

		c3, ch, err := c1.Simplify(c2)
		fmt.Fprintf(sb, "\nUSE        :\t%s\nTO SIMPLIFY:\t%s\n\t -- >\t", c1.String(), c2.String())
		if err != nil {
			fmt.Fprintf(sb, "ERROR : %v", err)
		} else {
			if !ch {
				fmt.Fprintf(sb, "NO CHANGE")
			} else {
				if len(c3) == 0 {
					fmt.Fprintf(sb, "REMOVE")
				} else {
					fmt.Fprintf(sb, "REPLACE WITH:\t  ")
					for k, r := range c3 {
						fmt.Fprintf(sb, "%s", r.String())
						if k != len(c3)-1 {
							fmt.Fprintf(sb, " and ")
						}
					}
				}
			}
		}

	}

	sb := new(strings.Builder)

	for i := 0; i < len(cc); i++ {
		for j := i; j < len(cc); j++ {
			fmt.Fprintln(sb)
			do(sb, i, j, cc)
			fmt.Fprintln(sb)
			do(sb, j, i, cc)

		}
	}
	return sb
}

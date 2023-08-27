package solver

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xavier268/myprolog/mytest"
)

var X, Y, Z = Variable{Name: "X", Nsp: 0}, Variable{Name: "Y", Nsp: 0}, Variable{Name: "Z", Nsp: 0}

var tConsSimplify1 = []Constraint{
	VarIsAtom{
		V: X,
		A: Atom{Value: "toto"},
	},
	VarIsAtom{
		V: Y,
		A: Atom{Value: "tata"},
	},
	VarIsVar{
		V: Y,
		W: Z,
	},
	VarIsVar{
		V: X,
		W: Z,
	},
}

var tConsSimplify2 = []Constraint{

	VarIsNumber{
		V:           X,
		Min:         Number{Num: 3, Den: 2, Normalized: false},
		Max:         Number{Num: 7, Den: 2, Normalized: false},
		IntegerOnly: false,
	},
	VarIsNumber{
		V:           Z,
		Min:         Number{Num: -3, Den: 2, Normalized: false},
		Max:         Number{Num: 1, Den: 2, Normalized: false},
		IntegerOnly: true,
	},
}

func TestSimplify(t *testing.T) {
	var sb *strings.Builder

	sb = run(t, tConsSimplify1)
	mytest.Verify(t, sb.String(), "simplify_test.1.want")

	sb = run(t, tConsSimplify2)
	mytest.Verify(t, sb.String(), "simplify_test.2.want")

	sb = run(t, append(tConsSimplify1, tConsSimplify2...))
	mytest.Verify(t, sb.String(), "simplify_test.1.2.want")

}

func run(t *testing.T, tCons []Constraint) *strings.Builder {
	sb := new(strings.Builder)
	for i, c1 := range tCons {
		fmt.Fprintln(sb)
		fmt.Fprintf(sb, "\n%d\tusing : %v", i, c1)
		for j, c2 := range tCons {
			fmt.Fprintf(sb, "\n%d-%d\t\t to simplify : %v    --> ", i, j, c2)
			cc, changed, err := c1.Simplify(c2)
			if err != nil {
				fmt.Fprintf(sb, "\n%d-%d\t\t\tERROR   : %v", i, j, err)
			} else {
				if changed {
					if len(cc) == 0 {
						fmt.Fprintf(sb, "\n%d-%d\t\t\tREMOVE", i, j)
					} else {
						fmt.Fprintf(sb, "\n%d-%d\t\t\tREPLACE BY", i, j)
						for _, c := range cc {
							fmt.Fprintf(sb, "\n%d-%d\t\t\t\t%v", i, j, c)
						}
					}
				} else {
					fmt.Fprintf(sb, "\n%d-%d\t\t\tNO CHANGE, KEEP AS IS", i, j)
				}
			}
		}
	}
	return sb
}

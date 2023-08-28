package solver

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xavier268/myprolog/mytest"
)

func TestSimplify(t *testing.T) {
	var sb *strings.Builder

	sb = run(t, TEST_VAR_IS_ATOM)
	mytest.Verify(t, sb.String(), "constraint_simplify_test.atoms")

	sb = run(t, TEST_VAR_IS_STRING)
	mytest.Verify(t, sb.String(), "constraint_simplify_test.strings")

	sb = run(t, append(TEST_VAR_IS_STRING, TEST_VAR_IS_ATOM...))
	mytest.Verify(t, sb.String(), "constraint_simplify_test.strings.atoms")

	sb = run(t, TEST_VAR_IS_VAR)
	mytest.Verify(t, sb.String(), "constraint_simplify_test.vars")

	sb = run(t, append(TEST_VAR_IS_ATOM, TEST_VAR_IS_VAR...))
	mytest.Verify(t, sb.String(), "constraint_simplify_test.atoms.vars")

	sb = run(t, append(TEST_VAR_IS_STRING, TEST_VAR_IS_VAR...))
	mytest.Verify(t, sb.String(), "constraint_simplify_test.strings.vars")

}

func run(t *testing.T, tCons []Constraint) *strings.Builder {

	sb := new(strings.Builder)

	// only test on checked, non nil, constraints
	tCons = cleanAllConstraints(tCons)

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

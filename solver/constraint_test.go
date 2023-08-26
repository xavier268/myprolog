package solver

import (
	"fmt"
	"strings"
	"testing"
)

// Data set of constraints to test
var cData = []Constraint{

	VarIsNumber{
		V:           Variable{Name: "X", Nsp: 0},
		Min:         Number{Num: 5, Den: 2, Normalized: false},
		Max:         Number{Num: 10, Den: 2, Normalized: false},
		IntegerOnly: false,
	},
}

func TestConstraintsCheck(t *testing.T) {
	var err error
	sb := new(strings.Builder)

	fmt.Fprintf(sb, "========= %s =========\n", t.Name())

	for i, c := range cData {
		fmt.Fprintf(sb, "\n\n%d\t(original\t%v\n", i, c.String())
		c, err = c.Check()
		if err != nil {
			fmt.Fprintf(sb, "%d\t%v\n", i, err)
		}
		fmt.Fprintf(sb, "%d\t(checked)\t%v\n", i, c)
	}

	verifyTest(t, sb.String(), "constraint_check_test.want")
}

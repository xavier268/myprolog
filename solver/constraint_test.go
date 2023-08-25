package solver

import (
	"fmt"
	"strings"
	"testing"
)

func TestConstraints(t *testing.T) {

	sb := new(strings.Builder)
	cdata := []Constraint{} // TODO: Add Constraint Data Here

	fmt.Fprintf(sb, "========= %s =========\n", t.Name())

	for i, c := range cdata {
		fmt.Fprintf(sb, "\n\n%d\t%v\n", i, c)
		// operations ...
	}

	verifyTest(t, sb.String(), "constraint_test.want")
}

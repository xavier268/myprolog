package solver

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xavier268/myprolog/mytest"
)

func TestConstraintsCheck(t *testing.T) {
	var err error
	sb := new(strings.Builder)

	fmt.Fprintf(sb, "\n========= Single constraint test =========\n")

	cData := append([]Constraint{}, TEST_EMPTY_CONSTRAINTS...)
	cData = append(cData, TEST_VAR_IS_ATOM...)
	cData = append(cData, TEST_VAR_IS_VAR...)

	for i, c := range cData {

		fmt.Fprintln(sb)

		// special nil constraints - nothing to see !
		if c == nil {
			fmt.Fprintf(sb, "\n%d\t<nil>", i)
			continue
		}

		fmt.Fprintf(sb, "\n%d\t(original)\t%v", i, c.String())
		fmt.Fprintf(sb, "\n%d\t(raw form)\t%#v", i, c)

		c, err = c.Check()
		fmt.Fprintf(sb, "\n%d\t(checked)\t%v", i, c)
		if err != nil {
			fmt.Fprintf(sb, ", error : %v", err)
		}
	}

	mytest.Verify(t, sb.String(), "constraint_test.check")
}

package solver

import (
	"testing"

	"github.com/xavier268/myprolog/mytest"
)

func TestNumberConstraints(t *testing.T) {

	cc := cleanAllConstraints(generateConstraintsFromNumbers(TEST_NUMBERS))

	sb := runConstraintSimplify2x2Test(t, cc)
	mytest.Verify(t, sb.String(), "constraint_number_test")
}

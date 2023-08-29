package solver

import (
	"strings"
	"testing"

	"github.com/xavier268/myprolog/mytest"
)

func TestSimplify(t *testing.T) {
	var sb *strings.Builder

	tatom := cleanAllConstraints(TEST_VAR_IS_ATOM)
	tvar := cleanAllConstraints(TEST_VAR_IS_VAR)
	tstring := cleanAllConstraints(TEST_VAR_IS_STRING)
	tempty := cleanAllConstraints(TEST_EMPTY_CONSTRAINTS)

	sb = runConstraintSimplify2x2Test(t, tempty)
	mytest.Verify(t, sb.String(), "constraint_simplify_test.atoms")

	sb = runConstraintSimplify2x2Test(t, tstring)
	mytest.Verify(t, sb.String(), "constraint_simplify_test.strings")

	sb = runConstraintSimplify2x2Test(t, append(tstring, tatom...))
	mytest.Verify(t, sb.String(), "constraint_simplify_test.strings.atoms")

	sb = runConstraintSimplify2x2Test(t, tvar)
	mytest.Verify(t, sb.String(), "constraint_simplify_test.vars")

	sb = runConstraintSimplify2x2Test(t, append(tatom, tvar...))
	mytest.Verify(t, sb.String(), "constraint_simplify_test.atoms.vars")

	sb = runConstraintSimplify2x2Test(t, append(tstring, tvar...))
	mytest.Verify(t, sb.String(), "constraint_simplify_test.strings.vars")

}

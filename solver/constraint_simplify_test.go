package solver

import (
	"strings"
	"testing"

	"github.com/xavier268/mytest"

	_ "github.com/xavier268/myprolog/parser"
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

	var c1, c2 Constraint

	// details zoom1
	c1 = VarIsCompoundTerm{ // X2 = [c,b,a]
		V: X2,
		T: CompoundTerm{
			Functor: "dot",
			Children: []Term{
				Atom{Value: "c"},
				CompoundTerm{
					Functor: "dot",
					Children: []Term{
						Atom{Value: "b"},
						CompoundTerm{
							Functor: "dot",
							Children: []Term{
								Atom{Value: "a"},
								CompoundTerm{
									Functor: "dot",
								},
							},
						},
					},
				},
			},
		},
	}
	c2 = VarIsCompoundTerm{ // Z = [Y2,X2]
		V: Z,
		T: CompoundTerm{
			Functor: "dot",
			Children: []Term{
				Y2, X2,
			},
		},
	}

	sb = runConstraintSimplify2x2Test(t, []Constraint{c1, c2})
	mytest.Verify(t, sb.String(), "constraint_simplify_test.zoom.1")

	// detail zoom2
	c1 = VarIsCompoundTerm{ // X2 = 4
		V: X2,
		T: Number{
			Num:        4,
			Den:        1,
			Normalized: true,
		},
	}

	c2 = VarIsVar{ // Z = X2
		V: Z,
		W: X2,
	}

	sb = runConstraintSimplify2x2Test(t, []Constraint{c1, c2})
	mytest.Verify(t, sb.String(), "constraint_simplify_test.zoom.2")

	// detail zoom3
	c1 = VarIsCompoundTerm{ // X2 = 4
		V: X2,
		T: Number{
			Num:        4,
			Den:        1,
			Normalized: true,
		},
	}

	c2 = VarIsVar{ // X2 = Z
		V: X2,
		W: Z,
	}
	sb = runConstraintSimplify2x2Test(t, []Constraint{c1, c2})
	mytest.Verify(t, sb.String(), "constraint_simplify_test.zoom.3")

}

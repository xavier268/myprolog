package solver

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xavier268/myprolog/mytest"
	"github.com/xavier268/myprolog/parser"
)

var subTestData = []Term{
	Atom{Value: "a"},
	String{Value: "b"},
	Number{Num: 1, Den: 2, Normalized: false},
	Underscore{},

	Variable{Name: "X", Nsp: 0},
	Variable{Name: "X", Nsp: 1},
	Variable{Name: "Y", Nsp: 1},

	parser.MustParseString("a(X,c).", "test data")[0],
	parser.MustParseString("a(b,X).", "test data")[0],
	parser.MustParseString("a(X, dot(X,Y)).", "test data")[0],
	parser.MustParseString("a(Y,_, d).", "test data")[0],
}

func TestReplace(t *testing.T) {

	sb := new(strings.Builder)

	// Replacements
	Rep1 := Y2
	Rep2 := Underscore{}
	Rep3 := parser.MustParseString("a(b,c).", "test data")[0]
	Rep4 := parser.MustParseString("a(X,c).", "test data")[0] // positive occur chech - should never happen anyway, but should work here.

	var tt Term
	found := false

	fmt.Fprintf(sb, "============== Test Replace ==============\n")

	for i, d := range subTestData {

		fmt.Fprintln(sb)
		fmt.Fprintf(sb, "\n%d\tReplacing %s in : \t%s", i, X.String(), d.String())

		tt, found = ReplaceVar(X, d, Rep1)
		fmt.Fprintf(sb, "\n%d\t\tby %s \t\tresult: %s", i, Rep1.String(), tt.String())
		fmt.Fprintf(sb, "\n%d\t\t\tfound: %t", i, found)

		tt, found = ReplaceVar(X, d, Rep2)
		fmt.Fprintf(sb, "\n%d\t\tby %s \t\tresult: %s", i, Rep2.String(), tt.String())
		fmt.Fprintf(sb, "\n%d\t\t\tfound: %t", i, found)

		tt, found = ReplaceVar(X, d, Rep3)
		fmt.Fprintf(sb, "\n%d\t\tby %s \t\tresult: %s", i, Rep3.String(), tt.String())
		fmt.Fprintf(sb, "\n%d\t\t\tfound: %t", i, found)

		tt, found = ReplaceVar(X, d, Rep4)
		fmt.Fprintf(sb, "\n%d\t\tby %s \t\tresult: %s", i, Rep4.String(), tt.String())
		fmt.Fprintf(sb, "\n%d\t\t\tfound: %t", i, found)

	}
	mytest.Verify(t, sb.String(), "substitution_test.replace")
}

func TestFindVars(t *testing.T) {
	sb := new(strings.Builder)
	fmt.Fprintf(sb, "============== Test FindVars ==============\n")
	for i, d := range subTestData {
		fmt.Fprintf(sb, "\n%d\t%s contains ... %v", i, d.String(), FindVars(d))
	}

	mytest.Verify(t, sb.String(), "substitution_test.findvars")
}

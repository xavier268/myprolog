package parser

import (
	"fmt"
	"strings"
	"testing"
)

var testNumberData = []Number{

	{}, // Zero value (valid)
	NaN,
	MaxNumber,
	MinNumber,
	{Num: 5, Den: 1},
	{Num: 2, Den: 3},
	{Num: -6, Den: -1}, // denormalized
	{Num: -2, Den: 6},
	{Num: 10, Den: -6},
}

func TestNumber(t *testing.T) {

	sb := new(strings.Builder)

	fmt.Fprintf(sb, "Test name : %s\n", t.Name())
	fmt.Fprintf(sb, "Using a set of %d test numbers : %v\n", len(testNumberData), testNumberData)

	fmt.Fprintf(sb, "\n============= Single number test =============\n")
	for i, n := range testNumberData {
		fmt.Fprintln(sb)
		fmt.Fprintf(sb, "%d\t\t(original)      n: %v\n", i, n.String())
		fmt.Fprintf(sb, "%d\t\t(pretty)        n: %v\n", i, n.Pretty())
		fmt.Fprintf(sb, "%d\t\t(normalized)    n: %v\n", i, n.Normalize().Pretty())
		fmt.Fprintf(sb, "%d\t\t(ChSign)        n: %v\n", i, n.ChSign().Pretty())
		fmt.Fprintf(sb, "%d\t\t(IsZero)        n: %v\n", i, n.IsZero())
		fmt.Fprintf(sb, "%d\t\t(IsNaN)         n: %v\n", i, n.IsNaN())
		fmt.Fprintf(sb, "%d\t\t(IsInteger)     n: %v\n", i, n.IsInteger())
		fmt.Fprintf(sb, "%d\t\t(Floor)         n: %v\n", i, n.Floor().Pretty())
		fmt.Fprintf(sb, "%d\t\t(Ceil)          n: %v\n", i, n.Ceil().Pretty())
		if !n.IsNaN() { // avoid panic on NaN
			fmt.Fprintf(sb, "%d\t\t(ToInt)         n: %v\n", i, n.ToInt())
		}
	}

	fmt.Fprintf(sb, "\n============= Two numbers test =============\n")
	for i, n := range testNumberData {
		for j, m := range testNumberData {
			fmt.Fprintln(sb)
			fmt.Fprintf(sb, "%d-%d\t\t(pretty)      n,m: %v and %v\n", i, j, n.Pretty(), m.Pretty())
			fmt.Fprintf(sb, "%d-%d\t\t(Eq)          n,m: %v\n", i, j, n.Eq(m))
			fmt.Fprintf(sb, "%d-%d\t\t(Plus)        n,m: %v\n", i, j, n.Plus(m).Pretty())
			fmt.Fprintf(sb, "%d-%d\t\t(Minus)       n,m: %v\n", i, j, n.Minus(m).Pretty())
			fmt.Fprintf(sb, "%d-%d\t\t(Times)       n,m: %v\n", i, j, n.Times(m).Pretty())
			fmt.Fprintf(sb, "%d-%d\t\t(Less)        n,m: %v\n", i, j, n.Less(m))
		}
	}

	verifyTest(t, sb.String(), "definitions_number_test.want")

}

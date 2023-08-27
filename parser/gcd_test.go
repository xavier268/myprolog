package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xavier268/myprolog/mytest"
)

func TestGcd(t *testing.T) {
	data := []int{2, 3, 7, 3 * 7, -5, 67, 0, 1, 33, 6, 11 * 17, 11 * 18, 5 * 18, -6 * 18}
	sb := new(strings.Builder)

	fmt.Fprintln(sb, "Testing PGCD computations")

	for i, j := 0, 1; j < len(data); i, j = i+1, j+1 {
		a, b := data[i], data[j]
		p := Gcd(a, b)
		fmt.Fprintf(sb, "%d\t%d\t\t--> \t%d\n", a, b, p)
	}
	mytest.Verify(t, sb.String(), "gcd_test")

}

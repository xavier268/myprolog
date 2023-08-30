package repl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xavier268/myprolog/mytest"
	"github.com/xavier268/myprolog/parser"
	"github.com/xavier268/myprolog/solver"
)

// use as test harness to control there are no changes whean developping the code
func TestEndToEnd(t *testing.T) {

	inputData := []string{
		"a(b,c).  ?- a(X,Y).",
		"a(b,c). a(d,e).    ?- a(X,Y).",
		"a(b,c). a(d,e).    ?- a(Y,X).",
		"a(b,c).a(c,d).     ?- a(U,V),a(V,W).",
		"a(b,c). a(c,d).a(e,f). ?- a(X,Y).",
	}

	for i, input := range inputData { // one file per input

		sb := new(strings.Builder)

		// simple solution handler, prints all solutions sucessively, until nil state reached.
		sh := func(st *solver.State) *solver.State {
			if st == nil {
				fmt.Fprintf(sb, "\nsolution:\tnil state")
				return st
			} else {
				fmt.Fprintf(sb, "\nsolution:\t%v", st.Constraints)
				return st.Parent
			}
		}

		fmt.Fprintf(sb, "\n\nInput:\t%v", input)
		fmt.Printf("\nProcessing input:\t%v\n", input)
		tt, err := parser.ParseString(input, t.Name())
		if err != nil {
			t.Fatalf("Error parsing input: %v", err)
		}
		fmt.Fprintf(sb, "\nParsed:\t%v", tt)

		st := solver.NewState(nil)         // create new state
		st.Goals = append(st.Goals, tt...) // add input to goal list

		for {
			if st == nil || sb.Len() > 10_000 {
				break
			}
			st = solver.Solve(st, sh)
			fmt.Fprintf(sb, "\nState:\t%v", st)
		}
		mytest.Verify(t, sb.String(), fmt.Sprintf("endToEnd_test.%d", i))
	}
}

// Test to work in detail on a single expression
func TestEndToEndDetail(t *testing.T) {
	t.Skip()
	input := "a(X,Z) :- a(X,Y),a(Y,Z). a(b,c). a(c,d). a(d,e).  ?- a(X,Y)."

	sb := new(strings.Builder)

	// simple solution handler, prints all solutions sucessively, until nil state reached.
	sh := func(st *solver.State) *solver.State {

		if st == nil {
			fmt.Fprint(sb, "\n\n : Solution handler : nil state\n")
			return st
		} else {

			fmt.Fprintf(sb, "\nState:\t%v", st)
			fmt.Fprintf(sb, "\n\n=========> solution:\t%v\n", st.Constraints)

			return st.Parent
		}
	}

	fmt.Fprintf(sb, "\n==================\nInput:\t%v\n==================\n\n", input)
	st := solver.NewState(nil)                     // create new state
	tt, err := parser.ParseString(input, t.Name()) // parse input
	if err != nil {
		t.Fatalf("Error parsing input: %v", err)
	}
	st.Goals = append(st.Goals, tt...) // add input to goal list
	fmt.Fprintf(sb, "\nState:\t%v", st)
	for {
		if st == nil || sb.Len() > 5_000 {
			break
		}
		fmt.Fprintf(sb, "\nState:\t%v", st)
		st = solver.Solve(st, sh)
		fmt.Fprintf(sb, "\nState:\t%v", st)
	}
	mytest.Verify(t, sb.String(), "endToEnd_test.detail")

}

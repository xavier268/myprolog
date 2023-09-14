// Package example contains examples, that are also run as tests.
package examples

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xavier268/myprolog/parser"
	"github.com/xavier268/myprolog/solver"
	"github.com/xavier268/mytest"
)

// use as test harness to control there are no changes whean developping the code
func TestExamples(t *testing.T) {

	inputData := []string{
		`
		a(b,c).  
		?- a(X,Y).`,
		`
		a(b,c). 
		a(d,e).    
		?- a(X,Y).`,
		`
		a(b,c).
		a(d,e).
		?- a(Y,X).`,
		`
		a(b,c). 
		a(d,e).    
		// here, we make multiple queries at once - implicitely, we want all the queries to be simultaneously satisfied.
		?- a(X,Y).
		?-a(T,_).`,
		`
		a(b,c). 
		a(d,e).    
		// here, we make multiple queries at once - explicitely require all the queries to be simultaneously satisfied.
		?- a(X,Y),a(T,_).`,
		`
		a(b,c). 
		a(d,e).    
		// here, we make multiple queries at once - explicitely require one of them to besatisfied.
		?-  a(X,Y) ; a(T,_) .`,
		`
		a(b,c).a(c,d).     ?- a(U,V),a(V,W).`,
		`
		a(b,c). a(c,d).a(e,f). ?- a(X,Y).`,
		`
		a(b,c). a(c,d).a(e,f).  ?- a(_,Z).`,
		`
		a(b,c). a(c,d).a(e,f).  ?- a(T,_).`,
		`
		a(b,c). a(c,d).a(e,f).a(b,f).  ?- a(_,Z).`,
		`
		// define reverse a list 
		reverse_list(List, Reversed)                :-      reverse_list(List, [], Reversed).
		reverse_list([], Acc, Acc).
		reverse_list([Head|Tail], Acc, Reversed)    :-      reverse_list(Tail, [Head|Acc], Reversed).		
		// query
		?- reverse_list([a,b,c,d], Reversed).
		`,
		`
		// reverse a list with named variables
		reverse_list(List, Reversed)             :-      reverse_list(List, [], Reversed).
		reverse_list([], Acc, Acc).
		reverse_list([Head|Tail], Acc, Reversed) :-      reverse_list(Tail, [Head|Acc], Reversed).
		
		?- reverse_list([a,_,b,Z,d], Reversed).
		`,

		`
		// constraints disjonction
		a(a,b). 
		a(X,c) :- X=2 ; X=3 . // means X is 2 or 3, will generate 2 more solutions.
		                      // IMPORTANT : notice that the period (.) must never follow immediately a number, 
							  // because that would mean a float, and the end of phrase marker will not be found, 
							  // triggering a parser error.
		?- a(X,Y).`,
		`
		// constraints conjonction
		a(a,b).
		a(X,c) :- X=2 , X=3 . // means X is both 2 and 3, will always fail !.
		?- a(X,Y).`,

		`
		// this will create infinite recursion, that will need to be managed by depth control		
		a(b,c). 
		a(c,d).
		a(X,Y) :- a(X,V),a(V,Y).		
		?- a(A,B).`,
	}

	for i, input := range inputData { // one file per input

		sb := new(strings.Builder)

		// simple solution handler, prints all solutions sucessively, until nil state reached.
		sh := func(st *solver.State) *solver.State {
			if st == nil {
				fmt.Fprintf(sb, "\n----------\nSolution:\t%v", "no (more) solution.")
				return st
			} else {
				fmt.Fprintln(sb)
				fmt.Fprintf(sb, "\n----------\nSolution:\t%v", solver.FilterSolutions(st.Constraints))
				fmt.Fprintf(sb, "\nBecause of : \n%s\n", st.RulesHistory())
				return st.Parent
			}
		}

		fmt.Fprintf(sb, "\n\nInput:\n%v\n", input)
		fmt.Printf("\nProcessing input:\t%v\n", input)
		tt, err := parser.ParseString(input, t.Name())
		if err != nil {
			t.Fatalf("Error parsing input: %v", err)
		}
		// fmt.Fprintf(sb, "\nParsed:\t%v", tt)

		st := solver.NewState(nil)         // create new state
		st.Goals = append(st.Goals, tt...) // add input to goal list

		for {
			if st == nil || sb.Len() > 10_000 {
				break
			}
			st.Session.ForceDepthControl(5) // maintain depth control at 5, examples  are not that big.
			st = solver.Solve(st, sh)
			fmt.Fprintf(sb, "\nState:\t%v", st)
		}
		mytest.Verify(t, sb.String(), fmt.Sprintf("example.%03d", i))
	}
}

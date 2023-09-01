package solver

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xavier268/myprolog/mytest"
	"github.com/xavier268/myprolog/parser"
)

// use as test harness to control there are no changes whean developping the code
func TestEndToEnd(t *testing.T) {

	inputData := []string{
		"a(b,c).  ?- a(X,Y).",
		"a(b,c). a(d,e).    ?- a(X,Y).",
		"a(b,c). a(d,e).    ?- a(Y,X).",
		"a(b,c).a(c,d).     ?- a(U,V),a(V,W).",
		"a(b,c). a(c,d).a(e,f). ?- a(X,Y).",
		"a(b,c). a(c,d).a(e,f).  ?- a(_,Z).",
		"a(b,c). a(c,d).a(e,f).  ?- a(T,_).",
		"a(b,c). a(c,d).a(e,f).a(b,f).  ?- a(_,Z).",
		"a(b,c). a(c,d).  a(X,Y) :- a(X,V),a(V,Y).  ?- a(A,B).",
		` // define reverse a list
		reverse_list(List, Reversed) 				:-		reverse_list(List, [], Reversed).
		reverse_list([], Acc, Acc).
		reverse_list([Head|Tail], Acc, Reversed) 	:-		reverse_list(Tail, [Head|Acc], Reversed).		
		// query
		?- reverse_list([a,b,c,d], Reversed).
		`,
		` // define reverse a list
		reverse_list(List, Reversed) 				:-		reverse_list(List, [], Reversed).
		reverse_list([], Acc, Acc).
		reverse_list([Head|Tail], Acc, Reversed) 	:-		reverse_list(Tail, [Head|Acc], Reversed).		
		// query
		?- reverse_list([1,2,3,4], Reversed).
		`,
		` // define reverse a list
		reverse_list(List, Reversed) 				:-		reverse_list(List, [], Reversed).
		reverse_list([], Acc, Acc).
		reverse_list([Head|Tail], Acc, Reversed) 	:-		reverse_list(Tail, [Head|Acc], Reversed).		
		// query
		?- reverse_list(["a","b","c","d"], Reversed).
		`,
		` 	// reverse a list
		reverse_list(List, Reversed) :-
    	reverse_list(List, [], Reversed).
		reverse_list([], Acc, Acc).
		reverse_list([Head|Tail], Acc, Reversed) :-
    	reverse_list(Tail, [Head|Acc], Reversed).	
		?- reverse_list([a,_,b], Reversed).
		`,
		` 	// reverse a list
		reverse_list(List, Reversed) :-
    	reverse_list(List, [], Reversed).
		reverse_list([], Acc, Acc).
		reverse_list([Head|Tail], Acc, Reversed) :-
    	reverse_list(Tail, [Head|Acc], Reversed).	
		?- reverse_list([a,_,b,c,d], Reversed).
		`,

		` 	// reverse a list
		reverse_list(List, Reversed) :-
    	reverse_list(List, [], Reversed).
		reverse_list([], Acc, Acc).
		reverse_list([Head|Tail], Acc, Reversed) :-
    	reverse_list(Tail, [Head|Acc], Reversed).	
		?- reverse_list([a,X,Y], Reversed).
		`,

		// constraints =
		"a(a,b). a(X,c) :- X=2 ; X=3 .  ?- a(X,Y).",
		"a(a,b). a(X,c) :- X=2 , X=3 .  ?- a(X,Y).",
		"a(a,b). a(X,c) :- X='a' , X=3 .  ?- a(X,Y).",
	}

	for i, input := range inputData { // one file per input

		sb := new(strings.Builder)

		// simple solution handler, prints all solutions sucessively, until nil state reached.
		sh := func(st *State) *State {
			if st == nil {
				fmt.Fprintf(sb, "\nsolution:\tnil state")
				return st
			} else {
				fmt.Fprintf(sb, "\nsolution:\t%v", FilterSolutions(st.Constraints))
				fmt.Fprintf(sb, "\nRules applied : \n%s\n", st.RulesHistory())
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

		st := NewState(nil)                // create new state
		st.Goals = append(st.Goals, tt...) // add input to goal list

		for {
			if st == nil || sb.Len() > 10_000 {
				break
			}
			st = Solve(st, sh)
			fmt.Fprintf(sb, "\nState:\t%v", st)
		}
		mytest.Verify(t, sb.String(), fmt.Sprintf("endToEnd_test.%d", i))
	}
}

// Test to work in detail on a single expression
func TestEndToEndDetail(t *testing.T) {

	input := "a(a,b). a(X,c) :- X=2 ; X=3 .  ?- a(X,Y)."

	sb := new(strings.Builder)

	// simple solution handler, prints all solutions sucessively, until nil state reached.
	sh := func(st *State) *State {

		if st == nil {
			fmt.Fprint(sb, "\n\n : Solution handler : nil state\n")
			return st
		} else {

			fmt.Fprintf(sb, "\n=========> solution cleaned:\t%v", FilterSolutions(st.Constraints))
			fmt.Fprintf(sb, "\n=========> All constraints :\t%v", st.Constraints)
			fmt.Fprintf(sb, "\n=========> Known rules:\n%s", st.Session.ListRules())
			fmt.Fprintf(sb, "\n=========> Rules applied :\n%s", st.RulesHistory())
			fmt.Fprintf(sb, "\n State:\n%v", st)

			fmt.Fprintln(sb)

			return st.Parent
		}
	}

	fmt.Fprintf(sb, "\n==================\nInput:\t%v\n==================\n\n", input)
	st := NewState(nil)                            // create new state
	tt, err := parser.ParseString(input, t.Name()) // parse input
	if err != nil {
		t.Fatalf("Error parsing input: %v\n%s", err, input)
	}
	st.Goals = append(st.Goals, tt...) // add input to goal list
	fmt.Fprintf(sb, "\nState:\t%v", st)
	for {
		if st == nil || sb.Len() > 5_000 {
			break
		}
		fmt.Fprintf(sb, "\nState:\t%v", st)
		st = Solve(st, sh)
		fmt.Fprintf(sb, "\nState:\t%v", st)
	}
	mytest.Verify(t, sb.String(), "endToEnd_test.detail")

}

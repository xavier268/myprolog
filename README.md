
[![Go Reference](https://pkg.go.dev/badge/github.com/xavier268/myprolog.svg)](https://pkg.go.dev/github.com/xavier268/myprolog) 
[![Go Report Card](https://goreportcard.com/badge/github.com/xavier268/myprolog)](https://goreportcard.com/report/github.com/xavier268/myprolog)
# myprolog
A Prolog (the langage) implementation for educational purposes

Inspired by the description in https://www.metalevel.at/prolog

# How to use ?

Typing 'go run .' in the root directory will launch the interactive interpreter.

To display the menu option, type 'h'.

You can enter facts, rules or queries, in any order, always terminating with a '.'.

See examples in the test files named *_endToEnd_test.xxx*. For instance :

```
        // reverse a list
		reverse_list(List, Reversed) :-
    	reverse_list(List, [], Reversed).
		reverse_list([], Acc, Acc).
		reverse_list([Head|Tail], Acc, Reversed) :-
    	reverse_list(Tail, [Head|Acc], Reversed).	
		?- reverse_list([a,X,3/4], Reversed).

```

will generate the solution : Reversed = [ 3/4, X, a ]

When multiple solutions are possible, there are successively proposed.

If prolog detects an infinite tree, it stops its exploration and looks for other solutions first.
This allows some non optimal  recursion to still generate solutions, such as in :

```

        a(b,c). 
        a(c,d).  
        a(X,Y) :- a(X,V),a(V,Y).  
        
        ?- a(A,B). 

```
will produce the 3 solutions : 	[A = b B = c] , [ A = c B = d]  and [A = b B = d] and then stop.

# Work in progress

Facts, rules, and queries must be terminated by a '.' (period).

',' (comma) and ';' (semi-colon) cand be used to indicate conjunction or disjunction between terms.

Errors are raisonnably detailled and reported, and should not crash the solver. Error on entry is reported and ignored.
Types of objects can be "strings", atoms, Variable, compound(terms, "etc..). The '_' (underscore) is a special variable that can take any value.
Lists can be entered as bracket lists : [ 1,2,3,] , as pairs [ a | b ] , or as dot(a,dot(b,dot())).

All Numbers are handled exactly, internally represented as int64 based rationnals. Rational numbers can also be entered directly, as in : lt(X,3/2). Although overflow could happen, precision is garanteed, ie 1/3 + 1/3 will always exactly equals 2/3.

Interval artithmethic is performed exactly, taking into account if number is integer or not.

# Predefined predicates

Compound form predicates require parenthesis :
* rule ( ...) : defines a rule structure, same as infix ':-'
* query ( ...)  : defines query structure, same as infix '?-'
* dot ( a, b), dot(a), dot() : are used to define lists [ 1,2,3 ] and pairs [ a | b ]
* and ( ...): same as infix ','
* or (...) : same as infix ';'
* number (X) : true only if child is a number
* integer (X) : true only if child is an integer number
* load ("path","to","file.ext") : load a file and evaluate it
* print ("message", "or", 3/4, atom, ...) : print on the console

Atomic form predicates :
* ! : the cut predicate, prevents backtracking
* fail : always fails
* "rules" : print the list of known rules and facts (debugging)
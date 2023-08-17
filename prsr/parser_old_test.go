package prsr

/*

import (
	"fmt"
	"testing"

	"github.com/xavier268/myprolog/node"
	"github.com/xavier268/myprolog/tknz"
)

func TestParser0(t *testing.T) {

	tab := map[string]struct { // map sources to ...
		ok  bool   // true if no error
		exp string // parsed output
	}{
		// basic atoms
		"aa":   {true, " aa"},
		"22":   {true, " 22"},
		"X2":   {true, " X2"},
		"+":    {true, " +"},
		":-":   {true, " :-"},
		"a b ": {true, " a b"},

		// parenthesis - simple
		"aa()":  {true, " aa"},
		"X()":   {false, " X"},  // not a functor
		"2()":   {false, " 2"},  // not a functor
		"aa(":   {false, " aa"}, // unbalanced parenth
		"aa)":   {false, " aa"}, // unbalanced parenth
		"+()":   {true, " +"},
		"+(aa)": {true, " + ( aa )"},

		// parenthesis - complex
		"()":          {false, ""},
		"(bb)":        {false, ""},
		"aa(())":      {false, " aa"},
		"aa(f(f)))":   {false, " aa ( f ( f ) )"},
		"aa(f)(f)))":  {false, " aa ( f f )"},  // unbalanced
		"aa((f)(f))":  {false, " aa"},          // adding the first (f) to a nil node ...
		"aa(f)(f)":    {true, " aa ( f f )"},   // ok to split the parameter list in multiple groups ...
		"aa(bb)(cc)":  {true, " aa ( bb cc )"}, // ok to split the parameter list in multiple groups ...
		"aa(f()aa)()": {true, " aa ( f aa )"},  // ok to split the parameter list in multiple groups ...

		"a,b":  {true, " a b"},
		"a,,b": {true, " a b"},
		",a b": {true, " a b"},
		"a,b,": {true, " a b"},

		// nil - checks will be made later ...
		"nil":    {true, " nil"},  //ok
		"nil ()": {false, " nil"}, //nok, because not a functor
		"nil(a)": {false, " nil"}, //nok, because not a functor

		// quotes
		"\"a b \"()":  {true, " a b "},
		"\"a b \"(c)": {true, " a b  ( c )"},
	}
	for i, want := range tab {
		tzr := tknz.NewTokenizerString(i)
		root := node.NewNode("test")
		err := parse0(tzr, root, new(int))
		if err != nil {
			fmt.Printf("info : %20s -> %v\n", i, err)
		}
		if (err == nil) != want.ok {
			t.Fatalf("unexpected error : %v,\nfor test : %s -> %v", err, i, tab[i])
		}
		if want.exp != "" && root.String() != " test ("+want.exp+" )" {
			t.Fatalf("unexpected result	for test : %s\ngot : %s\nwant: %s", i, root, " test ("+want.exp+" )")
		}
		if want.exp == "" && root.String() != " test" {
			t.Fatalf("unexpected result	\ngot : >%s<\nwant: > test<", root)
		}
	}
}

func TestPreProcList(t *testing.T) {
	tab := map[string]struct { // map sources to ...
		ok  bool   // true if no error
		exp string // non idented value of parsed output
	}{

		// canonical
		"a":        {true, "a"},           // ok
		"nil":      {true, "nil"},         // ok
		"nil()":    {false, "nil"},        // ok
		"nil(a)":   {false, "nil"},        // nok, check on nil as functor !
		"dot(a,b)": {true, "dot ( a b )"}, // ok

		// bar lists
		"[|]":   {false, "[ | ]"},      // nok
		"[a|]":  {false, "[ a | ]"},    // nok
		"[|b]":  {false, "[ | b ]"},    // nok
		"[a|b]": {true, "dot ( a b )"}, // ok

		// bracket lists
		"[]":          {true, "dot ( nil nil )"},                         // ok
		"[nil]":       {true, "dot ( nil dot ( nil nil ) )"},             // ok
		"[a]":         {true, "dot ( a dot ( nil nil ) )"},               // ok
		"[ a b ]":     {true, "dot ( a dot ( b dot ( nil nil ) ) )"},     // ok
		"x [ a b ] c": {true, "x dot ( a dot ( b dot ( nil nil ) ) ) c"}, // ok

		// mixed
		"[a(X,y)]":      {true, "dot ( a ( X y ) dot ( nil nil ) )"},                                     // ok
		"[a(X,y)|b(T)]": {true, "dot ( a ( X y ) b ( T ) )"},                                             // ok
		"[a|[b c]]":     {true, "dot ( a dot ( b dot ( c dot ( nil nil ) ) ) )"},                         // ok
		"[a b c]":       {true, "dot ( a dot ( b dot ( c dot ( nil nil ) ) ) )"},                         // ok
		"[a [b] c]":     {true, "dot ( a dot ( dot ( b dot ( nil nil ) ) dot ( c dot ( nil nil ) ) ) )"}, // ok
		"[a [b c]]":     {true, "dot ( a dot ( dot ( b dot ( c dot ( nil nil ) ) ) dot ( nil nil ) ) )"}, // ok

	}

	for src, got := range tab {
		tzr := tknz.NewTokenizerString(src)
		root := node.NewNode("test")
		err := parse0(tzr, root, new(int))
		if err == nil {
			// no need to preProcBar if already previous error !
			err = preProcList(root)
		}
		if err != nil {
			fmt.Printf("info : %20s -> %v\n", src, err)
		}
		if (err == nil) != got.ok {
			t.Fatalf("unexpected error : %v,\nfor test : %s -> %v", err, src, tab[src])
		}
		if got.exp != "" && root.String() != " test ( "+got.exp+" )" {
			t.Fatalf("unexpected result	for %s\ngot : %s\nwant: %s", src, root, " test ( "+got.exp+" )")
		}
		if got.exp == "" && root.String() != " test " {
			t.Fatalf("unexpected result	\ngot : >%s<\nwant: > test <", root)
		}
	}
}

func TestPreProcRule(t *testing.T) {
	tab := map[string]struct { // map sources to ...
		ok  bool   // true if no error
		exp string // non idented value of parsed output
	}{
		// basic atoms, no alternatives
		":-":                  {false, ":-"},                                      // illegal empty canonical
		":-(a)":               {true, "rule ( a )"},                               // missing period
		":-(aa)":              {true, "rule ( aa )"},                              // valid canonical form
		"aa.":                 {true, "rule ( aa )"},                              // valid fact
		"X.":                  {false, "X ."},                                     // Variable as a rule
		"X:-a.":               {true, "rule ( X a )"},                             // No error here !
		"a":                   {false, "a"},                                       // missing final period
		"2.":                  {false, "2"},                                       // missing final period, because part of number !
		"2..":                 {false, "2 ."},                                     // one period for the number, the other for the rule, but a Number is not a rule on its own,
		"a:-":                 {false, "a :-"},                                    // missing final period
		"a:-b.":               {true, "rule ( a b )"},                             // ok
		"a:-.":                {true, "rule ( a )"},                               // ok
		"a.":                  {true, "rule ( a )"},                               // ok
		"a.b.":                {true, "rule ( a ) rule ( b )"},                    // ok
		"a.:-(b)":             {true, "rule ( a ) rule ( b )"},                    // ok
		":-(a)b.":             {true, "rule ( a ) rule ( b )"},                    // ok
		"a.:-(b c)":           {true, "rule ( a ) rule ( b c )"},                  // ok
		":-(a c)b.":           {true, "rule ( a c ) rule ( b )"},                  // ok
		"a:-b c. d :- e f.":   {true, "rule ( a b c ) rule ( d e f )"},            // ok
		"a:-b c. d :- e f.g.": {true, "rule ( a b c ) rule ( d e f ) rule ( g )"}, // ok
		"a:-b c. d :- e.g.":   {true, "rule ( a b c ) rule ( d e ) rule ( g )"},   // ok
		"a:-b c. d :- .g.":    {true, "rule ( a b c ) rule ( d ) rule ( g )"},     // ok

		// // compound terms, no alternative
		"a(X,Y).":           {true, "rule ( a ( X Y ) )"},             // ok
		"a(X,Y):-.":         {true, "rule ( a ( X Y ) )"},             // ok
		"a(X,Y):-b(Y,X).":   {true, "rule ( a ( X Y ) b ( Y X ) )"},   // ok
		"a(X,Y):-b(Y,X),X.": {true, "rule ( a ( X Y ) b ( Y X ) X )"}, // ok
		"a(X,_):-b(_,X),X.": {true, "rule ( a ( X _ ) b ( _ X ) X )"}, // ok

		// // alternatives

		"a:-c;d e.":    {true, "rule ( a c ) rule ( a d e )"},               // ok
		"a:-b;c;d e.":  {true, "rule ( a b ) rule ( a c ) rule ( a d e )"},  // ok
		"a:-b;c.":      {true, "rule ( a b ) rule ( a c )"},                 // ok
		"a:-b;c;d e":   {false, "a :- b ; c ; d e"},                         // missing period
		"a:-b;d e":     {false, "a :- b ; d e"},                             // missing period
		"a:-b;":        {false, "a :- b ;"},                                 // missing period
		"a;c.":         {false, "a ; c ."},                                  // missing  :-
		"a(X,Y):-b;c.": {true, "rule ( a ( X Y ) b ) rule ( a ( X Y ) c )"}, // ok
		"a:-b.;c.":     {false, "rule ( a b ) ; c ."},                       // syntax error
		"a:-b.d;c.":    {false, "rule ( a b ) d ; c ."},                     // syntax error
		"a:-;d.":       {true, "rule ( a ) rule ( a d )"},                   // ok
		"a;d.":         {false, "a ; d ."},                                  // nok

		// Queries
		"?q.":     {true, "query ( q )"},    // ok
		"?q r .":  {true, "query ( q r )"},  // ok
		"a?q r .": {false, "a ? q r ."},     // neither a query nor a rule
		"? a":     {false, "? a"},           // missing period
		"?? a":    {false, "? ? a"},         // syntax error
		"?q.b":    {false, "query ( q ) b"}, // nok

		// mixed rules - queries
		"?q.a.":   {true, "query ( q ) rule ( a )"},            // ok
		"b.?q.a.": {true, "rule ( b ) query ( q ) rule ( a )"}, // ok
		"b?q.a.":  {false, "b ? q . a ."},                      // missing period
		"b.?q.a":  {false, "rule ( b ) query ( q ) a"},         // missing period

	}

	for src, got := range tab {
		tzr := tknz.NewTokenizerString(src)
		root := node.NewNode("test")
		err := parse0(tzr, root, new(int))
		if err != nil {
			t.Fatalf("parse0 should not fail for %s with %e", src, err)
		}
		err = preProcRule(root)
		if err != nil {
			fmt.Printf("info : %20s -> %v\n", src, err)
		}
		if (err == nil) != got.ok {
			t.Fatalf("unexpected error : %v,\nfor test : %s -> %v", err, src, tab[src])
		}
		if got.exp != "" && root.String() != " test ( "+got.exp+" )" {
			t.Fatalf("unexpected result	for %s\ngot : %s\nwant: %s", src, root, " test ( "+got.exp+" )")
		}
		if got.exp == "" && root.String() != "got " {
			t.Fatalf("unexpected result	\ngot : >%s<\nwant: >got <", root)
		}
	}
}

*/

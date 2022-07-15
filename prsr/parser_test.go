package prsr

import (
	"fmt"
	"testing"

	"github.com/xavier268/myprolog/node"
	"github.com/xavier268/myprolog/tknz"
)

func TestParser0(t *testing.T) {

	tab := map[string]struct { // map sources to ...
		ok  bool   // true if no error
		exp string // non idented value of parsed output
	}{
		// basic atoms
		"aa":       {true, " aa"},
		"22":       {true, " 22"},
		"X2":       {true, " X2"},
		"+":        {true, " +"},
		":-":       {true, " :-"},
		"\"a b \"": {true, " \"a b \""},

		// parenthesis - simple
		"aa()":        {true, " aa"},
		"X()":         {false, " X"},  // not a functor
		"2()":         {false, " 2"},  // not a functor
		"aa(":         {false, " aa"}, // unbalanced parenth
		"aa)":         {false, " aa"}, // unbalanced parenth
		"+()":         {true, " +"},
		"+(aa)":       {true, " + ( aa )"},
		"\"a b \"()":  {true, " \"a b \""},
		"\"a b \"(c)": {true, " \"a b \" ( c )"},

		// parenthesis - complex
		"()":          {false, ""},
		"(bb)":        {false, ""},
		"aa(())":      {false, " aa"},
		"aa(f(f)))":   {false, " aa ( f ( f ) )"},
		"aa(f)(f)))":  {false, " aa ( f f )"},  // unbalanced
		"aa((f)(f))":  {false, " aa"},          // adding the firts (f) to a nil node ...
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

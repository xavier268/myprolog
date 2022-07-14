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
		"aa":       {true, "aa"},
		"22":       {true, "22"},
		"X2":       {true, "X2"},
		"+":        {true, "+"},
		":-":       {true, ":-"},
		"\"a b \"": {true, "\"a b \""},

		// parenthesis - simple
		"aa()":        {true, "aa"},
		"X()":         {false, "X"},
		"2()":         {false, "2"},
		"aa(":         {false, "aa"},
		"+()":         {true, "+"},
		"+(aa)":       {true, "+ ( aa )"},
		"\"a b \"()":  {true, "\"a b \""},
		"\"a b \"(c)": {true, "\"a b \" ( c )"},

		// parenthesis - complex
		"()":          {false, ""},
		"aa(())":      {false, "aa"},
		"aa)":         {false, "aa"},
		"aa(f(f)))":   {false, "aa ( f ( f ) )"},
		"aa(f)(f)))":  {false, "aa ( f )"},
		"aa(bb)(cc)":  {false, "aa ( bb )"},
		"aa(f()aa)()": {false, "aa ( f aa )"},

		"a,b":  {true, "a b"},
		"a,,b": {true, "a b"},
		",a b": {true, "a b"},
		"a,b,": {true, "a b"},

		// nil - checks will be made later ...
		"nil":    {true, "nil"},  //ok
		"nil ()": {false, "nil"}, //not a functor
		"nil(a)": {false, "nil"}, //not a functor

	}
	for src, got := range tab {
		tzr := tknz.NewTokenizerString(src)
		root := node.NewNode("got")
		err := parse0(tzr, root, new(int))
		if err != nil {
			fmt.Printf("info : %20s -> %v\n", src, err)
		}
		if (err == nil) != got.ok {
			t.Fatalf("unexpected error : %v,\nfor test : %s -> %v", err, src, tab[src])
		}
		if got.exp != "" && root.String() != "got ( "+got.exp+" ) " {
			t.Fatalf("unexpected result	for test : %s\ngot : %s\nwant: %s", src, root, "got ( "+got.exp+" ) ")
		}
		if got.exp == "" && root.String() != "got " {
			t.Fatalf("unexpected result	\ngot : >%s<\nwant: >got <", root)
		}
	}
}

package inter

import (
	"fmt"
	"testing"
)

func TestParse0Visual(t *testing.T) {
	//t.Skip()
	src := `a(b(1+2)c[5,6])`
	fmt.Println("__Source___")
	fmt.Println(src)
	pi := NewInter()
	tzr := NewTokenizerString(src)
	root := pi.nodeFor("parsed")
	err := pi.parse0(tzr, root, new(int))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("__Plain____")
	fmt.Println(root)
	fmt.Println("__Indent___")
	fmt.Println(root.StringPretty())
	root.DumpTree(true)
	pi.dumpSymt()
	fmt.Println("Testing if root is constant  ...")
	if root.isConstant() {
		pi.dumpSymt()
		fmt.Println("root IS constant/immutable")
	} else {
		pi.dumpSymt()
		fmt.Println("root IS NOT constant/immutable")
	}
	fmt.Println("Proactively marking constant for root ...")
	root.markConstant()
	pi.dumpSymt()
}

func TestParse0Table(t *testing.T) {

	tab := map[string]struct { // map sources to ...
		ok  bool   // true if no error
		exp string // non idented value of parsed output
	}{
		// basic atoms
		"aa":       {true, "aa"},
		"22":       {true, "22"},
		"X2":       {true, "X2"},
		"+":        {true, "+"},
		":-":       {true, ": -"},
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
	pi := NewInter()
	for src, got := range tab {
		tzr := NewTokenizerString(src)
		root := pi.nodeFor("got")
		err := pi.parse0(tzr, root, new(int))
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

func TestPreProcRule(t *testing.T) {
	tab := map[string]struct { // map sources to ...
		ok  bool   // true if no error
		exp string // non idented value of parsed output
	}{
		// basic atoms, no alternatives
		"~":                 {false, "~"},                              // illegal empty canonical
		"~(a)":              {true, "~ ( a )"},                         // missing period
		"~(aa)":             {true, "~ ( aa )"},                        // valid canonical form
		"aa.":               {true, "~ ( aa )"},                        // valid fact
		"X.":                {false, "X ."},                            // variable head error
		"a":                 {false, "a"},                              // missing final period
		"a~":                {false, "a ~"},                            // missing final period
		"a~b.":              {true, "~ ( a b )"},                       // ok
		"a~.":               {true, "~ ( a )"},                         // ok
		"a.":                {true, "~ ( a )"},                         // ok
		"a.b.":              {true, "~ ( a ) ~ ( b )"},                 // ok
		"a.~(b)":            {true, "~ ( a ) ~ ( b )"},                 // ok
		"~(a)b.":            {true, "~ ( a ) ~ ( b )"},                 // ok
		"a.~(b c)":          {true, "~ ( a ) ~ ( b c )"},               // ok
		"~(a c)b.":          {true, "~ ( a c ) ~ ( b )"},               // ok
		"a~b c. d ~ e f.":   {true, "~ ( a b c ) ~ ( d e f )"},         // ok
		"a~b c. d ~ e f.g.": {true, "~ ( a b c ) ~ ( d e f ) ~ ( g )"}, // ok
		"a~b c. d ~ e.g.":   {true, "~ ( a b c ) ~ ( d e ) ~ ( g )"},   // ok
		"a~b c. d ~ .g.":    {true, "~ ( a b c ) ~ ( d ) ~ ( g )"},     // ok

		// compound terms, no alternative
		"a(X,Y).":          {true, "~ ( a ( X Y ) )"},             // ok
		"a(X,Y)~.":         {true, "~ ( a ( X Y ) )"},             // ok
		"a(X,Y)~b(Y,X).":   {true, "~ ( a ( X Y ) b ( Y X ) )"},   // ok
		"a(X,Y)~b(Y,X),X.": {true, "~ ( a ( X Y ) b ( Y X ) X )"}, // ok
		"a(X,_)~b(_,X),X.": {true, "~ ( a ( X _ ) b ( _ X ) X )"}, // ok

		// alternatives
		"a~b;c;d e.":  {true, "~ ( a b ) ~ ( a c ) ~ ( a d e )"},     // ok
		"a~b;c.":      {true, "~ ( a b ) ~ ( a c )"},                 // ok
		"a~b;c;d e":   {false, "a ~ b ; c ; d e"},                    // missing period
		"a~b;d e":     {false, "a ~ b ; d e"},                        // missing period
		"a~b;":        {false, "a ~ b ;"},                            // missing period
		"a;c.":        {false, "a ; c ."},                            // missing  ~
		"a(X,Y)~b;c.": {true, "~ ( a ( X Y ) b ) ~ ( a ( X Y ) c )"}, // ok
		"a~b.;c.":     {false, "~ ( a b ) ; c ."},                    // syntax error
		"a~b.d;c.":    {false, "~ ( a b ) d ; c ."},                  // syntax error
		"a~;d.":       {true, "~ ( a ) ~ ( a d )"},                   // ok

	}
	pi := NewInter()
	for src, got := range tab {
		tzr := NewTokenizerString(src)
		root := pi.nodeFor("got")
		err := pi.parse0(tzr, root, new(int))
		if err != nil {
			t.Fatalf("parse0 should not fail for %s with %e", src, err)
		}
		err = pi.preProcRule(root)
		if err != nil {
			fmt.Printf("info : %20s -> %v\n", src, err)
		}
		if (err == nil) != got.ok {
			t.Fatalf("unexpected error : %v,\nfor test : %s -> %v", err, src, tab[src])
		}
		if got.exp != "" && root.String() != "got ( "+got.exp+" ) " {
			t.Fatalf("unexpected result	for %s\ngot : %s\nwant: %s", src, root, "got ( "+got.exp+" ) ")
		}
		if got.exp == "" && root.String() != "got " {
			t.Fatalf("unexpected result	\ngot : >%s<\nwant: >got <", root)
		}
	}
}

func TestFunctor(t *testing.T) {
	if isFunctor("nil") {
		t.Fatal("nil is not a functor")
	}
	if isFunctor("_") {
		t.Fatal("_ is not a functor")
	}
	if isFunctor("22") {
		t.Fatal("22 is not a functor")
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
		"[]":      {true, "dot ( nil nil )"},                     // ok
		"[a]":     {true, "dot ( a dot ( nil nil ) )"},           // ok
		"[ a b ]": {true, "dot ( a dot ( b dot ( nil nil ) ) )"}, // ok

		// mixed
		"[a(X,y)]":      {true, "dot ( a ( X y ) dot ( nil nil ) )"},                                     // ok
		"[a(X,y)|b(T)]": {true, "dot ( a ( X y ) b ( T ) )"},                                             // ok
		"[a|[b c]]":     {true, "dot ( a dot ( b dot ( c dot ( nil nil ) ) ) )"},                         // ok
		"[a b c]":       {true, "dot ( a dot ( b dot ( c dot ( nil nil ) ) ) )"},                         // ok
		"[a [b] c]":     {true, "dot ( a dot ( dot ( b dot ( nil nil ) ) dot ( c dot ( nil nil ) ) ) )"}, // ok
		"[a [b c]]":     {true, "dot ( a dot ( dot ( b dot ( c dot ( nil nil ) ) ) dot ( nil nil ) ) )"}, // ok

	}
	pi := NewInter()
	for src, got := range tab {
		tzr := NewTokenizerString(src)
		root := pi.nodeFor("got")
		err := pi.parse0(tzr, root, new(int))
		if err == nil {
			// no need to preProcBar if already previous error !
			err = pi.preProcList(root)
		}
		if err != nil {
			fmt.Printf("info : %20s -> %v\n", src, err)
		}
		if (err == nil) != got.ok {
			t.Fatalf("unexpected error : %v,\nfor test : %s -> %v", err, src, tab[src])
		}
		if got.exp != "" && root.String() != "got ( "+got.exp+" ) " {
			t.Fatalf("unexpected result	for %s\ngot : %s\nwant: %s", src, root, "got ( "+got.exp+" ) ")
		}
		if got.exp == "" && root.String() != "got " {
			t.Fatalf("unexpected result	\ngot : >%s<\nwant: >got <", root)
		}
	}
}

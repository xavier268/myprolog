package inter

import (
	"fmt"
	"testing"
)

func TestParseVisual(t *testing.T) {

	src := `f(1 555 X Y gggggg(deux f ( Z 666 _ ) 5 _ 5)`
	fmt.Println("__Source___")
	fmt.Println(src)
	pi := NewInter()
	tzr := NewTokenizerString(src)
	root := pi.nodeFor("parsed")
	err := pi.Parse(tzr, root)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("__Plain____")
	fmt.Println(root)
	fmt.Println("__Indent___")
	fmt.Println(root.StringIndent())
	root.dump()
	pi.dumpSymt()
}

func TestParseTable(t *testing.T) {

	tab := map[string]struct { // map sources to ...
		ok  bool   // true if no error
		exp string // non idented value of parsed output
	}{
		// basic atoms
		"aa": {true, "aa"},
		"22": {true, "22"},
		"X2": {true, "X2"},

		// parenthesis
		"aa()":       {true, "aa"},
		"aa(":        {false, "aa"},
		"()":         {false, ""},
		"aa(())":     {false, "aa"},
		"aa)":        {false, "aa"},
		"aa(f(f)))":  {false, "aa ( f ( f ) )"},
		"aa(f)(f)))": {false, "aa ( f f )"},

		"a,b":  {true, "a b"},
		"a,,b": {true, "a b"},
		",a b": {true, "a b"},
		"a,b,": {true, "a b"},
	}
	pi := NewInter()
	for src, got := range tab {
		tzr := NewTokenizerString(src)
		root := pi.nodeFor("got")
		err := pi.Parse(tzr, root)
		if err != nil {
			fmt.Printf("info : %s -> %v\n", src, err)
		}
		if (err == nil) != got.ok {
			t.Fatalf("unexpected error : %v,\nfor test : %s -> %v", err, src, tab[src])
		}
		if got.exp != "" && root.String() != "got ( "+got.exp+" ) " {
			t.Fatalf("unexpected result	\ngot : %s\nwant: %s", root, "got ( "+got.exp+" ) ")
		}
		if got.exp == "" && root.String() != "got " {
			t.Fatalf("unexpected result	\ngot : >%s<\nwant: >got <", root)
		}
	}

}

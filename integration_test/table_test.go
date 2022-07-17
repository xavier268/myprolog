package integration

import (
	"fmt"
	"testing"

	"github.com/xavier268/myprolog/repl"
)

type testDataType struct {
	src string // input program
	ok  bool   // ok if no error
	res string //resulting constraints
}

var testData = []testDataType{

	// Basic, simple solutions
	{"toto(a,b).	?toto(X,Y).", true, "[X = a,  Y = b, ]"}, // ultra simple
	{"?toto(X,Y).	toto(a,b).", true, "[X = a,  Y = b, ]"}, // query comes first ..
	{"?toto(X,Y).	toto(a,a).", true, "[X = a,  Y = a, ]"}, // query comes first ..

	// Basic, no solution
	{"toto.			?toto(X,Y).", false, "[]"}, // arity does not match !
	{"toto(X).		?toto(X,Y).", false, "[]"}, // arity does not match !

	// Underscore in queries
	{"toto(a,b).	?toto(_,Y).", true, "[Y = b, ]"}, //
	{"toto(a,b).	?toto(X,_).", true, "[X = a, ]"}, //
	{"toto(a,b).	?toto(_,_).", true, "[]"}, //

	// underscore in rules
	{"toto(_).		?toto(X,Y).", false, "[]"}, // arity does not match !
	{"toto(a,_).	?toto(X,Y).", true, "[X = a, ]"}, //

	// Misc.
	{"toto(a,b,c).toto(a,X,X).    ?toto(X,X,c).", false, "[]"}, //

	// Syntax issues
	{"toto(a,b) 	?toto(X,Y).", false, "<nil>"}, // missing period

}

func TestRunStringTable(t *testing.T) {

	for i, d := range testData {
		pc, err := repl.RunString(d.src)
		if (err == nil) != d.ok {
			fmt.Println(pc.StringDetailed())
			fmt.Printf("ON TEST CASE #%d\n%v\n", i, d)
			t.Fatalf("unexpected error outcome : %v", err)
		}
		if d.res != pc.ResultString() {
			fmt.Println(pc.StringDetailed())
			fmt.Printf("ON TEST CASE #%d\nGOT :%v\nWANT:%v\n", i, pc.ResultString(), d.res)
			t.Fatalf("unexpected error outcome")
		}
	}

}

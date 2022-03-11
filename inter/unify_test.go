package inter

import (
	"fmt"
	"testing"
)

func TestUnify(t *testing.T) {

	fmt.Println("\n           ========= TestUnify=========")

	type tty struct {
		g  string
		h  string
		ok bool
	}
	tab := []tty{

		{"a b", "a b", true},
		{"a", "b", false},

		{"X", "a", true},
		{"X", "Y", true},
		{"a(Y,b)", "a(X,b)", true},

		{"a(X)", "a(b)", true},
		{"a(X)", "a(X)", true},
		{"a(X,b)", "a(X,b)", true},
		{"a(b,X)", "a(b,c)", true},
		{"a(b,X)", "a(b,c(f,g))", true},

		{"X Y", "b c", true},
		{"a(b,X)", "a(Y,c)", true},

		{"a(b,X)", "a(X,c)", false},
		{"a(b,X)", "a(X,c)", false},
		{"a(b,X)", "a(c)", false},
		{"a(X)", "X", false},

		{"a(X,X,X)", "a(a,Y,c)", false},
		{"a(X,X,X)", "a(b,Y,b)", true},
		{"a(X,X,X)", "a(Y,g(Y),b)", false},

		{"a(_,X)", "a(b,c(f,g))", true},
		{"a(_,X,c)", "a(X,_,g(X))", false},
		{"a(Y,X,c)", "a(X,_,X)", true}, // <- NOT CLEAN, check algo ?!
		{"a(_,X,c)", "a(X,_,X)", true},

		/*  */
	}

	in := NewInter()
	for i, tt := range tab {
		g, h := in.n(tt.g), in.n(tt.h)
		dotest(t, -i, h, g, tt.ok)
		g, h = in.n(tt.g), in.n(tt.h)
		dotest(t, i, g, h, tt.ok)

	}
}

func dotest(t *testing.T, i int, h *Node, g *Node, ok bool) {
	fmt.Printf("INFO %d: \t%s and \t%s, expect %v :\t", i, h, g, ok)
	ctx, err := unify(NewPContext(), h, g)
	fmt.Println(ctx)
	if (err == nil) != ok {
		//ctx.dump()
		t.Fatalf("Unexpected unification result line %d:\n %s and %s\n%s", i, h, g, err)
	}

	//fmt.Printf("INFO %d: %s and %s unified\n", i, h, g)
	//ctx.dump()

}

func (in *Inter) n(src string) *Node {
	tzr := NewTokenizerString(src)
	n := &Node{
		name: "n",
	}
	err := in.parse0(tzr, n, new(int))
	if err != nil {
		panic("invalid source : " + src)
	}
	return n
}

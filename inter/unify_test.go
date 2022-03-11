package inter

import (
	"fmt"
	"testing"
)

func TestUnify(t *testing.T) {

	type tty struct {
		g  string
		h  string
		ok bool
	}
	tab := []tty{

		{"a(b,X)", "a(Y,c)", true}, // TODO - Fix infinite loop here !!
		{"a(b,X)", "a(b,c)", true},
		{"a(b,X)", "a(b,c(f,g))", true},

		{"a(_,X)", "a(b,c(f,g))", true},

		{"a(b,X)", "a(X,c)", false},
		{"a(b,X)", "a(X,c)", false},
		{"a(b,X)", "a(c)", false},
	}

	in := NewInter()

	for i, tt := range tab {

		g, h := in.n(tt.g), in.n(tt.h)
		ctx, err := unify(NewPContext(), h, g)
		if (err == nil) != tt.ok {
			ctx.dump()
			t.Fatalf("Unexpected unification result line %d:\n %s and %s\n%s", i, h, g, err)
		}
		if err == nil {
			fmt.Printf("INFO %d: %s and %s unified\n", i, h, g)
			ctx.dump()
		}
		// exchange values ...
		ctx, err = unify(NewPContext(), g, h)
		if (err == nil) != tt.ok {
			ctx.dump()
			t.Fatalf("Unexpected unification result line %d:\n %s and %s\n%s", i, g, h, err)
		}
		if err == nil {
			fmt.Printf("INFO %d: %s and %s unified\n", i, g, h)
			ctx.dump()
		}
	}
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

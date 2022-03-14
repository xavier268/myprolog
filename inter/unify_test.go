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
		{"a(Y,X,c)", "a(X,_,X)", true},
		{"a(_,X,c)", "a(X,_,X)", true},

		{"a(_, X,f(Y),c)", "a(X,_,X,Y)", true},

		{"f(X)", "c", false},
		{"a(f(Y),c)", "a(f(f(X)),Y)", false},
		{"a(U,f(X,Y),f(X,Z))", "a(Z,f(_,Z),f(c,X))", true},
		{"a(U,f(X,Y),f(X,Z))", "a(Z,f(_,Z),f(c,X))", true},

		{"a(k(k(U,U),X),f(X,Y),f(X,Z))", "a(Z,f(_,Z),f(c,X))", false},
		{"a(k(k(i,j),X),f(X,Y),f(X,Z))", "a(Z,f(_,Z),f(c,X))", false},

		{"k(k(i,j),_)", "k(Z,Z)", true},
		{"k(k(i,j),X)", "k(Z,Z)", true},

		{"k(U,k(U,666)),U", "Z,555", true},

		/*  test lists in bracket or bar form ! */

		{"[ 1 2 3]", "[ X | Y ]", true},
		{"[ 1 _ 3]", "[ X | Y ]", true},

		{"[ 1 X 3]", "[ Z 4 5 ]", false},
		{"[ 1 X 3]", "[ Z 4 3 ]", true},

		{"[1 ]", "[X|Y]", true},
		{"[1 ]", "[X Y]", false},
		{"[ ]", "X", true},
	}

	in := NewInter()
	for i, tt := range tab {
		g, h := in.n(tt.g), in.n(tt.h)
		dotest(t, in, i, g, h, tt.ok)
		g, h = in.n(tt.g), in.n(tt.h)
		dotest(t, in, -i, h, g, tt.ok)

	}
}

// which test case line should provide detailled output in dotest ?
var detail = map[int]bool{
	// 4: true,
	// -26: true,
	-31: true,
}

func dotest(t *testing.T, in *Inter, i int, h *Node, g *Node, ok bool) {
	fmt.Printf("INFO %d: \t%s and \t%s, expect %v :\t", i, h, g, ok)
	ctx, err := in.unify(NewPContext(), h, g)
	fmt.Println(ctx)
	if (err == nil) != ok {
		//ctx.dump()
		t.Fatalf("Unexpected unification result line %d:\n %s and %s\n%s", i, h, g, err)
	}

	//fmt.Printf("INFO %d: %s and %s unified\n", i, h, g)
	if detail[i] { // select which test line to detail ?
		in.Dump()
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
	err = in.preProcList(n)
	if err != nil {
		panic("invalid source (list syntax) : " + src)
	}
	return n
}

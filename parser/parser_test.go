package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {

	var tdata = []string{

		// rules implicit
		"un(deux,trois).",
		"un.",
		"un(deux).",
		"un(2,3).",
		"empty().",
		"un(deux(),trois).",

		// rules explicit
		"a(X):-un(X).",
		"a(X):-un(X,Y).",
		"a(X):-un(X,Y,Z);b(Y);c(X).",
		"a(X):-un(X,Y,Z),b(Y),c(X).",
		"a(X):-un(X,Y,Z);b(Y),c(X).",
		"a(X):-un(X,Y,Z),b(Y);c(X).",
		"a(X):-.",

		// queries
		"?- 3.", // invalid syntax
		"?- test.",
		"?- un(deux,X).",
		"?- un(deux,X), trois(X).",
		"?- un(deux,X); trois(X).",
		"?- un(deux,X); trois(X),quatre(Y,_).",
		"?- un(deux,X), trois(X);quatre(Y,_).",

		// list
		"[].",
		"[2].",
		"[2,3].",
		"[2,3,4].",

		// sublist
		"[[2,3],4].",
		"[2,[3,4]].",

		// list and pairs
		"[2|3].",
		"[4|].",
		"[4|X].",

		// non list
		"dot(1,dot(2,3)).",                   // not a list
		"dot(1,dot(dot(4,dot(5,dot())),3)).", // not a list, but contains a list

		// underscore
		"un(_,_,X,2).",
		"un(_,_,X,2,3).",
		"un(_,_,X,2):-deux(X,_,5.).",
		"?-un(_,_,X,2,3),deux(X,_,5.).",
		"?- _.",
		"_ .",

		// invalid
		"un,deux.",
		"2(a).",
		" [|2].",
		"[|]",
		"a(b,,).",
		":-.",
	}

	res := run(tdata)

	verify(t, res, "parser_test.wanted")

}

func run(tdata []string) string {
	sb := new(strings.Builder)
	for i, d := range tdata {

		r, err := ParseString(d, fmt.Sprintf("test # %d <%s>", i, d))
		fmt.Fprintf(sb, "\n%d \t\t<%s>\n", i, d)
		fmt.Fprintf(sb, "%d\t\terr=%v\n", i, err)
		for _, v := range r {
			fmt.Fprintf(sb, "%d\t\t(string)    %s\n", i, v.String())
			fmt.Fprintf(sb, "%d\t\t(pretty)    %s\n", i, v.Pretty())
		}
	}
	return sb.String()
}

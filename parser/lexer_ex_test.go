package parser

import (
	"fmt"
	"strings"
)

func Example() {

	lx := newLexerString(`hello world |
		X23 _ , . 
		:- ?- 
		12 0x12 0b111 
		1001e-3 3.14
		"a normal string" '&'  /* a comment */
		// another comment
		555.55  ` + " `a raw \"string`  " +
		" a(b,c(d,e,f,g))")
	vtok := new(mySymType)
	for tk := lx.Lex(vtok); tk != eof; tk = lx.Lex(vtok) {
		if len(vtok.list) > 0 {
			fmt.Printf("the 'list' value should not be set, but was set to %v\n", vtok.list)
		}
		fmt.Printf("token type: %d, ( %q ), ", tk, tk)
		if vtok.value == nil {
			fmt.Printf("lvalue is nil\n")
		} else {
			fmt.Printf("string representation: %s\n", vtok.value.String())
		}

	}
	// Output:
	// token type: 57348, ( '\ue004' ), string representation: hello
	// token type: 57348, ( '\ue004' ), string representation: world
	// token type: 124, ( '|' ), lvalue is nil
	// token type: 57352, ( '\ue008' ), string representation: X23
	// token type: 95, ( '_' ), string representation: _
	// token type: 44, ( ',' ), lvalue is nil
	// token type: 46, ( '.' ), lvalue is nil
	// token type: 57346, ( '\ue002' ), lvalue is nil
	// token type: 57347, ( '\ue003' ), lvalue is nil
	// token type: 57350, ( '\ue006' ), string representation: 12
	// token type: 57350, ( '\ue006' ), string representation: 18
	// token type: 57350, ( '\ue006' ), string representation: 7
	// token type: 57351, ( '\ue007' ), string representation: 1.001e+00
	// token type: 57351, ( '\ue007' ), string representation: 3.140e+00
	// token type: 57349, ( '\ue005' ), string representation: "a normal string"
	// token type: 57349, ( '\ue005' ), string representation: "&"
	// token type: 57351, ( '\ue007' ), string representation: 5.555e+02
	// token type: 57349, ( '\ue005' ), string representation: "a raw \"string"
	// token type: 57348, ( '\ue004' ), string representation: a
	// token type: 40, ( '(' ), lvalue is nil
	// token type: 57348, ( '\ue004' ), string representation: b
	// token type: 44, ( ',' ), lvalue is nil
	// token type: 57348, ( '\ue004' ), string representation: c
	// token type: 40, ( '(' ), lvalue is nil
	// token type: 57348, ( '\ue004' ), string representation: d
	// token type: 44, ( ',' ), lvalue is nil
	// token type: 57348, ( '\ue004' ), string representation: e
	// token type: 44, ( ',' ), lvalue is nil
	// token type: 57348, ( '\ue004' ), string representation: f
	// token type: 44, ( ',' ), lvalue is nil
	// token type: 57348, ( '\ue004' ), string representation: g
	// token type: 41, ( ')' ), lvalue is nil
	// token type: 41, ( ')' ), lvalue is nil

}

// newLexerString from string
func newLexerString(src string) *myLex {
	ps := NewLexer(strings.NewReader(src), "string")
	return ps
}

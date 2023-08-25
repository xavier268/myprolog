package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {

	sb := new(strings.Builder)

	data := `hello world |
	X23 _ , . 
	:- ?- 
	12 0x12 0b111 
	1001e-3 3.14
	"a normal string" '&'  /* a comment */
	// another comment
	555.55  ` + " `a raw \"string`  " +
		" a(b,c(d,e,f,g))"

	fmt.Fprintf(sb, "Lexer test for :\n(quoted)-------\n%q\n(unquoted)--------\n%s\n--------\n", data, data)
	lx := newLexerString(data)
	vtok := new(mySymType)
	for tk := lx.Lex(vtok); tk != eof; tk = lx.Lex(vtok) {
		if len(vtok.list) > 0 {
			fmt.Fprintf(sb, "the 'list' value should not be set, but was set to %v\n", vtok.list)
		}
		fmt.Fprintf(sb, "token type: %d, ( %q ), ", tk, tk)
		if vtok.value == nil {
			fmt.Fprintf(sb, "lvalue is nil\n")
		} else {
			fmt.Fprintf(sb, "string representation: %s\n", vtok.value.String())
		}

	}

	verify(t, sb.String(), "lexer_test.wanted")

}

// newLexerString from string
func newLexerString(src string) *myLex {
	ps := NewLexer(strings.NewReader(src), "string")
	return ps
}

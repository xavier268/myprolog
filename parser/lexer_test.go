package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {

	sb := new(strings.Builder)

	dataset := []string{
		"12e3",
		"1.2e3",
		"13/4",
		"-12",
		"-3/4",
		"3/-4",
		`hello world |
	X23 _ , . 
	:- ?- 
	12 0x12 0b111 
	1001e-3 3.14
	"a normal string" '&'  /* a comment */
	// another comment
	555.55  `,
		" `a raw \"string`  ",

		" a(b,c(d,e,f,g))\n",
		" 12", " 12.3", " 12/4 ", " 12e4", " 12\n",
		" -12", " -12.3", " -12.2 ", "-12/4", " -12e4", " -12\n",
	}

	for _, data := range dataset {
		fmt.Fprintf(sb, "\n\nLexer test for :\n(quoted)-------\n%q\n(unquoted)--------\n%s\n--------\n", data, data)
		lx := newLexerString(data)
		vtok := new(mySymType)
		for tk := lx.Lex(vtok); tk != eof; tk = lx.Lex(vtok) {
			if len(vtok.list) > 0 {
				fmt.Fprintf(sb, "the 'list' value should not be set, but was set to %v\n", vtok.list)
			}
			fmt.Fprintf(sb, "token type: %d, ( %q )\n", tk, tk)
			if vtok.value == nil {
				fmt.Fprintf(sb, "\tlvalue is nil\n")
			} else {
				fmt.Fprintf(sb, "\tstring representation: %s\n", vtok.value.String())
				fmt.Fprintf(sb, "\tpretty representation: %s\n", vtok.value.Pretty())
			}
		}
		for i, e := range lx.LastErr {
			fmt.Fprintf(sb, "\nLexer errors %d\t: %s\n", i, e)
		}
	}

	verifyTest(t, sb.String(), "lexer_test.wanted")

}

// newLexerString from string
func newLexerString(src string) *myLex {
	ps := NewLexer(strings.NewReader(src), "string")
	return ps
}

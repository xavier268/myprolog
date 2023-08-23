package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"
)

// Lexer for prolog

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type myLex struct {
	s scanner.Scanner // golang scanner
}

// Required to satisfy interface.
func (lx *myLex) Error(s string) {
	fmt.Printf("LexerError in %s, line %d : %v\n", lx.s.Filename, lx.s.Line, s)
}

// The myLexer interface is implemented by myLex, defined by generated parser.
var _ MyLexer = new(myLex)

// NewLexer from io.Reader
func NewLexer(input io.Reader) *myLex {
	ps := new(myLex)
	ps.s.Mode = scanner.GoTokens
	ps.s.Init(input)
	ps.s.Filename = "io.Reader"
	return ps
}

// NewLexerString from string
func NewLexerString(src string) *myLex {
	ps := NewLexer(strings.NewReader(src))
	ps.s.Filename = "string"
	return ps
}

// NewLexerFile from file name. File is left open.
func NewLexerFile(fileName string) (*myLex, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	ps := NewLexer(bufio.NewReader(f))
	ps.s.Filename = fileName
	return ps, nil
}

// The parser calls this method to get each new token.
// This implementation returns token extracted from scanner
func (lx *myLex) Lex(lval *MySymType) int {

	lval.list = []Term{} // lexer never return a list
	lval.value = nil     // reset the value to nil

	tk := lx.s.Scan()
	switch tk {
	case scanner.EOF:
		return eof // 0 for the lexer
	case scanner.Ident:
		txt := lx.s.TokenText()
		if txt == "_" {
			lval.value = &Underscore{}
			return '_'
		}
		if txt[0] >= 'A' && txt[0] <= 'Z' {
			lval.value = &Variable{
				Name: txt,
				Nsp:  0,
			}
			return VARIABLE
		}
		if len(txt) > 1 && txt[0] == '_' {
			lx.Error(fmt.Sprintf("Identifier cannot begin with underscore : %v", txt))
			return eof
		}
		lval.value = &Atom{
			Value: txt,
		}
		return ATOM

	case scanner.Comment: // ignore comments
		return lx.Lex(lval)

	case scanner.RawString, scanner.String:
		lval.value = &String{
			Value: lx.s.TokenText()[1 : len(lx.s.TokenText())-1],
		}
		return STRING

	case scanner.Char: // handled as a STRING with 1 char.
		lval.value = &String{
			Value: lx.s.TokenText()[1 : len(lx.s.TokenText())-1],
		}
		return STRING

	case scanner.Float:
		var z float64
		fmt.Sscanf(lx.s.TokenText(), "%v", &z)
		lval.value = &Float{
			Value: z,
		}
		return FLOAT

	case scanner.Int:
		var z int
		fmt.Sscanf(lx.s.TokenText(), "%v", &z)
		lval.value = &Integer{
			Value: z,
		}
		return INTEGER

	case '(', ')', '[', ']', ',', ';', '.': // single char tokens recognized by parser, cannot begin a multichar operaotor.
		// yylval is not set for these.
		return int(tk)

	case ':': // possibly multichar token.
		if lx.s.Peek() == '-' {
			lx.s.Scan()
			return OPRULE
		}
		// a single ':' is recognized as an atom.
		lval.value = &Atom{
			Value: ":",
		}
		return ATOM

	case '?': // possibly a query operator.
		if lx.s.Peek() == '-' {
			lx.s.Scan()
			return OPQUERY
		}
		// a single '?' is recognized as an atom.
		lval.value = &Atom{
			Value: "?",
		}
		return ATOM

	default:
		lx.Error(fmt.Sprintf("Unknown token type %v :  %v", tk, lx.s.TokenText()))
		return eof
	}

}

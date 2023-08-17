package myyacc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"

	"github.com/xavier268/myprolog/config"
)

// Lexer for prolog

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type myLex struct {
	s scanner.Scanner // scanner
}

// Required to statisfy interface.
func (*myLex) Error(s string) {
	fmt.Println("DEBUG : Error : ", s)
}

// The myLexer interface is implemented by myLex.
var _ myLexer = new(myLex)

// newLexer from io.Reader
func newLexer(c *config.Config, input io.Reader) *myLex {
	ps := new(myLex)
	ps.s.Mode = c.ScannerMode
	ps.s.Init(input)
	ps.s.Filename = "<io.Reader>"
	return ps
}

// newLexerString from string
func newLexerString(c *config.Config, src string) *myLex {
	ps := newLexer(c, strings.NewReader(src))
	ps.s.Filename = "<string>"
	return ps
}

// newLexerFile from file name. File is left open.
func newLexerFile(c *config.Config, fileName string) (*myLex, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	ps := newLexer(c, bufio.NewReader(f))
	ps.s.Filename = fileName
	return ps, nil
}

// The parser calls this method to get each new token. This
// implementation returns token extracted from scanner
func (ps *myLex) Lex(yylval *mySymType) int {

	tk := ps.s.Scan()
	switch tk {
	case scanner.EOF:
		return eof // 0 for the lexer
		// TODO : shouldn't we exploit further the type of token as analysed by the lib ?
	}
	yylval.name = ps.s.TokenText()
	n := yylval.name[0]
	switch n {
	case ')', '(', '.', ',', ';':
		return int(n)
	case '"':
		return STRING
	default:
		return ATOM
	}
}

package inter

import (
	"bufio"
	"io"
	"os"
	"strings"
	"text/scanner"
)

// Tokenizer scanns the token from an input stream.
// Token follow the go syntax, including comments.
// Combined tokens, such as :- or <= appear as two separate tokens.
type Tokenizer interface {
	Next() string // Get next token, empty string on EOF
}

type pscanner struct {
	s scanner.Scanner
}

func (ps *pscanner) Next() string {
	tk := ps.s.Scan()
	if tk == scanner.EOF {
		return ""
	}
	n := ps.s.TokenText()

	// handle double tokens
	switch n {
	case ":": // :-
		if ps.s.Peek() == '-' {
			ps.s.Scan()
			return n + ps.s.TokenText()
		} else {
			return n
		}

	case "<", ">", "!": // <=, != >=
		if ps.s.Peek() == '=' {
			ps.s.Scan()
			return n + ps.s.TokenText()
		} else {
			return n
		}

	default:
		return n
	}

}

// NewTokenizer from io.Reader
func NewTokenizer(input io.Reader) Tokenizer {
	ps := new(pscanner)
	//ps.s.Mode = scanner.GoTokens
	ps.s.Init(input)
	ps.s.Filename = "io.Reader"
	return ps
}

// NewTokenizerString from srource string
func NewTokenizerString(src string) Tokenizer {
	return NewTokenizer(strings.NewReader(src))
}

// NewTokenizerFile from file name
func NewTokenizerFile(fileName string) Tokenizer {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	return NewTokenizer(bufio.NewReader(f))
}

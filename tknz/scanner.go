// Packahe tknz contains the scanner and utilities.
package tknz

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"
)

// Tokenizer scans the token from an input stream.
type Tokenizer struct {
	s scanner.Scanner
}

// Next provides the next token.
// Token follow the go syntax, including comments.
// Combined tokens  ( :- , <= , >= , != ) are reconstructed as a single token.
// Upon EOF, an empty string is returned.
func (ps *Tokenizer) Next() string {
	tk := ps.s.Scan()
	if tk == scanner.EOF {
		return ""
	}
	n := ps.s.TokenText()

	// handle quoted strings.
	if len(n) >= 2 && n[0] == '"' {
		fmt.Sscanf(n, "%q", &n)
	}

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
func NewTokenizer(input io.Reader) *Tokenizer {
	ps := new(Tokenizer)
	//ps.s.Mode = scanner.GoTokens
	ps.s.Init(input)
	ps.s.Filename = "io.Reader"
	return ps
}

// NewTokenizerString from string
func NewTokenizerString(src string) *Tokenizer {
	return NewTokenizer(strings.NewReader(src))
}

// NewTokenizerFile from file name
func NewTokenizerFile(fileName string) (*Tokenizer, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return NewTokenizer(bufio.NewReader(f)), nil
}

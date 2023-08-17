// Package tknz contains the scanner and utilities.
package tknz

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"

	"github.com/xavier268/myprolog/config"
)

// Tokenizer scans the token from an input stream.
type Tokenizer struct {
	s       scanner.Scanner // scanner
	pending []string        // pending tokens to send later
}

// Next provides the next token.
// Token follow the go syntax, including comments.
// Combined tokens  ( :- , <= , >= , != , ?- ) are reconstructed as a single token.
// quotes are not removed from quted strings.
// Upon EOF, an empty string is returned.
func (ps *Tokenizer) Next() string {

	tk := ps.s.Scan()
	if tk == scanner.EOF {
		return ""
	}
	n := ps.s.TokenText()

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

	case "?": // ?-
		if ps.s.Peek() == '-' {
			ps.s.Scan()
			return n + ps.s.TokenText()
		} else {
			return n
		}

	case ".":
		return n

	default:
		if strings.Contains(n, ".") {
			fmt.Println("dot not allowed in token : " + n)
			return strings.Split(n, ".")[0] // drop the decimal part, keep an int
		}
		return n
	}

}

// NewTokenizer from io.Reader
// Preprocessing is performed using the preprocessor transformer.
func NewTokenizer(c *config.Config, input io.Reader) *Tokenizer {
	ps := new(Tokenizer)
	ps.s.Mode = c.ScannerMode
	ps.s.Init(input)
	ps.s.Filename = "<io.Reader>"
	return ps
}

// NewTokenizerString from string
func NewTokenizerString(c *config.Config, src string) *Tokenizer {
	ps := NewTokenizer(c, strings.NewReader(src))
	ps.s.Filename = "<string>"
	return ps
}

// NewTokenizerFile from file name. File is left open.
func NewTokenizerFile(c *config.Config, fileName string) (*Tokenizer, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	ps := NewTokenizer(c, bufio.NewReader(f))
	ps.s.Filename = fileName
	return ps, nil
}

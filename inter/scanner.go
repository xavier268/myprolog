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
	return ps.s.TokenText()
}

// NewTokenizer from io.Reader
func NewTokenizer(input io.Reader) Tokenizer {
	ps := new(pscanner)
	ps.s.Init(input)
	ps.s.Filename = "io.Reader"
	return ps
}

// NewTokenizerString from srource string
func NewTokenizerString(src string) Tokenizer {
	ps := new(pscanner)
	ps.s.Init(strings.NewReader(src))
	ps.s.Filename = "string"
	return ps
}

// NewTokenizerFile from file name
func NewTokenizerFile(fileName string) Tokenizer {
	ps := new(pscanner)
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	ps.s.Init(bufio.NewReader(f))
	ps.s.Filename = fileName
	return ps
}

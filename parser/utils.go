package parser

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func Parse(rdr io.Reader, sourcename string) ([]Term, error) {
	var err error
	lx := NewLexer(rdr, sourcename)
	p := myNewParser()
	e := p.Parse(lx)
	if e != 0 {
		err = fmt.Errorf("Parse error")
	}
	r := append([]Term{}, lastParseResult...) // copy slice before returning to avoid mutating the original slice
	lastParseResult = nil
	return r, err
}

func ParseFile(filename string) ([]Term, error) {
	var err error
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f, filename)
}

func ParseString(input string, sourcename string) ([]Term, error) {
	return Parse(strings.NewReader(input), sourcename)
}

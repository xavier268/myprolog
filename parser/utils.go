package parser

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	RED    = "\x1b[31m"
	GREEN  = "\x1b[32m"
	BLUE   = "\x1b[34m"
	YELLOW = "\x1b[33m"
	CYAN   = "\x1b[36m"

	RESET = "\x1b[0m"
)

func Parse(rdr io.Reader, sourcename string) ([]Term, error) {
	var err error
	lx := NewLexer(rdr, sourcename)
	p := myNewParser()
	e := p.Parse(lx)
	if e != 0 {
		err = fmt.Errorf("error : %v", lx.LastErr)
	}
	r := append([]Term{}, lx.LastResult...) // copy slice before returning to avoid mutating the original slice
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

func MustParseString(input string, sourcename string) []Term {
	r, err := ParseString(input, sourcename)
	if err != nil {
		panic(err)
	}
	return r
}

// The error, if any, is added to the (partially) parsed list
// Useful mainly for debugging.
func noFailParseString(input string, sourcename string) []Term {
	r, err := ParseString(input, sourcename)
	if err != nil {
		r = append(r, CompoundTerm{
			Functor:  "error",
			Children: []Term{String{Value: fmt.Sprintf("error : %v", err)}},
		})
	}
	return r
}

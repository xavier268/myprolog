package myyacc

import (
	"fmt"
	"testing"

	"github.com/xavier268/myprolog/config"
)

func TestParser(t *testing.T) {

	data := " un , deux ; troix , quatre ."
	lx := newLexerString(config.New(), data)

	pp := myNewParser()
	res := pp.Parse(lx)
	fmt.Println("Parse results ( 0 = no error ): ", res)
}

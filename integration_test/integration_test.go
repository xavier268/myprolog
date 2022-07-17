package integration

import (
	"testing"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/repl"
)

func TestRunFile(*testing.T) {
	config.FlagDebug = false
	config.FlagVerbose = false
	repl.RunFile("test.pl")
}

func TestRunString(*testing.T) {
	config.FlagDebug = false
	config.FlagVerbose = false
	repl.RunString(`?print("hello world\n").`)
}

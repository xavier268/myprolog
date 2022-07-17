package integration

import (
	"testing"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/repl"
)

func Example_runFile1() {
	config.FlagDebug = false
	config.FlagVerbose = false
	repl.RunFile("test1.pl")
	// Output:
	// Queries : [ query ( queries )  query ( goals )  query ( queries )  query ( goals )  query ( print ( hello word ) )  query ( queries )  goals]
	// Queries : [ query ( queries )  query ( goals )  query ( queries )  query ( goals )  query ( print ( hello word ) )  queries] hello word
	// Queries : [ query ( queries )  query ( goals )  query ( queries )  goals]
	// Queries : [ query ( queries )  query ( goals )  queries]
	// Queries : [ query ( queries )  goals]
	// Queries : [ queries]
}

func Example_runString1() {
	config.FlagDebug = false
	config.FlagVerbose = config.FlagDebug
	repl.RunString(`?print("hello\nworld").`)
	// Output:
	// hello
	// world
}

func Test_runString2(*testing.T) {
	config.FlagDebug = false
	config.FlagVerbose = config.FlagDebug
	repl.RunString(`toto(a,b).?toto(a,X).`)

}

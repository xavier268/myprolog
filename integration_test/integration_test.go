package integration

import (
	"fmt"
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

func Example_runString2() {
	config.FlagDebug = false
	config.FlagVerbose = config.FlagDebug
	pc, err := repl.RunString(`toto(a,b).?toto(a,X).`)
	fmt.Println(err)
	fmt.Println(pc)
	// Output:
	// <nil>
	// Constrt : [X = b, ]
	// Goals   : []
	// Root    : false
}

func Example_runString3() {
	config.FlagDebug = false
	config.FlagVerbose = config.FlagDebug
	pc, err := repl.RunString(`toto(a,b).?toto(a,X).`)
	fmt.Println(err)
	fmt.Println(pc)
	// Output:
	// <nil>
	// Constrt : [X = b, ]
	// Goals   : []
	// Root    : false
}

func Example_runString4() {
	config.FlagDebug = false
	config.FlagVerbose = config.FlagDebug
	pc, err := repl.RunString(`?toto(a,X).`)
	fmt.Println(err)
	fmt.Println(pc)
	// output:
	// no solution, unknown keyword
	// Constrt : []
	// Goals   : [ toto ( a X )]
	// Root    : true
}

func Example_runString5() {
	config.FlagDebug = false
	config.FlagVerbose = config.FlagDebug
	pc, err := repl.RunString(`?toto(a,X).toto(b,c).`)
	fmt.Println(err)
	fmt.Println(pc)
	// output:
	// no solution, unknown keyword
	// Constrt : []
	// Goals   : [ toto ( a X )]
	// Root    : true
}
func Example_runString6() {
	config.FlagDebug = false
	config.FlagVerbose = config.FlagDebug
	pc, err := repl.RunString(`?toto(X,X).toto(b,c).`)
	fmt.Println(err)
	fmt.Println(pc)
	// output:
	// no solution, unknown keyword
	// Constrt : []
	// Goals   : [ toto ( X X )]
	// Root    : true
}
func Example_runString7() {
	config.FlagDebug = false
	config.FlagVerbose = config.FlagDebug
	pc, err := repl.RunString(`?toto(X,Y).toto(b,c).`)
	fmt.Println(err)
	fmt.Println(pc)
	// output:
	// <nil>
	// Constrt : [X = b,  Y = c, ]
	// Goals   : []
	// Root    : false
}
func TestTT(t *testing.T) {
	//func Example_runString5() {
	config.FlagDebug = true
	config.FlagVerbose = config.FlagDebug
	pc, err := repl.RunString(`?toto(X,Y).toto(b,c).`)
	fmt.Println(err)
	fmt.Println(pc)
	// output:
	// no solution, unknown keyword
	// Constrt : []
	// Goals   : [ toto ( a X )]
	// Root    : true
}

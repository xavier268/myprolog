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
	// hello world from test1.pl
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
	// no solution
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
	// no solution
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
	// no solution
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

func Example_runString8() {
	config.FlagDebug = false
	config.FlagVerbose = config.FlagDebug
	pc, err := repl.RunString(`?println(X yyy " " zzz ).`)
	fmt.Println(err)
	fmt.Println(pc)
	// output:
	// X yyy   zzz
	// <nil>
	// Constrt : []
	// Goals   : []
	// Root    : true
}

func TestHelpVisual(_ *testing.T) {
	config.FlagDebug = false
	config.FlagVerbose = config.FlagDebug
	repl.RunString(`?help().`)
}

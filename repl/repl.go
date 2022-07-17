package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/node"
	"github.com/xavier268/myprolog/pcontext"
	"github.com/xavier268/myprolog/prsr"
	"github.com/xavier268/myprolog/tknz"
)

// REPL launch the intercative "Request-Execute-Print-Loop" main loop.
func REPL() {

	fmt.Printf("%s\nVersion : %s - %s\n", config.WELCOME, config.VERSION, config.COPYRIGHT)

	pg := node.NewProgramNode()
	pc := pcontext.NewPContext(pg)

	fmt.Println("\nOk.")
	for {

		reader := bufio.NewReader(os.Stdin)
		src, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		if config.FlagDebug {
			fmt.Println("DEBUG REPL : before parsing", pc)
		}
		prog, err := prsr.Parse(tknz.NewTokenizerString(src))
		if config.FlagDebug {
			fmt.Println("DEBUG REPL : new inputs parsed ", prog)
		}
		pc.AddGoals(prog)
		pc.AddRules(prog)
		if config.FlagDebug {
			fmt.Println("DEBUG REPL : after parsing", pc)
		}

		if err != nil {
			fmt.Println(err)
			continue
		}
		pc, err = pc.Run()
		if err != nil {
			fmt.Println(err)
			continue
		}
		pc.Display() // Display results
		fmt.Println("\nOk.")
	}

}

// RUN a non intercative program with rules and goals.
func RUN(filename string) error {

	tk, err := tknz.NewTokenizerFile(filename)
	if err != nil {
		return err
	}
	prog, err := prsr.Parse(tk)
	if err != nil {
		return err
	}

	pc := pcontext.NewPContext(prog)
	_, err = pc.Run()
	return err

}

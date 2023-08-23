// Package repl contains the high-level entry points to run programs.
package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/node"
	"github.com/xavier268/myprolog/pcontext"
	"github.com/xavier268/myprolog/prsr"
	"github.com/xavier268/myprolog/tknz"
)

// REPL launch the intercative "Request-Execute-Print-Loop" main loop.
func REPL() {

	config.PrintFullWelcome()

	pg := node.NewProgramNode()
	pc := pcontext.NewPContext(pg)
	fmt.Println("Type a blank line to stop input and launch execution.")
	fmt.Print("\nOk>")
	for {
		// Capture src until a blank line is input.
		var src string
		for {
			reader := bufio.NewReader(os.Stdin)
			s, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				break
			}
			src = src + s
			if len(strings.TrimSpace(s)) == 0 {
				break
			}

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
			fmt.Println("Error : ", err)
			continue
		}
		fmt.Printf("\nResults : %s\nOk> ", pc.ResultString())
	}

}

// RunFile a non interactive program with rules and goals from file.
func RunFile(filename string) (*pcontext.PContext, error) {

	tk, err := tknz.NewTokenizerFile(filename)
	if err != nil {
		return nil, err
	}
	prog, err := prsr.Parse(tk)
	if err != nil {
		return nil, err
	}

	pc := pcontext.NewPContext(prog)
	return pc.Run()

}

// RunString a non interactive program with rules and goals from string.
func RunString(s string) (*pcontext.PContext, error) {

	tk := tknz.NewTokenizerString(s)
	prog, err := prsr.Parse(tk)
	if err != nil {
		return nil, err
	}

	pc := pcontext.NewPContext(prog)
	return pc.Run()

}

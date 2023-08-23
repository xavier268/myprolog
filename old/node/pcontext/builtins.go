package pcontext

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/node"
	"github.com/xavier268/myprolog/prsr"
	"github.com/xavier268/myprolog/tknz"
)

// DoBuiltin execute root of the goal node if it is a builtin keyword.
// Zero, one or more modified Node are returned.
// No-op if the root of the node is not a keyword.
func (pc *PContext) DoBuiltin(goal *node.Node) ([]*node.Node, error) {

	if pc == nil || goal == nil {
		return nil, nil
	}

	switch kw := goal.GetLoad().(type) {

	case node.Keyword: // Handle actual keywords

		// First, check arity
		if node.Keywords[kw.String()].Arity() >= 0 && goal.NbChildren() != node.Keywords[kw.String()].Arity() {
			return nil, fmt.Errorf("unknown keyword with such arity")
		}

		switch kw.String() {

		case "query": // unpack the query content
			if config.FlagDebug {
				fmt.Println("\nDEBUG preparing to execute : ", goal)
			}
			return goal.GetChildren(), nil

		case "print": // print all the children nodes and remove node
			pc.doPrint(goal.GetChildren()...)
			return nil, nil
		case "println": // print all the children nodes and remove node
			pc.doPrintln(goal.GetChildren()...)
			return nil, nil
		case "halt": // pause program
			pc.doHalt()
			return nil, nil

		case "exit": // exit program
			pc.doExit()
			return nil, nil

		case "help": // print help
			pc.doHelp(goal.GetChildren()...)
			return nil, nil

		case "verbose": // toggle the verbose switch
			pc.doVerbose(goal.GetChildren()...)
			return nil, nil

		case "debug": // toggle the debug switch
			pc.doVerbose(goal.GetChildren()...)
			pc.doDebug(goal.GetChildren()...)
			return nil, nil

		case "list": // list current rule
			pc.doList()
			return nil, nil

		case "load": // load more rules from an external file.
			// For the moment : ( TODO)
			// only use the first child, ignore the rest.
			// Use the name of the load, do not indirect variables.
			file := fmt.Sprint(goal.GetChild(0).GetLoad())
			tk, err := tknz.NewTokenizerFile(file)
			if err != nil {
				abs, _ := filepath.Abs(file)
				fmt.Printf("\nWARNING (cannot load %s) : %v", abs, err)
				return nil, nil
			}
			rr, err := prsr.Parse(tk)
			if err != nil {
				fmt.Printf("\nWARNING (cannot load %s) : %v\n", file, err)
				return nil, nil
			}
			pc.AddRules(rr)
			return nil, nil
			/*
				case "rules": // dump rules
					fmt.Printf("\nRules : %v", pc.rules)
					return nil, nil

				case "queries", "goals": // dump queries, also called goals
					fmt.Printf("\nQueries : %v", pc.goals)
					return nil, nil
			*/

		default: // keyword has no effect, ignore but keep - see unknown !
			return []*node.Node{goal}, nil
		}

	default: // not a keyword, remove, signal an error since rules have already been tried
		return nil, fmt.Errorf("no solution")
	}
}
func (pc *PContext) doDebug(children ...*node.Node) {
	v, ok := children[0].GetLoad().(node.Number)
	if ok {
		if v.GetValue() == 0 {
			config.FlagDebug = false
		} else {
			config.FlagDebug = true
		}
	} // ignore non number.
	fmt.Println("Debug flag is now ", config.FlagDebug)
}
func (pc *PContext) doVerbose(children ...*node.Node) {
	v, ok := children[0].GetLoad().(node.Number)
	if ok {
		if v.GetValue() == 0 {
			config.FlagVerbose = false
		} else {
			config.FlagVerbose = true
		}
	} // ignore non number.
	fmt.Println("Verbose flag is now ", config.FlagVerbose)
}
func (pc *PContext) doExit() {
	fmt.Println("\nProgram stopped")
	os.Exit(0)
}

func (pc *PContext) doHalt() {
	fmt.Println("\nProgram paused. Type ^C to exit, or <Enter> to continue")
	fmt.Scanln(new(string))
	fmt.Println("\nProgram resumed.")
}

func (pc *PContext) doPrint(children ...*node.Node) {
	for _, c := range children {
		fmt.Print(c.String())
	}
}
func (pc *PContext) doPrintln(children ...*node.Node) {
	pc.doPrint(children...)
	fmt.Println()
}

func (pc *PContext) doList() {
	fmt.Printf("List of current rules :\n")
	for i, r := range pc.rules {
		fmt.Printf("%3d)\t%s\n", i, r.String())
	}
}

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

// DoBuiltin execute root of the goal node if it is a keyword.
// Zero, one or more modified Node are returned.
// No-op if the root of the node is not a keyword.
func (pc *PContext) DoBuiltin(goal *node.Node) ([]*node.Node, error) {

	if pc == nil || goal == nil {
		return nil, nil
	}

	switch kw := goal.GetLoad().(type) {

	case node.Keyword: // Handle actual keywords

		switch kw.String() {
		case "print": // print all the children nodes and remove node
			for _, c := range goal.GetChildren() {
				fmt.Print(pc.StringContent(c))
			}
			return nil, nil
		case "halt": // pause program
			fmt.Println("\nProgram paused. Type ^C to exit, or <Enter> to continue")
			fmt.Scanln(new(string))
			fmt.Println("\nProgram resumed.")
			return nil, nil

		case "exit": // exit program
			fmt.Println("\nProgram stopped")
			os.Exit(0)
			return nil, nil

		case "query": // unpack the query content
			if config.FlagDebug {
				fmt.Println("\nDEBUG preparing to execute : ", goal)
			}
			return goal.GetChildren(), nil

		case "verbose": // toggle the verbose switch
			config.FlagVerbose = !config.FlagVerbose
			fmt.Println("\nVerbose flag is now", config.FlagVerbose)
			return nil, nil

		case "debug": // toggle the debug switch
			config.FlagDebug = !config.FlagDebug
			fmt.Println("\nDebug flag is now", config.FlagDebug)
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

		case "rules": // dump rules
			fmt.Printf("\nRules : %v", pc.rules)
			return nil, nil

		case "queries", "goals": // dump queries, also called goals
			fmt.Printf("\nQueries : %v", pc.goals)
			return nil, nil

		default: // keyword has no effect, ignore but keep
			return []*node.Node{goal}, nil
		}

	default: // not a keyword, ignore but keep
		return []*node.Node{goal}, nil
	}
}

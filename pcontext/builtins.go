package pcontext

import (
	"fmt"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/node"
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
		case "halt", "exit": // remove goal and exit program
			fmt.Println("\nProgram paused. Type ^C to exit, or <Enter> to continue")
			fmt.Scanln(new(string))
			fmt.Println("Program resumed.")
			return nil, nil

		case "query": // unpack the query content
			return goal.GetChildren(), nil

		case "trace": // toggle the verbose switch
			config.FlagVerbose = !config.FlagVerbose
			fmt.Println("Verbose flag is now", config.FlagVerbose)
			return nil, nil

		case "debug": // toggle the debug switch
			config.FlagDebug = !config.FlagDebug
			fmt.Println("Debug flag is now", config.FlagDebug)
			return nil, nil

		default: // keyword has no effect, ignore but keep
			return []*node.Node{goal}, nil
		}

	default: // not a keyword, ignore but keep
		return []*node.Node{goal}, nil
	}
}

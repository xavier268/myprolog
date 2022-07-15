// Package inter contains the top level interpreter andthe repl loop.
package inter

import (
	"github.com/xavier268/myprolog/node"
	"github.com/xavier268/myprolog/pcontext"
)

// Execute will attempt to erase the goal, using the rules defined in the program.
func Execute(prog *node.Node, goal *node.Node) (*pcontext.PContext, error) {

	ctx := pcontext.NewPContext(prog, goal)
	for ctx != nil {
		err := execute(ctx)
		if err == nil {
			ctx.DisplaySolutions()
		} else {
			// Backtrack
			ctx = ctx.Parent()
		}
	}

}

// execute is the main execution loop on the context.
func execute(ctx *pcontext.PContext) error {
	panic("todo")
}

package pcontext

import (
	"fmt"
	"testing"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/node"
)

var (
	// Variables
	X, Y, Z = node.NewVariable("X", 0), node.NewVariable("Y", 0), node.NewVariable("Z", 0)
	X1      = node.NewVariable("X", 1)

	// Nodes

	na, nb, nc, nd  = node.NewStringNode("a"), node.NewStringNode("b"), node.NewStringNode("c"), node.NewStringNode("d")
	nX, nX1, nY, nZ = node.NewVariableNode(X), node.NewVariableNode(X1), node.NewVariableNode(Y), node.NewVariableNode(Z)

	n1 = na.Clone().Add(nX, nX1)
	n2 = nb.Clone().Add(n1.Clone(), nY.Clone(), nX.Clone())
	n3 = na.Clone().Add(nc, nd, nZ)
)

var data = []Cons{
	{X, ConsEQ, na},
	{Y, ConsEQ, nX},
	{X, ConsEQ, nY},
	{Z, ConsEQ, n3},
	{X1, ConsEQ, n2},
}

func TestTableSimplifyVisual(t *testing.T) {

	config.FlagDebug = false
	for i := 0; i <= len(data); i++ {

		test := data[:i]
		fmt.Print(test, " ---> ")
		fmt.Println(Simplify(test))
	}

}

package pcontext

import (
	"fmt"
	"testing"

	"github.com/xavier268/myprolog/node"
)

func TestConEqual(t *testing.T) {

	cXY := NewConEqual(node.NewVariable("X", 0), node.NewVariableNode(node.NewVariable("Y", 0)))
	// cYX := NewConEqual(node.NewVariable("Y", 0), node.NewVariableNode(node.NewVariable("X", 0)))
	// cXX := NewConEqual(node.NewVariable("X", 0), node.NewVariableNode(node.NewVariable("X", 0)))
	cXa := NewConEqual(node.NewVariable("X", 0), node.NewStringNode("a"))

	pc := NewPContext(nil, nil)
	pc.SetConstraint(cXY)
	fmt.Println(pc.ResultString())
	pc.SetConstraint(cXa)
	fmt.Println(pc.ResultString())

}

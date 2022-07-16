package pcontext

import (
	"fmt"

	"github.com/xavier268/myprolog/node"
)

// Equality between a Variable and a (sub)tree
type ConEqual struct {
	v node.Variable
	t *node.Node
}

var _ Constraint = NewConEqual(node.Variable{}, nil)

func NewConEqual(v node.Variable, t *node.Node) ConEqual {
	c := ConEqual{v, t}
	return c
}

// checks obvious X=X
func (c ConEqual) Obvious() bool {
	return c.t.GetLoad() == c.v || c.t.GetLoad() == node.Underscore{}
}

// Verify checks for positive occurs error.
func (c ConEqual) Verify() error {
	return c.t.Walk(
		func(n *node.Node) error {
			if n != nil && n.GetLoad() == c.v {
				return fmt.Errorf("positive occurs check")
			}
			return nil
		})
}

func (c ConEqual) Merge(old Constraint) (remove bool, newcon []Constraint, err error) {

	if c.Obvious() {
		return false, nil, nil
	}
	if err := c.Verify(); err != nil {
		return false, nil, err
	}

	// TODO replace X in all previous occurence in the rhs ?

	// default ..
	return false, []Constraint{c}, nil

}

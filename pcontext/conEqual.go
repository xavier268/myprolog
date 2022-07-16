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

func (ConEqual) Cons() {}

func (c ConEqual) String() string {
	return fmt.Sprintf("%s = %s", c.v.String(), c.t.String())
}

func NewConEqual(v node.Variable, t *node.Node) ConEqual {
	c := ConEqual{v, t}
	return c
}

// checks obvious X=X
func (c ConEqual) IsObvious() bool {
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

// Update the cc constraints that contain a reference the Variable in c in their rhs, by replacing this Variable by the c rhs.
// Return nil if no update required.
func (c ConEqual) Update(cc Constraint) (upcc Constraint) {

	if cc == nil {
		return nil
	}

	switch v := cc.(type) {
	case ConEqual:

		vnode := node.NewVariableNode(v.v)
		tree, changed := v.t.ReplaceSubTree(vnode, v.t)
		if changed {
			return ConEqual{
				v: v.v,
				t: tree,
			}
		} else {
			return nil
		}

	default:
		panic("type of constraint not implemented")
	}

}

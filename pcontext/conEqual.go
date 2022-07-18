package pcontext

import (
	"fmt"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/node"
)

// Equality between a Variable and a (sub)tree
type ConEqual struct {
	v node.Variable
	t *node.Node
}

func (ConEqual) Cons() {}

func (c ConEqual) String() string {
	return fmt.Sprintf("%s =%s, ", c.v.String(), c.t.String())
}

func NewConEqual(v node.Variable, t *node.Node) ConEqual {
	if config.FlagDebug {
		fmt.Println("DEBUG : CREATING CONSTRAINT EQUAL ", v, t)
	}
	c := ConEqual{v, t}
	return c
}

// checks obvious X=X
func (c ConEqual) IsObvious() bool {
	return c.t.GetLoad() == c.v || c.t.GetLoad() == node.Underscore{}
}

// VerifyPosOcc checks for positive occur error X = f(a,X,...)
func (c ConEqual) VerifyPosOcc() error {
	return c.t.Walk(
		func(n *node.Node) error {
			if n != nil && n.GetLoad() == c.v {
				return fmt.Errorf("positive occurs check")
			}
			return nil
		})
}

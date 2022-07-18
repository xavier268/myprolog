package cons

import (
	"fmt"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/node"
)

// Generic constraint structure on a Variable.
type Cons struct {
	variable node.Variable // the variable upon which the constraint aplies
	relation RelType       // relation
	tree     *node.Node    // a tree defining the constraint
	// possibly other fields for further constraints ...
}

// Type of constraint relation
type RelType int

const (
	ConsEQ RelType = iota
)

func (c *Cons) String() string {
	if c == nil {
		return fmt.Sprint(nil)
	}
	switch c.relation {
	case ConsEQ:
		return fmt.Sprintf("%s =%s, ", c.variable.String(), c.tree.String())
	default:
		panic("Unimplemented constraint")
	}

}

func NewConEqual(v node.Variable, t *node.Node) Cons {
	if config.FlagDebug {
		fmt.Println("DEBUG : CREATING CONSTRAINT EQUAL ", v, t)
	}
	return Cons{v, ConsEQ, t}
}

// Checks the validity of the constraint.
// Invaid constraints implies no solution and will trigger backtracking.
func (c *Cons) IsValid() bool {
	if c == nil {
		return true
	}
	return !c.isPosOccurErr()
}

// Checks the relevance of the constraint.
// Irrevelant constraint can be safely ignored.
// (avoid X = X and X = _  ; _ = xx cannot happen since _ is not a Variable)
func (c *Cons) IsRelevant() bool {
	if c == nil {
		return false
	}
	switch c.relation {
	case ConsEQ:
		return c.tree.GetLoad() != c.variable && c.tree.GetLoad() != node.Underscore{}
	default:
		panic("unimplemented constraint type")
	}
}

// isPosOccurErr checks if constraints has a positive occur error, ie X = f(a,X,...)
func (c Cons) isPosOccurErr() bool {
	switch c.relation {
	case ConsEQ:
		return c.tree.Walk(
			func(n *node.Node) error {
				if n != nil && n.GetLoad() == c.variable {
					return fmt.Errorf("positive occurs check")
				}
				return nil
			}) != nil
	default:
		panic("unimplemented constraint type")
	}
}

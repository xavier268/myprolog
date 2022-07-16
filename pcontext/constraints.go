package pcontext

import "github.com/xavier268/myprolog/node"

// A Constraints is the main information we have about the potential solutions.
// PContext will continuously try to combine the constraints together to extract the best possible description of the solution.
// Once created, constraints are immutable.
type Constraint interface {
	// Merge current Constraint with old previous constraint, returning new constraints, possibly replacing old.
	// Remove boolean indicates if the previous constraint needs to be removed.
	// Error is returned if an impossibility is detected, and backtracking is required.
	Merge(old Constraint) (remove bool, newconst []Constraint, err error)
}

type ConsEqual struct {
	v node.Variable
	n *node.Node
}

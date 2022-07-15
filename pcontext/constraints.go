package pcontext

// A Constraints is the main information we have about the potential solutions.
// PContext will continuously try to combine the constraints together to extract the best possible description of the solution.
type Constraint interface {
	// Merge the constraint with another one
	// Error is returned if an impossibility is detected, and backtracking will be required.
	Merge(cc *Constraint) error
}

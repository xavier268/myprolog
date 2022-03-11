package inter

import "fmt"

// ErrSolve is set when backtracking is required.
var ErrSolve = fmt.Errorf("backtracking is required")

// Solve will attempt to solve the provided goal, using the rules previously loaded.
// Solve should have no side effect on Inter and be indempotent,
// excpet for the symbol table, necessary to garantee unicity of variable nodes.
// Solve returns no context on error.
func (in *Inter) Solve(ctx *PContext, goal *Node) (*PContext, error) {

	if goal == nil {
		return ctx, nil
	}

	for _, r := range in.rules.args { // iterate on each known original rule ...

		// compare functor and arity of goal vs rule HEAD, to abort rapidly
		if goal.name != r.args[0].name || len(goal.args) != len(r.args[0].name) {
			continue
		}

		// Lets rescope the rule head.
		// Use a different suffix per rule, since different rules can name their
		// variables using the same symbol.
		suffix := in.Uid()
		head := in.Rescope(r.args[0], suffix)

		// try to unify head with goal, keeping the new context with the new constraints.
		newctx, err := unify(ctx, goal, head)
		if err == nil {
			// unification of head suceeded.
			// rescope rule body and attempt to solve each one,
			// generating new contexts.
			for i := 1; err == nil && i < len(r.args); i++ {
				b := in.Rescope(r.args[i], suffix) // rescoped body, same suffix as head !
				newctx, err = in.Solve(newctx, b)
			}
			if err == nil {
				//solving succeeded with rule r !
				return newctx, err
			} else {
				// try another rule ...
				newctx = nil // helps garbage collector ...
				continue     // iterate on next rules
			}
		}
	}

	// all rules failed, report error.
	return nil, ErrSolve
}

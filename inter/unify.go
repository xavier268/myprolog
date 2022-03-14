package inter

import "fmt"

// unify attempts to unify a goal g with a rescoped head h.
// Upon success a new pcontext is returned.
// Upon error, NO context is returned, to facilitate garbage collection.
func (in *Inter) unify(ctx *PContext, g *Node, h *Node) (*PContext, error) {
	//ctx.dump()
	//fmt.Printf("DEBUG : entering unify : %s\t and \t%s\n", g, h)

	switch {

	// --- exclude nil -----
	case h == nil || g == nil:
		//fmt.Printf("DEBUG :  trying to unify %s and %s\n", g.String(), h.String())
		panic("cannot unify a nil goal or a nil head")

	// --- handle _  on either side ----
	case g.name == "_" || h.name == "_": // _ = ... or ... = _
		return ctx, nil

	// --- handle G = G ---- # RULE 2
	case isVariable(g.name) && isVariable(h.name) && (g.name == h.name):
		return ctx, nil // ignore

	// --- Handle G = H ---
	case isVariable(g.name) && isVariable(h.name) && g.name != h.name:
		// TODO - do not overwite previous X ! <<<<<<<<<<<<<<<<<<<<<<<????????????TOODOOOO !
		// test if G=u already exists ?
		u := ctx.Get(EQ, g)
		if u != nil {
			return in.unify(ctx, h, u)
		}
		// test if H=u exists ?
		u = ctx.Get(EQ, h)
		if u != nil {
			return in.unify(ctx, g, u)
		}
		// neither H nor G were already set ...
		return ctx.Set(EQ, g, h), nil

	// --- handle G = h -----
	case isVariable(g.name) && !isVariable(h.name):

		if h.contains(g) { // # RULE 7
			fmt.Println("\nWARNING : Positive error check")
			return ctx, ErrPosOcc // fail, positive occur check (ie : a looping tree would be needed ...)
		}

		// Test if G already in LHS of previous equations.
		// test if G=u exists already ? (and u does not contains G, by construction).
		u := ctx.Get(EQ, g) // TODO : generalize to different operators ...
		if u != nil {
			// X = h and X = u both exist.
			ctxu, err := in.unify(ctx, u, h) // check if u = h ?
			if err != nil {
				return ctx, ErrSolve // fail to unify u & h !
			} else {
				// ctxu can unify u & h
				return ctxu, nil
			}
		}

		// Test if G already in RHS of previous equations ?
		// If it is, we need to substitute G with h everywhere.
		lhs, rhs := ctx.rhsGet(EQ, g)
		if rhs == nil && lhs == nil { // no prior use of g.
			return ctx.Set(EQ, g, h), nil // store G = h and continue
		}

		// subsitute h in all G occurences.
		rhs = in.Rescope(rhs, "") // clone rhs, before modifying it !
		if rhs == g {             // we had LHS = G as an equation already.
			rhs = h
			nctx := ctx.Set(EQ, lhs, rhs) // set LHS = h, since G = h
			return in.unify(nctx, g, h)   // unify again with new context
		}
		// here, we have LHS = f(G)
		rhs.substitute0(g, h)
		nctx := ctx.Set(EQ, lhs, rhs) // set LHS = f(h), since G = h
		return in.unify(nctx, g, h)   // unify again with new context

	case !isVariable(g.name) && isVariable(h.name): // g = H
		return in.unify(ctx, h, g)

	// ---- handle a == b , same arity = 0 --- # RULE 1.1
	case !isVariable(g.name) && !isVariable(h.name) && len(g.args) == 0 && len(h.args) == 0:
		if g.name == h.name {
			return ctx, nil // ignore
		} else {
			return ctx, ErrSolve // fail
		}

	// ---- handle g(a,b, ...) = h(a,b,...) --- # RULE 4, same arity
	case !isVariable(g.name) && !isVariable(h.name) && len(g.args) == len(h.args) && g.name == h.name:
		var err error
		c := ctx
		for i := range g.args {
			c, err = in.unify(c, g.args[i], h.args[i])
			if err != nil {
				return ctx, ErrSolve
			}
		}
		return c, nil

		// ---- handle g(a,b, ...) = h(a,b,...) --- # RULE 1.4, different arity
	case !isVariable(g.name) && !isVariable(h.name) && len(g.args) != len(h.args):
		return ctx, ErrSolve // fail

	default:
		fmt.Printf("DEBUG : default fallback, cannot unify %s and %s\n", g.String(), h.String())
		return ctx, ErrSolve
	}
}

// does node n contains node X ?
func (n *Node) contains(X *Node) bool {

	if n == X {
		return true
	}

	for _, a := range n.args {
		if a == X {
			return true
		}
	}
	return false
}

// Subsitute all occurence of x with y in n.args, recursively.
// if n ==x, it cannot be changed. Only its args will change.
// n is modified. No checks, no assumptions.
// Very dangerous !! You may want to rescope or clone n before.
func (n *Node) substitute0(x, y *Node) {
	if x == y {
		return
	}
	for i, a := range n.args {
		if a == x {
			n.args[i] = y
		} else {
			a.substitute0(x, y)
		}
	}
}

/*
// substituteVariables all variables in n that are known to ctx.
// n itself should NOT be a variable.
func (n *Node) substituteVariables(ctx *PContext) {
	if ctx == nil {
		return
	}
	for i, a := range n.args {
		if isVariable(a.name) {
			s := ctx.Get(EQ, a)
			if s != nil {
				n.args[i] = s
			}
		} else {
			a.substituteVariables(ctx)
		}
	}
}

*/

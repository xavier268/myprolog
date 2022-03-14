package inter

import "fmt"

// unify attempts to unify a goal g with a rescoped head h.
// Upon success a new pcontext is returned.
// Upon error, NO context is returned, to facilitate garbage collection.
func (in *Inter) unify(ctx *PContext, g *Node, h *Node) (*PContext, error) {
	//ctx.dump()
	//fmt.Printf("DEBUG : entering unify : %s\t and \t%s\n", g, h)

	// --- exclude nil -----
	if h == nil || g == nil {
		//fmt.Printf("DEBUG :  trying to unify %s and %s\n", g.String(), h.String())
		panic("cannot unify a nil goal or a nil head")
	}

	switch {

	// --- handle variables ----
	// TODO - verify handling of _ ?
	case g.name == "_" || h.name == "_": // _ = ... or ... = _
		return ctx, nil

	case isVariable(g.name) && isVariable(h.name): // G = H
		if g.name == h.name { // G = G
			return ctx, nil // ignore
		} else { // G = H
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
		}

	case isVariable(g.name) && !isVariable(h.name): // G == h
		if h.contains(g) {
			fmt.Println("WARNING : Positive error check")
			return ctx, ErrPosOcc // fail, positive occur check (ie : a looping tree would be needed ...)
		}
		// if a variable already known appears in X=t, then substitute it (including the X in lhs !)

		// if G=u exists already
		u := ctx.Get(EQ, g)
		if u != nil {
			// X=u exists
			ctxu, err := in.unify(ctx, u, h) // check if u = h ?
			if err != nil {
				return ctx, ErrSolve // fail to unify u & h !
			} else {
				// ctxu could unify u & h
				return ctxu, nil
			}
		}
		if h.isConstant() {
			//h being constant, it is now safe to Set G=h
			return ctx.Set(EQ, g, h), nil
		}
		// Before adding G = h, we need to substitute all former equations where the value G exists.
		h.substituteVariables(ctx)
		return ctx.Set(EQ, g, h), nil

	case !isVariable(g.name) && isVariable(h.name): // g = H
		return in.unify(ctx, h, g)

	// ---- handle no variables ----
	// same arity = 0
	case len(g.args) == 0 && len(h.args) == 0:
		if g.name == h.name {
			return ctx, nil // ignore
		} else {
			return ctx, ErrSolve // fail
		}

	// arity != 0, same functor, same arity, for non variables
	case len(g.args) == len(h.args) && g.name == h.name:
		var err error
		c := ctx
		for i := range g.args {
			c, err = in.unify(c, g.args[i], h.args[i])
			if err != nil {
				return ctx, ErrSolve
			}
		}
		return c, nil

	// different arity, but not variables
	case len(g.args) != len(h.args):
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

/*
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
*/

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

package pcontext

import (
	"fmt"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/cons"
	"github.com/xavier268/myprolog/node"
)

var errNoMatch = fmt.Errorf("no match")

// Unify attempts unification of the the goal (gh) with rule head (rh), updating context constraints.
// It is called recusively, trying to unify children with each other.
// As unification proceeds, the context contraints are updated (but not the goals).
func (pc *PContext) Unify(gh, rh *node.Node) error {

	if config.FlagDebug {
		fmt.Println("DEBUG UNIFYING : ", gh, rh)
	}

	if pc == nil {
		panic("cannot unify with a nil context")
	}
	if rh == nil && gh == nil {
		return nil
	}

	goalLoad := gh.GetLoad()
	ruleLoad := rh.GetLoad()

	switch goalLoad.(type) {

	case node.Number:
		if err := pc.unifyNumber(gh, rh); err != nil {
			return err
		}

	case node.String:
		// check heads
		if err := pc.unifyString(gh, rh); err != nil {
			return err
		}

		// recursion could be required on children ?
		switch ruleLoad.(type) {
		case node.String: // string/string requires recursion
			if ruleLoad != goalLoad || rh.NbChildren() != gh.NbChildren() {
				return errNoMatch
			}
			for i := range rh.GetChildren() {
				if err := pc.Unify(gh.GetChild(i), rh.GetChild(i)); err != nil {
					return err
				}
			}
		default: // no child checks for string vs anything else
		}
	case node.Underscore:
		// always match

	case node.Variable:
		if err := pc.unifyVariable(gh, rh); err != nil {
			return err
		}

	case node.Keyword:
		return fmt.Errorf("keywords never match unification - should be handled elsewhere")

	default:
		panic("internal error")
	}

	return nil
}

// The first node is a number
func (pc *PContext) unifyNumber(gh, rh *node.Node) error {

	switch rh.GetLoad().(type) {
	case node.Number:
		if rh.GetLoad() != gh.GetLoad() {
			return errNoMatch
		}
	case node.Variable:
		// unify in reverse order
		if err := pc.unifyVariable(rh, gh); err != nil {
			return err
		}
	case node.Underscore:
		// ignore
	case node.Keyword, node.String:
		return errNoMatch
	default:
		panic("internal error")
	}

	return nil
}

// unify a goal string with ...
func (pc *PContext) unifyString(gh, rh *node.Node) error {
	switch rh.GetLoad().(type) {
	case node.Number, node.Keyword:
		return errNoMatch
	case node.Variable:
		// reverse order
		if err := pc.unifyVariable(rh, gh); err != nil {
			return err
		}
	case node.Underscore:
		// ignore
	case node.String:
		if rh.GetLoad() != gh.GetLoad() {
			return errNoMatch
		}
		if rh.NbChildren() != gh.NbChildren() {
			return errNoMatch
		}
		for i := range gh.GetChildren() {
			if err := pc.Unify(gh.GetChild(i), rh.GetChild(i)); err != nil {
				return errNoMatch
			}
		}
	default:
		panic("internal error")
	}
	return nil
}

// SetConstraint tries to update context.
// If error is not nil, then a contradiction appeared.
// If no error, the context has been already simplified.
func (pc *PContext) SetConstraint(c cons.Cons) error {
	if !c.IsRelevant() {
		return nil // ok, but ignore
	}
	if !c.IsValid() {
		return errNoMatch // internal contradiction
	}
	newlist, changed, err := cons.Simplify(pc.cstr)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	pc.cstr = newlist
	return nil
}

// unify a node variable in the first position.
func (pc *PContext) unifyVariable(gh, rh *node.Node) error {

	switch rh.GetLoad().(type) {
	case node.Number, node.Variable, node.String:
		c := cons.NewConEqual(gh.GetLoad().(node.Variable), rh)
		err := pc.SetConstraint(c)
		if err != nil {
			return err // cannot unify !
		}
		return nil

	case node.Underscore:
		// ignore

	case node.Keyword:
		return fmt.Errorf("keywords cannot be unified")

	default:
		panic("internal error")
	}

	return nil
}

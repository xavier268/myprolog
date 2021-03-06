package pcontext

import (
	"fmt"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/node"
)

var errNoMatch = fmt.Errorf("no match")

// Unify attempts unification of the rule head (rh) with the goal (gh), updating context constraints.
// It is called recusively, trying to unify children with each other.
// As unification proceeds, the context contraints are updated (but not the goals).
func (pc *PContext) Unify(rh, gh *node.Node) error {

	if config.FlagDebug {
		fmt.Println("DEBUG UNIFYING : ", rh, gh)
	}

	if pc == nil {
		panic("cannot unify with a nil context")
	}
	if rh == nil && gh == nil {
		return nil
	}

	v1 := rh.GetLoad()
	v2 := gh.GetLoad()

	switch v1.(type) {

	case node.Number:
		if err := pc.unifyNumber(rh, gh); err != nil {
			return err
		}

	case node.String:
		// check heads
		if err := pc.unifyString(rh, gh); err != nil {
			return err
		}

		// recursion could be required on children ?
		switch v2.(type) {
		case node.String: // string/string requires recursion
			if v1 != v2 || rh.NbChildren() != gh.NbChildren() {
				return errNoMatch
			}
			for i := range rh.GetChildren() {
				if err := pc.Unify(rh.GetChild(i), gh.GetChild(i)); err != nil {
					return err
				}
			}
		default: // no child checks for string vs anything else
		}
	case node.Underscore:
		// always match

	case node.Variable:
		if err := pc.unifyVariable(rh, gh); err != nil {
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
func (pc *PContext) unifyNumber(rh, gh *node.Node) error {

	switch gh.GetLoad().(type) {
	case node.Number:
		if rh.GetLoad() != gh.GetLoad() {
			return errNoMatch
		}
	case node.Variable:
		if err := pc.unifyVariable(gh, rh); err != nil {
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

// unify a string with ...
func (pc *PContext) unifyString(rh, gh *node.Node) error {
	switch gh.GetLoad().(type) {
	case node.Number, node.Keyword:
		return errNoMatch
	case node.Variable:
		if err := pc.unifyVariable(gh, rh); err != nil {
			return err
		}
	case node.Underscore:
		// ignore
	case node.String:
		if rh.GetLoad() != gh.GetLoad() {
			return errNoMatch
		}
	default:
		panic("internal error")
	}
	return nil
}

// unify a node variable in the first position.
func (pc *PContext) unifyVariable(rh, gh *node.Node) error {

	switch gh.GetLoad().(type) {
	case node.Number, node.Variable, node.String:
		c := NewConEqual(rh.GetLoad().(node.Variable), gh)
		err := pc.SetConstraint(c)
		if err != nil {
			return err
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

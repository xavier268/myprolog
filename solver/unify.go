package solver

import "fmt"

var ErrUnificationImpossible = fmt.Errorf("unification impossible")

// Unify recursively unifies rule head and goal.
// We avoid state directly, by just appending to a list of constraints.
// Simplification will occur later.
// No predicates are handled here.
// No backtracking is done.
func Unify(cList []Constraint, head Term, goal Term) ([]Constraint, error) {
	var err error

	// special nil cases - do nothing if both are nil
	if head == nil && goal == nil {
		return cList, nil
	}
	// cannot unify nil with non nil
	if head == nil || goal == nil {
		return cList, ErrUnificationImpossible
	}

	// handle other non nil & non underscore cases
	switch head := head.(type) {

	case Number: // integer head

		switch goal := goal.(type) {
		case Underscore:
			return cList, nil
		case Number: // numbers can  unify with themselves
			if head.Eq(goal) {
				return cList, nil
			} else {
				return cList, ErrUnificationImpossible
			}
		case String, Atom: //  do not unify
			return cList, ErrUnificationImpossible
		case Variable: // goal is a variable
			c := VarEQ{
				V:     goal,
				Value: head,
			}
			return CheckAddConstraint(cList, c)
		case CompoundTerm:
			return cList, ErrUnificationImpossible
		default:
			panic("unreacheable code reached")
		}

	case String: // string head

		switch goal := goal.(type) {
		case Underscore:
			return cList, nil
		case String:
			if head.Value == goal.Value {
				return cList, nil
			} else {
				return cList, ErrUnificationImpossible
			}
		case Number, Atom: //  do not unify
			return cList, ErrUnificationImpossible
		case Variable: // goal is a variable
			c := VarIsString{
				V: goal,
				S: head,
			}
			return CheckAddConstraint(cList, c)
		case CompoundTerm:
			return cList, ErrUnificationImpossible
		default:
			panic("unreacheable code reached")
		}

	case Variable: // variable head

		switch goal := goal.(type) {
		case Underscore:
			return cList, nil
		case Variable:
			c := VarIsVar{
				V: goal,
				W: head,
			}
			return CheckAddConstraint(cList, c)
		case CompoundTerm:
			c := VarIsCompoundTerm{
				V: head, // head is the variable
				T: goal, // goal not a variable anymore
			}
			return CheckAddConstraint(cList, c)
		case Number:
			c := VarEQ{
				V:     head,
				Value: goal,
			}
			return CheckAddConstraint(cList, c)
		case String:
			c := VarIsString{
				V: head, // head is the variable
				S: goal,
			}
			return CheckAddConstraint(cList, c)
		case Atom:
			c := VarIsAtom{
				V: head, // head is the variable
				A: goal,
			}
			return CheckAddConstraint(cList, c)
		default:
			panic("unreacheable code reached")
		}

	case Underscore: // do nothing
		return cList, nil

	case Atom:
		switch goal := goal.(type) {
		case Underscore:
			return cList, nil
		case Variable:
			c := VarIsAtom{
				V: goal, // prefer goal when it is the variable, switch order
				A: head,
			}
			return CheckAddConstraint(cList, c)
		case CompoundTerm, String, Number:
			return cList, ErrUnificationImpossible
		case Atom:
			if head.Value == goal.Value {
				return cList, nil
			} else {
				return cList, ErrUnificationImpossible
			}
		default:
			panic("unreacheable code reached")
		}

	case CompoundTerm: // compound head
		switch goal := goal.(type) {
		case String, Number, Atom:
			return cList, ErrUnificationImpossible
		case Underscore:
			return cList, nil
		case Variable:
			c := VarIsCompoundTerm{
				V: goal,
				T: head,
			}
			return CheckAddConstraint(cList, c)
		case CompoundTerm:
			if goal.Functor != head.Functor || len(goal.Children) != len(head.Children) {
				return cList, ErrUnificationImpossible
			}
			for i, h := range head.Children {
				if cList, err = Unify(cList, h, goal.Children[i]); err != nil {
					return cList, ErrUnificationImpossible
				}
			}
			return cList, nil
		default:
			panic("unreacheable code reached")
		}
	default:
		panic("unreacheable code reached")
	}
}

// Check (individual), and Add a constraint to the list, return error for backtracking.
// Simplification not performed.
// Avoid later workload by notadding nil constraints.
func CheckAddConstraint(cc []Constraint, c Constraint) ([]Constraint, error) {
	clean, err := c.Check()
	if err == nil && clean != nil {
		return append(cc, clean), nil
	} else {
		return cc, err
	}
}

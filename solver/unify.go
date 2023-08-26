package solver

import "fmt"

var ErrUnificationImpossible = fmt.Errorf("Unification impossible")

// Unify recursively unifies rule head and goal.
// State is modified during unification, because constraints and/or goals are added as we unify.
// No backtracking is done.
func Unify(st *State, head Term, goal Term) (newstate *State, err error) {

	// special nil cases - do nothing if both are nil
	if head == nil && goal == nil {
		return st, nil
	}
	// cannot unify nil with non nil
	if head == nil || goal == nil {
		return st, ErrUnificationImpossible
	}

	// handle other non nil & non underscore cases
	switch head := head.(type) {

	case Number: // integer head

		switch goal := goal.(type) {
		case Underscore:
			return st, nil
		case Number: // numbers can  unify with themselves
			if head.Eq(goal) {
				return st, nil
			} else {
				return st, ErrUnificationImpossible
			}
		case String, Atom: //  do not unify
			return st, ErrUnificationImpossible
		case Variable: // goal is a variable
			c := VarIsNumber{
				V:           goal,
				Min:         head,
				Max:         head,
				IntegerOnly: false,
			}
			err := st.AddConstraint(c)
			return st, err
		case CompoundTerm:
			return st, ErrUnificationImpossible
		default:
			panic("unreacheable code reached")
		}

	case String: // string head

		switch goal := goal.(type) {
		case Underscore:
			return st, nil
		case String:
			if head.Value == goal.Value {
				return st, nil
			} else {
				return st, ErrUnificationImpossible
			}
		case Number, Atom: //  do not unify
			return st, ErrUnificationImpossible
		case Variable: // goal is a variable
			c := VarIsString{
				V: goal,
				S: head.Value,
			}
			err := st.AddConstraint(c)
			return st, err

		case CompoundTerm:
			return st, ErrUnificationImpossible
		default:
			panic("unreacheable code reached")
		}

	case Variable: // variable head

		switch goal := goal.(type) {
		case Underscore:
			return st, nil
		case Variable:
			c := VarIsVar{
				V: goal,
				W: head,
			}
			err := st.AddConstraint(c)
			return st, err
		case CompoundTerm:
			c := VarIsCompoundTerm{
				V: head, // head is the variable
				T: goal, // goal not a variable anymore
			}
			err := st.AddConstraint(c)

			return st, err
		case Number:
			c := VarIsNumber{
				V:           head,
				Min:         goal,
				Max:         goal,
				IntegerOnly: false,
			}
			err := st.AddConstraint(c)
			return st, err
		case String:
			c := VarIsString{
				V: head, // head is the variable
				S: goal.Value,
			}
			err := st.AddConstraint(c)
			return st, err
		case Atom:
			c := VarIsAtom{
				V: head, // head is the variable
				A: goal,
			}
			err := st.AddConstraint(c)
			return st, err
		default:
			panic("unreacheable code reached")
		}

	case Underscore: // do nothing
		return st, nil

	case Atom:
		switch goal := goal.(type) {
		case Underscore:
			return st, nil
		case Variable:
			c := VarIsAtom{
				V: goal, // prefer goal when it is the variable, switch order
				A: head,
			}
			err := st.AddConstraint(c)
			return st, err
		case CompoundTerm, String, Number:
			return st, ErrUnificationImpossible
		case Atom:
			if head.Value == goal.Value {
				return st, nil
			} else {
				return st, ErrUnificationImpossible
			}
		default:
			panic("unreacheable code reached")
		}

	case CompoundTerm: // compound head
		switch goal := goal.(type) {
		case String, Number, Atom:
			return st, ErrUnificationImpossible
		case Underscore:
			return st, nil
		case Variable:
			c := VarIsCompoundTerm{
				V: goal,
				T: head,
			}
			err := st.AddConstraint(c)
			return st, err
		case CompoundTerm:
			if goal.Functor != head.Functor || len(goal.Children) != len(head.Children) {
				return st, ErrUnificationImpossible
			}
			for i, h := range head.Children {
				if st, err = Unify(st, h, goal.Children[i]); err != nil {
					return st, ErrUnificationImpossible
				}
			}
			return st, nil
		default:
			panic("unreacheable code reached")
		}
	default:
		panic("unreacheable code reached")
	}
}

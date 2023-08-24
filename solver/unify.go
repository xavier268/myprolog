package solver

// Unify recursively unifies rule head and goal.
// State is modified during unification, because constraints are added while we unify.
// No backtracking is done.
func Unify(st *State, head Term, goal Term) (newstate *State, ok bool) {

	// special nil cases
	if head == nil && goal == nil {
		return st, true
	}
	if head == nil || goal == nil {
		return st, false
	}

	// handle other non nil & non underscore cases
	switch head := head.(type) {

	case *Integer: // integer head

		switch goal := goal.(type) {
		case *Underscore:
			return st, true
		case *Integer: // integers unify with themselves
			return st, head.Value == goal.Value
		case *Float, *String, *Char, *Atom: //  do not unify
			return st, false
		case *Variable: // goal is a variable
			c := VarEqCons{
				V: Clone(goal).(*Variable),
				T: head,
			}
			err := st.AddConstraint(c)
			if err != nil {
				return st, false
			} else {
				return st, true
			}
		case *CompoundTerm:
			return st, false
		default:
			panic("unreacheable code reached")
		}

	case *Float: // float head

		switch goal := goal.(type) {
		case *Underscore:
			return st, true
		case *Float:
			return st, head.Value == goal.Value
		case *Integer, *String, *Char, *Atom: //  do not unify
			return st, false
		case *Variable: // goal is a variable
			c := VarEqCons{
				V: goal,
				T: head,
			}
			err := st.AddConstraint(c)
			if err != nil {
				return st, false
			} else {
				return st, true
			}
		case *CompoundTerm:
			return st, false
		default:
			panic("unreacheable code reached")
		}

	case *String: // string head

		switch goal := goal.(type) {
		case *Underscore:
			return st, true
		case *String:
			return st, head.Value == goal.Value
		case *Integer, *Float, *Char, *Atom: //  do not unify
			return st, false
		case *Variable: // goal is a variable
			c := VarEqCons{
				V: goal,
				T: head,
			}
			err := st.AddConstraint(c)
			if err != nil {
				return st, false
			} else {
				return st, true
			}
		case *CompoundTerm:
			return st, false
		default:
			panic("unreacheable code reached")
		}

	case *Char: // char head

		switch goal := goal.(type) {
		case *Underscore:
			return st, true
		case *Char:
			return st, head.Char == goal.Char
		case *Integer, *Float, *String, *Atom: //  do not unify
			return st, false
		case *Variable: // goal is a variable
			c := VarEqCons{
				V: goal,
				T: head,
			}
			err := st.AddConstraint(c)
			if err != nil {
				return st, false
			} else {
				return st, true
			}
		case *CompoundTerm:
			return st, false
		default:
			panic("unreacheable code reached")
		}

	case *Variable: // variable head

		switch goal := goal.(type) {
		case *Underscore:
			return st, true
		case *Variable:
			c := VarEqCons{
				V: goal, // prefer goal when it is the variable, switch order
				T: head,
			}
			err := st.AddConstraint(c)
			if err != nil {
				return st, false
			} else {
				return st, true
			}

		case *CompoundTerm, *Atom, *String, *Char, *Integer, *Float:
			c := VarEqCons{
				V: head, // head is the variable
				T: goal, // goal not a variable anymore
			}
			err := st.AddConstraint(c)
			if err != nil {
				return st, false
			} else {
				return st, true
			}
		default:
			panic("unreacheable code reached")
		}

	case *Underscore:
		return st, true

	case *Atom:

		switch goal := goal.(type) {
		case *Underscore:
			return st, true
		case *Variable:
			c := VarEqCons{
				V: goal, // prefer goal when it is the variable, switch order
				T: head,
			}
			err := st.AddConstraint(c)
			if err != nil {
				return st, false
			} else {
				return st, true
			}

		case *CompoundTerm, *String, *Char, *Integer, *Float:
			return st, false

		case *Atom:
			return st, head.Value == goal.Value

		default:
			panic("unreacheable code reached")
		}

	case *CompoundTerm: // compound head

		switch goal := goal.(type) {
		case *Char, *String, *Float, *Integer, *Atom:
			return st, false
		case *Underscore:
			return st, true
		case *Variable:
			c := VarEqCons{
				V: goal, // prefer goal when it is the variable, switch order
				T: head,
			}
			err := st.AddConstraint(c)
			if err != nil {
				return st, false
			} else {
				return st, true
			}
		case *CompoundTerm:
			if goal.Functor != head.Functor || len(goal.Children) != len(head.Children) {
				return st, false
			}
			for i, h := range head.Children {
				if st, ok = Unify(st, h, goal.Children[i]); !ok {
					return st, false
				}
			}
			return st, true
		default:
			panic("unreacheable code reached")
		}
	default:
		panic("unreacheable code reached")
	}
}

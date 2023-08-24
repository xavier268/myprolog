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
		case *String, *Char, *Atom: //  do not unify
			return st, false
		case *Float: // integer can unif with float
			return st, float64(head.Value) == goal.Value
		case *Variable: // goal is a variable
			c := VarIsInteger{
				V:   goal,
				Min: head.Value,
				Max: head.Value,
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
		case *String, *Char, *Atom: //  do not unify
			return st, false
		case *Integer: // Float and Integer can unify
			return st, head.Value == float64(goal.Value)
		case *Variable: // goal is a variable
			c := VarIsFloat{
				V:   goal,
				Min: head.Value,
				Max: head.Value,
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
			c := VarIsString{
				V: goal,
				S: head.Value,
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
			c := VarIsChar{
				V: goal,
				C: head.Char,
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
			c := VarIsVar{
				V: goal,
				W: head,
			}
			err := st.AddConstraint(c)
			if err != nil {
				return st, false
			} else {
				return st, true
			}

		case *CompoundTerm, *Atom, *String, *Char, *Integer, *Float:
			c := VarIsCompoundTerm{
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
			c := VarIsAtom{
				V: goal, // prefer goal when it is the variable, switch order
				A: *head,
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
			c := VarIsCompoundTerm{
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

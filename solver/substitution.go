package solver

// Find if Variable exists in Term
func FindVar(v Variable, t Term) (found bool) {
	switch t := t.(type) {
	case Variable:
		return t.Name == v.Name && t.Nsp == v.Nsp
	case Number, String, Atom, Underscore:
		return false
	case CompoundTerm:
		for _, c := range t.Children {
			if FindVar(v, c) {
				return true
			}
		}
		return false
	default:
		panic("code should have been unreacheable")
	}
}

// Replace Variable v in Term with w.
// If replacement occurs, Term is cloned.
func ReplaceVar(oldVar Variable, t Term, newVar Variable) (tt Term, found bool) {

	found = FindVar(oldVar, t)
	if !found {
		return t, false
	}

	tt = t.Clone()
	replaceVarInPlace(oldVar, tt, newVar)
	return tt, true
}

// replace in place. No cloning.
func replaceVarInPlace(oldVar Variable, t Term, newVar Variable) {

	switch t := t.(type) {
	case Variable:
		if t.Name == oldVar.Name && t.Nsp == oldVar.Nsp {
			t.Name = newVar.Name
			t.Nsp = newVar.Nsp
			return
		}
		return

	case Number, String, Atom, Underscore:
		return
	case CompoundTerm:
		for _, c := range t.Children {
			replaceVarInPlace(oldVar, c, newVar)
		}
		return
	default:
		panic("code should have been unreacheable")
	}
}

package solver

// Find if Variable exists in Term
func FindVar(v Variable, t Term) (found bool) {
	switch t := t.(type) {
	case *Variable:
		return t.Name == v.Name && t.Nsp == v.Nsp
	case *Float, *Integer, *String, *Char, *Atom, *Underscore:
		return false
	case *CompoundTerm:
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
func ReplaceVar(v Variable, t Term, w Variable) (tt Term, found bool) {

	found = FindVar(v, t)

	if !found {
		return t, false
	}

	tt = Clone(t)
	replaceVarInPlace(v, tt, w)
	return tt, true
}

// replace in place. No cloning.
func replaceVarInPlace(v Variable, t Term, w Variable) {

	switch t := t.(type) {
	case *Variable:
		if t.Name == v.Name && t.Nsp == v.Nsp {
			t.Name = w.Name
			t.Nsp = w.Nsp
			return
		}
		return

	case *Float, *Integer, *String, *Char, *Atom, *Underscore:
		return
	case *CompoundTerm:
		for _, c := range t.Children {
			replaceVarInPlace(v, c, w)
		}
		return
	default:
		panic("code should have been unreacheable")
	}
}

package solver

// Find if Variable exists in Term
func FindVar(v Variable, t Term) (found bool) {
	switch t := t.(type) {
	case Variable:
		return v.Eq(t)
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

// Get the set of all vars in Term
func FindVars(t Term) map[Variable]bool {
	switch t := t.(type) {
	case Variable:
		return map[Variable]bool{t: true}
	case Number, String, Atom, Underscore:
		return nil
	case CompoundTerm:
		res := make(map[Variable]bool)
		for _, c := range t.Children {
			for v := range FindVars(c) {
				res[v] = true
			}
		}
		return res
	default:
		panic("code should have been unreacheable")
	}
}

// Replace Variable v with the term w in  t.
// If replacement occurs, Term is cloned and found flag is set.
func ReplaceVar(v Variable, t Term, w Term) (res Term, found bool) {

	switch tt := t.(type) {
	case Variable:
		if v.Eq(tt) { // match !
			return w, true
		}
		return t, false
	case Number, String, Atom, Underscore:
		return tt, false
	case CompoundTerm:
		var f bool
		res := CompoundTerm{
			Functor:  tt.Functor,
			Children: make([]Term, len(tt.Children)),
		}
		for i, c := range tt.Children {
			res.Children[i], f = ReplaceVar(v, c, w)
			found = found || f
		}
		if found {
			return res, true
		} else {
			return t, false
		}
	default:
		panic("code unreacheable")
	}
}

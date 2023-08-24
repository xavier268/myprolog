package solver

// make a deep clone of Term
// OBSOLETE - use Clone interface instead.
func Clone(t Term) Term {
	switch t := t.(type) {

	case *Variable:
		return &Variable{
			Name: t.Name,
			Nsp:  t.Nsp,
		}

	case *Atom:
		return &Atom{
			Value: t.Value,
		}

	case *Underscore,
		*String,
		*Char,
		*Float,
		*Integer:
		return t

	case *CompoundTerm:
		tt := &CompoundTerm{
			Functor:  t.Functor,
			Children: []Term{},
		}
		for _, c := range t.Children {
			tt.Children = append(tt.Children, Clone(c))
		}
		return tt

	default:
		panic("attempting to clone an unknown type")
	}
}

// clone and change the name space of all variables
func CloneNsp(t Term, nsp int) Term {
	switch t := t.(type) {

	case *Variable:
		return &Variable{
			Name: t.Name,
			Nsp:  nsp,
		}

	case *Atom:
		return &Atom{
			Value: t.Value,
		}

	case *Underscore,
		*String,
		*Char,
		*Float,
		*Integer:
		return t

	case *CompoundTerm:
		tt := &CompoundTerm{
			Functor:  t.Functor,
			Children: []Term{},
		}
		for _, c := range t.Children {
			tt.Children = append(tt.Children, CloneNsp(c, nsp))
		}
		return tt

	default:
		panic("attempting to clonensp an unknown type")
	}
}

package solver

import "github.com/xavier268/myprolog/parser"

// make a deep clone of Term
func Clone(t parser.Term) parser.Term {
	switch t := t.(type) {

	case *parser.Variable:
		return &parser.Variable{
			Name: t.Name,
			Nsp:  t.Nsp,
		}

	case *parser.Atom:
		return &parser.Atom{
			Value: t.Value,
		}

	case *parser.Underscore,
		*parser.String,
		*parser.Char,
		*parser.Float,
		*parser.Integer:
		return t

	case *parser.CompoundTerm:
		tt := &parser.CompoundTerm{
			Functor:  t.Functor,
			Children: []parser.Term{},
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
func CloneNsp(t parser.Term, nsp int) parser.Term {
	switch t := t.(type) {

	case *parser.Variable:
		return &parser.Variable{
			Name: t.Name,
			Nsp:  nsp,
		}

	case *parser.Atom:
		return &parser.Atom{
			Value: t.Value,
		}

	case *parser.Underscore,
		*parser.String,
		*parser.Char,
		*parser.Float,
		*parser.Integer:
		return t

	case *parser.CompoundTerm:
		tt := &parser.CompoundTerm{
			Functor:  t.Functor,
			Children: []parser.Term{},
		}
		for _, c := range t.Children {
			tt.Children = append(tt.Children, CloneNsp(c, nsp))
		}
		return tt

	default:
		panic("attempting to clonensp an unknown type")
	}
}

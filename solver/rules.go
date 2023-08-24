package solver

import (
	"fmt"

	"github.com/xavier268/myprolog/parser"
)

// Set of rules that can be applied in a state
type RuleSet struct {
	rules []*CompoundTerm
}

func (rs *RuleSet) AddRule(rule Term) {
	if rule == nil {
		panic("trying to add a nil rule")
	}
	r, ok := rule.(*CompoundTerm)
	if !ok {
		panic("Trying to add a Term as rule that is not a CompoundTerm")
	}
	if r.Functor != "rule" {
		panic("Trying to add a Term that is not a rule")
	}
	if len(r.Children) == 0 {
		fmt.Println(parser.START_RED, "WARNING : trying to add a rule with no children - ignored", parser.END_RED)
		return
	}
	rs.rules = append(rs.rules, r)
}

// Find the next rule applicable in the current context.
// Retuns an alpha-transformed rule.
// New State is created if many rules can be applied.
// Return nil if no rule available.
func FindNextRule(st *State) (*State, Term) {

	if len(st.Goals) == 0 {
		return st, nil // no rule can be found for a non existing goal
	}

	count := 0     // number of applicable rules. Needed to decide to fork state or not.
	selected := -1 // index of first applicable rule selected
	goal := st.Goals[0]

	// iterate over all rules, check which are applicable
	for i, rule := range st.Rules.rules {
		if SameStructure(rule.Children[0], goal) {
			if count == 0 {
				selected = i // remember only the first one found
			}
			count = count + 1
			if count >= 2 {
				break // multiple rules are applicable
			}
		}
	}

	if count <= 0 { // no rule found
		st.NextRule = len(st.Rules.rules)
		return st, nil //  will trigger backtracking ...
	}

	if count == 1 { // only one rule is applicable - do not fork state
		st.Uid = st.Uid + 1
		rule := CloneNsp(st.Rules.rules[selected], st.Uid)
		return st, rule
	}

	if count >= 2 { // more than one r=ule possible - need to fork state
		st.NextRule = selected + 1 // st becomes the old state we will backtrack into
		nst := NewState(st)        // fork a new state to continue with
		nst.NextRule = selected
		nst.Uid = st.Uid + 1
		rule := CloneNsp(st.Rules.rules[selected], st.Uid)
		return st, rule
	}
	panic("internal error - code should be unreacheable")
}

// Apply a rule to the state.
// No new state is created.
// Returns nil if no rule is applicable.
func ApplyRule(st *State, rule Term) (*State, error) {
	panic("not implemented")
}

// Does the structure of both terms match and should be considered,
// regardless of potential constraints . It may not be unifiable later.
// The goal is to eliminate early non match.
func SameStructure(t1, t2 Term) bool {

	// special nil cases
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}

	// handle underscores
	if _, ok := t1.(*Underscore); ok {
		return true
	}
	if _, ok := t2.(*Underscore); ok {
		return true
	}

	// handle variables
	if _, ok := t1.(*Variable); ok {
		return true
	}
	if _, ok := t2.(*Variable); ok {
		return true
	}

	// handle other term
	switch t1 := t1.(type) {

	case *Integer:
		switch t2 := t2.(type) {
		case *Integer:
			return t1.Value == t2.Value
		default:
			return false
		}

	case *Float:
		switch t2 := t2.(type) {
		case *Float:
			return t1.Value == t2.Value
		default:
			return false
		}

	case *String:
		switch t2 := t2.(type) {
		case *String:
			return t1.Value == t2.Value
		default:
			return false
		}

	case *Char:
		switch t2 := t2.(type) {
		case *Char:
			return t1.Char == t2.Char
		default:
			return false
		}

	case *Atom:
		switch t2 := t2.(type) {
		case *Atom:
			return t1.Value == t2.Value
		default:
			return false
		}

	case *CompoundTerm:
		switch t2 := t2.(type) {
		case *CompoundTerm:
			if (t1.Functor != t2.Functor) || len(t1.Children) != len(t2.Children) {
				return false
			}
			for i, c1 := range t1.Children {
				if !SameStructure(c1, t2.Children[i]) {
					return false
				}
			}
			return true
		default:
			return false

		}
	default:
		panic("internal error - code should be unreacheable")
	}
}

package solver

import (
	"fmt"

	"github.com/xavier268/myprolog/parser"
)

// Set of rules that can be applied in a state
type RuleSet struct {
	rules []CompoundTerm
}

func (rs *RuleSet) AddRule(rule Term) {
	if rule == nil {
		panic("trying to add a nil rule")
	}
	r, ok := rule.(CompoundTerm)
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
// Returns an alpha-transformed rule.
// New State is created if a choice was made between multiple applicable rules.
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
		rule := st.Rules.rules[selected].CloneNsp(st.Uid)
		st.NextRule = selected + 1 // to ensure no loop
		return st, rule
	}

	if count >= 2 { // more than one r=ule possible - need to fork state
		st.NextRule = selected + 1 // st becomes the old state we will backtrack into
		nst := NewState(st)        // fork a new state to continue with
		nst.NextRule = selected
		nst.Uid = st.Uid + 1
		rule := st.Rules.rules[selected].CloneNsp(nst.Uid)
		return nst, rule
	}
	panic("internal error - code should be unreacheable")
}

// Apply a rule to the state.
// No new state is created, no backtracking is performed here, an error is returned instead.
// Goals and constraints are potentially updated.
// Returns err!=nil if backtracking will be required.
func ApplyRule(st *State, rule Term) (*State, error) {
	var err error
	if rule == nil {
		panic("Trying to apply a nil rule")
	}
	crule, ok := rule.(CompoundTerm)
	if !ok || crule.Functor != "rule" {
		panic("Trying to apply a rule that is not a valid rule")
	}
	if len(crule.Children) == 0 {
		panic("Trying to apply a rule with no head")
	}
	// Try to unify the head with the first goal, adding constraints to the state
	st, err = Unify(st, crule.Children[0], st.Goals[0])
	if err != nil {
		return st, err
	}
	// Update goals after successfull unification with the rule rhs.
	st.Goals = append(crule.Children[1:], st.Goals[1:]...)
	st.NextRule = 0 // reset next rule because goals changed
	return st, nil

}

// Does the structure of both terms seem to match and should be considered,
// regardless of potential constraints . It may not be unifiable later.
// Do not check too deep into the structure, since predicates or expressions
// can be hidden inside that will disappear later.
// Float and Integer could unify if same float64 value.
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
	if _, ok := t1.(Underscore); ok {
		return true
	}
	if _, ok := t2.(Underscore); ok {
		return true
	}

	// handle variables
	if _, ok := t1.(Variable); ok {
		return true
	}
	if _, ok := t2.(Variable); ok {
		return true
	}

	// handle other non nil, non variable, non underscore
	switch t1 := t1.(type) {

	case Number:
		switch t2 := t2.(type) {
		case Number:
			return t1.Eq(t2)

		default:
			return false
		}

	case String:
		switch t2 := t2.(type) {
		case String:
			return t1.Value == t2.Value
		default:
			return false
		}

	case Atom:
		if AtomPredicate[t1.Value] {
			return true // predicates may modify the structure and are not a reliable compare.
		}
		switch t2 := t2.(type) {
		case Atom:
			if AtomPredicate[t1.Value] {
				return true // predicates may modify the structure and are not a reliable compare.
			}
			return t1.Value == t2.Value
		case CompoundTerm:
			_, ok := CompPredicateMap[t2.Functor]
			return ok // predicates may modify the structure and are not a reliable compare.
		default:
			return false
		}

	case CompoundTerm:
		if _, ok := CompPredicateMap[t1.Functor]; ok {
			return true // predicates may modify the structure and are not a reliable compare.
		}
		switch t2 := t2.(type) {
		case Atom:
			return AtomPredicate[t2.Value] // predicates may modify the structure and are not a reliable compare.
		case CompoundTerm:
			if _, ok := CompPredicateMap[t2.Functor]; ok {
				return true // predicates may modify the structure and are not a reliable compare.
			}

			if (t1.Functor != t2.Functor) || len(t1.Children) != len(t2.Children) {
				return false
			}
			// If there are nested predicates, they should match to be considered unifiable.
			// The thinking here is that predicates only produce an effct as a top level goal,
			// otherwise, they behave as a usual Atom/CompoundTerm.
			for i, c := range t1.Children {
				if !SameStructure(c, t2.Children[i]) {
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

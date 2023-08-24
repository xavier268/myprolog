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

	goal := st.Goals[0]
	count := 0     // number of applicable rules. Needed to decide to fork state or not.
	selected := -1 // index of first applicable rule selected

	for i := st.NextRule; i < len(st.Rules.rules); i++ {

		// TODO - set selected rule to the first rule matching the goal arity and functor.
		// Set the count to the number of matching rules (stop at 2)

	}

	if count == 0 { // no rule found
		st.NextRule = len(st.Rules.rules)
		return st, nil //  will trigger backtracking ...
	}

	if count == 1 { // only one rule is applicable
		st.Uid = st.Uid + 1
		rule := CloneNsp(st.Rules.rules[selected], st.Uid)
		return st, rule
	}

	if count >= 2 { // need to fork state
		st.NextRule = selected + 1 // st becomes the old state we will backtrack into
		nst := NewState(st)        // fork a new state to continue with
		nst.NextRule = selected
		nst.Uid = st.Uid + 1
		rule := CloneNsp(st.Rules.rules[selected], st.Uid)
		return st, rule
	}

	// TO DO
	panic("not fully implemented")

}

// Apply a rule to the state.
// No new state is created.
// Returns nil if no rule is applicable.
func ApplyRule(st *State, rule Term) (*State, error) {
	panic("not implemented")
}

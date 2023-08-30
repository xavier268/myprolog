package solver

import (
	"fmt"
	"strings"

	"github.com/xavier268/myprolog/parser"
)

// State maintains the current state of the puzzle current state, and is used for backtracking.
// A new state is created if and only if we have to make a choice and backtracking may be necessary.
type State struct {
	Rules       *RuleSet // pointer to applicable ruleset.
	Constraints []Constraint
	Goals       []Term // Goals we are trying to erase. They are ordered in the order they should be executed.
	Parent      *State // for backtracking
	Uid         int    // Unique id for alpha-conversion of rules
	NextRule    int    // index of the next rule to consider. It may not be applicable.
}

// Creates a new state, using provided parent.
// Constraints and goals cloned from parent.
// Rules are copied but not cloned.
// Uid is copied verbatim.
// Next rule is set to 0.
func NewState(parent *State) *State {

	if parent == nil {
		return new(State)
	}

	st := new(State)
	st.Rules = parent.Rules // default is to point to the same ruleset
	st.Parent = parent
	st.Uid = parent.Uid
	st.Constraints = append(st.Constraints, parent.Constraints...)
	st.Goals = append(st.Goals, parent.Goals...)
	return st
}

func (st *State) AddRule(rule ...CompoundTerm) {
	if st.Rules == nil {
		st.Rules = &RuleSet{rules: []parser.CompoundTerm{}}
	}
	st.Rules.rules = append(st.Rules.rules, rule...)
}

func (st *State) String() string {
	sb := new(strings.Builder)

	fmt.Fprintf(sb, "State: \n")
	fmt.Fprintf(sb, "Rules :\n%s", st.Rules)
	fmt.Fprintf(sb, "Constraints : %s\n", st.Constraints)
	fmt.Fprintf(sb, "Goals : %s\n", st.Goals)
	fmt.Fprintf(sb, "NextRule : %d\n", st.NextRule)
	fmt.Fprintf(sb, "Uid : %d\n", st.Uid)
	if st.Parent != nil {
		fmt.Fprintf(sb, "Parent : NO\n")
	} else {
		fmt.Fprintf(sb, "Parent : YES\n")
	}
	return sb.String()
}

// Partially resets the state. The state is modified.
func (st *State) Abort() {
	if st == nil {
		return
	}
	st.NextRule = 0
	st.Constraints = st.Constraints[:0]
	st.Goals = st.Goals[:0]
}

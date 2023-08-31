package solver

import (
	"fmt"
	"strings"
)

// State maintains the current state of the puzzle current state, and is used for backtracking.
// A new state is created if and only if we have to make a choice and backtracking may be necessary.
type State struct {
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
		st := new(State)
		st.Uid = 1
		return st
	}

	st := new(State)
	st.Parent = parent
	st.Uid = parent.Uid + 1
	st.Constraints = append(st.Constraints, parent.Constraints...) // deep copy the constraints
	st.Goals = append(st.Goals, parent.Goals...)                   // deep copy the goals
	return st
}

func (st *State) String() string {
	sb := new(strings.Builder)

	fmt.Fprintf(sb, "State: \n")
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

// a printable explanation of the rules applied to reach this state.
// rules are listed in the order they were applied.
func (st *State) RulesHistory() string {

	res := ""
	for s := st; s != nil; s = s.Parent {
		if s.NextRule <= 0 {
			continue
		}
		rule := MYDB.rules[s.NextRule-1]
		res = fmt.Sprintf("rule#%d\t%s\n%s", s.NextRule, rule.Pretty(), res)
	}

	return res

}

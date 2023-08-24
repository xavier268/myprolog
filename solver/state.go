package solver

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
	st.Constraints = make([]Constraint, 0, len(parent.Constraints))
	st.Goals = make([]Term, 0, len(parent.Goals))
	for _, g := range parent.Goals {
		st.Goals = append(st.Goals, Clone(g))
	}
	for _, c := range parent.Constraints {
		st.Constraints = append(st.Constraints, c.Clone())
	}
	return st
}

// Add a constraint to state and simplify immediately.
func (s *State) AddConstraint(c Constraint) (err error) {
	s.Constraints = append(s.Constraints, c)
	s.Constraints, err = SimplifyConstraints(s.Constraints)
	return err
}
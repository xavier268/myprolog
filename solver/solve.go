package solver

import (
	"github.com/xavier268/myprolog/parser"
)

// Aliases to parser types.
type Term = parser.Term

type CompoundTerm = parser.CompoundTerm
type Variable = parser.Variable
type Underscore = parser.Underscore
type Atom = parser.Atom
type String = parser.String
type Number = parser.Number

// useful numbers
var ZeroNumber, OneNumber = parser.ZeroNumber, parser.OneNumber

// Solve for a given state.
// Backtracking is managed only in this function.
func Solve(st *State, sh SolutionHandler) *State {
	var err error
	for {

		if st == nil { // we have no more solutions, stop !
			return nil
		}
		if len(st.Goals) == 0 { // we have found a solution, answer is available in the state constraints.
			st := sh(st)   // handle the solution and ensure the state moves backward
			if st == nil { // stop !
				return st
			} else {
				continue // loop again
			}
		}

		// Apply predicate or remove obvious goals on the first top goal, if possible
		// State MAY change.
		st, err = DoPredicate(st)
		if st == nil { // stop !
			return nil
		}
		if err != nil { // backtrack required
			st = st.Parent
			continue // loop again
		}
		if len(st.Goals) == 0 { // we have found a solution, answer is available in the state constraints.
			// it is the responsability of sh to replace st by st.Parent if we want to search next soltions
			st = sh(st)
			if st == nil { // stop solving !
				return st
			} else {
				continue // loop again, try to find more solutions ...
			}
		}

		// Find next rule to apply.
		// State is forked if a rule is returned to try.
		var rule Term
		st, rule = FindNextRule(st)
		if st == nil { // stop solver !
			return nil
		}
		if len(st.Goals) == 0 { // we have found a solution, answer is available in the state constraints.
			st = sh(st)
			if st == nil { // stop !
				return st
			} else {
				continue // loop again
			}
		}
		if rule == nil { // no rule to apply, backtrack
			st = st.Parent
			continue // loop again
		}

		// Apply selected rule, replacing the first goal with rule head.
		// No new state is created, backtracking is not performed inside ApplyRule.
		st, err = ApplyRule(st, rule)
		if st == nil {
			return nil
		}
		if err != nil {
			st = st.Parent
			continue
		}
		if len(st.Goals) == 0 { // we have found a solution, answer is available in the state constraints.
			st = sh(st)
			if st == nil { // stop !
				return st
			} else {
				continue // loop again
			}
		}

		// Iterate ...
	}
}

// Handles a solution. Return value indicates what to do next.
// The state MAY be modified.
// If returned state is nil, Solver stops.
type SolutionHandler func(st *State) *State

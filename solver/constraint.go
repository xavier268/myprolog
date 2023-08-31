package solver

import (
	"fmt"
)

// a Constraint is immutable
type Constraint interface {
	// deep copy
	Clone() Constraint
	// Check will check validity of constraint, clean it or simplify it.
	// It will return the constraint itself,  possibly modified, or nil if constraint should be ignored (
	// CAUTION : return can be nil, despite a nil error !
	// An error means the constraint is impossible to satisfy (e.g. positive occur check, empty number interval, ...)
	Check() (Constraint, error)
	// simplify c, taking into account the calling constraint (unchanged).
	// if error, other results are not significant and should not be used.
	// If changed, replace c by all of cc (possibly empty, to just suppress c).
	// If changed is false, cc has no meaning.
	// Input constraints MUST BE checked
	// Output constraints are checkedutput.
	Simplify(c Constraint) (cc []Constraint, changed bool, err error)
	String() string
	GetV() Variable // Get (a copy of) the main variable forthe constraint
}

// X = ...
var _ Constraint = VarIsVar{}          // X = Y
var _ Constraint = VarIsCompoundTerm{} // X = a(b,Y)

// X = constant atomic term
var _ Constraint = VarIsString{}
var _ Constraint = VarIsAtom{}

// constraint on numbers
var _ Constraint = VarINT{} // X is integer number
var _ Constraint = VarEQ{}  // X = n
var _ Constraint = VarLTE{} // X <= n
var _ Constraint = VarLT{}  // X < n
var _ Constraint = VarGT{}  // X > n
var _ Constraint = VarGTE{} // X >! n

var ErrInvalidConstraintNaN = fmt.Errorf("invalid constraint (NaN)")
var ErrInvalidConstraintEmptyRange = fmt.Errorf("invalid constraint, specified range is empty")
var ErrInvalidConstraintSimplify = fmt.Errorf("incompatible constraints detected when simplifying")
var ErrPositiveOccur = fmt.Errorf("positive occur check")
var ErrNaN = fmt.Errorf("not a number (NaN)")

var ErrNotImplemented = fmt.Errorf(RED + "not implemented" + RESET) // should be removed later on //TODO

// Attempt to simplify constraint list.
// Return error if an incompatibility was detected.
// Constraints are supposed to have been checked before calling this function.
// Its is a garantee that returned constraints are checked.
func SimplifyConstraints(constraints []Constraint) ([]Constraint, error) {

	for changed := true; changed; { // loop until no more changes are made ...
		changed = false
	hasChangedLoop: // loop again for each changed constraint ...
		for i, c := range constraints {
			for j, d := range constraints {
				if i != j && c != nil && d != nil { // filter out nil values
					dd, ch, err := c.Simplify(d)
					if err != nil {
						return constraints, err // will trigger backtracking ...
					}
					if ch { // update if requested to do so. Could be nil.
						changed = true
						if len(dd) == 0 { // suppress this  constraint
							constraints[j] = nil
						} else {
							constraints[j] = dd[0]                       // replace d
							constraints = append(constraints, dd[1:]...) // if more than 1, add at the end
						}
						break hasChangedLoop
					}
				}
			}
		}
	}

	// Now, that we went through the previous loop once without changing anything,
	// we should clean/remove remaining nil constraints
	cc := make([]Constraint, 0, len(constraints))
	for _, c := range constraints {
		if c != nil {
			cc = append(cc, c)
		}
	}
	return cc, nil
}

// Filter constraints to display, keeping only 0 Nsp.
func FilterSolutions(cc []Constraint) []Constraint {
	ncc := make([]Constraint, 0, len(cc))
	for _, c := range cc {

		if c == nil {
			continue
		}

		if c.GetV().Nsp == 0 {
			ncc = append(ncc, c)
		}
	}
	return ncc
}

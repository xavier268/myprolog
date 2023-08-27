package solver

import (
	"fmt"
)

// a Constraint is immutable
type Constraint interface {
	// deep copy
	Clone() Constraint
	// Check will check validity of constraint, clean it or simplify it.
	// It will return the constraint itself,  possibly modified, or nil if constraint should be ignored.
	// An error means the constraint is impossible to satisfy (e.g. positive occur check, empty number interval, ...)
	Check() (Constraint, error)
	// simplify c, taking into account the calling constraint (unchanged).
	// If changed, replace c by all of cc (possibly empty, to just suppress c).
	// It is assumed that input constraints are checked, it is garanteed that output constraints have been checked.
	Simplify(c Constraint) (cc []Constraint, changed bool, err error)
	String() string
}

var _ Constraint = VarIsCompoundTerm{}
var _ Constraint = VarIsString{}
var _ Constraint = VarIsNumber{}
var _ Constraint = VarIsVar{}
var _ Constraint = VarIsAtom{}

var ErrInvalidConstraintNaN = fmt.Errorf("invalid constraint (NaN)")
var ErrInvalidConstraintEmptyRange = fmt.Errorf("invalid constraint, specified range is empty")
var ErrInvalidConstraintSimplify = fmt.Errorf("incompatible constraints detected when simplifying")
var ErrPositiveOccur = fmt.Errorf("positive occur check")
var ErrNotImplemented = fmt.Errorf(RED + "not implemented" + RESET) // should be removed later on //TODO

// Attempt to simplify constraint list.
// Return error if an incompatibility was detected.
// Constraints are supposed to have been checked before calling this function.
// Its is a garantee that retuened constraints are checked.
func SimplifyConstraints(constraints []Constraint) ([]Constraint, error) {

hasChangedLoop: // loop again for each changed constraint ...
	for i, c := range constraints {
		for j, d := range constraints {
			if i != j && c != nil && d != nil {
				dd, ch, err := c.Simplify(d)
				if err != nil {
					return constraints, err // will trigger backtracking ...
				}
				if ch { // update if requested to do so. Could be nil.
					if len(dd) == 0 { // suppress this  constraint
						constraints[j] = nil
					} else {
						constraints[j] = dd[0] // replace with one or more new constraints
						constraints = append(constraints, dd[1:]...)
					}
					break hasChangedLoop
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

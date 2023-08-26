package solver

// Attempt to simplify constraint list.
// Return error if an incompatibility was detected.
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

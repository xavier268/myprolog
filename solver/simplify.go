package solver

// Attempt to simplify constratint list.
// Return error if an incompatibility was detected.
func SimplifyConstraints(constraints []Constraint) ([]Constraint, error) {

	changed := true
	for changed { // iterate while constraints change ...

		changed = false

		for i, c := range constraints {
			for j, d := range constraints {
				if i != j && c != nil && d != nil {
					dd, ch, err := c.Simplify(d)
					if err != nil {
						return constraints, err // will trigger backtracking ...
					}
					if ch { // update if requested to do so. Could be nil.
						constraints[j] = dd
					}
					changed = changed || ch
				}
			}
		}

		if changed { // remove nil constraints
			cc := make([]Constraint, 0, len(constraints))
			for _, c := range constraints {
				if c != nil {
					cc = append(cc, c)
				}
			}
			constraints = cc
		}
	}
	return constraints, nil
}

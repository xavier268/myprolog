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
					changed = changed || ch
					constraints[j] = dd
				}
			}
		}
	}

	return constraints, nil
}

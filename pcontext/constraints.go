package pcontext

type Constraint interface {
	Cons() // marquer for the interface. No-op.
	String() string
}

// Add a constraint to the current context.
// Error if impossibility is detected (backtracking will be required !)
func (pc *PContext) SetConstraint(cc Constraint) error {

	if cc == nil || pc == nil {
		return nil
	}

	switch c := cc.(type) {

	case ConEqual: // form X = xxx
		if c.IsObvious() {
			return nil
		}
		if err := c.Verify(); err != nil { // positive occur check
			return err
		}

		// Update old constr if needed, replacing previous X occurence by xxx
		for i, old := range pc.cstr {
			old2 := c.Update(old)
			if old2 != nil {
				pc.cstr[i] = old2
			}
		}

		// Look for previously existing X = yyy ?
		for _, old := range pc.cstr {
			switch old2 := old.(type) {
			case ConEqual:
				if old2.v == c.v { // same X = ....
					err := pc.Unify(old2.t, c.t)
					if err != nil {
						// we couldn't unify a and b in X=a & X=b !
						return err
					} else {
						// since we could unify a and b in X=a & X=b, do not add X=b, not required.
						return nil
					}
				}
			default: // ignore
			}
		}

		// by default, add this new constraint and return no error
		pc.cstr = append(pc.cstr, c)
		return nil

	default:
		panic("unknown constraint")

	}

	//return nil
}

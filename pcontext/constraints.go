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

	case ConEqual:
		if c.IsObvious() {
			return nil
		}
		if err := c.Verify(); err != nil { // positive occur check
			return err
		}
		for i, old := range pc.cstr {
			old2 := c.Update(old)
			if old2 != nil {
				pc.cstr[i] = old2
			}
		}

	default:
		panic("unknown constraint")

	}

	return nil
}

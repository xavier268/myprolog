package pcontext

type Constraint interface {
	Cons() // marquer for the interface. No-op.
	String() string
}

// Add a constraint to the current context.
// Error if impossibility is detected, due to the constraint it self.
// Not checks are performed vs previously existing constraints.
func (pc *PContext) SetConstraint(cc Constraint) error {

	if cc == nil || pc == nil {
		return nil
	}

	switch c := cc.(type) {

	case ConEqual: // form X = xxx
		if c.IsObvious() {
			return nil
		}
		if err := c.VerifyPosOcc(); err != nil { // positive occur check, don't add !
			return err
		}

		// by default, add this new constraint and return no error
		pc.cstr = append(pc.cstr, c)
		return nil

	default:
		panic("unknown constraint")

	}

	//return nil
}

// Simplify the constraint set
func (pc *PContext) Simplify() error {

}

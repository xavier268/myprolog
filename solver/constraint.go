package solver

type Constraint interface {
	Clone() Constraint
}

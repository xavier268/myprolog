// Package pcontext contains the constraint context.
package pcontext

import (
	"time"

	"github.com/xavier268/myprolog/node"
)

type PContext struct {
	goals   []*node.Node
	rules   []*node.Node
	current int          // index of the current rule being tried.
	cstr    []Constraint // Current list of constraints, can contain nils
	parent  *PContext
	start   time.Time
}

// Create a new PContext, using a program and a goal (ie, a query).
// The program may already contain goals that will be added first to the goals list.
func NewPContext(prog *node.Node, goal *node.Node) *PContext {
	pc := new(PContext)
	pc.current = 0
	kr, _ := node.NewKeyword("rule")
	kq, _ := node.NewKeyword("query")
	for _, r := range prog.GetChildren() {
		if r.GetLoad() == kr {
			pc.rules = append(pc.rules, r)
		}
		if r.GetLoad() == kq {
			pc.goals = append(pc.goals, r)
		}
	}
	pc.parent = nil
	pc.cstr = nil
	pc.goals = append(pc.goals, goal)
	pc.start = time.Now()
	return pc
}

func (pc *PContext) String() string {
	panic("todo")
}

func (pc *PContext) DisplaySolutions() {
	panic("todo")
}

// Freeze will remove the parent context, making backtracking impossible.
func (pc *PContext) Freeze() {
	pc.parent = nil
}

// Creates a new PContext from the provided parent pcontext.
func (pc *PContext) Clone() *PContext {
	n := new(PContext)
	n.parent = pc
	n.rules = append(n.rules, pc.rules...)
	n.goals = append(n.goals, pc.goals...)
	n.current = pc.current
	// special handling for constraints to cleanup nil constraints
	for _, c := range pc.cstr {
		if c != nil {
			n.cstr = append(n.cstr, c)
		}
	}
	n.start = time.Now()
	return n
}

// Parent returns to the previous state, or nil if no backtracking possible.
func (pc *PContext) Parent() *PContext {
	return pc.parent
}

// Add a constraint to the current context.
// Error if impossibility is detected (backtracking will be required !)
func (pc *PContext) Set(cc Constraint) error {

	if cc == nil || pc == nil {
		return nil
	}

	for i := 0; i < len(pc.cstr); i++ {
		c := pc.cstr[i]
		if c == nil {
			continue
		}
		remove, nc, err := cc.Merge(c)
		if err != nil {
			return err
		}
		if remove {
			pc.cstr[i] = nil
		}
		pc.cstr = append(pc.cstr, nc...)
	}
}

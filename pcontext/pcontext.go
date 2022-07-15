// Package pcontext contains the constraint context.
package pcontext

import (
	"time"

	"github.com/xavier268/myprolog/node"
)

type PContext struct {
	goals  []*node.Node
	rules  []*node.Node
	next   int          // index of the next rule to try
	cstr   []Constraint // Current list of constraints
	parent *PContext
	start  time.Time
}

func NewPContext(prog *node.Node, goal *node.Node) *PContext {
	pc := new(PContext)
	pc.goals = append(pc.goals, goal)
	pc.next = 0
	kr, _ := node.NewKeyword("rule")
	for _, r := range prog.GetChildren() {
		if r.GetLoad() == kr {
			pc.rules = append(pc.rules, r)
		}
	}
	pc.parent = nil
	pc.cstr = nil
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
	n.cstr = append(n.cstr, pc.cstr...)
	n.rules = append(n.rules, pc.rules...)
	n.goals = append(n.goals, pc.goals...)
	n.next = pc.next
	n.start = time.Now()
	return n
}

// Parent returns to the previous state, or nil if no backtracking possible.
func (pc *PContext) Parent() *PContext {
	return pc.parent
}

// Add a constraint to the current context.
// Error if impossibility is detected (backtracking will be required !)
func (pc *PContext) Set(cc *Constraint) error {

	for _, c := range pc.cstr {
		err := c.Merge(cc)
		if err != nil {
			return err
		}
	}
	return nil
}

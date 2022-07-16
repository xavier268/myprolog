// Package pcontext contains the constraint context.
package pcontext

import (
	"fmt"
	"strings"
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
	uid     *int // unique id generator sherd accross context
}

func (pc *PContext) UID() int {
	*pc.uid = *pc.uid + 1
	return *pc.uid
}

// Create a new initial PContext, using a program and optionnal goals (ie, a query).
// The program may already contain goals that will be added first to the goals list.
func NewPContext(prog *node.Node, goals ...*node.Node) *PContext {

	pc := new(PContext)
	pc.current = 0
	pc.AddRules(prog)
	pc.AddGoals(prog)
	pc.parent = nil
	pc.cstr = nil
	pc.goals = append(pc.goals, goals...) // additionnal goals will run after the ones already in the program.
	pc.start = time.Now()
	pc.uid = new(int)
	return pc
}

// AddRules append new goals from new program.
func (pc *PContext) AddRules(prog *node.Node) {
	ruleLoad, _ := node.NewKeyword("rule")
	for _, r := range prog.GetChildren() {
		if r.GetLoad() == ruleLoad {
			pc.rules = append(pc.rules, r)
		}
	}
}

// AddGoals append new goals from new program.
func (pc *PContext) AddGoals(prog *node.Node) {
	queryLoad, _ := node.NewKeyword("query")
	for _, r := range prog.GetChildren() {
		if r.GetLoad() == queryLoad {
			pc.goals = append(pc.goals, r)
		}
	}
}

func (pc *PContext) String() string {
	var sb strings.Builder
	fmt.Fprintln(&sb, "------ pcontext --")
	fmt.Fprintf(&sb, "\tConstraints :\n%v\n", pc.cstr)
	fmt.Fprintf(&sb, "\tGoals :\n%v\n", pc.goals)
	fmt.Fprintf(&sb, "\tParent :\t%v\n", pc.parent)
	fmt.Fprintf(&sb, "\tCurrent rule :\t%v\n", pc.current)
	fmt.Fprintf(&sb, "\tRules :\n%v\n", pc.rules)
	fmt.Fprintf(&sb, "\tStart :\t%v\n", pc.start)
	fmt.Fprintf(&sb, "\tUID :\t%d\n", *pc.uid)

	return sb.String()
}

func (pc *PContext) Display() {
	fmt.Println("TODO : Better results display needed ...")
	fmt.Println(pc) // TODO - better solution display that a simple dump !
}

// Freeze will remove the parent context, making backtracking impossible.
func (pc *PContext) Freeze() {
	pc.parent = nil
}

// Creates a new PContext from the provided parent pcontext.
func (pc *PContext) Push() *PContext {
	n := new(PContext)
	n.parent = pc
	n.uid = pc.uid
	*n.uid = *n.uid + 1
	n.rules = append(n.rules, pc.rules...)
	n.goals = append(n.goals, pc.goals...)
	n.current = 0 // start looking at all rules again !
	// special handling for constraints to cleanup nil constraints
	for _, c := range pc.cstr {
		if c != nil {
			n.cstr = append(n.cstr, c)
		}
	}
	n.start = time.Now()
	return n
}

// Pop returns to the previous state, or nil if no backtracking possible.
func (pc *PContext) Pop() *PContext {
	if pc == nil {
		return nil
	}
	return pc.parent
}

// Add a constraint to the current context.
// Error if impossibility is detected (backtracking will be required !)
func (pc *PContext) SetConstraint(cc Constraint) error {

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

	panic("todo")
}

// TODO
func (pc *PContext) StringContent(n *node.Node) string {
	return "PContext.StringContent not implemented - using String instead :\n" + n.String()
}

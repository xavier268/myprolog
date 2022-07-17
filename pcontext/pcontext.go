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
	uid     *int // unique id generator shared across context
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
	fmt.Fprintf(&sb, "Constrt : %v\n", pc.cstr)
	fmt.Fprintf(&sb, "Goals   : %v\n", pc.goals)
	fmt.Fprintf(&sb, "Root    : %v\n", pc.parent == nil)

	return sb.String()
}

// Like String, but more detailed.
func (pc *PContext) StringDetailed() string {
	if pc == nil {
		return fmt.Sprintf("\n------ pcontext --\n%v", nil)
	}
	var sb strings.Builder
	fmt.Fprintln(&sb, "\n------ pcontext --")
	fmt.Fprintf(&sb, "\nConstraints :\n%v", pc.cstr)
	fmt.Fprintf(&sb, "\nGoals :\n%v", pc.goals)
	fmt.Fprintf(&sb, "\nRoot :\t%v", pc.parent == nil)
	fmt.Fprintf(&sb, "\nCurrent rule :\t%v\n", pc.current)
	fmt.Fprintf(&sb, "\nRules :\n%v", pc.rules)
	fmt.Fprintf(&sb, "\nStart :\t%v", pc.start)
	fmt.Fprintf(&sb, "\nUID :\t%d", *pc.uid)

	return sb.String()
}

// Results, if any.
func (pc *PContext) Display() {
	fmt.Println("Constraints : ", pc.ResultString())
}

// Results, if any.
func (pc *PContext) ResultString() string {
	if pc == nil {
		return fmt.Sprint(nil)
	}
	return fmt.Sprint(pc.cstr)
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

// TODO - redirect to variable content
func (pc *PContext) StringContent(n *node.Node) string {
	return n.String()
}

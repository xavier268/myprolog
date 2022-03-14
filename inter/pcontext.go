package inter

import (
	"fmt"
	"strings"
)

// PContext maintain the list of constraints needed to solve a query.
type PContext struct {
	// parent context. Used for backtracking.
	prt *PContext
	ctp CType
	lhs *Node
	rhs *Node
}

func NewPContext() *PContext {
	c := new(PContext)
	c.prt = nil
	return c
}

// CType is the type of constraint.
type CType int

const (
	EQ     CType = iota // X = toto(1,2)
	LT                  // X < 3
	GT                  // X < 3
	NEQ                 //X != b(a,c)
	NUMBER              // X is a number ...
)

func (ct CType) String() string {
	switch ct {
	case EQ:
		return "="
	default:
		return "?notIplm?"
	}
}

// Set the given constraint, returns the new modified context.
// Binding to nil is the same as not set.
func (ctx *PContext) Set(ctp CType, lhs *Node, rhs *Node) *PContext {
	//fmt.Printf("DEBUG : Setting %d : %s  --> %s\n", ctp, lhs, rhs)
	// if already set to the correct value, do nothing.
	if rhs == ctx.Get(ctp, lhs) {
		return ctx
	}
	// actual setting, even for nil !
	return &PContext{
		prt: ctx,
		ctp: ctp,
		lhs: lhs,
		rhs: rhs,
	}
}

// Get the rhs for a given lhs and type.
// Return nil if not found.
func (ctx *PContext) Get(ctp CType, lhs *Node) (rhs *Node) {
	/*if ctx.prt == ctx {
		panic("loop in pcontext !")
	}*/
	for c := ctx; c != nil; c = c.prt {
		if c.ctp == ctp && c.lhs == lhs {
			return c.rhs
		}
	}
	return nil
}

// GetRHSContains returns the FIRST rhs that contains the node x (typically a Variable, for further substitution), possible at a low nested level.
// The result returned is NOT rescoped, and probably should be before it is modified.
// Retun nil if not found.
func (ctx *PContext) getRHSContains(ctp CType, x *Node) *Node {

	dedup := make(map[*Node]bool)

	for c := ctx; c != nil; c = c.prt {
		if c.lhs == nil || c.ctp != ctp || dedup[c.lhs] {
			continue
		}
		dedup[c.lhs] = true
		if c.rhs == x || c.rhs.contains(x) {
			return c.rhs
		}
	}
	return nil
}

func (ctx *PContext) dump() {
	fmt.Printf("-------- dump context %p -----------\n", ctx)

	for c := ctx; c != nil; c = c.prt {
		fmt.Printf("%02d:\t%s   ---->  %s\n", c.ctp, c.lhs, c.rhs)
	}
	fmt.Println("String : ", ctx)
	fmt.Println("------------------------------------")
}

// String for a single line, deduped context.
func (ctx *PContext) String() string {

	var b strings.Builder
	dedup := make(map[*Node]bool)

	for c := ctx; c != nil; c = c.prt {
		if c.lhs != nil && !dedup[c.lhs] {
			dedup[c.rhs] = true
			fmt.Fprintf(&b, "%s %s %s, ", c.lhs, c.ctp, c.rhs)
		}
	}
	return b.String()
}

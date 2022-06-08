package inter

import (
	"fmt"
	"strings"
)

// PContext maintain the list of constraints needed to solve a query.
type PContext struct {
	// parent context. Used for backtracking.
	prt *PContext
	ctp CType  // the type of constraint
	lhs string // LHS can only be a Variable name to which the constraint applies. It should never be the empty string.
	rhs *Node  // any node will do, can be nil (ie of NUMBER ...)
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
	LT                  // X < 3, only for NUMBERs
	GT                  // X < 3, only for NUMBERs
	NUMBER              // X is a number ...
	INT                 // X is an integer, all INTs are NUMBERs
)

func (ct CType) String() string {
	switch ct {
	case EQ:
		return "EQ"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case NUMBER:
		return "NUMBER"
	case INT:
		return "INT"

	default:
		return "?notIplm?"
	}
}

// Set lhs for EQ.
// Actual change is made only if not alredy set.
func (ctx *PContext) SetEQ(lhs string, rhs *Node) *PContext {
	if lhs == "" || lhs == "_" || !isVariable(lhs) {
		panic("Trying to set an invalid EQ constraint, with lhs : " + lhs)
	}

	// if already set to the correct value, do nothing.
	prev := ctx.GetEQ(lhs)
	if prev == nil && rhs == nil {
		return ctx // no change
	}
	if prev != nil && prev.Equal(rhs) {
		return ctx // no change
	}
	// actual setting, even for nil !
	return &PContext{
		prt: ctx,
		ctp: EQ,
		lhs: lhs,
		rhs: rhs,
	}
}

// Get the rhs where lhs EQ rhs.
// Return nil if not found.
func (ctx *PContext) GetEQ(lhs string) (rhs *Node) {
	for c := ctx; c != nil; c = c.prt {
		if c.ctp == EQ && c.lhs == lhs {
			return c.rhs
		}
	}
	return nil
}

// SetNUMBER sets lhs to ba a number.
// Actual change is made only if not alredy set.
func (ctx *PContext) SetNUMBER(lhs string) *PContext {
	if lhs == "" || lhs == "_" || !isVariable(lhs) {
		panic("Trying to set an invalid NUMBER constraint, with lhs : " + lhs)
	}
	// If already a number, no change
	if ctx.IsInt(lhs) || ctx.IsNumber(lhs) {
		return ctx // not needed
	}
	// Otherwise, set it !
	return &PContext{
		prt: ctx,
		ctp: NUMBER,
		lhs: lhs,
		rhs: nil,
	}
}

// SetINT sets lhs to ba an  INT.
// Actual change is made only if not alredy set.
func (ctx *PContext) SetINT(lhs string) *PContext {
	if lhs == "" || lhs == "_" || !isVariable(lhs) {
		panic("Trying to set an invalid INT constraint, with lhs : " + lhs)
	}
	// If already a number, no change
	if ctx.IsInt(lhs) {
		return ctx // not needed
	}
	// Otherwise, set it !
	return &PContext{
		prt: ctx,
		ctp: INT,
		lhs: lhs,
		rhs: nil,
	}
}

// IsNumber tells if lhs should be a number.
// Default to false.
func (ctx *PContext) IsNumber(lhs string) bool {
	if lhs == "" {
		return false
	}
	if isNumber(lhs) || ctx.IsInt(lhs) {
		return true
	}
	// actual check
	for c := ctx; c != nil; c = c.prt {
		if c.ctp == NUMBER && c.lhs == lhs {
			return true
		}
	}
	return false
}

func (ctx *PContext) IsInt(lhs string) bool {
	if lhs == "" {
		return false
	}
	if isInt(lhs) {
		return true
	}
	// actual check
	for c := ctx; c != nil; c = c.prt {
		if c.ctp == INT && c.lhs == lhs {
			return true
		}
	}
	return false
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
	dedup := make(map[*Node]bool) // dedup lhs

	for c := ctx; c != nil; c = c.prt {
		if c.lhs != nil && !dedup[c.lhs] {
			dedup[c.lhs] = true
			fmt.Fprintf(&b, "%s %s %s, ", c.lhs, c.ctp, c.rhs)
		}
	}
	return b.String()
}

// rhsGet return the FIRST, VALID lhs which rhs, also returned, is or contains x.
// returns nil, nil if not found.
func (ctx *PContext) rhsGet(ctp CType, x *Node) (lhs *Node, rhs *Node) {

	if x == nil {
		panic("cannot call rhsGet with x==nil")
	}

	dedup := make(map[*Node]bool)

	for c := ctx; c != nil; c = c.prt {
		if dedup[c.lhs] || c.rhs == nil {
			continue
		}
		dedup[c.lhs] = true
		if c.rhs == x || c.rhs.contains(x) {
			return c.lhs, c.rhs
		}
	}
	return nil, nil
}

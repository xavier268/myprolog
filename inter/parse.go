package inter

import (
	"fmt"
	"math"
	"strings"
)

// a Variable is a leafnode that starts with a capital letter or an underscore.
func isVariable(name string) bool {
	return (name[0] >= 'A' && name[0] <= 'Z') || name[0] == '_'
}

// a number is a valid go token starting with a digit.
func isNumber(name string) bool {
	return name[0] >= '0' && name[0] <= '9'
}

// isFunctor to check if valid functor name
func isFunctor(name string) bool {
	return name != "" && !isVariable(name) && !isNumber(name) && !strings.Contains("()_[]|.", name) && name != "nil"
}

func (i *Inter) Parse(tzr Tokenizer, root *Node) error {
	return i.parse0(tzr, root, new(int))
}

// Parse tokens, adding them as children of the provided node.
// Only parenthesis and basic checks are managed.
func (i *Inter) parse0(tzr Tokenizer, root *Node, par *int) error {
	for tk := tzr.Next(); tk != ""; tk = tzr.Next() {
		switch tk {
		case ",": // ignore ','
		case "(":
			*par++
			n := root.lastArg()
			if n == nil {
				return fmt.Errorf("cannot parse an expression starting with '('")
			}
			if len(n.args) != 0 {
				return fmt.Errorf("a compound object cannot be a functor before parenthesis")
			}
			if isVariable(n.name) {
				return fmt.Errorf("the Variable %s cannot be a functor before parenthesis", n.name)
			}
			if isNumber(n.name) {
				return fmt.Errorf("the Number %s cannot be a functor before parenthesis", n.name)
			}
			if !isFunctor(n.name) {
				return fmt.Errorf("the name %s cannot be a functor, appearing before parenthesis", n.name)
			}
			err := i.parse0(tzr, n, par)
			if err != nil {
				return err
			}
		case ")":
			*par--
			if *par < 0 {
				return fmt.Errorf("closing a parenthesis that was not open before")
			}
			return nil
		case "":
			panic("unexpected attempt to parse empty token")
		default:
			n := i.nodeFor(tk)
			root.args = append(root.args, n)
		}
	}

	if *par != 0 { // capture opened, not closed parenthesis
		return fmt.Errorf("parenthesis do not match")
	}
	return nil
}

// isConstant recursively and lazily tells if the tree is constant, ie, contains NO variable.
// Any node flagged constant should be considered immutable and should never change, not be given additionnal children.
// New nodes start as non constant.
func (n *Node) isConstant() bool {

	if n.constant { // already defined as constant
		return true
	}
	if isVariable(n.name) { // Variable are NEVER constant.
		n.constant = false
		return false
	}
	/*if len(n.args) == 0 { // all atomic that are not variable are constant.
		n.constant = true
		return true
	}*/
	for _, a := range n.args { // now, all atomic with NO children will be constant.
		if !a.isConstant() {
			return false // lazily return. Some subtrees might not have been marked.
		}
	}
	n.constant = true
	return true
}

// markConstant walk the tree and recomputes all constant values.
// Needed everytime a variable could gave been added somewhere, or the tree could have changed.
// It returns the constant status.
func (n *Node) markConstant() bool {

	//fmt.Printf("DEBUG - exploring : %s\n", n)

	if isVariable(n.name) { // Variable are NEVER constant.
		n.constant = false
		return false
	}
	if len(n.args) == 0 { // atomic and not variable are ALWAYS constant.
		n.constant = true
		//fmt.Printf("DEBUG - marking constant %s\n", n)
		return true
	}
	/*
		// WARNING : This looks attractive, but the compiler overoptimizes and takes short cuts, not covering all nodes !
			cc := true
			for _, a := range n.args {
				cc = cc && a.markConstant() // compiler optimzes and break on false !!!
			}
			n.constant = cc
			return n.constant
	*/

	// Alternative, using a trick to prevent compiler over-optimization ...
	nbc := 0
	for _, a := range n.args {
		if a.markConstant() {
			nbc++
		} // no early break here !
	}
	n.constant = (nbc == len(n.args))
	return n.constant
}

const MaxInt int = int(math.MaxInt32)

// preProcRule pre-processes rules. It is idempotent.
// It will turn postfix rules into prefix rules, using the "~" functor, and checking rule syntax.
// It handles facts and alternative (semi-colon) rules.
// Rules are supposed to be the children of n. Not recursion to look for them below that.
func (in *Inter) preProcRule(n *Node) error {

	if n == nil || len(n.args) == 0 {
		return nil
	}

	rstart := 0 // points to the first child we want to process
	for {       // manual loop on rstart
		if rstart >= len(n.args) {
			break // done
		}

		// reset rule internal pointers
		rperiod := MaxInt // points to first valid .
		rtilda := MaxInt  // points to first ~
		rsemi := MaxInt   // points to first semi, if before period and after tilda.

		// set rule pointers
		for i := rstart; i < len(n.args); i++ {
			if i < rperiod && n.args[i].name == "." {
				rperiod = i
				break // do not update after !
			}
			if i < rtilda && i <= rstart+1 && n.args[i].name == "~" { // tilda can only appear prefix or postfix.
				rtilda = i
			}
			if i < rsemi && rsemi > rtilda && n.args[i].name == ";" { // tilda required before semi
				rsemi = i
			}
		}

		//fmt.Printf("DEBUG : $=%d/%d, .=%d, ~=%d, ;=%d\n", rstart, len(n.args), rperiod, rtilda, rsemi)

		// check syntax and process
		if rtilda == rstart { // canonical from
			if len(n.args[rtilda].args) != 0 {
				// valid canonical form.
				// Ignore and continue.
				rstart++
				continue
			}
			// invalid canonical form
			return fmt.Errorf("canonical form rule has no head")
		}
		if rperiod == MaxInt {
			return fmt.Errorf("rule is missing the final period")
		}
		if len(n.args[rperiod].args) != 0 {
			return fmt.Errorf("the period cannot be a functor")
		}
		if rperiod == rstart {
			return fmt.Errorf("empty fact rule")
		}

		if rperiod == rstart+1 { // fact
			if isVariable(n.args[rstart].name) { // invalid fact
				return fmt.Errorf("head of rule cannot be a Variable")
			}
			// construct actual rule, reusing the period node.
			n.args[rperiod].name = "~"
			n.args[rperiod].args = append(n.args[rperiod].args, n.args[rstart])
			// remove head pointer
			n.args = append(n.args[:rstart], n.args[rperiod:]...)
			// proceed.
			rstart++ // jump to the node following the period.
			continue
		}

		if rtilda == rstart+1 && rsemi == MaxInt { // postfix rule (no alternative).
			head := n.args[rstart]
			n.args[rtilda].args = append(n.args[rtilda].args, head)
			n.args[rtilda].args = append(n.args[rtilda].args, n.args[rtilda+1:rperiod]...)

			// cleanup
			//fmt.Println("DEBUG : n.args before cleanup:", n.args)
			if rperiod < len(n.args) {
				n.args = append(n.args[:rtilda+1], n.args[rperiod+1:]...)
			} else {
				n.args = n.args[:rtilda+1]
			}
			n.args = append(n.args[:rstart], n.args[rtilda:]...)
			//fmt.Println("DEBUG : n.args after cleanup:", n.args)
			rstart++
			continue
		}

		if rtilda == rstart+1 && rsemi != MaxInt { // postfix rule (with alternative).
			head := n.args[rstart]
			n.args[rtilda].args = append(n.args[rtilda].args, head)
			n.args[rtilda].args = append(n.args[rtilda].args, n.args[rtilda+1:rsemi]...)

			// cleanup
			//fmt.Println("DEBUG : n.args before cleanup:", n.args)
			n.args[rstart], n.args[rtilda] = n.args[rtilda], n.args[rstart]
			n.args[rsemi].name = "~"
			n.args = append(n.args[:rtilda+1], n.args[rsemi:]...)
			//fmt.Println("DEBUG : n.args after cleanup:", n.args)
			rstart++
			continue
		}

		return fmt.Errorf("unknown syntax")
	}
	return nil
}

// preProcBar recursively pre-processes bar-formed lists, transforming non canonical forms into canonical forms.
// The canonical form for list of a,b and c uses the dot operator, as in :
// dot(a dot(b dot(c)))//
// The bar form is :
// a | RestOfList or a | nil
// The | operator is just the postfix version of dot/2.
func (in *Inter) preProcBar(n *Node) error {

	if n == nil || len(n.args) == 0 {
		return nil // no child, done.
	}
	//n has children ...
	if !isFunctor(n.name) {
		return fmt.Errorf("%s is not a valid functor", n.name)
	}
	i := 0
	for {
		// manual loop on n.args
		if i >= len(n.args) {
			break
		}

		// recurse lower level
		err := in.preProcBar(n.args[i])
		if err != nil {
			return err
		}

		if n.args[i].name == "|" {
			if i-1 >= 0 && i+1 <= len(n.args)-1 { // do we have a child before and after ?
				n.args[i].name = "dot"                             // reuse | node as dot node
				n.args[i].args = []*Node{n.args[i-1], n.args[i+1]} // copy children
				// cleanup
				n.args[i-1], n.args[i] = n.args[i], n.args[i-1] // swap to prefix
				if i+2 < len(n.args) {
					n.args = append(n.args[:i], n.args[i+2:]...)
				} else {
					n.args = n.args[:i]
				}
				// no need to increment. i is now pointing on the next node already.
				continue
			}
			// no child before and after ...
			return fmt.Errorf("the | operator should have exactly 2 postfix arguments")
		}
		i++
	}
	return nil
}

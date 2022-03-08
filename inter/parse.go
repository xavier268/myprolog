package inter

import (
	"fmt"
)

type buildin struct { // descriptor for buildin symbol or operator
	// arity of the symbol, < 0 means it can have any arity
	arity int
	//symbol is postfixed if it should appear immediately after the first Atom, as in :
	// 3 symbol 4 5
	// Symbol is not postfixed (default) if it appears as in :
	// symbol(3 4 5)
	postfix bool
}

// Definition of 'built in' operators and symbols.
var BuildIns = map[string]buildin{
	":-": {-1, true},
	"!=": {2, true},
	"=":  {2, true},
	"/":  {0, false},
}

// a Variable is a leafnode that starts with a capital letter or is just an underscore.
func isVariable(name string) bool {
	return (name[0] >= 'A' && name[0] <= 'Z') || name == "_"
}

// a number is a valid go token starting with a digit.
func isNumber(name string) bool {
	return name[0] >= '0' && name[0] <= '9'
}

func (i *Inter) Parse(tzr Tokenizer, root *Node) error {
	return i.parse0(tzr, root, 0, new(int))
}

// Parse tokens, adding them as children of the provided node.
// Only parenthesis and basic checks are managed.
func (i *Inter) parse0(tzr Tokenizer, root *Node, level int, par *int) error {
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
			err := i.parse0(tzr, n, level+1, par)
			if err != nil {
				return err
			}
		case ")":
			*par--
			if *par < 0 {
				return fmt.Errorf("closing a parenthesis too early")
			}
			if level == 0 {
				return fmt.Errorf("closing a parenthesis that was never opened")
			}
			return nil
		case "":
			panic("unexpected attempt to parse empty token")
		default:
			n := i.nodeFor(tk)
			root.args = append(root.args, n)
		}
	}

	if level == 0 && *par != 0 {
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

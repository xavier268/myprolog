package inter

import (
	"fmt"
	"math"
)

// a Variable starts with a capital letter or an underscore.
func isVariable(name string) bool {
	return (name[0] >= 'A' && name[0] <= 'Z') || name[0] == '_'
}

// a number is a valid go token starting with a digit.
func isNumber(name string) bool {
	return name[0] >= '0' && name[0] <= '9'
}

// isFunctor to check if valid functor name.
// It should be neither a variable, nor a number, nor named 'nil'
func isFunctor(name string) bool {
	return name != "" && !isVariable(name) && !isNumber(name) && name != "nil"
}

// ParseRules will load the rules (or additionnal rules) to the rules node.
// New rules are appended to the previous ones.
func (in *Inter) ParseRules(tzr Tokenizer) error {
	err := in.parse0(tzr, in.rules, new(int))
	if err != nil {
		return err
	}
	err = in.preProcList(in.rules)
	if err != nil {
		return err
	}
	err = in.preProcRule(in.rules)
	if err != nil {
		return err
	}
	return nil
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
			n := nodeFor(tk)
			root.args = append(root.args, n)
		}
	}

	if *par != 0 { // capture opened, not closed parenthesis
		return fmt.Errorf("parenthesis do not match")
	}
	return nil
}

const MaxInt int = int(math.MaxInt32)

// preProcRule pre-processes rules. It is idempotent.
// It will turn postfix rules into prefix rules, using the ":-" functor, and checking rule syntax.
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
		rarrow := MaxInt  // points to first :-
		rsemi := MaxInt   // points to first semi, if before period and after tilda.

		// set rule pointers
		for i := rstart; i < len(n.args); i++ {
			if i < rperiod && n.args[i].name == "." {
				rperiod = i
				break // do not update after !
			}
			if rarrow != MaxInt && n.args[i].name == ":-" {
				return fmt.Errorf("there can only be one arrow :- per valid rule. Did you forget a period ?")
			}
			if i < rarrow && i <= rstart+1 && n.args[i].name == ":-" { // arrow can only appear prefix or postfix.
				rarrow = i
			}
			if i < rsemi && rsemi > rarrow && n.args[i].name == ";" { // arrow required before semi
				rsemi = i
			}

		}

		// check syntax and process
		if rarrow == rstart { // canonical from
			if len(n.args[rarrow].args) != 0 {
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
				return fmt.Errorf("a Variable is not a rul on its own")
			}
			// construct actual rule, reusing the period node.
			n.args[rperiod].name = ":-"
			n.args[rperiod].args = append(n.args[rperiod].args, n.args[rstart])
			// remove head pointer
			n.args = append(n.args[:rstart], n.args[rperiod:]...)
			// proceed.
			rstart++ // jump to the node following the period.
			continue
		}

		if rarrow == rstart+1 && rsemi == MaxInt { // postfix rule (no alternative).
			head := n.args[rstart]
			n.args[rarrow].args = append(n.args[rarrow].args, head)
			n.args[rarrow].args = append(n.args[rarrow].args, n.args[rarrow+1:rperiod]...)

			// cleanup
			//fmt.Println("DEBUG : n.args before cleanup:", n.args)
			if rperiod < len(n.args) {
				n.args = append(n.args[:rarrow+1], n.args[rperiod+1:]...)
			} else {
				n.args = n.args[:rarrow+1]
			}
			n.args = append(n.args[:rstart], n.args[rarrow:]...)
			//fmt.Println("DEBUG : n.args after cleanup:", n.args)
			rstart++
			continue
		}

		if rarrow == rstart+1 && rsemi != MaxInt { // postfix rule (with alternative).
			head := n.args[rstart]
			n.args[rarrow].args = append(n.args[rarrow].args, head)
			n.args[rarrow].args = append(n.args[rarrow].args, n.args[rarrow+1:rsemi]...)

			// cleanup
			//fmt.Println("DEBUG : n.args before cleanup:", n.args)
			n.args[rstart], n.args[rarrow] = n.args[rarrow], n.args[rstart]
			n.args[rsemi].name = ":-"
			n.args = append(n.args[:rarrow+1], n.args[rsemi:]...)
			//fmt.Println("DEBUG : n.args after cleanup:", n.args)
			rstart++
			continue
		}

		return fmt.Errorf("unknown syntax")
	}
	return nil
}

// preProcList recursively pre-processes bracket lists, transforming non canonical forms into canonical forms.
// The canonical form for list of a,b and c uses the dot operator, as in :
// dot(a dot(b dot(c)))//
// The bracket form is :
// [ a b c ] or [ a | [ b c ]]
func (in *Inter) preProcList(n *Node) error {

	if n == nil || len(n.args) == 0 {
		return nil // nothing to do
	}
	// Now, n has children.
	if !isFunctor(n.name) {
		return fmt.Errorf("%s cannot have children, it is not a valid functor", n.name)
	}

	// while loop until all lists are handled
	for {
		// find latest open
		open, close, bar := -1, -1, -1
		for i, a := range n.args {
			if a.name == "[" {
				open = i
			}
			if open < 0 && a.name == "]" {
				return fmt.Errorf("missing opening bracket before closing")
			}
			if open < 0 && a.name == "|" {
				return fmt.Errorf("the | symbol must be enclosed in brackets")
			}
		}
		if open < 0 {
			break // no more open, and no hanging close or bar.
		}
		// find earliest close and bar AFTER open
		for i := open; i < len(n.args); i++ {
			a := n.args[i]
			if a.name == "|" {
				if bar < 0 {
					bar = i
				} else {
					return fmt.Errorf("illegal multiple | in the same bracket list")
				}
			}
			if a.name == "]" && close < 0 {
				close = i
				break
			}
		}
		if close < 0 {
			return fmt.Errorf("missing closing bracket")
		}
		if bar > close {
			bar = -1
		}

		// now, open, close and bar are valid and consitent.
		if bar < 0 { // standard bracket list, [ a b c ] with no bar to worry about

			// reuse the open node
			list := n.args[open]
			list.name = "dot"
			list.args = []*Node{nodeFor("nil"), nodeFor("nil")}

			for p := open + 1; n.args[p].name != "]"; p++ { // iterate on the inner list
				nn := nodeFor("dot")
				nn.args = []*Node{nodeFor("nil"), nodeFor("nil")}
				list.args = []*Node{n.args[p], nn}
				list = list.args[1]
			}

			// cleanup - suppressing n.args nodes from open+1 included to close included.
			if close+1 < len(n.args) {
				n.args = append(n.args[:open+1], n.args[close+1:]...)
			} else {
				n.args = n.args[:open+1]
			}
		}

		if bar > 0 { // bracket list in the form [ a | b ]
			// check syntax
			if bar-open != 2 || close-bar != 2 {
				return fmt.Errorf("wrong number of arguments for the [x|y] operator : %s %s %s", n.args[bar-1].name, "|", n.args[bar+1].name)
			}
			// reuse the open node
			n.args[open].name = "dot"
			n.args[open].args = []*Node{n.args[open+1], n.args[open+3]}
			// cleanup : suppress from open excluded to close included
			if close+1 < len(n.args) {
				n.args = append(n.args[:open+1], n.args[close+1:]...)
			} else {
				n.args = n.args[:open+1]
			}
		}
		// go search the next list !
	}

	// now, recurse on children ...
	for _, a := range n.args {
		err := in.preProcList(a)
		if err != nil {
			return err
		}
	}

	return nil
}

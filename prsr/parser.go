package prsr

import (
	"fmt"

	"github.com/xavier268/myprolog/node"
	"github.com/xavier268/myprolog/tknz"
)

// Parse tokens, adding them as children of the provided node.
// Only parenthesis and basic checks are managed.
func parse0(tzr *tknz.Tokenizer, root *node.Node, par *int) error {
	for tk := tzr.Next(); tk != ""; tk = tzr.Next() {
		switch tk {
		case ",": // ignore ','
		case "(":
			*par++
			n := root.LastChild()
			if n != nil && n.ChildrenAllowed() {
				err := parse0(tzr, n, par)
				if err != nil {
					return err
				}
			} else { // n == nil or not a functor !
				return fmt.Errorf("cannot add children to such a node : %#v", n)
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
			root.Add(node.NewNode(tk))
		}
	}

	if *par != 0 {
		return fmt.Errorf("parenthesis do not match")
	}
	return nil
}

// preProcList recursively pre-processes bracket lists, transforming non canonical forms into canonical forms.
// The canonical form for list of a,b and c uses the dot operator, as in :
// dot(a dot(b dot(c)))//
// The bracket form is :
// [ a b c ] or [ a | [ b c ]]
func preProcList(n *node.Node) error {

	if n == nil || n.NbChildren() == 0 {
		return nil // nothing to do
	}
	// Now, n has children.
	if !n.ChildrenAllowed() {
		return fmt.Errorf("Node cannot have children, it is not a valid functor, %#v", n)
	}

	// while loop until all lists are handled
	for {
		// find latest open
		open, close, bar := -1, -1, -1
		for i, a := range n.GetChildren() {
			if a.GetLoad() == node.String("[") {
				open = i
			}
			if open < 0 && a.GetLoad() == node.String("]") {
				return fmt.Errorf("missing opening bracket before closing")
			}
			if open < 0 && a.GetLoad() == node.String("|") {
				return fmt.Errorf("the | symbol must be enclosed in brackets")
			}
		}
		if open < 0 {
			break // no more open, and no hanging close or bar.
		}
		// find earliest close and bar AFTER open
		for i := open; i < n.NbChildren(); i++ {
			a := n.GetChild(i)
			if a.GetLoad() == node.String("|") {
				if bar < 0 {
					bar = i
				} else {
					return fmt.Errorf("illegal multiple | in the same bracket list")
				}
			}
			if a.GetLoad() == node.String("]") && close < 0 {
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

		// now ... open, close and bar are valid and consitent.
		if bar < 0 { // standard bracket list, [ a b c ] with no bar to worry about

			// reuse the open node
			list := n.GetChild(open)
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

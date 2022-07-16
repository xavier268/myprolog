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
		return fmt.Errorf("node cannot have children, it is not a valid functor, %#v", n)
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
		fmt.Println("Debug : open, close, bar =", open, close, bar)

		if bar < 0 { // standard bracket list, [ a b c ] with no bar to worry about

			// replace the open node
			dot := node.NewDotNode()
			dot.Add(nil, nil)
			n.ReplaceChild(open, dot)

			for n.GetChild(open+1).GetLoad() != node.String("]") { // iterate on the inner list
				nn := node.NewDotNode()
				nn.Add(nil, nil)
				dot.ReplaceChild(0, n.GetChild(open+1))
				dot.ReplaceChild(1, nn)
				dot = dot.GetChild(1)
				n.RemoveChild(open + 1) // cleanup - suppressing n.args nodes from open+1 included to close excluded.
			}
			n.RemoveChild(open + 1) // clean close
		}

		if bar > 0 { // bracket list in the form [ a | b ]
			// check syntax
			if bar-open != 2 || close-bar != 2 {
				return fmt.Errorf("wrong number of arguments for the [x|y] operator : %s %s %s", n.GetChild(bar-1).GetLoad(), "|", n.GetChild(bar+1).GetLoad())
			}

			// handle the first part
			dot := node.NewDotNode()
			dot.Add(n.GetChild(open+1), n.GetChild(open+3))
			n.RemoveChild(open + 1)   // remove car
			n.RemoveChild(open + 1)   // remove |
			n.RemoveChild(open + 1)   // remove cdr
			n.RemoveChild(open + 1)   // remove ]
			n.ReplaceChild(open, dot) // replace [ with dot-car, cdr)

		}
	}

	// now, recurse on children ...
	for _, a := range n.GetChildren() {
		err := preProcList(a)
		if err != nil {
			return err
		}
	}

	return nil
}

// preProcRule pre-processes rules and queries. It is idempotent.
// It will turn infix rules into prefix rules, using the ":-" functor, changed into the 'rule' keyword, and checking rule syntax.
// It handles facts and alternative (semi-colon) rules.
// Rules are supposed to be the children of input Node. Not recursion to look for them below that.
// Queries start with a question mark (?) and finish with a period (.).
func preProcRule(n *node.Node) error {

	if n == nil || n.NbChildren() == 0 {
		return nil
	}

	// Define useful  constants and references
	const large = 100_000_000 // less than max signed int 32
	qk, err := node.NewKeyword("query")
	if err != nil {
		panic("internal error")
	}
	rk, err := node.NewKeyword("rule")
	if err != nil {
		panic("internal error")
	}

	rstart := 0 // points to the first child we want to process
	for {       // manual loop on rstart
		if rstart >= n.NbChildren() {
			break // done
		}

		// reset rule internal pointers
		rperiod := large // points to first valid .
		rarrow := large  // points to first :-
		rsemi := large   // points to first semi, if before period and after tilda.
		rquery := large

		// set rule pointers
		for i := rstart; i < n.NbChildren(); i++ {

			// skip canonical query/rule forms
			if n.GetChild(rstart).GetLoad() == rk || n.GetChild(rstart).GetLoad() == qk {
				rstart++ // skip prexisting queries or rules allready processed.
				continue
			}

			if rquery == large && rperiod == large && n.GetChild(i).GetLoad() == node.String("?") {
				rquery = i
			}
			if i < rperiod && n.GetChild(i).GetLoad() == node.String(".") {
				rperiod = i
				break // do not update further beyond the first period !
			}
			if rarrow != large && n.GetChild(i).GetLoad() == node.String(":-") {
				return fmt.Errorf("there can only be one arrow :- per valid rule. Did you forget a period ?")
			}
			if i < rarrow && i <= rstart+1 && n.GetChild(i).GetLoad() == node.String(":-") { // arrow can only appear prefix or postfix.
				rarrow = i
			}
			if i < rsemi && rsemi > rarrow && n.GetChild(i).GetLoad() == node.String(";") { // arrow required before semi
				rsemi = i
			}
		}
		fmt.Println("DEBUG : ", n)
		fmt.Println("DEBUG : rstart, rquery, rarrow, rsemi, rperiod:", rstart, rquery, rarrow, rsemi, rperiod)

		// process query if needed
		if rquery != large {
			if rquery != rstart {
				return fmt.Errorf("missing a period before query ?")
			}
			if rperiod == large {
				return fmt.Errorf("missing a period after query ?")
			}
			nq := node.NewQueryNode()
			for n.GetChild(rstart+1).GetLoad() != node.String(".") {
				nq.Add(n.GetChild(rstart + 1))
				n.RemoveChild(rstart + 1)
			}
			n.RemoveChild(rstart + 1) // remove "."
			n.ReplaceChild(rstart, nq)

			rstart++
			continue
		}

		// check syntax and process rule
		if rarrow == rstart { // canonical form
			if n.GetChild(rarrow).NbChildren() != 0 {
				// valid canonical form.
				// Change :- into 'rule' and continue
				rule := node.NewRuleNode()
				rule.Add(n.GetChild(rstart).GetChildren()...)
				n.ReplaceChild(rstart, rule)
				rstart++
				continue
			}
			// invalid canonical form
			return fmt.Errorf("canonical form rule has no head")
		}
		if rperiod == large {
			return fmt.Errorf("rule is missing the final period")
		}
		if n.GetChild(rperiod).NbChildren() != 0 {
			return fmt.Errorf("the period cannot be a functor")
		}
		if rperiod == rstart {
			return fmt.Errorf("empty fact rule")
		}

		if rperiod == rstart+1 { // single fact, with period.
			// invalid fact
			if _, ok := n.GetChild(rstart).GetLoad().(node.Variable); ok {
				return fmt.Errorf("a Variable is not a rule on its own")
			}
			if _, ok := n.GetChild(rstart).GetLoad().(node.Underscore); ok {
				return fmt.Errorf("a Underscore is not a rule on its own")
			}
			if _, ok := n.GetChild(rstart).GetLoad().(node.Number); ok {
				return fmt.Errorf("a Number is not a rule on its own")
			}
			// construct actual rule, reusing the period node.
			rule := node.NewRuleNode()
			rule.Add(n.GetChild(rstart))
			n.ReplaceChild(rstart, rule)
			n.RemoveChild(rperiod)

			// proceed.
			rstart++ // jump to the node following the period.
			continue
		}

		if rarrow == rstart+1 && rsemi == large { // postfix rule (no alternative).
			rule := node.NewRuleNode()
			rule.Add(n.GetChild(rstart))
			for n.GetChild(rarrow+1).GetLoad() != node.String(".") {
				rule.Add(n.GetChild(rarrow + 1))
				n.RemoveChild(rarrow + 1)
			}
			n.RemoveChild(rarrow + 1)
			n.RemoveChild(rarrow)
			n.ReplaceChild(rstart, rule)

			rstart++
			continue
		}

		if rarrow == rstart+1 && rsemi != large { // postfix rule (with alternative).

			rule := node.NewRuleNode()
			rule.Add(n.GetChild(rstart))
			for n.GetChild(rarrow+1).GetLoad() != node.String(";") {
				rule.Add(n.GetChild(rarrow + 1))
				n.RemoveChild(rarrow + 1)
			}

			n.ReplaceChild(rstart+2, node.NewStringNode(":-"))
			n.ReplaceChild(rstart+1, n.GetChild(rstart))
			n.ReplaceChild(rstart, rule)

			rstart++
			continue
		}

		return fmt.Errorf("unknown syntax")
	}
	return nil
}

// Parse a Prolog program.
// Valid syntax is a set of valid rules.
// Root node is called program.
func Parse(tk *tknz.Tokenizer) (*node.Node, error) {
	p := node.NewProgramNode()
	err := ParseMore(p, tk)
	return p, err
}

// Update the provided program with the rules and queries parsed from tokenizer.
func ParseMore(p *node.Node, tk *tknz.Tokenizer) error {
	err := parse0(tk, p, new(int))
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = preProcList(p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = preProcRule(p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

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

// Parse tokens, adding them as children of the provided node.
func (i *Inter) Parse(tzr Tokenizer, root *Node) error {
	return i.parse(tzr, root, 0, new(int))
}

func (i *Inter) parse(tzr Tokenizer, root *Node, level int, par *int) error {
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
			err := i.parse(tzr, n, level+1, par)
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

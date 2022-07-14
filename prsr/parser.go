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
			n := root.LastArg()
			err := parse0(tzr, n, par)
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
			root.Add(node.NewNode(tk))
		}
	}

	if *par != 0 {
		return fmt.Errorf("parenthesis do not match")
	}
	return nil
}

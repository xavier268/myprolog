// TODO - WARNING - SIMPLIFY NEEDS TO BE ABLE TO CALL UNIFY !!!
// Reconsider pulling Unify xx into the cons package ?
package cons

import (
	"fmt"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/node"
)

// Simplify a given list of constraints.
// Simplification means :
// 	* remove duplicates
//	* substiture k,own variables by their definitions
// It is assumed that, individually, each constraint is already both Valid ans Relevant.
// If trees are modified, they are always cloned first.

func Simplify(oldlist []Cons) (newlist []Cons, changed bool, err error) {

	changed = false
	err = nil
	newlist = append(newlist, oldlist...)

	if config.FlagDebug {
		fmt.Println("DEBUG SIMPLIFY CONS : ", oldlist)
	}

	loop := true
	for loop {
		loop = false

		// Look for duplicates X = .. and X = ...
		for i := 0; i < len(newlist); i++ {
			for j := i + 1; j < len(newlist); j++ {

				if newlist[i].variable == newlist[j].variable &&
					newlist[i].relation == newlist[j].relation &&
					newlist[i].tree.Equal(newlist[j].tree) {
					// remove j
					if j+1 < len(newlist) {
						newlist = append(newlist[:j], newlist[j+1:]...)
					} else {
						newlist = newlist[:j]
					}
					j-- // look again same index, different content
					changed = true
					loop = true
				}
				// no change
			}
		}

		// replace X by its value in right hand sides, whatever the relation
		for i := range newlist {
			for j := range newlist {

				newtree, repl := newlist[j].tree.ReplaceSubTree(node.NewVariableNode(newlist[i].variable), newlist[i].tree)
				if repl {
					if i == j {
						return nil, false, fmt.Errorf("positive error check") // X = f(... X ...) - never replace !
					}
					changed = true
					loop = true
					newlist[j].tree = newtree
				}
			}
		}

	}
	if config.FlagDebug {
		fmt.Println("DEBUG SIMPLIFY - new, chg, err =", newlist, changed, err)
	}
	if changed {
		return newlist, changed, nil
	} else {
		newlist = nil
		return oldlist, changed, nil
	}

}

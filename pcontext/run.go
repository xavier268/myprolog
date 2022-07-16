package pcontext

import (
	"fmt"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/node"
)

// Run  the pcontext, until all goals are erased or an error has occured.
// Only rule selection trigger a context push/pop.
func (pc *PContext) Run() (*PContext, error) {

	for {

		if config.FlagDebug {
			fmt.Println("DEBUG RUN : ", pc)
		}

		if pc == nil {
			return nil, nil
		}
		if len(pc.goals) == 0 { // done, but we keep the context for further solutions
			return pc, nil
		}

		goal := pc.goals[len(pc.goals)-1]
		for goal == nil { // trim nil goals
			pc.goals = pc.goals[:len(pc.goals)-1]
			if len(pc.goals) > 0 {
				goal = pc.goals[len(pc.goals)-1]
			} else { // done, no more goals, but we keep the context for further solutions
				return pc, nil
			}
		}

		rule := pc.FindRule(goal)
		if rule != nil {
			// push context, and try to confirm this rule
			pc = pc.Push()
			err := pc.Unify(rule.GetChild(0), goal)
			if err == nil { // success !
				// update goals, erasing the initial goal used, and adding the body of the localized rule.
				if len(pc.goals) > 1 {
					pc.goals = pc.goals[:len(pc.goals)-2]
				} else {
					pc.goals = pc.goals[:0]
				}
				if rule.NbChildren() > 1 {
					pc.goals = append(pc.goals, rule.GetChildren()[1:]...)
				}

				continue // loop with the new context.
			}
			pc = pc.Pop() // on failure, pop context and continue
			continue

		} else { // no rule matches - verify if builtin applies ?
			gg, err := pc.DoBuiltin(goal)
			if err != nil {
				// backtracking required
				return pc, err
			}
			// update the new goals
			if gg == nil {
				pc.goals = pc.goals[:len(pc.goals)-1]
			} else {
				pc.goals = append(pc.goals[:len(pc.goals)-1], gg...)
			}
			continue
		}

	}
}

// Try to find an applicable rule, starting from the current rule pointer.
// If found, return a localized version and update the current rule pointer of the context.
// If no found, return nil. It is just a quickcheck, that will have to be confirmed later.
// No change to the context at this stage.
func (pc PContext) FindRule(goal *node.Node) (localizedRule *node.Node) {
	//found := -1 // found rule
	for i := pc.current; i < len(pc.rules); i++ {
		head := pc.rules[i].GetChild(0)                                                 // rule head
		if head.GetLoad() == goal.GetLoad() && head.NbChildren() == goal.NbChildren() { // found a candidate
			// Because neither can be a Variable here, it is ok to just compare payload and arity.
			pc.current = i // update curent, so that next time the next rule will be selected.
			return pc.rules[i].CloneLocal(pc.UID())
		}
	}
	return nil
}

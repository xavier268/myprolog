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
			// attempt to use this rule !
			// TODO - push a new context and attempt to use this rule.
			// If it fails, pop the context back.
		} else {
			// verify if builtin applies ?
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
// If found, retun a localized version, and update the current rule pointer.
// If no found, return nil.
func (pc PContext) FindRule(goal *node.Node) (localizedRule *node.Node) {
	// TODO
	fmt.Println("WARNING : FindRule not implemented - TODO")
	return nil
}

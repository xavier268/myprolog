package solver

import "fmt"

// Known compound predicates with their imperative arity.
// A compund predicate MUST be a CompoundTerm.
// Arity set to -1 means it can vary.
// Assumes only canonical form.
// Constant values are eliminated (numbers, strings ...) as well as underscore.
var CompPredicateMap = map[string]int{
	"rule":  -1, // rule definition
	"query": -1, // query definition
	"dot":   -1, // actually, between 0 and 2 - special check enforced in code
	"and":   2,  // both children must erase
	"or":    2,  // at least one child must erase
}

// Known atomic predicates.
var AtomPredicate = map[string]bool{
	"!": true, // the 'cut' predicate will prevent backtracking from now on.
}

// Execute predicates, if possible, using the functor of the first goal.
// Includes removing underscore. Constant values (numbers, strings ...) are kept.
// State MAY change, if  alternative must be considered.
func DoPredicate(st *State) (*State, error) {
	for {
		if len(st.Goals) == 0 {
			return st, nil
		}
		goal := st.Goals[0]
		switch g := goal.(type) {
		case Underscore: // remove underscore.
			st.Goals = st.Goals[1:]
			st.NextRule = 0 // when goal change, reset the next rule pointer ...
			continue
		case Variable: // Variables are not predicates. Leave as is.
			return st, nil
		case Number, String: // keep number and constants
			return st, nil
		case Atom:
			if !AtomPredicate[g.Value] { // not a predicate, leave as is.
				return st, nil
			}
			switch g.Value {
			case "!": // cut
				st.Parent = nil // cut, so no parent
				st.Goals = st.Goals[1:]
				st.NextRule = 0 // when goal change, reset the next rule pointer ...
				continue
			default: // should not happen
				panic("internal error : unknown Atom predicate : " + g.Value)
			}
		case CompoundTerm:
			functor := g.Functor
			arity, ok := CompPredicateMap[functor]
			if !ok { // not in the predicate map, ignore it
				return st, nil
			}
			if arity > 0 && arity != len(g.Children) { // arity check fails
				return st, fmt.Errorf("predicate arity check failed: %s", g.String())
			}
			if functor == "dot" && len(g.Children) > 2 { // special arity check for dot
				return st, fmt.Errorf("'dot' arity check failed: %s", g.String())
			}
			switch functor {
			case "rule":
				continue // don't handle rule. As a predicate, it has no effect.
			case "query": // remove query, and add its children as goals instead.
				st.Goals = append(g.Children, st.Goals[1:]...)
				st.NextRule = 0 // when goal change, reset the next rule pointer ...
				continue
			case "dot":
				continue // don't handle dot yet.
			case "and": // remove and add children as goals.
				st.Goals = append(g.Children, st.Goals[1:]...)
				st.NextRule = 0 // when goal change, reset the next rule pointer ...
				continue
			case "or": // fork state.
				// st becomes parent state and  will handle the second alternative,
				// new state handles the first and becomes current.
				nst := NewState(st)

				st.Goals = append(g.Children[1:2], st.Goals[1:]...)
				nst.Goals = append(g.Children[0:1], st.Goals[1:]...)

				st.NextRule = 0 // when goal change, reset the next rule pointer ...
				nst.NextRule = 0
				return nst, nil

			default: // should be unreacheable ...
				panic("internal error for predicate : " + g.String())
			}

		default: // should be unreacheable ...
			panic("internal error : should be unreacheable type")
		}
	}
}

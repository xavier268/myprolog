package solver

import (
	"fmt"
	"strings"

	"github.com/xavier268/myprolog/parser"
)

// this file contains stic global variables to store persistent data for the lifetime of the application.

// Global varial containing the persisent (memory) database : rules, statistiques, ...
var MYDB = Database{
	rules:      make([]CompoundTerm, 0, 10),
	depthStats: 0,
}

// Set of rules that can be applied in a state
type Database struct {
	rules      []CompoundTerm // known usable, non deduplicated, rules.
	depthStats int            // count of truncated searches
}

// Count known rules in DB
func CountDBRules() int {
	return len(MYDB.rules)
}

// Reset the global DB
func ResetDB() {
	MYDB = Database{
		rules:      make([]CompoundTerm, 0, 10),
		depthStats: 0,
	}
}

func ListDBRules() string {
	sb := new(strings.Builder)
	for _, rule := range MYDB.rules {
		fmt.Fprintln(sb, rule.Pretty())
	}
	return sb.String()
}

func AddDBRule(rules ...CompoundTerm) {
	if len(rules) == 0 {
		return
	}
	for _, r := range rules {
		if r.Functor != "rule" {
			panic("Trying to add a Term that is not a rule")
		}
		if len(r.Children) == 0 {
			fmt.Println(parser.RED, "WARNING : trying to add a rule with no children - ignored", parser.RESET)
			return
		}
		MYDB.rules = append(MYDB.rules, r)
	}
}

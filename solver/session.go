package solver

import (
	"fmt"
	"strings"

	"github.com/xavier268/myprolog/parser"
)

// The Session object contains rules and various session constants.

// Set of rules that can be applied in a state
type Session struct {
	rules      []CompoundTerm // known usable, non deduplicated, rules.
	depthStats int            // count of truncated searches
}

func NewSession() *Session {
	db := new(Session)
	db.ResetSession()
	return db
}

// Count known rules in DB
func (sess Session) CountRules() int {
	return len(sess.rules)
}

// Reset the global session
func (sess *Session) ResetSession() {
	sess.rules = make([]CompoundTerm, 0, 10)
	sess.depthStats = 0
}

func (sess Session) ListRules() string {
	sb := new(strings.Builder)
	for i, rule := range sess.rules {
		fmt.Fprintf(sb, "rule#%d>\t%s", i+1, rule.Pretty())
		if i < len(sess.rules)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// AddRule, no dedup.
func (sess *Session) AddRule(rules ...CompoundTerm) {
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
		sess.rules = append(sess.rules, r)
	}
}

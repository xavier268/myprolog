package node

import (
	"fmt"
	"strings"

	"github.com/xavier268/myprolog/config"
)

type KeywordHelp struct {
	arity int    // required nb of children, -1 for any
	short string // short, one line description
	long  string // more lines of description
}

func (kwh KeywordHelp) Arity() int {
	return kwh.arity
}

// contains reserved keywords, and their arity. Arity of -1 means undefined.
var Keywords = map[string]KeywordHelp{
	"rule":    {-1, "defines a rule in a program.", ""},
	"query":   {-1, "defines a query in a program.", ""},
	"program": {-1, "defines the top level node for a program.", ""},
	"freeze":  {0, "prevents backtracking from here.", ""},
	"halt":    {0, "pauses execution.", "Execution will resume upon user request."},
	"exit":    {0, "immediately terminates program execution.", ""},
	"print":   {-1, "prints the string representation of each of the children node.", "Parameters can be anything, including complex nodes. Variable X will appear as 'X', bound values are not used."},
	"println": {-1, "prints the string representation of each of the children node.", "Parameters can be anything, including complex nodes. Variable X will appear as 'X', bound values are not used."},
	"load":    {-1, "loads new rules from files.", "Successively load the provided files. Parameters should only be strings. Only rules are loaded, queries are ignored."},
	"verbose": {1, "set verbose flag.", "Parameter should be 0 or 1. Any other value has no effect."},
	"debug":   {1, "set debug flag.", "Parameter should be 0 or 1. Any other value has no effect."},
	"help":    {-1, "displays help information.", "If no parameter, it displays full help. If string or keyword parameters are provided, it displays detailed help."},
	"list":    {0, "list current rules.", ""},
	"unknown": {-1, "has no effect.", "However, because unknown is a keyword, although it has no effect, it will never bind to anything. An unknown goal cannot be erased."},
}

type Keyword struct {
	name string
}

func (kw *Keyword) String() string {
	if kw == nil {
		return " nil"
	}
	return kw.name
}

// Get help on a keywoard. No-op if not a keyword.
func HelpKeyword(kw string) (string, error) {

	var sb strings.Builder
	h, ok := Keywords[kw]
	if !ok {
		return "", errNotFound
	}
	a := fmt.Sprint(h.arity)
	if h.arity < 0 {
		a = "any"
	}
	fmt.Fprintf(&sb, "%s(%s) is a keyword that %s\n", kw, a, h.short)
	lines := config.CutString(h.long, config.TEXT_WIDTH)
	for _, line := range lines {
		if len(line) > 0 {
			fmt.Fprintln(&sb, config.TEXT_PREFIX, line)
		}
	}
	return sb.String(), nil
}

var errNotFound = fmt.Errorf("keyword not found")

func NewKeyword(kw string) (Keyword, error) {
	_, ok := Keywords[kw]
	if !ok {
		return Keyword{}, errNotFound
	}
	return Keyword{kw}, nil
}

func NewDotNode() *Node {
	n := new(Node)
	n.load = String("dot")
	return n
}

func NewRuleNode() *Node {
	n := new(Node)
	n.load = Keyword{"rule"}
	return n
}
func NewSlashNode() *Node {
	n := new(Node)
	n.load = Keyword{"freeze"}
	return n
}

func NewProgramNode() *Node {
	n := new(Node)
	n.load = Keyword{"program"}
	return n
}

func NewQueryNode() *Node {
	n := new(Node)
	n.load = Keyword{"query"}
	return n
}

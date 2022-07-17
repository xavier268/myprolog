package node

import "fmt"

// contains reserved keywords
var keywords = []string{
	"rule",    // define a rule in a program
	"query",   // define a query in a program
	"program", // contains rule and query children
	"freeze",  // prevents backtracking from now on.
	"halt",    // pause program.
	"exit",    // exit program.
	"print",   // print message or node
	"load",    // load a program (rules only)
	"reset",   // restart - TODO
	"verbose", // toggle verbose flag
	"debug",   // toggle debug flag
	"rules",   // dump current rules
	"queries", // dump current queries/goals
	"goals",   // same - dump current queries/goals
	// "help",    // TODO
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

var errNotFound = fmt.Errorf("keyword not found")

func NewKeyword(kw string) (Keyword, error) {
	for _, k := range keywords {
		if k == kw {
			return Keyword{kw}, nil
		}
	}
	return Keyword{}, errNotFound
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

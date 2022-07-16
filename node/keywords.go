package node

import "fmt"

// contains reserved keywords
var keywords = []string{
	"rule",
	"query",
	"program",
	"slash",
	"halt",
	"print",
	"import",
	"reset",
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
	n.load = Keyword{"slash"}
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

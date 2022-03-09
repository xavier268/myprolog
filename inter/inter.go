package inter

import "fmt"

// Inter contains the global interpreter context.
type Inter struct {
	// symt maps symbols of variables and numbers to their (unique) pointer.
	// We can reuse the same object, because variables and numbers CANNOT HAVE CHILDREN ; they cannot be compound objects.
	// This is enforced by the parser.
	// However, they are not necessary "constant", ie "immutable" (Variables are NOT "contant")
	symt map[string]*Node
}

func NewInter() *Inter {
	return &Inter{
		symt: map[string]*Node{
			"_": {name: "_", constant: false, args: []*Node{}},
		},
	}
}

// nodeFor provides a node to represent the given text.
// If a number or a variable, then the same node is always returned.
// Otherwise, a new node is created.
func (i *Inter) nodeFor(text string) *Node {
	if text == "" {
		panic("name is required")
	}
	n := i.symt[text]
	if n == nil { // not yet in symbol table
		n = &Node{name: text}
		if isVariable(text) || isNumber(text) {
			i.symt[text] = n
		}

	}
	return n
}

// dumpSymt dumps symbol table.
func (i *Inter) dumpSymt() {

	fmt.Print("\n   Symbol       ")
	dumpHeader()
	for s, n := range i.symt {
		fmt.Printf("%10s --> ", s)
		n.DumpTree(false) // no header
	}
}

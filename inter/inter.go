package inter

import "fmt"

// Inter contains the global interpreter context.
type Inter struct {
	// symt maps symbols of variables and numbers to their (unique) pointer.
	// We can reuse the same object, because variables and numbers cannot have children ; they cannot be compound objects.
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
// If a number or a variable (having no children), then the same node is always returned.
// Otherwise, a new node is created.
func (i *Inter) nodeFor(text string) *Node {
	if text == "" {
		panic("name is required")
	}
	n := i.symt[text]
	if n == nil { // not yet in symbol table
		n = &Node{name: text}
		if isVariable(text) || isNumber(text) { // should it be in symbol table ?
			n.constant = true
			i.symt[text] = n
		}
	}
	return n
}

// dumpSymt dumps symbol table.
func (i *Inter) dumpSymt() {

	fmt.Println("\n   Symbol       pointer   const?    name  :  args ....")
	for s, n := range i.symt {
		fmt.Printf("%10s --> ", s)
		n.dump()
	}

}

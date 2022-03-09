package inter

import "fmt"

// Inter contains the global interpreter context.
type Inter struct {
	// symt maps symbols of variables and numbers to their (unique) pointer.
	// We can reuse the same object, because variables and numbers CANNOT HAVE CHILDREN ; they cannot be compound objects.
	// This is enforced by the parser.
	// However, symbols themslves are invariant, but they are not "constant", ie "immutable" (Variables are NOT "contant").
	// When a rule is rescoped, new symbols will be needed for those reprensting Variable.
	symt map[string]*Node
	// Uinique id, used to generate unique variable names.
	uid int
	// the rules node is names 'rules' and contains the rules ;-)
	rules *Node
}

func NewInter() *Inter {
	return &Inter{
		symt: map[string]*Node{
			"_":   {name: "_", constant: false, args: []*Node{}},  // predefined, match-all variable.
			"nil": {name: "nil", constant: true, args: []*Node{}}, // predefined. Cannot be a functor.
		},
	}
}

// Uid generate a unique suffix, such as "_125".
// Typically, to be added to a Symbol variable to create a new  distinct local variable.
func (i *Inter) Uid() string {
	i.uid++
	return fmt.Sprintf("#%d", i.uid)
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

func (in *Inter) Dump() {
	in.dumpRules()
	in.dumpSymt()
}

// dumpSymt dumps symbol table.
func (in *Inter) dumpSymt() {

	fmt.Println("\n------------ Symbole table --------------")
	fmt.Print("   Symbol       ")
	dumpHeader()
	for s, n := range in.symt {
		fmt.Printf("%10s --> ", s)
		n.DumpTree(false) // no header
	}
	fmt.Println("-----------------------------------------")

}

// dump the rules loaded
func (in *Inter) dumpRules() {
	fmt.Printf("\n --------------- %s ------------\n", in.rules.name)
	for i, r := range in.rules.args {
		fmt.Printf("%4d: ", i)
		fmt.Println(r.String())
	}
	fmt.Println("-----------------------------------")
}

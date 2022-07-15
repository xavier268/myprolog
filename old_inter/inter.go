package inter

import "fmt"

// Inter contains the global interpreter context.
type Inter struct {
	// Unique id, used to generate unique variable names when rescoping rules
	uid int
	// the rules node is named 'rules' and it contains the rules ;-)
	rules *Node
}

func NewInter() *Inter {
	return new(Inter).Reset()
}

// Reset the internals of the Interpreter, also resetting the rules.
func (in *Inter) Reset() *Inter {
	in.rules = nodeFor("rules")
	in.uid = 0
	return in
}

// Uid generate a unique suffix, such as "#125".
// Typically, to be added to a Symbol variable to create a new  distinct local variable.
func (i *Inter) Uid() string {
	i.uid++
	return fmt.Sprintf("#%d", i.uid)
}

// nodeFor provides a node to represent the given text.
func nodeFor(text string) *Node {
	if text == "" {
		panic("name is required")
	}
	return &Node{name: text}
}

func (in *Inter) Dump() {
	in.dumpRules()
}

// dump the rules loaded
func (in *Inter) dumpRules() {
	if in.rules == nil {
		fmt.Println("\n---------- NO RULES ------------")
		return
	}
	fmt.Printf("\n --------------- %s ------------\n", in.rules.name)
	for i, r := range in.rules.args {
		fmt.Printf("%4d: ", i)
		fmt.Println(r.String())
	}
	fmt.Println("-----------------------------------")
}

// --------------------- Rescoping nodes ------------------------------------

// Rescope recursiveley the provided node, by adding the provided suffix to each of its variable names.
// The _ Variable is NEVER rescoped, because it never binds.
func (in *Inter) Rescope(n *Node, suffix string) *Node {

	if n == nil || n.name == "_" {
		return n
	}

	// duplicate root node
	r := nodeFor(n.name)

	// if variable, adjust name
	if isVariable(n.name) {
		n.name = n.name + suffix
	}

	// duplicate children ...
	for _, a := range n.args {
		r.args = append(r.args, in.Rescope(a, suffix))
	}

	return r
}

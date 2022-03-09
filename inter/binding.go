package inter

import (
	"fmt"
)

type Binder interface {
	// Get the binded tree attached to the variable. If unbound, gets nil.
	// To unset, set to nil.
	Get(variable *Node) *Node
	// Set the binded tree attached to the symbol.
	Set(variable *Node, value *Node)

	// Push pushes a new binding frame, where all bindings are now registered,
	// combining with all the previous frames.
	Push()
	// Pop the current frame, restoring the bindings as they were before
	// PushFrame was called. Error means there are no more frames to pop.
	Pop() error
	// Slash() prevents further poping, freezing the frames as they are.
	// New frames can be pushed, though. Typically results from the / in prolog,
	// prenventing backtracing.
	Slash()
}

type binder struct {
	frames []map[*Node]*Node
	top    int // frames[top] is the current frame.
	bottom int // lowest legally accessible frame
}

// Compiler check
var _ Binder = new(binder)

func NewBinder() Binder {
	b := new(binder)
	b.Push()
	return b
}

func (b *binder) Push() {
	m := make(map[*Node]*Node, 4)
	b.frames = append(b.frames, m)
	b.top = len(b.frames) - 1
}

// EOS is End of Stack error.
var ErrBinder = fmt.Errorf("cannot pop, binder frame stack already at bottom")

func (b *binder) Pop() error {
	if b.top <= b.bottom {
		return ErrBinder
	}
	b.frames = b.frames[:b.top]
	b.top = len(b.frames) - 1
	return nil
}

func (b *binder) Set(v *Node, data *Node) {
	b.frames[b.top][v] = data
}

func (b *binder) Get(v *Node) *Node {
	for i := b.top; i >= 0; i-- {
		r := b.frames[i][v]
		if r != nil {
			return r
		}
	}
	return nil
}

func (b *binder) Slash() {
	b.bottom = b.top
}

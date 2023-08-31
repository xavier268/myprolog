package solver

import (
	"fmt"

	"github.com/xavier268/myprolog/mytest"
)

const MIN_DEPTH = 30 // minimum depth control for state/uids

var (
	depthStats      int // count of truncated searches
	ErrDepthControl = fmt.Errorf("maximum allowed nesting depth reached - truncating search tree")
)

// Check if maxdepth is reached.
func (st State) CheckDepth() error {
	// Prefer usuing this rather than doing a direct compare,
	// since this will allow to detect too frequent hits and dynamically increase maxdepth if needed.

	md := MIN_DEPTH + st.Rules.Count()*3
	if st.Uid >= md {
		depthStats += 1
		err := ErrDepthControl
		if depthStats < 10 ||
			(depthStats < 100 && depthStats%10 == 0) ||
			(depthStats < 1000 && depthStats%100 == 0) ||
			(depthStats < 10000 && depthStats%1000 == 0) { // limit volume and frequency of information displayed !
			err = fmt.Errorf("reached max allowed depth (%d) - forcing backtracking (%d times)", md, depthStats)
			fmt.Printf("%sWARNING : %v%s\n", mytest.CYAN, err, mytest.RESET)
		}
		return err // in large volume, don't spend CPU to refpormat an error each time !
	}
	return nil
}

package solver

import (
	"fmt"

	"github.com/xavier268/myprolog/mytest"
)

const MIN_DEPTH = 30 // minimum depth control for state/uids

// Check if maxdepth is reached.

func (st State) CheckDepth() error {
	// Prefer usuing this rather than doing a direct compare,
	// since this will allow to detect too frequent hits and dynamically increase maxdepth if needed.

	// naive heuristic, to be improved later...
	md := MIN_DEPTH + st.Rules.Count()*3
	if st.Uid >= md {
		err := fmt.Errorf("reached max allowed depth (%d) - forcing backtracking", md)
		fmt.Printf("%sWARNING : %v%s\n", mytest.CYAN, err, mytest.RESET)
		return err
	}
	return nil
}

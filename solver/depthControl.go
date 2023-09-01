package solver

import (
	"fmt"

	"github.com/xavier268/myprolog/mytest"
)

const MIN_DEPTH = 30 // minimum depth control for state/uids

var ErrDepthControl = fmt.Errorf("maximum allowed nesting depth reached - truncating search tree")

// Check if maxdepth is reached.
func (st State) CheckDepth() error {
	// Prefer usuing this rather than doing a direct compare,
	// since this will allow to detect too frequent hits and dynamically increase maxdepth if needed.
	db := st.Session
	md := MIN_DEPTH + db.CountRules()*3
	if st.Uid >= md {
		db.depthStats += 1
		err := ErrDepthControl
		if db.depthStats < 10 ||
			(db.depthStats < 100 && db.depthStats%10 == 0) ||
			(db.depthStats < 1000 && db.depthStats%100 == 0) ||
			(db.depthStats < 10000 && db.depthStats%1000 == 0) { // limit volume and frequency of information displayed !
			err = fmt.Errorf("reached max allowed depth (%d) - forcing backtracking (%d times)", md, db.depthStats)
			fmt.Printf("%sWARNING : %v%s\n", mytest.CYAN, err, mytest.RESET)
		}
		return err // in large volume, don't spend CPU to refpormat an error each time !
	}
	return nil
}

package solver

import (
	"fmt"

	"github.com/xavier268/myprolog/mytest"
)

const MIN_DEPTH = 30  // minimum depth control for state/uids
const MAX_DEPTH = 300 // maximum depth control for state/uids

var ErrDepthControl = fmt.Errorf("maximum allowed nesting depth reached - truncating search tree")

// Check if maxdepth is reached.
func (st State) CheckDepth() error {
	// Prefer using this rather than doing a direct compare,
	// since this will allow to detect too frequent hits and dynamically increase maxdepth if needed.
	sess := st.Session
	if st.Uid >= st.Session.depthLimit {
		sess.depthCount += 1
		err := ErrDepthControl
		if sess.depthCount < 10 ||
			(sess.depthCount < 100 && sess.depthCount%10 == 0) ||
			(sess.depthCount < 1000 && sess.depthCount%100 == 0) ||
			(sess.depthCount < 10000 && sess.depthCount%1000 == 0) { // limit volume and frequency of information displayed !
			err = fmt.Errorf("reached max allowed depth (%d) - forcing backtracking (%d times)", sess.depthLimit, sess.depthCount)
			fmt.Printf("%sWARNING : %v%s\n", mytest.CYAN, err, mytest.RESET)
		}
		return err // in large volume, don't spend CPU to refpormat an error each time !
	}
	return nil
}

// Force the limit beyond which search trees will be truncated.
// Use at your own risks !
func (s *Session) ForceDepthControl(limit int) {
	s.depthLimit = limit
	s.depthCount = 0
	fmt.Println("\nDepth control set to", limit, "steps")
}

// heuristic to adjust depth control
// called when adding new rules.
// multiple calls will increase the limit ...
func (s *Session) AdjustDepthControl() {
	// naive heuristic , based on number of rules.
	md := s.depthLimit + s.CountRules()
	if md > MAX_DEPTH {
		md = MAX_DEPTH
	}
	s.ForceDepthControl(md)
}

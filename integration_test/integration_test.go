package integration

import (
	"testing"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/repl"
)

func TestIntegration(t *testing.T) {
	config.FlagDebug = false
	config.FlagVerbose = false
	repl.RUN("test.pl")
}

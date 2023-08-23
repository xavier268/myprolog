package pcontext

import (
	"fmt"
	"strings"

	"github.com/xavier268/myprolog/config"
)

type PVariableHelp struct {
	short string // short description
	long  string // long description
}

var PVariables = map[string]PVariableHelp{
	"PI": {"represents the PI number.", ""},
}

// Return help on pseudo variable. No-op if not a pseudo variable.
func HelpPVariable(t string) (string, error) {

	var sb strings.Builder

	th, ok := PVariables[t]
	if !ok {
		return "", fmt.Errorf("no a pseudovariable")
	}

	fmt.Fprintf(&sb, "%s is a pseudo-variable that %s\n", t, th.short)
	lines := config.CutString(th.long, config.TEXT_WIDTH)
	for _, line := range lines {
		if len(line) > 0 {
			fmt.Fprintln(&sb, config.TEXT_PREFIX, line)
		}
	}
	return sb.String(), nil
}

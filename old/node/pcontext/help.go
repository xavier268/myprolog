package pcontext

import (
	"fmt"
	"sort"
	"strings"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/node"
)

type TopicHelp struct {
	short, long string
}

var Topics = map[string]TopicHelp{
	"keyword":        {"is a reserved word.", "Keywords are executed when they are goals. They cannot unify with anything, preventing accidental redefinition."},
	"pseudovariable": {"is a predefined variable.", "Pseudovariables are defined and bound at the time they are met. Examples of pseudo variables are PI, or RULES."},
}

// Return help on topic. No-op if not a topic.
func HelpTopic(t string) (string, error) {

	var sb strings.Builder

	th, ok := Topics[t]
	if !ok {
		return "", fmt.Errorf("no a topic")
	}

	fmt.Fprintf(&sb, "A %s %s\n", t, th.short)
	lines := config.CutString(th.long, config.TEXT_WIDTH)
	for _, line := range lines {
		if len(line) > 0 {
			fmt.Fprintln(&sb, config.TEXT_PREFIX, line)
		}
	}
	return sb.String(), nil
}

func HelpString(s string) string {
	var h string
	var err error

	h, err = HelpTopic(s)
	if err == nil {
		return h
	}
	h, err = HelpPVariable(s)
	if err == nil {
		return h
	}
	h, err = node.HelpKeyword(s)
	if err == nil {
		return h
	}
	return ""
}

func (pc *PContext) doHelp(children ...*node.Node) {

	if len(children) != 0 {
		for _, c := range children {
			s := strings.TrimSpace(c.String())
			h := HelpString(s)
			if len(h) != 0 {
				fmt.Println(h)
			}
		}
		return
	}

	// otherwise, print full help.
	hh := make([]string, 0, len(Topics)+len(PVariables)+len(node.Keywords))
	for h := range Topics {
		hh = append(hh, h)
	}
	for h := range PVariables {
		hh = append(hh, h)
	}
	for h := range node.Keywords {
		hh = append(hh, h)
	}
	sort.Strings(hh)

	// Actual printing
	fmt.Println("-------- Full Help ------------")
	for _, h := range hh {
		fmt.Println(HelpString(h))
	}
	fmt.Println("-------------------------------")

}

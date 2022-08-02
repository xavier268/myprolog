// Package config contains shared configuration package flags and constants.
package config

import (
	"fmt"
	"strings"
)

const VERSION = "v0.5.1"
const WELCOME = "Prolog interpreter in Go"
const GITHUBNAME = "https://github.com/xavier268/myprolog"
const COPYRIGHT = "(c) 2022 Xavier Gandillot (aka xavier268)"

const TEXT_WIDTH = 40           // default width of paragraphs
const TEXT_PREFIX = "         " // default prefix for indentation

var FlagVerbose = false
var FlagDebug = false

func PrintFullWelcome() {
	fmt.Printf(
		"--------------------------------------------\n"+
			"%s - Version %s\n"+
			"- see %s\n"+
			"%s\n"+
			"--------------------------------------------\n",
		WELCOME, VERSION, GITHUBNAME, COPYRIGHT)
}

// CutString is a utility to intelligently cut the string in pieces with a max width.
func CutString(s string, width int) (lines []string) {

	words := strings.Split(s, " ")
	line := ""
	for _, w := range words {
		w = strings.TrimSpace(w)
		if len(w)+len(line) <= width {
			line = line + " " + w
		} else {
			lines = append(lines, strings.TrimSpace(line))
			line = w
		}
	}
	line = strings.TrimSpace(line)
	if len(line) != 0 {
		lines = append(lines, line)
	}

	return lines
}

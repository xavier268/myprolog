// Package config contains shared configuration package flags and constants.
package config

import "fmt"

const VERSION = "v0.3.4"
const WELCOME = "Prolog interpreter in Go"
const GITHUBNAME = "https://github.com/xavier268/myprolog"
const COPYRIGHT = "(c) 2022 Xavier Gandillot (aka xavier268)"

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

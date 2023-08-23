// Package config contains shared configuration package flags and constants.
package config

import (
	"fmt"
	"text/scanner"
)

const (
	FLOATFORMAT = "%0.2f"
	INTFORMAT   = "%d"
)

type Config struct {
	Version     string // version
	FlagVerbose bool
	FlagDebug   bool
	TextWidth   int    // default width of paragraphs
	TextPrefix  string // default prefix for indentation
	Welcome     string // welcome message
	GitHubName  string // link to github repository
	CopyRight   string // copyright message

	ReservedWords []struct { // list of reserved words, and their arity and meaning. Arity of -1 means no constraint.
		key     string
		arity   int
		meaning string
	}
	ScannerMode uint // default scanner mode when scanning tokens
}

func New() *Config {
	c := new(Config)
	c.TextWidth = 40
	c.TextPrefix = "         "
	c.Version = "v0.6.0"
	c.Welcome = "Prolog interpreter in Go"
	c.GitHubName = "https://github.com/xavier268/myprolog"
	c.CopyRight = "(c) Xavier Gandillot (aka xavier268) 2022,2023"
	c.ReservedWords = []struct {
		key     string
		arity   int
		meaning string
	}{ // list of reserved words, and their arity and meaning. Arity of -1 means no constraint.
		{"dot", 2, "canonical functor for lists"},
		{"-:", -1, "-: operator for rule definition"},
		{"?-", -1, "-: operator for query definition"},
		{".", 0, ". final marker for rules, facts or queries"},
		{"true", 0, "true literal"},
		{"false", 0, "false literal"},
	}

	c.ScannerMode = scanner.ScanInts | scanner.ScanIdents | scanner.ScanStrings | scanner.ScanComments | scanner.SkipComments

	return c
}

func (c *Config) PrintFullWelcome() {
	fmt.Printf(
		"--------------------------------------------\n"+
			"%s - Version %s\n"+
			"- see %s\n"+
			"%s\n"+
			"--------------------------------------------\n",
		c.Welcome, c.Version, c.GitHubName, c.CopyRight)
}
